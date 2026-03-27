package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Loboo34/travel/auth"
	"github.com/Loboo34/travel/database"
	"github.com/Loboo34/travel/handlers"
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
		24*time.Hour,
	)

	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Auth setup
	userRepo := auth.NewUserRepo(db)
	userService := auth.NewUserService(userRepo, jwtManager)
	userHandler := handlers.NewUserHandler(userService)

	//auth
	r.HandleFunc("/auth/register", userHandler.Register)
	r.HandleFunc("/auth/login", userHandler.Login)
	r.HandleFunc("/auth/profile", userHandler.GetProfile)

	//user

	//admin

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	utils.Logger.Info("Server starting on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
