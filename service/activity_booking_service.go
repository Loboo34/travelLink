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

type ActivityBookingService struct {
	activityRepo *repository.ActivityRepo
	bookingRepo  *repository.ActivityBookingRepo
	payment      payment.Provider
}

func NewActivityBookingService(activityRepo *repository.ActivityRepo,
	bookingRepo *repository.ActivityBookingRepo,
	payment payment.Provider) *ActivityBookingService {
	return &ActivityBookingService{
		activityRepo: activityRepo,
		bookingRepo:  bookingRepo,
		payment:      payment,
	}
}

type ActivityBookingResult struct {
	BookingID  primitive.ObjectID  `json:"bookingID"`
	Status     model.BookingStatus `json:"status"`
	AmountPaid int64               `json:"amountPaid"`
	Currency   string              `json:"currency"`
}

func (s *ActivityBookingService) Book(ctx context.Context, userID primitive.ObjectID, req model.ActivityBookingRequest) (*ActivityBookingResult, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation error")
	}

	slot, err := s.bookingRepo.CheckAndReserv(ctx, req.ActivityID, req.TimeSlotID, req.Participants)
	if err != nil {
		return nil, fmt.Errorf("error rserving: %w", err)
	}

	totalPrice := slot.PricePerPerson * int64(req.Participants)

	bookingID := primitive.NewObjectID()
	now := time.Now()

	booking := &model.ActivityBooking{
		ID:           bookingID,
		ActivityID:   req.ActivityID,
		TimeSlotID:   req.TimeSlotID,
		UserID:       userID,
		UserDetails:  req.ParticipantDetails,
		Participants: req.Participants,
		RefundStatus: model.RefundStatusNone,
		Status: model.BookingStatusPending,
		Payment: model.Payment{
			PaymentMethod: req.PaymentMethod,
			TotalAmount:   totalPrice,
			Currency:      req.Currency,
			Status:        model.PaymentPending,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.bookingRepo.CreateBooking(ctx, booking); err != nil {
		if releaseErr := s.bookingRepo.ReleaseReservation(ctx,  req.TimeSlotID, req.Participants); releaseErr != nil {
			return nil, fmt.Errorf("error releasing reservation after booking failiure")
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
			"type":      "activity",
		},
	})
	if err != nil {
		utils.Logger.Warn("payment failed")

		if releaseErr := s.bookingRepo.ReleaseReservation(ctx,  req.TimeSlotID, req.Participants); releaseErr != nil {
			utils.Logger.Error("failed to release reservation after payment failure")
		}

		_ = s.bookingRepo.UpdateBooking(ctx, bookingID, model.BookingStatusFailed, nil)

		return nil, fmt.Errorf("payment failed")
	}

	confirmPayment := &model.Payment{
		PaymentMethod:    req.PaymentMethod,
		TotalAmount:      totalPrice,
		Currency:         req.Currency,
		Status:           model.PaymentCompleted,
		PaymentReference: paymentResult.Reference,
		PaidAt:           &now,
	}

	if err := s.bookingRepo.UpdateBooking(ctx, bookingID, model.BookingStatusConfirmed, confirmPayment); err != nil {
		utils.Logger.Error("payment confired but failed to confirm booking")
		return nil, fmt.Errorf("confirming booking: %w", err)
	}

	return &ActivityBookingResult{
		BookingID:  bookingID,
		Status:     model.BookingStatusConfirmed,
		AmountPaid: totalPrice,
		Currency:   req.Currency,
	}, nil
}



func (s *ActivityBookingService) Cancel(ctx context.Context, userID primitive.ObjectID, req model.Cancellation) (*CancellationResult, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation error")
	}

	booking, err := s.bookingRepo.GetBooking(ctx, req.BookingID)
	if err != nil {
		return nil, fmt.Errorf("error getting activity: %w", err)
	}

	if booking.UserID != userID {
	return  nil, fmt.Errorf("")
	}

	if booking.Status != model.BookingStatusConfirmed{
		return nil, fmt.Errorf("")
	}

	if err := s.bookingRepo.CancelFlight(ctx, req.BookingID, req.Reason); err != nil {
		return nil, fmt.Errorf("error cancellign activity")
	}

	if releaseErr := s.bookingRepo.ReleaseReservation(ctx, booking.TimeSlotID, booking.Participants); releaseErr != nil {
		utils.Logger.Error("")
	}

	//payment
	 if booking.Payment.PaymentReference != "" {
        if err := s.payment.Refund(
            ctx, booking.Payment.PaymentReference, booking.AmountPaid,
        ); err != nil {
            utils.Logger.Error("")
           
        }
    }
	return &CancellationResult{
		BookingID:    req.BookingID,
		Status:       model.BookingStatusCanceled,
		RefundStatus: model.RefundStatusPending,
	}, nil

}
