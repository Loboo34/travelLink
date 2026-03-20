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

type FlightBookingService struct {
	flightRepo  *repository.FlightRepo
	bookingRepo *repository.FlightBookingRepo
	payment     payment.Provider
}

func NewFlightBookingService(flightRepo *repository.FlightRepo,
	bookingRepo *repository.FlightBookingRepo,
	payment payment.Provider) *FlightBookingService {
	return &FlightBookingService{
		flightRepo:  flightRepo,
		bookingRepo: bookingRepo,
		payment:     payment,
	}
}

type FlightBookingResults struct {
	BookingID  primitive.ObjectID  `json:"bookingID"`
	Status     model.BookingStatus `json:"status"`
	AmountPaid int64               `json:"amountPaid"`
	Currency   string              `json:"currency"`
}

func (s *FlightBookingService) Book(ctx context.Context, userID primitive.ObjectID, req model.FlightRequest) (*FlightBookingResults, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation error")
	}

	seats := len(req.PassengerDetails)

	offer, err := s.bookingRepo.CheckAndReserv(ctx, req.FlightID, seats)
	if err != nil {
		return nil, fmt.Errorf("Error reserving: %w", err)
	}

	bookingID := primitive.NewObjectID()
	now := time.Now()

	booking := &model.FlightBooking{
		ID:               bookingID,
		UserID:           userID,
		OfferReference:   req.FlightID,
		PassengerDetails: req.PassengerDetails,
		Status:           model.BookingStatusPending,
		SelectedSeat:     req.SelectedSeats,
		NumberOfSeats:    uint(seats),
		Payment: model.Payment{
			PaymentMethod: req.PaymentMethod,
			TotalAmount:   offer.PriceTotal,
			Currency:      req.Currency,
			Status:        model.PaymentPending,
		},
		RefundStatus: model.RefundStatusNone,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.bookingRepo.CreateBooking(ctx, booking); err != nil {
		if releaseErr := s.bookingRepo.ReleaseReservation(ctx, req.FlightID, seats); releaseErr != nil {
			utils.Logger.Error("failed to release reservation after booking creation failure")
		}

		return nil, fmt.Errorf("error creating booking: %w", err)
	}

	//payment

	paymentResult, err := s.payment.Charge(ctx, payment.ChargeRequest{
		Amount:   offer.PriceTotal,
		Currency: req.Currency,
		Method:   req.PaymentMethod,
		UserID:   userID.Hex(),
		Metadata: map[string]string{
			"bookingID": bookingID.Hex(),
			"type":      "flight",
		},
	})

	//if payment fails
	if err != nil {
		utils.Logger.Warn("payment failed")

		if releaseErr := s.bookingRepo.ReleaseReservation(ctx, req.FlightID, seats); releaseErr != nil {
			utils.Logger.Error("failed to release reservation after payment failure")
		}

		_ = s.bookingRepo.UpdateBooking(ctx, bookingID, model.BookingStatusFailed, nil)

		return nil, fmt.Errorf("payment failed")
	}

	confirmedPayment := &model.Payment{
		PaymentMethod:    req.PaymentMethod,
		TotalAmount:      offer.PriceTotal,
		Currency:         req.Currency,
		Status:           model.PaymentCompleted,
		PaymentReference: paymentResult.Reference,
		PaidAt:           &now,
	}

	if err := s.bookingRepo.UpdateBooking(ctx, req.FlightID, model.BookingStatusConfirmed, confirmedPayment); err != nil {
		utils.Logger.Error("payment succeeded but booking confirmation failed")

		return nil, fmt.Errorf("confirming booking: %w", err)

	}

	return &FlightBookingResults{
		BookingID:  bookingID,
		Status:     model.BookingStatusConfirmed,
		AmountPaid: offer.PriceTotal,
		Currency:   req.Currency,
	}, nil

}
