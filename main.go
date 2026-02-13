package main

import (
	"fmt"
	"log"
	"travel/database"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil{
		log.Fatal("Failed to load env")
	}
}

func main(){

	db := database.DB
	fmt.Println("DbName:", db.Name())
}