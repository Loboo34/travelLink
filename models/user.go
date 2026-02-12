package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitepty" json:"id"`
	FirstName   string             `bson:"firstName" json:"firstName"`
	LastName    string             `bson:"lasttName" json:"lasttName"`
	DateOfBirth time.Time          `bson:"dateofBirth" json:"dateofbirth"`
	Nationality string             `bson:"nationality" json:"nationality"`
	PhoneNumber string             `bson:"phoneNumber" json:"phoneNumber"`
	Email       string             `bson:"email" json:"email"`
	Password    string             `bson:"password" json:"password"`
	ProfilePic  string             `bson:"profilepic" json:"profilepic"`
	IsActive    bool               `bson:"isActive" json:"isActive"`
}

type ID struct {
	ID       primitive.ObjectID `bson:"_id,omitepty" json:"id"`
	User     User               `bson:"user" json:"user"`
	IDNumber string             `bson:"idNumber" json:"idNumber"`
}

type Passport struct {
	ID             primitive.ObjectID `bson:"_id,omitepty" json:"id"`
	User           User               `bson:"user" json:"user"`
	PassportNumber string             `bson:"passportNumber" json:"passportNumber"`
	ExpiryDate     time.Time          `bson:"expiryDate" json:"expiryDate"`
	IssuingCountry string             `bson:"issuingCountry" json:"issuingCountry"`
}

type TravelPreferences struct {
	ID                  primitive.ObjectID `bson:"_id,omitepty" json:"id"`
	User                User               `bson:"user" json:"user"`
	UserID              uint               `bson:"userID" json:"userID"`
	Interests           []string           `bson:"interests" json:"interests"`
	DietaryRestrictions []string           `bson:"dietaryRestrictions" json:"dietaryRestrictions"`
	AccessibilityNeeds  string             `bson:"accessibilityNeeds" json:"accessibilityNeeds"`
	BudgetRangeMin      float64            `bson:"budgetRangemin" json:"budgetRangemin"`
	BudgetRangeMax      float64            `bson:"bugetRangemax" json:"bugetRangemax"`
}

type BookingHistory struct {
	ID      primitive.ObjectID `bson:"_id,omitepty" json:"id"`
	User    User               `bson:"user" json:"user"`
	Type    string             `bson:"type" json:"type"`
	Status  string             `bson:"status" json:"status"`
	Amount  float64            `bson:"amount" json:"amount"`
	Details string             `bson:"details" json:"details"`
}
