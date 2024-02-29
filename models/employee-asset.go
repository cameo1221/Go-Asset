package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// EmployeeAsset represents the mapping between an employee and an asset
type EmployeeAsset struct {
	ID         uuid.UUID `json:"id"`
	AssetID    uuid.UUID `json:"asset_id"`
	EmployeeID uuid.UUID `json:"employee_id"`
	CreatedAt  time.Time `json:"created_at"`
	ArchivedAt time.Time `json:"archive_at,omitempty"`
}

// EmployeeAssetModel represents the model for employee-asset mapping operations
type EmployeeAssetModel struct {
	DB *sql.DB
}

// CreateEmployeeAsset creates a new employee-asset mapping in the database
func (eam *EmployeeAssetModel) CreateEmployeeAsset(employeeAsset *EmployeeAsset) error {
	query := `
		INSERT INTO employee_asset_mapping (id, asset_id, employee_id, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := eam.DB.QueryRow(query, employeeAsset.ID, employeeAsset.AssetID, employeeAsset.EmployeeID, employeeAsset.CreatedAt).Scan(&employeeAsset.ID)
	if err != nil {
		return err
	}

	return nil
}

// UpdateEmployeeAsset updates an existing employee-asset mapping in the database
func (eam *EmployeeAssetModel) UpdateEmployeeAsset(employeeAsset *EmployeeAsset) error {
	// Implement database operation to update the employee-asset mapping in the database
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

// ArchiveEmployeeAsset archives an existing employee-asset mapping in the database
func (eam *EmployeeAssetModel) ArchiveEmployeeAsset(id uuid.UUID) error {
	// Implement database operation to archive the employee-asset mapping in the database
	query := `
		UPDATE employee_asset_mapping
		SET archive_at = NOW()
		WHERE id = $1
	`

	_, err := eam.DB.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

// GetEmployeeAssetByID retrieves an employee-asset mapping from the database by its ID
func (eam *EmployeeAssetModel) GetEmployeeAssetByID(id uuid.UUID) (*EmployeeAsset, error) {
	// Implement database operation to retrieve the employee-asset mapping by ID from the database
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

// GetAllEmployeeAssets retrieves all employee-asset mappings from the database
func (eam *EmployeeAssetModel) GetAllEmployeeAssets() ([]*EmployeeAsset, error) {
	// Implement database operation to retrieve all employee-asset mappings from the database
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
