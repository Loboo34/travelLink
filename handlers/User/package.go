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
// Get packages/package
func GetPackages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET allowed")
		return
	}

	packageCollection := database.DB.Collection("packags")

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	cursor, err := packageCollection.Find(ctx, bson.M{})
	if err != nil {
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

	utils.RespondWithJson(w, http.StatusOK,  packages)

}

// get flight
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
	utils.RespondWithJson(w, http.StatusOK,  map[string]interface{}{"package": packageID})

}

// search for packages
type PackageHandler struct {
	packageService *service.PackageService
}

func NewPackageHandler(packageService *service.PackageService) *PackageHandler {
	return &PackageHandler{packageService: packageService}
}

func (h *PackageHandler) PackageSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET allowed")
		return
	}

	params, err := parsePackageParams(r.URL.Query())
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Search param error")
		return
	}

	results, err := h.packageService.Search(r.Context(), params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error making search")
		utils.Logger.Warn("Failed to make search")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, results)
}

func parsePackageParams(q url.Values) (model.PackageSearchParams, error) {
	var params model.PackageSearchParams

	params.Destination = q.Get("destination")

	startDate, err := time.Parse("2006-01-02", q.Get("startDate"))
	if err != nil {
		return params, errors.New("wrong date format")
	}
	params.StartDate = startDate

	end, err := time.Parse("2006-01-02", q.Get("endDate"))
	if err != nil {
		return params, errors.New("wrong date format")
	}
	params.EndDate = end

	adults, err := strconv.Atoi(q.Get("adults"))
	if err != nil || adults < 1 {
		return params, errors.New("adults must be a number greater than 0")
	}
	params.Travelers.Adults = adults

	if c := q.Get("children"); c != "" {
		children, err := strconv.Atoi(c)
		if err != nil {
			return params, errors.New("children must be a valid number")
		}
		params.Travelers.Children = children
	}

	if i := q.Get("infants"); i != "" {
		infants, err := strconv.Atoi(i)
		if err != nil {
			return params, errors.New("infants must be a valid number")
		}
		params.Travelers.Infants = infants
	}

	if b := q.Get("maxBudget"); b != "" {
		budget, err := strconv.ParseInt(b, 10, 64)
		if err != nil {
			return params, errors.New("maxBudget must be a valid number")
		}
		params.MaxBudget = budget
	}

	for _, t := range q["tag"] {
		params.Tags = append(params.Tags, model.PackageTag(t))
	}

	for _, c := range q["component"] {
		params.Components = append(params.Components, model.Component(c))
	}
	params.Components = []model.Component{model.Component(q.Get("components"))}
	params.SortBy = model.PackageSortOption(q.Get("sortBy"))
	page, _ := strconv.Atoi(q.Get("page"))
	params.Page = page

	pageSize, _ := strconv.Atoi(q.Get("pageSize"))
	params.PageSize = pageSize

	return params, nil
}

//get reccs
