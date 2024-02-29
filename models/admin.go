package models

import (
	"database/sql"
	"time"
	"fmt"

	"github.com/google/uuid"
)

// Admin represents an admin user in the system
type Admin struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Password   string    `json:"-"`
	CreatedAt  time.Time `json:"created_at"`
	ArchivedAt time.Time `json:"archived_at,omitempty"`
}

// AdminModel represents the model for admin operations
type AdminModel struct {
	DB *sql.DB
}

// CreateAdmin creates a new admin user in the database
func (am *AdminModel) CreateAdmin(admin *Admin) error {
	// Generate a new UUID for the asset ID
	admin.ID = uuid.New()

	// Write the SQL statement for inserting an asset into the database
	err := am.DB.QueryRow("INSERT INTO admin (id, name, email, Password, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id",admin.ID, admin.Name, admin.Email, admin.Password, time.Now()).Scan(&admin.ID)

	// Execute the SQL statement and check for errors

	if err != nil {
		return fmt.Errorf("error creating admin: %w", err)
	}

	return nil
}


// UpdateAdmin updates an existing admin user in the database
func (am *AdminModel) UpdateAdmin(admin *Admin) error {
	query := `
		UPDATE admin
		SET name = $1, email = $2, password = $3, archived_at = $4
		WHERE id = $5
	`
	
	_, err := am.DB.Exec(query, admin.Name, admin.Email, admin.Password, admin.ArchivedAt, admin.ID)
	return err
}

// ArchiveAdmin archives an existing admin user in the database
func (am *AdminModel) ArchiveAdmin(id uuid.UUID) error {
	query := `
		UPDATE admin
		SET archived_at = $1
		WHERE id = $2
	`

	_, err := am.DB.Exec(query, time.Now(), id)
	return err
}

// GetAdminByID retrieves an admin user from the database by its ID
func (am *AdminModel) GetAdminByID(id uuid.UUID) (*Admin, error) {
	query := `
		SELECT id, name, email, password, created_at, archived_at
		FROM admin
		WHERE id = $1
	`

	admin := &Admin{}
	err := am.DB.QueryRow(query, id).Scan(&admin.ID, &admin.Name, &admin.Email, &admin.Password, &admin.CreatedAt, &admin.ArchivedAt)
	if err != nil {
		return nil, err
	}

	return admin, nil
}

// GetAllAdmins retrieves all admin users from the database
func (am *AdminModel) GetAllAdmins() ([]*Admin, error) {
	query := `
		SELECT id, name, email, password, created_at, archived_at
		FROM admin
	`

	rows, err := am.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var admins []*Admin
	for rows.Next() {
		admin := &Admin{}
		err := rows.Scan(&admin.ID, &admin.Name, &admin.Email, &admin.Password, &admin.CreatedAt, &admin.ArchivedAt)
		if err != nil {
			return nil, err
		}
		admins = append(admins, admin)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return admins, nil
}
