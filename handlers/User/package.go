package handlers_user

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

// search for packages
type PackageSearchHandler struct {
	packageService *service.PackageSearchService
}

func NewSearchPackageHandler(packageService *service.PackageSearchService) *PackageSearchHandler {
	return &PackageSearchHandler{packageService: packageService}
}

func (h *PackageSearchHandler) PackageSearch(w http.ResponseWriter, r *http.Request) {
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
