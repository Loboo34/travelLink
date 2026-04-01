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

type FlightBookingService struct {
	flightRepo  *repository.FlightSearchRepo
	bookingRepo *repository.FlightBookingRepo
	payment     payment.Provider
}

func NewFlightBookingService(flightRepo *repository.FlightSearchRepo,
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
		return nil, &model.ValidationError{Message: err.Error()}
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
		NumberOfSeats:    seats,
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
			utils.Logger.Error("failed to release reservation after booking creation failure",
				zap.String("flighID", req.FlightID.Hex()),
				zap.Error(releaseErr))
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
		utils.Logger.Warn("payment failed", zap.String("bookingID", bookingID.Hex()),
			zap.Error(err))

		if releaseErr := s.bookingRepo.ReleaseReservation(ctx, req.FlightID, seats); releaseErr != nil {
			utils.Logger.Error("failed to release reservation after payment failure",
				zap.String("bookingID", bookingID.Hex()),
				zap.Error(releaseErr),
			)
		}

		_ = s.bookingRepo.UpdateBooking(ctx, bookingID, model.BookingStatusFailed, nil)

		return nil, &model.PaymentError{Message: "payment processing failed"}
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
		utils.Logger.Error("payment succeeded but booking confirmation failed", zap.String("bookingID", bookingID.Hex()),
			zap.Error(err))

		return nil, fmt.Errorf("confirming booking: %w", err)

	}

	return &FlightBookingResults{
		BookingID:  bookingID,
		Status:     model.BookingStatusConfirmed,
		AmountPaid: offer.PriceTotal,
		Currency:   req.Currency,
	}, nil

}

type CancellationResult struct {
	BookingID    primitive.ObjectID  `json:"bookingID"`
	Status       model.BookingStatus `json:"status"`
	RefundStatus model.RefundStatus  `json:"refundStatus"`
}

func (s *FlightBookingService) Cancel(ctx context.Context, userID primitive.ObjectID, req model.Cancellation) (*CancellationResult, error) {
	if err := req.Validate(); err != nil {
		return nil, &model.ValidationError{Message: err.Error()}
	}

	booking, err := s.bookingRepo.GetBooking(ctx, req.BookingID)
	if err != nil {
		return nil, fmt.Errorf("error getting flight: %w", err)
	}

	if booking.UserID != userID {
		return nil, &model.AuthError{Message: "unauthorized to cancel booking"}
	}

	if booking.Status != model.BookingStatusConfirmed {
		return nil, &model.ConflictError{Message: fmt.Sprintf("cannot cancel booking with status %q", booking.Status)}
	}

	if err := s.bookingRepo.Cancel(ctx, req.BookingID, req.Reason); err != nil {
		return nil, fmt.Errorf("error cancellign flight")
	}

	if releaseErr := s.bookingRepo.ReleaseReservation(ctx, booking.OfferReference, booking.NumberOfSeats); releaseErr != nil {
		utils.Logger.Error("failed to release seats", zap.String("bookingID", req.BookingID.Hex()), zap.String("flighID", booking.OfferReference.Hex()), zap.Error(releaseErr))
	}

	//payment
	if booking.Payment.PaymentReference != "" {
		if err := s.payment.Refund(
			ctx, booking.Payment.PaymentReference, booking.AmountPaid,
		); err != nil {
			utils.Logger.Error("refund failed after cancellation",
                zap.String("bookingID", req.BookingID.Hex()),
                zap.String("paymentReference", booking.Payment.PaymentReference),
                zap.Error(err),
            )

		}
	}
	return &CancellationResult{
		BookingID:    req.BookingID,
		Status:       model.BookingStatusCanceled,
		RefundStatus: model.RefundStatusPending,
	}, nil

}
