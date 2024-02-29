package models

import (
	"database/sql"
	"time"
	"fmt"

	"github.com/google/uuid"
)

type Admin struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Password   string    `json:"-"`
	CreatedAt  time.Time `json:"created_at"`
	ArchivedAt *time.Time `json:"archive_at,omitempty"`
}

type AdminModel struct {
	DB *sql.DB
}

func (am *AdminModel) CreateAdmin(admin *Admin) error {
	admin.ID = uuid.New()

	err := am.DB.QueryRow("INSERT INTO admin (id, name, email, Password, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id",admin.ID, admin.Name, admin.Email, admin.Password, time.Now()).Scan(&admin.ID)


	if err != nil {
		return fmt.Errorf("error creating admin: %w", err)
	}

	return nil
}


func (am *AdminModel) UpdateAdmin(admin *Admin) error {
	query := `
		UPDATE admin
		SET name = $1, email = $2, password = $3, archive_at = $4
		WHERE id = $5
	`
	
	_, err := am.DB.Exec(query, admin.Name, admin.Email, admin.Password, admin.ArchivedAt, admin.ID)
	return err
}

func (am *AdminModel) ArchiveAdmin(id uuid.UUID) error {
	query := `
		UPDATE admin
		SET archive_at = $1
		WHERE id = $2
	`

	_, err := am.DB.Exec(query, time.Now(), id)
	return err
}

func (am *AdminModel) GetAdminByID(id uuid.UUID) (*Admin, error) {
	query := `
		SELECT id, name, email, password, created_at, archive_at
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

func (am *AdminModel) GetAllAdmins() ([]*Admin, error) {
	query := `
		SELECT id, name, email, password, created_at, archive_at
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
