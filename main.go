package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Loboo34/travel/database"
	"github.com/Loboo34/travel/utils"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load env")
	}
}

func main() {

	db := database.DB
	fmt.Println("DbName:", db.Name())

	if err := utils.InitCloudinary(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	); err != nil {
		log.Fatalf("cloudinary setup failed: %v", err)
	}
}
