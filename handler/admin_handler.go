package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux" // Import gorilla mux

	"github.com/cameo1221/Go-Asset/models" // Update with your actual package path
)

// AdminHandler represents the handler for managing assets
type AdminHandler struct {
	AdminModel *models.AdminModel
}

// NewAdminHandler creates a new instance of AdminHandler
func NewAdminHandler(adminModel *models.AdminModel) *AdminHandler {
	return &AdminHandler{AdminModel: adminModel}
}

// createAsset is a helper function for handling asset creation logic
func (ah *AdminHandler) createAdmin(w http.ResponseWriter, r *http.Request) {
	// Parse request body to get Admin data
	var admin models.Admin
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return 
	}

	if len(body) > 0 {
		err = json.Unmarshal(body, &admin)
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

	// Call the model method to create the Admin in the database
	err = ah.AdminModel.CreateAdmin(&admin)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating admin: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with success message and status code
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Admin created successfully")
	
}

func (ah *AdminHandler) getAllAdmins(w http.ResponseWriter, r *http.Request) {
    // Call the model method to get all Admins from the database
    admins, err := ah.AdminModel.GetAllAdmins()
    if err != nil {
        http.Error(w, fmt.Sprintf("Error getting admins: %v", err), http.StatusInternalServerError)
        return
    }

    // Convert admins slice to JSON
    adminsJSON, err := json.Marshal(admins)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error marshaling admins to JSON: %v", err), http.StatusInternalServerError)
        return
    }

    // Set content type header and write response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _, err = w.Write(adminsJSON)
    if err != nil {
        log.Printf("Error writing response: %v\n", err)
    }
}
// getadmin is a helper function for handling admin retrieval logic
func (ah *AdminHandler) getAdmin(w http.ResponseWriter, r *http.Request) {
	// Extract Admin ID from URL path using gorilla mux
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Parse the ID string into a UUID
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid admin ID", http.StatusBadRequest)
		return
	}

	// Call the model method to retrieve the admin from the database
	admin, err := ah.AdminModel.GetAdminByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving admin: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the retrieved admin
	json.NewEncoder(w).Encode(admin)
}

// updateadmin is a helper function for handling admin update logic
// updateadmin is a helper function for handling admin update logic
func (ah *AdminHandler) updateAdmin(w http.ResponseWriter, r *http.Request) {
	// Extract Admin ID from URL path using gorilla mux
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Parse the ID string into a UUID
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid admin ID", http.StatusBadRequest)
		return
	}

	// Parse request body to get updated Admin data
	var updatedAdmin models.Admin
	err = json.NewDecoder(r.Body).Decode(&updatedAdmin)
	if err != nil {
		http.Error(w, "Unable to decode request body", http.StatusBadRequest)
		return
	}

	// Set the ID of the updated Admin
	updatedAdmin.ID = id

	// Call the model method to update the Admin in the database
	err = ah.AdminModel.UpdateAdmin(&updatedAdmin)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating admin: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with success message and status code
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Admin updated successfully")
}


func (ah *AdminHandler) deleteAdmin(w http.ResponseWriter, r *http.Request) {
	// Extract Admin ID from URL path using gorilla mux
	vars := mux.Vars(r)
	id := vars["id"]

	// Convert the id variable to the appropriate type (int)
	adminID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid admin ID", http.StatusBadRequest)
		return
	}

	// Call the model method to delete the admin from the database
	err = ah.AdminModel.ArchiveAdmin(adminID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting admin: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with success message and status code
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Admin deleted successfully")

}


// RegisterRoutes registers all admin related routes on the provided router
func RegisterAdminRoutes(router *mux.Router, ah *AdminHandler) {
	router.Use(jsonContentTypeMiddleware)
	router.Use(loggingMiddleware)

	router.HandleFunc("/admins", ah.createAdmin).Methods("POST")
	router.HandleFunc("/admins", ah.getAllAdmins).Methods("Get")
	router.HandleFunc("/admins/{id}", ah.getAdmin).Methods("GET")
	router.HandleFunc("/admins/{id}", ah.updateAdmin).Methods("PUT")
	router.HandleFunc("/admins/{id}", ah.deleteAdmin).Methods("DELETE")
}
