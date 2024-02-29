package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID         uuid.UUID `json:"id"`
	AdminID    uuid.UUID `json:"admin_id"`
	Archive_at time.Time `json:"archive_at"`
	CreatedAt  time.Time `json:"created_at"`
}

type SessionModel struct {
	DB *sql.DB
}

func (sm *SessionModel) GetAllSessions() ([]*Session, error){
	query := `Select * FROM admin_session`
	rows, err := sm.DB.Query(query)
	if err != nil{
		return nil, err

	}
	defer rows.Close()

	var sessions []*Session
	for rows.Next(){
		session := &Session{}
		err := rows.Scan(&session.ID, &session.AdminID, &session.Archive_at, &session.CreatedAt)
		if err != nil{
			return nil, err
		}
		sessions = append(sessions, session)

	}
	if err:= rows.Err(); err != nil{
		return nil, err

	}
	return sessions, nil
}
func (sm *SessionModel) CreateSession(session *Session) error {
	archive_at := time.Now().Add(time.Minute * 24)
	query := "INSERT INTO admin_session (id, admin_id, archive_at, created_at) VALUES ($1, $2, $3, $4) RETURNING id"

    session.ID = uuid.New()
	session.CreatedAt = time.Now()

	err := sm.DB.QueryRow(query, session.ID, session.AdminID, archive_at, session.CreatedAt).Scan(&session.ID)

	if err != nil {
		return err
	}

	return nil
}

func (sm *SessionModel) GetSessionByID(id uuid.UUID) (*Session, error) {
	query := `
		SELECT admin_id,  archive_at, created_at
		FROM admin_session
		WHERE id = $1
	`

	session := &Session{}
	err := sm.DB.QueryRow(query, id).Scan(&session.AdminID,&session.Archive_at, &session.CreatedAt)
	if err != nil {
		return nil, err
	}

	session.ID = id
	return session, nil
}
func (sm *SessionModel) UpdateSession(session *Session)error{
	query := `
	    UPDATE admin_session
		SET archive_at = $1 
		WHERE id = $2
	`
	_, err := sm.DB.Exec(query,session.Archive_at,session.ID)
	return err
}

func (sm *SessionModel) ArchiveSession(id uuid.UUID) error{
	query := `UPDATE admin_session SET archive_at = $1 WHERE id = $2`
	_, err := sm.DB.Exec(query,time.Now(),id)
	if err != nil{
		return err
	}
	return nil
}
func (sm *SessionModel) DeleteSession(id uuid.UUID) error {
	query := `
		DELETE FROM admin_session
		WHERE id = $1
	`

	_, err := sm.DB.Exec(query, id)
	return err
}
