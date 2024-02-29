
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

type SessionHandler struct {
	SessionModel *models.SessionModel
}

func NewSessionHandler(sessionModel *models.SessionModel) *SessionHandler {
	return &SessionHandler{SessionModel: sessionModel}
}

func (ah *SessionHandler) createSession(w http.ResponseWriter, r *http.Request) {
	var session models.Session
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return 
	}

	if len(body) > 0 {
		err = json.Unmarshal(body, &session)
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

	err = ah.SessionModel.CreateSession(&session)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating session: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "session created successfully")
	
}

func (ah *SessionHandler) getAllSessions(w http.ResponseWriter, r *http.Request) {
    sessions, err := ah.SessionModel.GetAllSessions()
    if err != nil {
        http.Error(w, fmt.Sprintf("Error getting Sessions: %v", err), http.StatusInternalServerError)
        return
    }

    sessionsJSON, err := json.Marshal(sessions)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error marshaling sessions to JSON: %v", err), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _, err = w.Write(sessionsJSON)
    if err != nil {
        log.Printf("Error writing response: %v\n", err)
    }
}
func (ah *SessionHandler) getSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid Session ID", http.StatusBadRequest)
		return
	}

	session, err := ah.SessionModel.GetSessionByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving session: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(session)
}

func (ah *SessionHandler) updateSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid Session ID", http.StatusBadRequest)
		return
	}

	var updatedSession models.Session
	err = json.NewDecoder(r.Body).Decode(&updatedSession)
	if err != nil {
		http.Error(w, "Unable to decode request body", http.StatusBadRequest)
		return
	}

	updatedSession.ID = id

	err = ah.SessionModel.UpdateSession(&updatedSession)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating Session: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Session updated successfully")
}


func (ah *SessionHandler) deleteSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	sessionID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	err = ah.SessionModel.ArchiveSession(sessionID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting session: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "session deleted successfully")

}


func RegisterSessionRoutes(router *mux.Router, ah *SessionHandler) {
    router.Use(middleware.JSONContentTypeMiddleware)
    router.Use(middleware.LoggingMiddleware)

	router.HandleFunc("/sessions", ah.createSession).Methods("POST")
	router.HandleFunc("/sessions", ah.getAllSessions).Methods("Get")
	router.HandleFunc("/sessions/{id}", ah.getSession).Methods("GET")
	router.HandleFunc("/sessions/{id}", ah.updateSession).Methods("PUT")
	router.HandleFunc("/sessions/{id}", ah.deleteSession).Methods("DELETE")
}
