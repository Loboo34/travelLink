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

	//user := auth.Authenticate(jwtManager)

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

	//auth
	r.HandleFunc("/auth/register", userHandler.Register)
	r.HandleFunc("/auth/login", userHandler.Login)
	r.HandleFunc("/auth/profile", userHandler.GetProfile)

	//user

	//admin
	admin := auth.RequireAdmin(jwtManager)

	//flights
	r.Handle("/flight/add", admin(http.HandlerFunc(flightHandler.AddFlight)))
	r.Handle("/flight/update/{flightID}", admin(http.HandlerFunc(flightHandler.UpdateFight)))
	r.Handle("/flight/status/{flightID}", admin(http.HandlerFunc(flightHandler.UpdateFlightStatus)))
	r.Handle("/flight/delete/{flightID}", admin(http.HandlerFunc(flightHandler.DeleteFlight)))

	//offers
	r.HandleFunc("/offer/create", flightHandler.FlightOffer)
	r.HandleFunc("/offer/update", flightHandler.UpdateOffer)
	r.HandleFunc("/offer/isactive", flightHandler.IsActive)
	r.HandleFunc("/offer/delete", flightHandler.DeleteOffer)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	utils.Logger.Info("Server starting on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
