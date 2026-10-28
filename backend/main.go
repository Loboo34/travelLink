package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Loboo34/travel/database"
	"github.com/Loboo34/travel/middleware"
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
	r.Use(middleware.Cors())
	//r.Use()

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

	useRoutes(r, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	utils.Logger.Info("Server starting on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
