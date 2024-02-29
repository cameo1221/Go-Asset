package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Session represents a session in the system
type Session struct {
	ID         uuid.UUID `json:"id"`
	AdminID    uuid.UUID `json:"admin_id"`
	Token      string    `json:"token"`
	Expiration time.Time `json:"expiration"`
	CreatedAt  time.Time `json:"created_at"`
}

// SessionModel represents the model for session operations
type SessionModel struct {
	DB *sql.DB
}

// CreateSession creates a new session in the database
func (sm *SessionModel) CreateSession(session *Session) error {
	query := `
		INSERT INTO admin_session (id, admin_id, token, expiration, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := sm.DB.QueryRow(query, session.ID, session.AdminID, session.Token, session.Expiration, session.CreatedAt).Scan(&session.ID)
	if err != nil {
		return err
	}

	return nil
}

// GetSessionByID retrieves a session from the database by its ID
func (sm *SessionModel) GetSessionByID(id uuid.UUID) (*Session, error) {
	query := `
		SELECT admin_id, token, expiration, created_at
		FROM admin_session
		WHERE id = $1
	`

	session := &Session{}
	err := sm.DB.QueryRow(query, id).Scan(&session.AdminID, &session.Token, &session.Expiration, &session.CreatedAt)
	if err != nil {
		return nil, err
	}

	session.ID = id
	return session, nil
}

// DeleteSession deletes a session from the database by its ID
func (sm *SessionModel) DeleteSession(id uuid.UUID) error {
	query := `
		DELETE FROM admin_session
		WHERE id = $1
	`

	_, err := sm.DB.Exec(query, id)
	return err
}
