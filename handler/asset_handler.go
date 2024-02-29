package handler

import (
	
	"Documents/Go_Asset/middleware"
	"encoding/json"
	"fmt"
	"io"
	
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux" // Import gorilla mux

	"github.com/cameo1221/Go-Asset/models" // Update with your actual package path
)

// AssetHandler represents the handler for managing assets
type AssetHandler struct {
	AssetModel *models.AssetModel
}

// NewAssetHandler creates a new instance of AssetHandler
func NewAssetHandler(assetModel *models.AssetModel) *AssetHandler {
	return &AssetHandler{AssetModel: assetModel}
}

// createAsset is a helper function for handling asset creation logic
func (ah *AssetHandler) createAsset(w http.ResponseWriter, r *http.Request) {
	// Parse request body to get asset data
	var asset models.Asset
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return 
	}

	if len(body) > 0 {
		err = json.Unmarshal(body, &asset)
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

	// Call the model method to create the asset in the database
	err = ah.AssetModel.CreateAsset(&asset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating asset: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with success message and status code
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Asset created successfully")
	
}

func (ah *AssetHandler) getAllAssets(w http.ResponseWriter, r *http.Request) {
    // Call the model method to get all assets from the database
    assets, err := ah.AssetModel.GetAllAssets()
    if err != nil {
        http.Error(w, fmt.Sprintf("Error getting assets: %v", err), http.StatusInternalServerError)
        return
    }

    // Convert assets slice to JSON
    assetsJSON, err := json.Marshal(assets)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error marshaling assets to JSON: %v", err), http.StatusInternalServerError)
        return
    }

    // Set content type header and write response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _, err = w.Write(assetsJSON)
    if err != nil {
        log.Printf("Error writing response: %v\n", err)
    }
}
// getAsset is a helper function for handling asset retrieval logic
func (ah *AssetHandler) getAsset(w http.ResponseWriter, r *http.Request) {
	// Extract asset ID from URL path using gorilla mux
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Parse the ID string into a UUID
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid asset ID", http.StatusBadRequest)
		return
	}

	// Call the model method to retrieve the asset from the database
	asset, err := ah.AssetModel.GetAssetByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving asset: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the retrieved asset
	json.NewEncoder(w).Encode(asset)
}

// updateAsset is a helper function for handling asset update logic
// updateAsset is a helper function for handling asset update logic
func (ah *AssetHandler) updateAsset(w http.ResponseWriter, r *http.Request) {
	// Extract asset ID from URL path using gorilla mux
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Parse the ID string into a UUID
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid asset ID", http.StatusBadRequest)
		return
	}

	// Parse request body to get updated asset data
	var updatedAsset models.Asset
	err = json.NewDecoder(r.Body).Decode(&updatedAsset)
	if err != nil {
		http.Error(w, "Unable to decode request body", http.StatusBadRequest)
		return
	}

	// Set the ID of the updated asset
	updatedAsset.Id = id

	// Call the model method to update the asset in the database
	err = ah.AssetModel.UpdateAsset(&updatedAsset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating asset: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with success message and status code
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Asset updated successfully")
}


func (ah *AssetHandler) deleteAsset(w http.ResponseWriter, r *http.Request) {
	// Extract asset ID from URL path using gorilla mux
	vars := mux.Vars(r)
	id := vars["id"]

	// Convert the id variable to the appropriate type (int)
	assetID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid asset ID", http.StatusBadRequest)
		return
	}

	// Call the model method to delete the asset from the database
	err = ah.AssetModel.ArchiveAsset(assetID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting asset: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with success message and status code
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Asset deleted successfully")

}



// RegisterRoutes registers all asset related routes on the provided router
func RegisterAssetRoutes(router *mux.Router, ah *AssetHandler) {
	router.Use(middleware.jsonContentTypeMiddleware)
	router.Use(middleware.loggingMiddleware)
	router.HandleFunc("/assets", ah.createAsset).Methods("POST")
	router.HandleFunc("/assets", ah.getAllAssets).Methods("Get")
	router.HandleFunc("/assets/{id}", ah.getAsset).Methods("GET")
	router.HandleFunc("/assets/{id}", ah.updateAsset).Methods("PUT")
	router.HandleFunc("/assets/{id}", ah.deleteAsset).Methods("DELETE")
	
}
