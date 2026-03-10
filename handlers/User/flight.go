package user

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Loboo34/travel/database"
	model "github.com/Loboo34/travel/models"
	"github.com/Loboo34/travel/service"
	"github.com/Loboo34/travel/utils"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// user
// get flights
func GetFlights(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET allowed")
		return
	}

	flightCollection := database.DB.Collection("flights")

	// use request context as parent
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	cursor, err := flightCollection.Find(ctx, bson.M{})
	if err != nil {
		utils.Logger.Warn("Failed to fetch flights")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error finding flights")
		return
	}
	defer cursor.Close(ctx)

	var flights []model.Flight
	if err := cursor.All(ctx, &flights); err != nil {
		utils.Logger.Warn("Failed decoding flights")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error decoding flights")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, "Fetched flights", flights)
}

func GetFlight(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET allowed")
		return
	}

	vars := mux.Vars(r)
	flightIDStr := vars["flightID"]
	if flightIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing flight ID")
		return
	}

	flightID, err := primitive.ObjectIDFromHex(flightIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid flight ID")
		return
	}

	flightCollection := database.DB.Collection("flights")
	var flight model.Flight

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = flightCollection.FindOne(ctx, bson.M{"_id": flightID}).Decode(&flight)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error finding flight")
		utils.Logger.Warn("Failed to find flight")
		return
	}

	utils.Logger.Info("Fetched flight")
	utils.RespondWithJson(w, http.StatusOK, "Flight found", map[string]interface{}{"flight": flightID})

}

type FlightHandler struct {
	flightService *service.FlightService
}

func NewFlightHandler(flightService *service.FlightService) *FlightHandler {
	return &FlightHandler{flightService: flightService}
}

// search flights
func (h *FlightHandler) SearchFlight(w http.ResponseWriter, r *http.Request) {
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

	utils.RespondWithJson(w, http.StatusOK, "", result)
}

func parseSearchParams(q url.Values) (model.FlightSearch, error) {
	var params model.FlightSearch

	params.OriginCode = q.Get("originCode")
	params.DestinationCode = q.Get("destinationCode")

	depTime, err := time.Parse("", q.Get("depatureTime"))
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
		params.Passengers.Infant = infants
	}

	// optional fields with defaults handled in Validate()
	params.CabinClass = model.CabinClassType(q.Get("cabinClass"))
	params.SortBy = model.SortOptions(q.Get("sortBy"))

	page, _ := strconv.Atoi(q.Get("page"))
	params.Page = page

	pageSize, _ := strconv.Atoi(q.Get("pageSize"))
	params.PageSize = pageSize

	return params, nil
}
