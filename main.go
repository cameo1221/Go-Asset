package main

import (
	"log"

	"github.com/cameo1221/Go-Asset/db" // Import your db package path
)

func main() {
	// Connect to the PostgreSQL database
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer database.Close()

	log.Println("Successfully connected to the database")
}
