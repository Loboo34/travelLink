package service

import (
	"context"
	"fmt"
	"time"

	model "github.com/Loboo34/travel/models"
	"github.com/Loboo34/travel/payment"
	"github.com/Loboo34/travel/repository"
	"github.com/Loboo34/travel/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccommodationBookingService struct {
	accommodationRepo *repository.AccommodationRepo
	bookingRepo       *repository.AccommodationBookingRepo
	payment           payment.Provider
}

func NewAccommodationBookingService(accommodationRepo *repository.AccommodationRepo,
	bookingRepo *repository.AccommodationBookingRepo,
	payment payment.Provider) *AccommodationBookingService {
	return &AccommodationBookingService{
		accommodationRepo: accommodationRepo,
		bookingRepo:       bookingRepo,
		payment:           payment,
	}
}

type AccommodationBookingResult struct {
	BookingID  primitive.ObjectID  `json:"bookingID"`
	Status     model.BookingStatus `json:"status"`
	AmountPaid int64               `json:"amountPaid"`
	Currency   string              `json:"currency"`
}

func (s *AccommodationBookingService) Book(ctx context.Context, userID primitive.ObjectID, req model.AccommodationBookingRequest) (*AccommodationBookingResult, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation error")
	}

	nights := int(req.CheckOut.Sub(req.CheckIn).Hours() / 24)

	totalPrice, err := s.accommodationRepo.GetTotalPrice(ctx, req.AccommodationID, req.RoomTypeID, req.CheckIn, req.CheckOut)
	    if err != nil {
        return nil, fmt.Errorf("fetching accommodation price: %w", err)
    }
    totalPrice = totalPrice * int64(req.Rooms)

	if err := s.bookingRepo.CheckAndReserv(ctx, req.AccommodationID, req.RoomTypeID, req.CheckIn, req.CheckOut, req.Rooms); err != nil {
		return nil, fmt.Errorf("error reserving: %w", err)
	}

	bookingID := primitive.NewObjectID()
	now := time.Now().UTC()

	booking := &model.AccommodationBooking{
		ID:              bookingID,
		UserID:          userID,
		AccommodationID: req.AccommodationID,
		UserDetails:     req.GuestDetails,
		CheckIn:         req.CheckIn,
		Checkout:        req.CheckOut,
		Nights:          nights,
		RoomTypeID:        req.RoomTypeID,
		Payment: model.Payment{
			PaymentMethod: req.PaymentMethod,
			TotalAmount:   totalPrice,
			Currency:      req.Currency,
			Status:        model.PaymentPending,
		},
		Status: model.BookingStatusPending,
		RefundStatus: model.RefundStatusNone,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.bookingRepo.CreateBooking(ctx, booking); err != nil {
		if releaseErr := s.bookingRepo.ReleaseReservation(ctx, req.AccommodationID, req.RoomTypeID, req.CheckIn, req.CheckOut, req.Rooms); releaseErr != nil {
			utils.Logger.Error("failed to relaease room reservation")
		}

		return nil, fmt.Errorf("error creating booking: %w", err)
	}


	//payment
	paymentResult, err := s.payment.Charge(ctx, payment.ChargeRequest{
		Amount:   totalPrice,
		Currency: req.Currency,
		Method:   req.PaymentMethod,
		UserID:   userID.Hex(),
		Metadata: map[string]string{
			"bookingID": bookingID.Hex(),
			"type":      "accommodation",
		},
	})
	if err != nil {
		utils.Logger.Warn("payment failed")
		if releaseErr := s.bookingRepo.ReleaseReservation(ctx, req.AccommodationID, req.RoomTypeID, req.CheckIn, req.CheckOut, req.Rooms); releaseErr != nil {
			utils.Logger.Error("failed to release reservation")
		}

		_ = s.bookingRepo.UpdateBooking(ctx, bookingID, model.BookingStatusFailed, nil)

		return nil, fmt.Errorf("payment failed: %w", err)
	}

	confirmedPayment := &model.Payment{
		PaymentMethod:    req.PaymentMethod,
		TotalAmount:      totalPrice,
		Currency:         req.Currency,
		Status:           model.PaymentCompleted,
		PaymentReference: paymentResult.Reference,
		PaidAt:           &now,
	}

	if err := s.bookingRepo.UpdateBooking(ctx, bookingID, model.BookingStatusConfirmed, confirmedPayment); err != nil {
		utils.Logger.Error("payment confired but failed to confirm booking")
		return nil, fmt.Errorf("confirming booking: %w", err)
	}

	return &AccommodationBookingResult{
		BookingID:  bookingID,
		Status:     model.BookingStatusConfirmed,
		AmountPaid: totalPrice,
		Currency:   req.Currency,
	}, nil

}
