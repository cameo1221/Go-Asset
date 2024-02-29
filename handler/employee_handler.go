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

	"github.com/cameo1221/Go-Asset/models" )

// EmployeeHandler represents the handler for managing assets
type EmployeeHandler struct {
	EmployeeModel *models.EmployeeModel
}

// NewEmployeeHandler creates a new instance of EmployeeHandler
func NewEmployeeHandler(employeeModel *models.EmployeeModel) *EmployeeHandler {
	return &EmployeeHandler{EmployeeModel: employeeModel}
}

// createAsset is a helper function for handling asset creation logic
func (ah *EmployeeHandler) createEmployee(w http.ResponseWriter, r *http.Request) {
	// Parse request body to get employee data
	var employee models.Employee
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return 
	}

	if len(body) > 0 {
		err = json.Unmarshal(body, &employee)
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

	// Call the model method to create the Employee in the database
	err = ah.EmployeeModel.CreateEmployee(&employee)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating employee: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with success message and status code
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Employee created successfully")
	
}

func (ah *EmployeeHandler) getAllEmployees(w http.ResponseWriter, r *http.Request) {
    // Call the model method to get all Employees from the database
    employees, err := ah.EmployeeModel.GetAllEmployees()
    if err != nil {
        http.Error(w, fmt.Sprintf("Error getting employees: %v", err), http.StatusInternalServerError)
        return
    }

    // Convert employees slice to JSON
    employeesJSON, err := json.Marshal(employees)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error marshaling employees to JSON: %v", err), http.StatusInternalServerError)
        return
    }

    // Set content type header and write response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _, err = w.Write(employeesJSON)
    if err != nil {
        log.Printf("Error writing response: %v\n", err)
    }
}
// getemployee is a helper function for handling employee retrieval logic
func (ah *EmployeeHandler) getEmployee(w http.ResponseWriter, r *http.Request) {
	// Extract Employee ID from URL path using gorilla mux
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Parse the ID string into a UUID
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	// Call the model method to retrieve the employee from the database
	employee, err := ah.EmployeeModel.GetEmployeeByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving employee: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the retrieved employee
	json.NewEncoder(w).Encode(employee)
}

// updateemployee is a helper function for handling employee update logic
func (ah *EmployeeHandler) updateEmployee(w http.ResponseWriter, r *http.Request) {
	// Extract Employee ID from URL path using gorilla mux
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Parse the ID string into a UUID
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	// Parse request body to get updated employee data
	var updatedEmployee models.Employee
	err = json.NewDecoder(r.Body).Decode(&updatedEmployee)
	if err != nil {
		http.Error(w, "Unable to decode request body", http.StatusBadRequest)
		return
	}

	// Set the ID of the updated Employee
	updatedEmployee.ID = id

	// Call the model method to update the Employee in the database
	err = ah.EmployeeModel.UpdateEmployee(&updatedEmployee)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating employee: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with success message and status code
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Employee updated successfully")
}


func (ah *EmployeeHandler) deleteEmployee(w http.ResponseWriter, r *http.Request) {
	// Extract Employee ID from URL path using gorilla mux
	vars := mux.Vars(r)
	id := vars["id"]

	// Convert the id variable to the appropriate type (int)
	employeeID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	// Call the model method to delete the employee from the database
	err = ah.EmployeeModel.ArchiveEmployee(employeeID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting employee: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with success message and status code
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Employee deleted successfully")

}


// RegisterRoutes registers all Employee related routes on the provided router
func RegisterEmployeeRoutes(router *mux.Router, ah *EmployeeHandler) {
    router.Use(middleware.JSONContentTypeMiddleware)
    router.Use(middleware.LoggingMiddleware)

	router.HandleFunc("/employees", ah.createEmployee).Methods("POST")
	router.HandleFunc("/employees", ah.getAllEmployees).Methods("Get")
	router.HandleFunc("/employees/{id}", ah.getEmployee).Methods("GET")
	router.HandleFunc("/employees/{id}", ah.updateEmployee).Methods("PUT")
	router.HandleFunc("/employees/{id}", ah.deleteEmployee).Methods("DELETE")
}
