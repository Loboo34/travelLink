package service

import (
	"context"
	"fmt"

	model "github.com/Loboo34/travel/models"
	"github.com/Loboo34/travel/repository"
)

type FlightSearchService struct {
	FlightSearchRepo  *repository.FlightSearchRepo
	AirportRepo *repository.AirportRepo
}

type FlightSearchResults struct {
	Outbound []model.FlightOffer `json:"outbound"`
	Inbound  []model.FlightOffer `json:"inbound"`
	IsReturn bool                `json:"isReturn"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"pageSize"`
}

func NewFlightSearchService(flightSearchRepo *repository.FlightSearchRepo, AirportRepo *repository.AirportRepo) *FlightSearchService {
	return &FlightSearchService{
		FlightSearchRepo: flightSearchRepo,
		AirportRepo: AirportRepo,
	}
}

func (s *FlightSearchService) Search(ctx context.Context, params model.FlightSearch) (*FlightSearchResults, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("invalid search params: %w", err)
	}

	originID, err := s.AirportRepo.FindIDByCode(ctx, params.OriginCode)
	if err != nil {
		return nil, fmt.Errorf("Depature airport: %w", err)
	}

	destinationID, err := s.AirportRepo.FindIDByCode(ctx, params.DestinationCode)
	if err != nil {
		return nil, fmt.Errorf("destination airport: %w", err)
	}

	filter := repository.FlightFilter{
		OriginID:      originID,
		DestinationID: destinationID,
		DepartureTime: params.DepartureTime,
		CabinClass:    params.CabinClass,
		MinSeats:      params.Passengers.Total(),
		SortBy:        params.SortBy,
		Page:          params.Page,
		PageSize:      params.PageSize,
	}

	if params.ReturnDate == nil {
		return s.oneWaySearch(ctx, params, filter)
	}

	return s.roundTripSearch(ctx, params, filter)
}

func (s *FlightSearchService) oneWaySearch(ctx context.Context, params model.FlightSearch, filter repository.FlightFilter) (*FlightSearchResults, error) {
	offers, err := s.FlightSearchRepo.SearchOffers(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("Outboud search: %w", err)
	}

	return &FlightSearchResults{
		Outbound: offers,
		IsReturn: false,
		Page:     params.Page,
		PageSize: params.PageSize,
	}, nil

}

func (s *FlightSearchService) roundTripSearch(ctx context.Context, params model.FlightSearch, outboundFilter repository.FlightFilter) (*FlightSearchResults, error) {

	inboundFilter := repository.FlightFilter{
		OriginID:      outboundFilter.DestinationID,
		DestinationID: outboundFilter.OriginID,
		DepartureTime: *params.ReturnDate,
		CabinClass:    outboundFilter.CabinClass,
		MinSeats:      outboundFilter.MinSeats,
		SortBy:        outboundFilter.SortBy,
		Page:          outboundFilter.Page,
		PageSize:      outboundFilter.PageSize,
	}

	type result struct {
		offers []model.FlightOffer
		err    error
	}

	outboundCh := make(chan result, 1)
	inboundCh := make(chan result, 1)

	go func() {
		offers, err := s.FlightSearchRepo.SearchOffers(ctx, outboundFilter)
		outboundCh <- result{offers, err}
	}()

	go func() {
		offers, err := s.FlightSearchRepo.SearchOffers(ctx, inboundFilter)
		inboundCh <- result{offers, err}
	}()

	outbound := <-outboundCh
	inbound := <-inboundCh

	if outbound.err != nil {
		return nil, fmt.Errorf("outbound search: %w", outbound.err)
	}
	if inbound.err != nil {
		return nil, fmt.Errorf("inbound search: %w", inbound.err)
	}

	return &FlightSearchResults{
		Outbound: outbound.offers,
		Inbound:  inbound.offers,
		IsReturn: true,
		Page:     params.Page,
		PageSize: params.PageSize,
	}, nil
}
