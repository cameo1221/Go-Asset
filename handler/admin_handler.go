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


type AdminHandler struct {
	AdminModel *models.AdminModel
}

func NewAdminHandler(adminModel *models.AdminModel) *AdminHandler {
	return &AdminHandler{AdminModel: adminModel}
}

func (ah *AdminHandler) createAdmin(w http.ResponseWriter, r *http.Request) {
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

	err = ah.AdminModel.CreateAdmin(&admin)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating admin: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Admin created successfully")
	
}

func (ah *AdminHandler) getAllAdmins(w http.ResponseWriter, r *http.Request) {
    admins, err := ah.AdminModel.GetAllAdmins()
    if err != nil {
        http.Error(w, fmt.Sprintf("Error getting admins: %v", err), http.StatusInternalServerError)
        return
    }

    adminsJSON, err := json.Marshal(admins)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error marshaling admins to JSON: %v", err), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _, err = w.Write(adminsJSON)
    if err != nil {
        log.Printf("Error writing response: %v\n", err)
    }
}
func (ah *AdminHandler) getAdmin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid admin ID", http.StatusBadRequest)
		return
	}

	admin, err := ah.AdminModel.GetAdminByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving admin: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(admin)
}

func (ah *AdminHandler) updateAdmin(w http.ResponseWriter, r *http.Request) {
	
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid admin ID", http.StatusBadRequest)
		return
	}

	var updatedAdmin models.Admin

	err = json.NewDecoder(r.Body).Decode(&updatedAdmin)
	if err != nil {
		http.Error(w, "Unable to decode request body", http.StatusBadRequest)
		return
	}

	updatedAdmin.ID = id

	err = ah.AdminModel.UpdateAdmin(&updatedAdmin)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating admin: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Admin updated successfully")
}


func (ah *AdminHandler) deleteAdmin(w http.ResponseWriter, r *http.Request) {
	
	vars := mux.Vars(r)
	id := vars["id"]
	adminID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid admin ID", http.StatusBadRequest)
		return
	}

	err = ah.AdminModel.ArchiveAdmin(adminID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting admin: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Admin deleted successfully")

}


func RegisterAdminRoutes(router *mux.Router, ah *AdminHandler) {
    router.Use(middleware.JSONContentTypeMiddleware)
    router.Use(middleware.LoggingMiddleware)

	router.HandleFunc("/admins", ah.createAdmin).Methods("POST")
	router.HandleFunc("/admins", ah.getAllAdmins).Methods("Get")
	router.HandleFunc("/admins/{id}", ah.getAdmin).Methods("GET")
	router.HandleFunc("/admins/{id}", ah.updateAdmin).Methods("PUT")
	router.HandleFunc("/admins/{id}", ah.deleteAdmin).Methods("DELETE")
}
