package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Loboo34/travel/auth"
	"github.com/Loboo34/travel/database"
	"github.com/Loboo34/travel/handlers"
	handlers_admin "github.com/Loboo34/travel/handlers/Admin"
	"github.com/Loboo34/travel/middleware"
	model "github.com/Loboo34/travel/models"
	"github.com/Loboo34/travel/repository"
	"github.com/Loboo34/travel/service"
	"github.com/Loboo34/travel/utils"
	"github.com/gorilla/mux"

	"github.com/joho/godotenv"
)

// loading envs
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load env")
	}
}

func main() {

	r := mux.NewRouter()
	r.Use()

	//db
	db := database.ConnectDB()
	fmt.Println("DbName:", db.Name())

	//logger
	utils.InitLogger(os.Getenv("ENV") == "production")
	defer utils.Logger.Sync()

	//cloudinary
	if err := utils.InitCloudinary(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	); err != nil {
		log.Fatalf("cloudinary setup failed: %v", err)
	}

	jwtManager := auth.NewJWTManager(
		os.Getenv("JWT_SECRET"),
		720*time.Hour, // 30 days
	)

	user := auth.Authenticate(jwtManager)

	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Auth setup
	userRepo := auth.NewUserRepo(db)
	userService := auth.NewUserService(userRepo, jwtManager)
	userHandler := handlers.NewUserHandler(userService)

	// Create admin
	adminService := middleware.NewAdminHandler(userRepo)
	if err := adminService.Create(context.Background()); err != nil {
		var conflictErr *model.ConflictError
		if errors.As(err, &conflictErr) {
			utils.Logger.Info("admin user already exists, skipping initialization")
		} else {
			log.Fatalf("failed to create admin user on startup: %v", err)
		}
	}

	//flight
	flightRepo := repository.NewFlightRepo(db)
	flightService := service.NewFlightService(flightRepo)
	flightHandler := handlers_admin.NewFlightHandler(flightService)

	accommodatioRepo := repository.NewAccommodationRepo(db)
	accommodationService := service.NewAccommodationService(accommodatioRepo)
	accommodationHandler := handlers_admin.NewAccommodationHandler(accommodationService)

	activityRepo := repository.NewActivityRepo(db)
	activityService := service.NewActivityService(activityRepo)
	activityHandler := handlers_admin.NewActivityHandler(activityService)

	packageRepo := repository.NewPackageRepo(db)
	packageService := service.NewPackageService(packageRepo)
	packageHandler := handlers_admin.NewPackageHandler(packageService)

	//Routes

	//auth
	r.HandleFunc("/auth/register", userHandler.Register)
	r.HandleFunc("/auth/login", userHandler.Login)
	r.Handle("/auth/profile", user(http.HandlerFunc(userHandler.GetProfile)))

	//admin
	admin := auth.RequireAdmin(jwtManager)

	//flights
	r.Handle("/flight/add", admin(http.HandlerFunc(flightHandler.AddFlight)))
	r.Handle("/flight/update/{flightID}", admin(http.HandlerFunc(flightHandler.UpdateFight)))
	r.Handle("/flight/status/{flightID}", admin(http.HandlerFunc(flightHandler.UpdateFlightStatus)))
	r.Handle("/flight/delete/{flightID}", admin(http.HandlerFunc(flightHandler.DeleteFlight)))

	//offers
	r.Handle("/offer/create", admin(http.HandlerFunc(flightHandler.FlightOffer)))
	r.Handle("/offer/update/{offerID}", admin(http.HandlerFunc(flightHandler.UpdateOffer)))
	r.Handle("/offer/isactive/{offerID}", admin(http.HandlerFunc(flightHandler.IsActive)))
	r.Handle("/offer/delete/{offerID}", admin(http.HandlerFunc(flightHandler.DeleteOffer)))

	//accommodation
	r.Handle("/accommodation/add", admin(http.HandlerFunc(accommodationHandler.AddAccommodation)))
	r.Handle("/accommodation/update/{accommodationID}", admin(http.HandlerFunc(accommodationHandler.Update)))
	r.Handle("/accommodation/delete/{accommodationID}", admin(http.HandlerFunc(accommodationHandler.DeleteAccommodation)))
	r.Handle("/availability/add", admin(http.HandlerFunc(accommodationHandler.Availability)))
	r.Handle("/availability/status/{availabilityID}", admin(http.HandlerFunc(accommodationHandler.IsActive)))
	r.Handle("/availability/remove/{availabilityID}", admin(http.HandlerFunc(accommodationHandler.RemoveAvailability)))

	//activities
	r.Handle("/activity/add", admin(http.HandlerFunc(activityHandler.CreateActivity)))
	r.Handle("/activity/update/{activityID}", admin(http.HandlerFunc(activityHandler.UpdateActivity)))
	r.Handle("/activity/delete/{activityID}", admin(http.HandlerFunc(activityHandler.DeleteActivity)))
	r.Handle("/activity/{activityID}/timeslot", admin(http.HandlerFunc(activityHandler.CreateTimeSlot)))

	//packages
	r.Handle("/package/add", admin(http.HandlerFunc(packageHandler.CreatePackage)))
	r.Handle("/package/update/{packageID}", admin(http.HandlerFunc(packageHandler.UpdatePackage)))
	r.Handle("/package/status/{packageID}", admin(http.HandlerFunc(packageHandler.SetActivePackage)))
	r.Handle("/package/delete/{packageID}", admin(http.HandlerFunc(packageHandler.DeletePackage)))


	//user
	//fetch
	r.HandleFunc("/flights", flightHandler.GetFlights)
	r.HandleFunc("/flight/{flightID}", flightHandler.GetFlight)

	r.HandleFunc("/offers", flightHandler.GetOffers)
	r.HandleFunc("/offer/{offerID}", flightHandler.GetOffer)

	r.HandleFunc("/accommodations", accommodationHandler.GetAcommodations)
	r.HandleFunc("/accommodation/{accommodationID}", accommodationHandler.GetAccommodation)

	r.HandleFunc("/available", accommodationHandler.GetAvailabilities)
	r.HandleFunc("/available/{availilableID}", accommodationHandler.GetAvailability)

	r.HandleFunc("/activitties", activityHandler.GetActivities)
	r.HandleFunc("/activity/{activityID}", activityHandler.GetActivity)

	r.HandleFunc("/timeslots", activityHandler.GetTimeslots)
	r.HandleFunc("/timeslot/{timeslotID}", activityHandler.GetTimeslot)
	//search

	//booking

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	utils.Logger.Info("Server starting on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
