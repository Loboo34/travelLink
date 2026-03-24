package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	model "github.com/Loboo34/travel/models"
	"github.com/Loboo34/travel/payment"
	"github.com/Loboo34/travel/repository"
	"github.com/Loboo34/travel/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PackageBookingService struct {
	flightRepo        *repository.FlightBookingRepo
	accommodationRepo *repository.AccommodationBookingRepo
	accomRepo         *repository.AccommodationRepo
	activityRepo      *repository.ActivityBookingRepo
	packageRepo       *repository.PackageBookingRepo
	payment           payment.Provider
}

func NewPackageBookingRepo(flightRepo *repository.FlightBookingRepo,
	accommodationRepo *repository.AccommodationBookingRepo,
	accomRepo *repository.AccommodationRepo,
	activityRepo *repository.ActivityBookingRepo,
	packageRepo *repository.PackageBookingRepo,
	payment payment.Provider) *PackageBookingService {
	return &PackageBookingService{
		flightRepo:        flightRepo,
		accommodationRepo: accommodationRepo,
		accomRepo:         accomRepo,
		activityRepo:      activityRepo,
		packageRepo:       packageRepo,
		payment:           payment,
	}
}

type PackageBookingResults struct {
	BookingID  primitive.ObjectID  `json:"bookingID"`
	Status     model.BookingStatus `json:"status"`
	AmountPaid int64               `json:"amountPaid"`
	Currency   string              `json:"currency"`
}

func (s *PackageBookingService) Book(ctx context.Context, userID primitive.ObjectID, req model.PackageBookingRequest) (*PackageBookingResults, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation error")
	}

	_, err := s.packageRepo.GetPackage(ctx, req.PackageID)
	if err != nil {
		return nil, fmt.Errorf("package does not exist: %w", err)
	}

	travelers := len(req.TravelersDetails)
	if err := s.packageRepo.ReserveSlot(ctx, req.PackageID, travelers); err != nil {
		return nil, fmt.Errorf("error reserving %w", err)
	}

	type componentResult struct {
		selection model.ComponentSelection
		price     int64
		err       error
	}

	resultCh := make(chan componentResult, len(req.ComponentSelection))

	for _, sel := range req.ComponentSelection {
		sel := sel
		go func() {
			price, err := s.reserveComponent(ctx, req, sel)
			resultCh <- componentResult{sel, price, err}
		}()
	}

	var successfulReservations []model.ComponentSelection
	var bookedComponents []model.BookedComponent
	var totalPrice int64
	var reservationErr error

	for range req.ComponentSelection {
		r := <-resultCh
		if r.err != nil {
			reservationErr = r.err
		} else {
			successfulReservations = append(successfulReservations, r.selection)
			totalPrice += r.price
			bookedComponents = append(bookedComponents, model.BookedComponent{
				ComponentType: r.selection.ComponentType,
				ReferenceID:   r.selection.ReferenceID,
				Price:         r.price,
			})
		}
	}

	if reservationErr != nil {
		s.rollBack(ctx, req.PackageID, successfulReservations, req, travelers)
		if errors.Is(reservationErr, reservationErr) {
			return nil, fmt.Errorf("reserving package: %w", reservationErr)
		}
	}

	bookingID := primitive.NewObjectID()
	now := time.Now()

	booking := &model.PackageBooking{
		ID:               bookingID,
		UserID:           userID,
		PackageID:        req.PackageID,
		UserDetails:      req.TravelersDetails,
		BookedComponents: bookedComponents,
		Status:           model.BookingStatusPending,
		RefundStatus:     model.RefundStatusNone,
		Payment: model.Payment{
			PaymentMethod: req.PaymentMethod,
			TotalAmount:   totalPrice,
			Currency:      req.Currency,
			Status:        model.PaymentPending,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.packageRepo.CreateBooking(ctx, booking); err != nil {
		s.rollBack(ctx, req.PackageID, successfulReservations, req, travelers)
		return nil, fmt.Errorf("creating package booking: %w", err)
	}

	//payment
	paymentResult, err := s.payment.Charge(ctx, payment.ChargeRequest{
		Amount:   totalPrice,
		Currency: req.Currency,
		Method:   req.PaymentMethod,
		UserID:   userID.Hex(),
		Metadata: map[string]string{
			"bookingID": bookingID.Hex(),
			"type":      "package",
		},
	})
	if err != nil {
		utils.Logger.Warn("payment failed")
		s.rollBack(ctx, req.PackageID, successfulReservations, req, travelers)
		_ = s.packageRepo.UpdateBookingStatus(ctx, bookingID, model.BookingStatusFailed, nil)
		return nil, fmt.Errorf("error releasing reservation: %w", err)
	}

	confirmedPayment := &model.Payment{
		PaymentMethod:    req.PaymentMethod,
		TotalAmount:      totalPrice,
		Currency:         req.Currency,
		Status:           model.PaymentCompleted,
		PaymentReference: paymentResult.Reference,
		PaidAt:           &now,
	}

	if err := s.packageRepo.UpdateBookingStatus(ctx, bookingID, model.BookingStatusCompleted, confirmedPayment); err != nil {
		utils.Logger.Error("payment confirem but failed to confirm booking")
		return nil, fmt.Errorf("confirming booking: %w", err)
	}

	return &PackageBookingResults{
		BookingID:  bookingID,
		Status:     model.BookingStatusConfirmed,
		AmountPaid: totalPrice,
		Currency:   req.Currency,
	}, nil

}

func (s *PackageBookingService) rollBack(ctx context.Context, packageID primitive.ObjectID, selections []model.ComponentSelection, req model.PackageBookingRequest, travelers int) {

	if err := s.packageRepo.ReleaseSlot(ctx, packageID, travelers); err != nil {
		utils.Logger.Error("failed to release package")
	}

	for _, sel := range selections {
		switch sel.ComponentType {
		case model.ComponentFlight:

			if err := s.flightRepo.ReleaseReservation(ctx, sel.ReferenceID, travelers); err != nil {
				utils.Logger.Error("")
			}

		case model.ComponentAccommodation:
			if sel.RoomTypeID == nil || sel.CheckIn == nil || sel.CheckOut == nil {
				continue
			}
			rooms := sel.Rooms
			if rooms < 0 {
				rooms = 1
			}
			if err := s.accommodationRepo.ReleaseReservation(ctx, sel.ReferenceID, *sel.RoomTypeID, *sel.CheckIn, *sel.CheckOut, rooms); err != nil {
				utils.Logger.Error("")
			}

		case model.ComponentActivity:
			if sel.TimeslotID == nil {
				continue
			}
			participants := sel.Participants
			if participants < 0 {
				participants = travelers
			}
			if err := s.activityRepo.ReleaseReservation(ctx, *sel.TimeslotID, participants); err != nil {
				utils.Logger.Error("")
			}
		}
	}
}

func (s *PackageBookingService) reserveComponent(ctx context.Context, req model.PackageBookingRequest, sel model.ComponentSelection) (int64, error) {
	travelers := len(req.TravelersDetails)

	switch sel.ComponentType {
	case model.ComponentFlight:
		offer, err := s.flightRepo.CheckAndReserv(ctx, sel.ReferenceID, travelers)
		if err != nil {
			return 0, err
		}

		return offer.PriceTotal, nil
	case model.ComponentAccommodation:
		if sel.RoomTypeID == nil || sel.CheckIn == nil || sel.CheckOut == nil {
			return 0, fmt.Errorf("accommodation selection missing roomTypeID, checkIn, checkOut")
		}

		rooms := sel.Rooms
		if rooms < 0 {
			rooms = 1
		}
		if err := s.accommodationRepo.CheckAndReserv(ctx, sel.ReferenceID, *sel.RoomTypeID, *sel.CheckIn, *sel.CheckOut, rooms); err != nil {
			return 0, err
		}

		totalPrice, err := s.accomRepo.GetTotalPrice(ctx, sel.ReferenceID, *sel.RoomTypeID, *sel.CheckIn, *sel.CheckOut)
		if err != nil {
			return 0, err
		}

		return totalPrice, nil

	case model.ComponentActivity:
		if sel.TimeslotID == nil {
			return 0, fmt.Errorf("activity selection missing timeslot ID")
		}
		participants := sel.Participants
		if participants < 0 {
			participants = 1
		}
		slot, err := s.activityRepo.CheckAndReserv(ctx, sel.ReferenceID, *sel.TimeslotID, participants)
		if err != nil {
			return 0, err
		}

		totalPrice := slot.PricePerPerson * int64(sel.Participants)

		return totalPrice, nil
	default:
		return 0, fmt.Errorf("invalid component")
	}

}
