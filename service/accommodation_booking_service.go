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
	"go.uber.org/zap"
)

type AccommodationBookingService struct {
	accommodationRepo *repository.AccommodationSearchRepo
	bookingRepo       *repository.AccommodationBookingRepo
	payment           payment.Provider
}

func NewAccommodationBookingService(accommodationRepo *repository.AccommodationSearchRepo,
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
		return nil, &model.ValidationError{Message: err.Error()}
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
		Rooms:           req.Rooms,
		RoomTypeID:      req.RoomTypeID,
		Payment: model.Payment{
			PaymentMethod: req.PaymentMethod,
			TotalAmount:   totalPrice,
			Currency:      req.Currency,
			Status:        model.PaymentPending,
		},
		Status:       model.BookingStatusPending,
		RefundStatus: model.RefundStatusNone,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.bookingRepo.CreateBooking(ctx, booking); err != nil {
		if releaseErr := s.bookingRepo.ReleaseReservation(ctx, req.AccommodationID, req.RoomTypeID, req.CheckIn, req.CheckOut, req.Rooms); releaseErr != nil {
			utils.Logger.Error("failed to relaease room reservation", zap.String("accomID", req.AccommodationID.Hex()), zap.Error(releaseErr))
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
		utils.Logger.Warn("payment failed", zap.String("bookingID", bookingID.Hex()),
			zap.Error(err))


		if releaseErr := s.bookingRepo.ReleaseReservation(ctx, req.AccommodationID, req.RoomTypeID, req.CheckIn, req.CheckOut, req.Rooms); releaseErr != nil {
			utils.Logger.Error("failed to release reservation", zap.String("bookingID", bookingID.Hex()), zap.Error(releaseErr))
		}

		_ = s.bookingRepo.UpdateBooking(ctx, bookingID, model.BookingStatusFailed, nil)

		return nil, &model.PaymentError{Message: "payment processing failed"}
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
		utils.Logger.Error("payment confired but failed to confirm booking", zap.String("bookingId", bookingID.Hex()), zap.Error(err))
		return nil, fmt.Errorf("confirming booking: %w", err)
	}

	return &AccommodationBookingResult{
		BookingID:  bookingID,
		Status:     model.BookingStatusConfirmed,
		AmountPaid: totalPrice,
		Currency:   req.Currency,
	}, nil

}

func (s *AccommodationBookingService) Cancel(ctx context.Context, userID primitive.ObjectID, req model.Cancellation) (*CancellationResult, error) {
	if err := req.Validate(); err != nil {
		return nil, &model.ValidationError{Message: err.Error()}
	}

	booking, err := s.bookingRepo.GetBooking(ctx, req.BookingID)
	if err != nil {
		return nil, fmt.Errorf("error getting accommodation: %w", err)
	}

	if booking.UserID != userID {
		return nil, &model.AuthError{Message: "unauthorized to cancel booking"}
	}

	if booking.Status != model.BookingStatusConfirmed {
		return nil, &model.ConflictError{Message: fmt.Sprintf("cannot cancel booking with status %q", booking.Status)}
	}

	if err := s.bookingRepo.Cancel(ctx, req.BookingID, req.Reason); err != nil {
		return nil, fmt.Errorf("error cancellign accommodation")
	}

	if releaseErr := s.bookingRepo.ReleaseReservation(ctx, booking.AccommodationID, booking.RoomTypeID, booking.CheckIn, booking.Checkout, booking.Rooms); releaseErr != nil {
		utils.Logger.Error("failed to release rooms", zap.String("bookingId", booking.ID.Hex()), zap.String("accommodationID", booking.AccommodationID.Hex()), zap.Error(releaseErr))
	}

	//payment
	if booking.Payment.PaymentReference != "" {
		if err := s.payment.Refund(
			ctx, booking.Payment.PaymentReference, booking.AmountPaid,
		); err != nil {
			utils.Logger.Error("refund failed after cancellation",   zap.String("bookingID", req.BookingID.Hex()),
                zap.String("paymentReference", booking.Payment.PaymentReference),
                zap.Error(err),)

		}
	}
	return &CancellationResult{
		BookingID:    req.BookingID,
		Status:       model.BookingStatusCanceled,
		RefundStatus: model.RefundStatusPending,
	}, nil

}
