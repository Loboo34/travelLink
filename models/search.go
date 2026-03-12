package model

import (
	"errors"
	"time"
)

type SearchOptions struct {
	Page     int `bson:"page" json:"page"`
	PageSize int `bson:"pageSize" json:"pageSize"`
}

type FlightSearch struct {
	OriginCode      string            `bson:"originCode" json:"originCode"`
	DestinationCode string            `bson:"destinationCode" json:"destinationCode"`
	DepartureTime   time.Time         `bson:"departureTime" json:"departureTime"`
	ReturnDate      *time.Time        `bson:"returnDate" json:"returnDate"` //nil for one way
	Passengers      PassengerCount    `bson:"passengers" json:"passengers"`
	CabinClass      CabinClassType    `bson:"cabinClass" json:"cabinClass"`
	SortBy          FlightSortOptions `bson:"sirtBy" json:"sortBy"`
	Page            int               `bson:"page" json:"page"`
	PageSize        int               `bson:"pageSize" json:"pageSize"`
}

type PassengerCount struct {
	Adults   int `json:"adults"`
	Children int `json:"children"`
	Infant   int `json:"infant"`
}

func (p PassengerCount) Total() int {
	return p.Adults + p.Children + p.Infant
}

type FlightSortOptions string

const (
	SortByPrice   FlightSortOptions = "Price"
	SortByAirline FlightSortOptions = "Airline"
	SortByStops   FlightSortOptions = "Stops"
)

func (p *FlightSearch) Validate() error {
	if p.OriginCode == "" {
		return errors.New("Origin airport code is required")
	}
	if p.DestinationCode == "" {
		return errors.New("Destination airport code is required")
	}

	if p.OriginCode == p.DestinationCode {
		return errors.New("Origin and Destination airpot can not be the same")
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	if p.DepartureTime.UTC().Before(today) {
		return errors.New("Departure time can not be in the past")
	}

	if p.ReturnDate != nil {
		if !p.ReturnDate.UTC().After(p.DepartureTime.UTC()) {
			return errors.New("return date must be after departure date")
		}
	}

	if p.Passengers.Adults < 1 {
		return errors.New("Atleast one adult passanger is needed")
	}
	if p.Passengers.Infant > p.Passengers.Adults {
		return errors.New("Infants must not exceed the number of adults")
	}

	if p.Passengers.Total() > 8 {
		return errors.New("Only mux of 8 passengers allowed to book at once")
	}

	switch p.CabinClass {
	case CabinClassFirst, CabinClassBusiness, CabinClassEconomy:
	default:
		p.CabinClass = CabinClassEconomy
	}

	switch p.SortBy {
	case SortByPrice, SortByAirline, SortByStops:
	default:
		p.SortBy = SortByPrice
	}

	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 || p.PageSize > 50 {
		p.PageSize = 20
	}

	return nil
}

type AccommodationSearch struct {
	Location     LocationSearch          `json:"location"`
	ChecKInDate  time.Time               `json:"checkInDate"`
	CheckOutDate time.Time               `json:"checkOutDate"`
	Guests       GuestCount              `json:"Guests"`
	PropertyType PropertyType            `json:"propertyType"`
	TotalRooms   int                     `json:"totalRooms"`
	SortBy       AccommodationSortOption `json:"sortBy"`
	Page         int                     `json:"page"`
	PageSize     int                     `json:"pageSize"`
}

type LocationSearch struct {
	City      string  `json:"city"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	RadiusKm  float64 `json:"radiusKm,omitempty"`
}

type GuestCount struct {
	Adults   int `json:"adults"`
	Children int `json:"children"`
	Infants  int `json:"infants"`
}

type AccommodationSortOption string

const (
	SortAccommodationByPrice  AccommodationSortOption = "price"
	SortAccommodationByRating AccommodationSortOption = "rating"
)

func (a *AccommodationSearch) Validate() error {
	if a.Location.City == "" && a.Location.Latitude == 0 && a.Location.Longitude == 0 {
		return errors.New("Location can not be empty")
	}

	if a.Location.Latitude != 0 || a.Location.Longitude != 0 {
		if a.Location.RadiusKm <= 0 {
			a.Location.RadiusKm = 10
		}
	}

	today := time.Now().Local().Truncate(24 * time.Hour)
	if a.ChecKInDate.Local().Before(today) {
		return errors.New("Check in date can not be in the past")
	}

	if !a.CheckOutDate.After(a.ChecKInDate) {
		return errors.New("Checkout date must be after check in date")
	}

	if a.Guests.Adults < 1 {
		return errors.New("At least one guest is required")
	}

	switch a.PropertyType {
	case PropertyTypeHotel, PropertyTypeVilla, PropertyTypeAirBnb, PropertyTypeGuesthouse, PropertyTypeResort:
	case "":
	default:
		return errors.New("Invalid property Type")
	}

	if a.PropertyType == PropertyTypeHotel && a.TotalRooms < 1 {
		a.TotalRooms = 1
	}

	switch a.SortBy {
	case SortAccommodationByPrice, SortAccommodationByRating:
	default:
		a.SortBy = SortAccommodationByPrice
	}

	if a.Page < 1 {
		a.Page = 1
	}
	if a.PageSize < 1 || a.PageSize > 50 {
		a.PageSize = 20
	}

	return nil
}

type ActivitySearch struct {
	Location     LocationSearch      `json:"location"`
	Participants ParticipantCount    `json:"participants"`
	ForAllAges   bool                `json:"forAllAges"` //should this be a search thing or filter??
	Category     ActivityCategory    `json:"category"`
	Date         time.Time           `json:"date"`
	SortBy       ActivitySortOptions `json:"sortBy"`
	Page         int                 `json:"page"`
	PageSize     int                 `json:"pageSize"`
}

type ParticipantCount struct {
	Adults   int `json:"adults"`
	Children int `json:"children"`
	Infants  int `json:"infants"`
}

type ActivitySortOptions string

const (
	SortActivityByPrice  ActivitySortOptions = "price"
	SortActivityByRating ActivitySortOptions = "rating"
)

func (a *ActivitySearch) Validate() error {
	if a.Location.City == "" && a.Location.Latitude == 0 && a.Location.Longitude == 0 {
		return errors.New("Location can not be empty")
	}

	if a.Location.Latitude != 0 || a.Location.Longitude != 0 {
		if a.Location.RadiusKm <= 0 {
			a.Location.RadiusKm = 10
		}
	}

	today := time.Now().Local().Truncate(24 * time.Hour)
	if a.Date.Local().Before(today) {
		return errors.New("Date can not be in the past")
	}

	if a.Participants.Adults < 1 {
		return errors.New("At least one participant is needed")
	}

	switch a.Category {
	case ActivityCategoryAdventure, ActivityCategoryCultural, ActivityCategoryFood, ActivityCategoryNature, ActivityCategoryNightlife, ActivityCategorySightseeing, ActivityCategoryWater, ActivityCategoryWellness, "":
	default:
		return errors.New("Invalid category")
	}

	switch a.SortBy {
	case SortActivityByPrice, SortActivityByRating:
	default:
		a.SortBy = SortActivityByPrice
	}

	if a.Page < 1 {
		a.Page = 1
	}
	if a.PageSize < 1 || a.PageSize > 50 {
		a.PageSize = 20
	}

	return nil
}
