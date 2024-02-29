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

type AssetHandler struct {
	AssetModel *models.AssetModel
}

func NewAssetHandler(assetModel *models.AssetModel) *AssetHandler {
	return &AssetHandler{AssetModel: assetModel}
}

func (ah *AssetHandler) createAsset(w http.ResponseWriter, r *http.Request) {
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

	err = ah.AssetModel.CreateAsset(&asset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating asset: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Asset created successfully")
	
}

func (ah *AssetHandler) getAllAssets(w http.ResponseWriter, r *http.Request) {
    assets, err := ah.AssetModel.GetAllAssets()
    if err != nil {
        http.Error(w, fmt.Sprintf("Error getting assets: %v", err), http.StatusInternalServerError)
        return
    }

    assetsJSON, err := json.Marshal(assets)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error marshaling assets to JSON: %v", err), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _, err = w.Write(assetsJSON)
    if err != nil {
        log.Printf("Error writing response: %v\n", err)
    }
}
func (ah *AssetHandler) getAsset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid asset ID", http.StatusBadRequest)
		return
	}

	asset, err := ah.AssetModel.GetAssetByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving asset: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(asset)
}

func (ah *AssetHandler) updateAsset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid asset ID", http.StatusBadRequest)
		return
	}

	var updatedAsset models.Asset
	err = json.NewDecoder(r.Body).Decode(&updatedAsset)
	if err != nil {
		http.Error(w, "Unable to decode request body", http.StatusBadRequest)
		return
	}

	updatedAsset.Id = id

	err = ah.AssetModel.UpdateAsset(&updatedAsset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating asset: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Asset updated successfully")
}


func (ah *AssetHandler) deleteAsset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	assetID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid asset ID", http.StatusBadRequest)
		return
	}

	err = ah.AssetModel.ArchiveAsset(assetID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting asset: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Asset deleted successfully")

}




func RegisterAssetRoutes(router *mux.Router, ah *AssetHandler) {
	router.Use(middleware.JSONContentTypeMiddleware)
    router.Use(middleware.LoggingMiddleware)
	router.HandleFunc("/assets", ah.createAsset).Methods("POST")
	router.HandleFunc("/assets", ah.getAllAssets).Methods("Get")
	router.HandleFunc("/assets/{id}", ah.getAsset).Methods("GET")
	router.HandleFunc("/assets/{id}", ah.updateAsset).Methods("PUT")
	router.HandleFunc("/assets/{id}", ah.deleteAsset).Methods("DELETE")
	
}
