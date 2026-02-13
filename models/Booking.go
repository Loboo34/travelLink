package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FlightBooking struct {
	ID                 primitive.ObjectID `bson:"_id" json:"id"`
	PassengerDetails   []UserDetails      `bson:"passengerDetails" json:"passengerDetails"`
	OfferReference     string             `bson:"offerReference" json:"offerReference"`
	Status             string             `bson:"status" json:"status"`
	SelectedSeat       []string           `bson:"selectedSeat" json:"selectedSeat"`
	NumberOfSeats      uint               `bson:"numberOfSeats" json:"numberOfSeats"`
	AmountPaid         float64            `bson:"amountPaid" json:"amountPaid"`
	CreatedAt          time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt          time.Time          `bson:"updatedAt" json:"updatedAt"`
	Payment            Payment            `bson:"payment" json:"payment"`
	CancellationDate   time.Time          `bson:"cancellationDate" json:"cancellationDate"`
	CancellationReason string             `bson:"cancellationReason" json:"cancellationReason"`
	RefundStatus       string             `bson:"refundStatus" json:"refundStatus"`
}

type AccommodationBooking struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AccommodationID    primitive.ObjectID `bson:"accommodationID" json:"accommodationID"`
	UserDetails        []UserDetails      `bson:"userDetails" json:"userDetails"`
	CheckIn            time.Time          `bson:"checkin" json:"checkin"`
	Checkout           time.Time          `bson:"checkout" json:"checkout"`
	Nights             uint               `bson:"nights" json:"nights"`
	RoomType           string             `bson:"roomType" json:"roomType"`
	Status             string             `bson:"status" json:"status"`
	AmountPaid         float64            `bson:"amountPaid" json:"amountPaid"`
	CreatedAt          time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt          time.Time          `bson:"updateAt" json:"updateAt"`
	Payment            Payment            `bson:"payment" json:"payment"`
	CancellationDate   time.Time          `bson:"cancellationDate" json:"cancellationDate"`
	CancellationReason string             `bson:"cancellationReason" json:"cancellationReason"`
	RefundStatus       string             `bson:"refundStatus" json:"refundStatus"`
}

type ActivityBooking struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ActivityID         primitive.ObjectID `bson:"activityID" json:"activityID"`
	TimeSlotID         primitive.ObjectID `bson:"timeslotID" json:"timeslotID"`
	UserDetails        []UserDetails      `bson:"userDetails" json:"userDetails"`
	StartTime          time.Time          `bson:"startTime" json:"startTime"`
	Participants       uint               `bson:"participants" json:"participants"`
	AmountPaid         float64            `bson:"amountPaid" json:"amountPaid"`
	CreatedAt          time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt          time.Time          `bson:"updateAt" json:"updateAt"`
	Payment            Payment            `bson:"payment" json:"payment"`
	CancellationDate   time.Time          `bson:"cancellationDate" json:"cancellationDate"`
	CancellationReason string             `bson:"cancellationReason" json:"cancellationReason"`
	RefundStatus       string             `bson:"refundStatus" json:"refundStatus"`
}

type PackageBooking struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PackageID          primitive.ObjectID `bson:"packageID" json:"packageID"`
	UserDetails        []UserDetails      `bson:"userDetails" json:"userDetails"`
	AmountPaid         float64            `bson:"amountPaid" json:"amountPaid"`
	CreatedAt          time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt          time.Time          `bson:"updateAt" json:"updateAt"`
	Payment            Payment            `bson:"payment" json:"payment"`
	CancellationDate   time.Time          `bson:"cancellationDate" json:"cancellationDate"`
	CancellationReason string             `bson:"cancellationReason" json:"cancellationReason"`
	RefundStatus       string             `bson:"refundStatus" json:"refundStatus"`
}

type Itinerary struct {
	ID          primitive.ObjectID `bson:"-id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"userID" json:"userID"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	StartDate   time.Time          `bson:"startDate" json:"startDate"`
	EndDate     time.Time          `bson:"endDate" json:"endDate"`
	Destination []string           `bson:"destination" json:"destination"`
	Status      string             `bson:"status" json:"status"`
	IsPublic    bool               `bson:"isPublic" json:"isPublic"`

	FlightBooking        []primitive.ObjectID `bson:"flightBooking" json:"flightBooking"`
	AccommodationBooking []primitive.ObjectID `bson:"accommodationBooking" json:"accommodationBooking"`
	ActivityrBooking          []primitive.ObjectID `bson:"activityBooking" json:"activityBooking"`

	DailyPlan   []DayPlan `bson:"dailyPlan" json:"dailyPlan"`
	TotalBudget float64   `bson:"budget" json:"budget"`
	AmountSpent float64   `bson:"amountSpent" json:"amountSpent"`
	CreatedAt   time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt" json:"updatedAt"`
}

type UserDetails struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"userID" json:"userID"`
	FirstName   string             `bson:"firstName" json:"firstName"`
	LastName    string             `bson:"lastName" json:"lastName"`
	UserPhone   string             `bson:"userPhone" json:"userPhone"`
	Gender      string             `bson:"gender" json:"gender"`
	Nationality string             `bson:"nationality" json:"nationality"`
}

type Payment struct {
	PaymentMethod string  `bson:"paymentMethod" json:"paymentMethod"`
	TotalAmount   float64 `bson:"totalAmount" json:"totalAmount"`
	Currency      string  `bson:"currency" json:"currency"`
	Status        string  `bson:"status" json:"status"`
}

type DayPlan struct {
	Day           int                  `bson:"day" json:"day"`
	Date          time.Time            `bson:"date" json:"date"`
	Location      string               `bson:"location" json:"location"`
	Activities    []primitive.ObjectID `bson:"activities" json:"activities"`
	Accommodation primitive.ObjectID   `bson:"accommodation" json:"accommodation"`
	Notes         string               `bson:"notes" json:"notes"`
}
