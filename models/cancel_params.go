package model

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cancellation struct {
	BookingID primitive.ObjectID `json:"bookingID"`
	Reason   string             `json:"reason"`
}

func (r *Cancellation) Validate() error{
	if r.BookingID.IsZero() {
		return errors.New("bookignID is required")
	}

	return nil
}


