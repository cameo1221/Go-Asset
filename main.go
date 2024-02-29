package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cameo1221/Go-Asset/db"
	"github.com/cameo1221/Go-Asset/handler"
	"github.com/cameo1221/Go-Asset/models"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize your database connection
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Initialize your asset model with the database connection
	assetModel := &models.AssetModel{DB: database.Conn}
	adminModel := &models.AdminModel{DB: database.Conn}
	employeeModel := &models.EmployeeModel{DB: database.Conn}
	employeeAssetModel := &models.EmployeeAssetModel{DB: database.Conn}
	sessionModel := &models.SessionModel{DB: database.Conn}

	// Initialize your asset handler with the asset model
	assetHandler := handler.NewAssetHandler(assetModel)
	adminHandler := handler.NewAdminHandler(adminModel)
	employeeHandler := handler.NewEmployeeHandler(employeeModel)
	employeeAssetHandler := handler.NewEmployeeassetHandler(employeeAssetModel)
	sessionHandler := handler.NewSessionHandler(sessionModel)



	// Initialize a new mux router
	router := mux.NewRouter()

	// Register asset routes with the router
	handler.RegisterAssetRoutes(router, assetHandler)
	handler.RegisterAdminRoutes(router, adminHandler)
	handler.RegisterEmployeeRoutes(router, employeeHandler)
	handler.RegisterEmployeeassetRoutes(router, employeeAssetHandler)
	handler.RegisterSessionRoutes(router, sessionHandler)


	// Start the HTTP server
	port := ":8080"
	fmt.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
