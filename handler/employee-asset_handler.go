package handler

import (
	"github.com/cameo1221/Go-Asset/middleware"
	
	"encoding/json"
	"fmt"
	"io"
	
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/cameo1221/Go-Asset/models" 
)


type EmployeeassetHandler struct {
	EmployeeassetModel *models.EmployeeAssetModel
}

func NewEmployeeassetHandler(EmployeeassetModel *models.EmployeeAssetModel) *EmployeeassetHandler {
	return &EmployeeassetHandler{EmployeeassetModel: EmployeeassetModel}
}

func (ah *EmployeeassetHandler) createEmployeeasset(w http.ResponseWriter, r *http.Request) {
	var employeeasset models.EmployeeAsset
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return 
	}

	if len(body) > 0 {
		err = json.Unmarshal(body, &employeeasset)
		if err != nil {
			return 
		}
	}
	if err != nil {
		if err == io.EOF {
			http.Error(w, "Request body is empty", http.StatusBadRequest)
			return
		}
		http.Error(w, fmt.Sprintf("Error decoding request body: %v", err), http.StatusBadRequest)
		return
	}

	err = ah.EmployeeassetModel.CreateEmployeeAsset(&employeeasset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating employeeasset: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "employeeasset created successfully")
	
}

func (ah *EmployeeassetHandler) getAllEmployeeassets(w http.ResponseWriter, r *http.Request) {
    employeeassets, err := ah.EmployeeassetModel.GetAllEmployeeAssets()
    if err != nil {
        http.Error(w, fmt.Sprintf("Error getting Employeeassets: %v", err), http.StatusInternalServerError)
        return
    }

    employeeassetsJSON, err := json.Marshal(employeeassets)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error marshaling employeeassets to JSON: %v", err), http.StatusInternalServerError)
        return
    }

    w.Header().Set("ContentType", "application/json")
    w.WriteHeader(http.StatusOK)
    _, err = w.Write(employeeassetsJSON)
    if err != nil {
        log.Printf("Error writing response: %v\n", err)
    }
}
func (ah *EmployeeassetHandler) getEmployeeasset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid Employeeasset ID", http.StatusBadRequest)
		return
	}

	employeeasset, err := ah.EmployeeassetModel.GetEmployeeAssetByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving Employeeasset: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(employeeasset)
}

func (ah *EmployeeassetHandler) updateEmployeeasset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid Employeeasset ID", http.StatusBadRequest)
		return
	}

	var updatedEmployeeasset models.EmployeeAsset

	err = json.NewDecoder(r.Body).Decode(&updatedEmployeeasset)
	if err != nil {
		http.Error(w, "Unable to decode request body", http.StatusBadRequest)
		return
	}

	updatedEmployeeasset.ID = id
	err = ah.EmployeeassetModel.UpdateEmployeeAsset(&updatedEmployeeasset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating Employeeasset: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Employeeasset updated successfully")
}


func (ah *EmployeeassetHandler) deleteEmployeeasset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	employeeassetID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid employeeasset ID", http.StatusBadRequest)
		return
	}

	err = ah.EmployeeassetModel.ArchiveEmployeeAsset(employeeassetID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting employeeasset: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "employeeasset deleted successfully")

}

func RegisterEmployeeassetRoutes(router *mux.Router, ah *EmployeeassetHandler) {
    router.Use(middleware.JSONContentTypeMiddleware)
    router.Use(middleware.LoggingMiddleware)

	router.HandleFunc("/employeeassets", ah.createEmployeeasset).Methods("POST")
	router.HandleFunc("/employeeassets", ah.getAllEmployeeassets).Methods("Get")
	router.HandleFunc("/employeeassets/{id}", ah.getEmployeeasset).Methods("GET")
	router.HandleFunc("/employeeassets/{id}", ah.updateEmployeeasset).Methods("PUT")
	router.HandleFunc("/employeesassets/{id}", ah.deleteEmployeeasset).Methods("DELETE")
}
