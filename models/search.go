package model

import "time"

type SearchOptions struct {
	Page     int `bson:"page" json:"page"`
	PageSize int `bson:"pageSize" json:"pageSize"`
}

type FlightSearch struct {
	OriginCode      string         `bson:"originCode" json:"originCode"`
	DestinationCode string         `bson:"destinationCode" json:"destinationCode"`
	Passengers      PassengerCount `bson:"passengers" json:"passengers"`
	CabinClass      CabinClassType `bson:"cabinClass" json:"cabinClass"`
	DepartureTime   time.Time      `bson:"departureTime" json:"departureTime"`
	ReturnDate      *time.Time     `bson:"returnDate" json:"returnDate"`
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
	SSortByPrice  SortOptions = "Price"
	SortByAirline SortOptions = "Airline"
	SortByStops   SortOptions = "Stops"
)
