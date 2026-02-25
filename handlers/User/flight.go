package user

import (
	"context"
	"net/http"
	"time"

	"github.com/Loboo34/travel/database"
	model "github.com/Loboo34/travel/models"
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

func getFlight(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet{
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

//search flights
//get flight bookings
//get flight by routes
//get flight availability/seats
