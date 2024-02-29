package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type EmployeeAsset struct {
	ID         uuid.UUID `json:"id"`
	AssetID    uuid.UUID `json:"asset_id"`
	EmployeeID uuid.UUID `json:"employee_id"`
	CreatedAt  time.Time `json:"created_at"`
	ArchivedAt *time.Time `json:"archive_at,omitempty"`
}

type EmployeeAssetModel struct {
	DB *sql.DB
}

func (eam *EmployeeAssetModel) CreateEmployeeAsset(employeeAsset *EmployeeAsset) error {
	query := `
		INSERT INTO employee_asset_mapping (id, asset_id, employee_id, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
    employeeAsset.ID = uuid.New()

	err := eam.DB.QueryRow(query, employeeAsset.ID, employeeAsset.AssetID, employeeAsset.EmployeeID, employeeAsset.CreatedAt).Scan(&employeeAsset.ID)
	if err != nil {
		return err
	}

	return nil
}

func (eam *EmployeeAssetModel) UpdateEmployeeAsset(employeeAsset *EmployeeAsset) error {
	query := `
		UPDATE employee_asset_mapping
		SET asset_id = $2, employee_id = $3, created_at = $4
		WHERE id = $1
	`

	_, err := eam.DB.Exec(query, employeeAsset.ID, employeeAsset.AssetID, employeeAsset.EmployeeID, employeeAsset.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (eam *EmployeeAssetModel) ArchiveEmployeeAsset(id uuid.UUID) error {
	query := `
		UPDATE employee_asset_mapping
		SET archive_at = $1
		WHERE id = $2
	`

	_, err := eam.DB.Exec(query, time.Now(), id)
	if err != nil {
		return err
	}

	return nil
}

func (eam *EmployeeAssetModel) GetEmployeeAssetByID(id uuid.UUID) (*EmployeeAsset, error) {
	query := `
		SELECT id, asset_id, employee_id, created_at, archive_at
		FROM employee_asset_mapping
		WHERE id = $1
	`

	employeeAsset := &EmployeeAsset{}
	err := eam.DB.QueryRow(query, id).Scan(&employeeAsset.ID, &employeeAsset.AssetID, &employeeAsset.EmployeeID, &employeeAsset.CreatedAt, &employeeAsset.ArchivedAt)
	if err != nil {
		return nil, err
	}

	return employeeAsset, nil
}

func (eam *EmployeeAssetModel) GetAllEmployeeAssets() ([]*EmployeeAsset, error) {
	query := `
		SELECT id, asset_id, employee_id, created_at, archive_at
		FROM employee_asset_mapping
	`

	rows, err := eam.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employeeAssets []*EmployeeAsset

	for rows.Next() {
		var employeeAsset EmployeeAsset
		err := rows.Scan(&employeeAsset.ID, &employeeAsset.AssetID, &employeeAsset.EmployeeID, &employeeAsset.CreatedAt, &employeeAsset.ArchivedAt)
		if err != nil {
			return nil, err
		}
		employeeAssets = append(employeeAssets, &employeeAsset)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return employeeAssets, nil
}
