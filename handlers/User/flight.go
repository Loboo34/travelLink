package handlers

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	model "github.com/Loboo34/travel/models"
	"github.com/Loboo34/travel/service"
	"github.com/Loboo34/travel/utils"
)

// user

type FLightSearchHandler struct {
	flightService *service.FlightSearchService
}

func NewFlightHandler(flightService *service.FlightSearchService) *FLightSearchHandler {
	return &FLightSearchHandler{flightService: flightService}
}

// search flights
func (h *FLightSearchHandler) SearchFlight(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET allowed")
		return
	}

	params, err := parseSearchParams(r.URL.Query())
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "search params error")
		return
	}

	result, err := h.flightService.Search(r.Context(), params)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Error making search")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, result)
}

func parseSearchParams(q url.Values) (model.FlightSearch, error) {
	var params model.FlightSearch

	params.OriginCode = q.Get("originCode")
	params.DestinationCode = q.Get("destinationCode")

	depTime, err := time.Parse("2006-01-02", q.Get("depatureTime"))
	if err != nil {
		return params, errors.New("")
	}

	params.DepartureTime = depTime

	if retStr := q.Get("returnDate"); retStr != "" {
		retDate, err := time.Parse("", retStr)
		if err != nil {

			return params, errors.New("")
		}
		params.ReturnDate = &retDate
	}

	adults, err := strconv.Atoi(q.Get("adults"))
	if err != nil || adults < 1 {
		return params, errors.New("adults must be a number greater than 0")
	}
	params.Passengers.Adults = adults

	if c := q.Get("children"); c != "" {
		children, err := strconv.Atoi(c)
		if err != nil {
			return params, errors.New("children must be a number")
		}
		params.Passengers.Children = children
	}

	if i := q.Get("infants"); i != "" {
		infants, err := strconv.Atoi(i)
		if err != nil {
			return params, errors.New("infants must be a number")
		}
		params.Passengers.Infants = infants
	}

	// optional fields with defaults handled in Validate()
	params.CabinClass = model.CabinClassType(q.Get("cabinClass"))
	params.SortBy = model.FlightSortOptions(q.Get("sortBy"))

	page, _ := strconv.Atoi(q.Get("page"))
	params.Page = page

	pageSize, _ := strconv.Atoi(q.Get("pageSize"))
	params.PageSize = pageSize

	return params, nil
}
