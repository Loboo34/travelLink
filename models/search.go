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
	OriginCode      string         `bson:"originCode" json:"originCode"`
	DestinationCode string         `bson:"destinationCode" json:"destinationCode"`
	DepartureTime   time.Time      `bson:"departureTime" json:"departureTime"`
	ReturnDate      *time.Time     `bson:"returnDate" json:"returnDate"` //nil for one way
	Passengers      PassengerCount `bson:"passengers" json:"passengers"`
	CabinClass      CabinClassType `bson:"cabinClass" json:"cabinClass"`
	SortBy          SortOptions    `bson:"sirtBy" json:"sortBy"`
	Pagination      SearchOptions  `bson:"pagination" json:"pagination"`
}

type PassengerCount struct {
	Adults   int `json:"adults"`
	Children int `json:"children"`
	Infant   int `json:"infant"`
}

func (p PassengerCount) Total() int {
	return p.Adults + p.Children + p.Infant
}

type SortOptions string

const (
	SortByPrice   SortOptions = "Price"
	SortByAirline SortOptions = "Airline"
	SortByStops   SortOptions = "Stops"
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

	if p.Pagination.Page < 1 {
		p.Pagination.Page = 1
	}
	if p.Pagination.PageSize < 1 || p.Pagination.PageSize > 50 {
		p.Pagination.PageSize = 20
	}

	return nil
}
