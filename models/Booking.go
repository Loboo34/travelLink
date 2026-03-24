package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FlightBooking struct {
	ID                 primitive.ObjectID `bson:"_id" json:"id"`
	UserID             primitive.ObjectID `bson:"userID" json:"userID"`
	OfferReference     primitive.ObjectID `bson:"offerReference" json:"offerReference"`
	PassengerDetails   []UserDetails      `bson:"passengerDetails" json:"passengerDetails"`
	Status             BookingStatus      `bson:"status" json:"status"`
	SelectedSeat       []string           `bson:"selectedSeat" json:"selectedSeat"`
	NumberOfSeats      uint               `bson:"numberOfSeats" json:"numberOfSeats"`
	AmountPaid         int64              `bson:"amountPaid" json:"amountPaid"`
	Payment            Payment            `bson:"payment" json:"payment"`
	CancellationDate   *time.Time         `bson:"cancellationDate" json:"cancellationDate"`
	CancellationReason string             `bson:"cancellationReason" json:"cancellationReason"`
	RefundStatus       RefundStatus       `bson:"refundStatus" json:"refundStatus"`
	CreatedAt          time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt          time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type AccommodationBooking struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID             primitive.ObjectID `bson:"userID" json:"userID"`
	AccommodationID    primitive.ObjectID `bson:"accommodationID" json:"accommodationID"`
	UserDetails        []UserDetails      `bson:"userDetails" json:"userDetails"`
	CheckIn            time.Time          `bson:"checkin" json:"checkin"`
	Checkout           time.Time          `bson:"checkout" json:"checkout"`
	Nights             int                `bson:"nights" json:"nights"`
	RoomTypeID         primitive.ObjectID `bson:"roomTypeID" json:"roomTypeID"`
	AmountPaid         int64              `bson:"amountPaid" json:"amountPaid"`
	Payment            Payment            `bson:"payment" json:"payment"`
	CancellationDate   *time.Time         `bson:"cancellationDate" json:"cancellationDate"`
	CancellationReason string             `bson:"cancellationReason" json:"cancellationReason"`
	RefundStatus       RefundStatus       `bson:"refundStatus" json:"refundStatus"`
	Status             BookingStatus      `bson:"status" json:"status"`
	CreatedAt          time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt          time.Time          `bson:"updateAt" json:"updateAt"`
}

type ActivityBooking struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID       primitive.ObjectID `bson:"userID" json:"userID"`
	ActivityID   primitive.ObjectID `bson:"activityID" json:"activityID"`
	TimeSlotID   primitive.ObjectID `bson:"timeslotID" json:"timeslotID"`
	UserDetails  []UserDetails      `bson:"userDetails" json:"userDetails"`
	Participants int                `bson:"participants" json:"participants"`
	//StartTime          time.Time          `bson:"startTime" json:"startTime"`
	AmountPaid         int64         `bson:"amountPaid" json:"amountPaid"`
	Payment            Payment       `bson:"payment" json:"payment"`
	CancellationDate   *time.Time    `bson:"cancellationDate" json:"cancellationDate"`
	CancellationReason string        `bson:"cancellationReason" json:"cancellationReason"`
	RefundStatus       RefundStatus  `bson:"refundStatus" json:"refundStatus"`
	Status             BookingStatus `bson:"status" json:"status"`
	CreatedAt          time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt          time.Time     `bson:"updateAt" json:"updateAt"`
}

type BookedComponent struct {
	ComponentType Component          `bson:"componentType" json:"componentType"`
	ReferenceID   primitive.ObjectID `bson:"referenceID" json:"referenceID"`
	Price         int64              `bson:"price" json:"price"`
}
type PackageBooking struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID             primitive.ObjectID `bson:"userID" json:"userID"`
	PackageID          primitive.ObjectID `bson:"packageID" json:"packageID"`
	UserDetails        []UserDetails      `bson:"userDetails" json:"userDetails"`
	BookedComponents   []BookedComponent  `bson:"bookedComponents" json:"bookedComponents"`
	AmountPaid         int64              `bson:"amountPaid" json:"amountPaid"`
	Payment            Payment            `bson:"payment" json:"payment"`
	CancellationDate   *time.Time         `bson:"cancellationDate" json:"cancellationDate"`
	CancellationReason string             `bson:"cancellationReason" json:"cancellationReason"`
	RefundStatus       RefundStatus       `bson:"refundStatus" json:"refundStatus"`
	Status             BookingStatus      `bson:"status" json:"status"`
	CreatedAt          time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt          time.Time          `bson:"updateAt" json:"updateAt"`
}

type Itinerary struct { //creating personal itinarary
	ID                   primitive.ObjectID   `bson:"-id,omitempty" json:"id"`
	UserID               primitive.ObjectID   `bson:"userID" json:"userID"`
	Title                string               `bson:"title" json:"title"`
	Description          string               `bson:"description" json:"description"`
	StartDate            time.Time            `bson:"startDate" json:"startDate"`
	EndDate              time.Time            `bson:"endDate" json:"endDate"`
	Destination          []string             `bson:"destination" json:"destination"`
	Status               BookingStatus        `bson:"status" json:"status"`
	IsPublic             bool                 `bson:"isPublic" json:"isPublic"`
	FlightBooking        []primitive.ObjectID `bson:"flightBooking" json:"flightBooking"`
	AccommodationBooking []primitive.ObjectID `bson:"accommodationBooking" json:"accommodationBooking"`
	ActivityrBooking     []primitive.ObjectID `bson:"activityBooking" json:"activityBooking"`
	DailyPlan            []DayPlan            `bson:"dailyPlan" json:"dailyPlan"`
	TotalBudget          int64                `bson:"budget" json:"budget"`
	AmountSpent          int64                `bson:"amountSpent" json:"amountSpent"`
	CreatedAt            time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt            time.Time            `bson:"updatedAt" json:"updatedAt"`
}

type UserDetails struct {
	UserID      primitive.ObjectID `bson:"userID" json:"userID"`
	FirstName   string             `bson:"firstName" json:"firstName"`
	LastName    string             `bson:"lastName" json:"lastName"`
	DateOfBirth time.Time          `bson:"dateOfBirth" json:"dateOfBirth"`
	UserPhone   string             `bson:"userPhone" json:"userPhone"`
	Gender      string             `bson:"gender" json:"gender"`
	Nationality string             `bson:"nationality" json:"nationality"`
	Passport    string             `bson:"passport" json:"passport"`
}

type Payment struct {
	PaymentMethod    PaymentMethod `bson:"paymentMethod" json:"paymentMethod"`
	TotalAmount      int64         `bson:"totalAmount" json:"totalAmount"`
	Currency         string        `bson:"currency" json:"currency"`
	Status           PaymentStatus `bson:"status" json:"status"`
	PaidAt           *time.Time    `bson:"paidAt" json:"paidAt"`
	PaymentReference string        `bson:"paymentReference" json:"paymentReference"`
}

type DayPlan struct {
	Day           int                  `bson:"day" json:"day"`
	Date          time.Time            `bson:"date" json:"date"`
	Location      string               `bson:"location" json:"location"`
	Activities    []primitive.ObjectID `bson:"activities" json:"activities"`
	Accommodation primitive.ObjectID   `bson:"accommodation" json:"accommodation"`
	Notes         string               `bson:"notes" json:"notes"`
}

type PaymentStatus string

const (
	PaymentCompleted PaymentStatus = "completed"
	PaymentPending   PaymentStatus = "pending"
	PaymentFailed    PaymentStatus = "failed"
	PaymentRefund    PaymentStatus = "refund"
)

type PaymentMethod string

const (
	PaymentMethodCard  PaymentMethod = "card"
	PaymentMethodMpesa PaymentMethod = "mpesa"
	PaymentMethodBank  PaymentMethod = "bank"
)

type BookingStatus string

const (
	BookingStatusCompleted BookingStatus = "completed"
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusFailed    BookingStatus = "Failed"
	BookingStatusConfirmed BookingStatus = "confirmed"
)

type RefundStatus string

const (
	RefundStatusNone      RefundStatus = "none"
	RefundStatusPending   RefundStatus = "pending"
	RefundStatusProcessed RefundStatus = "processed"
	RefundStatusFailed    RefundStatus = "failed"
)
