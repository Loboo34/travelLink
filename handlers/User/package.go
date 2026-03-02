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

//user
//Get packages/package
func GetPackages(w http.ResponseWriter, r *http.Request){
   if r.Method != http.MethodGet {
      utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET allowed")
      return 
   }

   packageCollection := database.DB.Collection("packags")

   ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
   defer cancel()

   cursor, err := packageCollection.Find(ctx, bson.M{})
   if err != nil{
      utils.Logger.Warn("Failed to fetch packags")
      utils.RespondWithError(w, http.StatusInternalServerError, "Error finding package")
      return 
   }
   defer cursor.Close(ctx)

	var packages []model.Package
	if err := cursor.All(ctx, &packages); err != nil {
		utils.Logger.Warn("Failed decoding packages")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error decoding packages")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, "Fetched packages", packages)

}

//get flight
func GetPackage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET allowed")
		return
	}

	vars := mux.Vars(r)
	packageIDStr := vars["packageID"]
	if packageIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing package ID")
		return
	}

	packageID, err := primitive.ObjectIDFromHex(packageIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid package ID")
		return
	}

	flightCollection := database.DB.Collection("packages")
	var pack model.Package

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = flightCollection.FindOne(ctx, bson.M{"_id": packageID}).Decode(&pack)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error finding package")
		utils.Logger.Warn("Failed to find package")
		return
	}

	utils.Logger.Info("Fetched package")
	utils.RespondWithJson(w, http.StatusOK, "Package found", map[string]interface{}{"package": packageID})

}
//search for packages
//get reccs
