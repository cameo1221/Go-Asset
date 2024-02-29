package models

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/google/uuid"
)

type Asset struct {
	Id         uuid.UUID  `json:"Id,omitempty" db:"Id"`
	Model      string     `json:"Model,omitempty" db:"Model"`
	Company    string     `json:"Company,omitempty" db:"Company"`
	CreatedAt  time.Time  `json:"createdAt,omitempty" db:"created_at"`
	ArchivedAt *time.Time `json:"archivedAt,omitempty" db:"archive_at"`
}

// AssetModel represents the model for asset operations
type AssetModel struct {
	DB *sql.DB
}

// CreateAsset creates a new asset in the database
func (am *AssetModel) CreateAsset(asset *Asset) error {
	// Generate a new UUID for the asset ID
	asset.Id = uuid.New()

	// Write the SQL statement for inserting an asset into the database
	err := am.DB.QueryRow("INSERT INTO asset (id, model, company, created_at) VALUES ($1, $2, $3, $4) RETURNING id",asset.Id, asset.Model, asset.Company,time.Now()).Scan(&asset.Id)

	// Execute the SQL statement and check for errors

	if err != nil {
		return fmt.Errorf("error creating asset: %w", err)
	}

	return nil
}

// UpdateAsset updates an existing asset in the database
func (am *AssetModel) UpdateAsset(asset *Asset) error {
	// Write the SQL statement for updating an asset in the database
	stmt := `UPDATE asset SET model = $1, company = $2 WHERE id = $3`

	// Execute the SQL statement
	_, err := am.DB.Exec(stmt, asset.Model, asset.Company, asset.Id)
	if err != nil {
		return err
	}

	return nil
}

// ArchiveAsset archives an existing asset in the database
func (am *AssetModel) ArchiveAsset(id uuid.UUID) error {
	// Write the SQL statement for archiving an asset in the database
	stmt := `UPDATE asset SET archive_at = $1 WHERE id = $2`

	// Execute the SQL statement
	_, err := am.DB.Exec(stmt, time.Now(), id)
	if err != nil {
		return err
	}

	return nil
}

// GetAssetByID retrieves an asset from the database by its ID
func (am *AssetModel) GetAssetByID(id uuid.UUID) (*Asset, error) {
	// Write the SQL statement for retrieving an asset by ID from the database
	stmt := `SELECT id, model, company, created_at, archive_at FROM asset WHERE id = $1`

	// Execute the SQL statement
	row := am.DB.QueryRow(stmt, id)

	// Create a new Asset struct to hold the retrieved data
	var asset Asset
	err := row.Scan(&asset.Id, &asset.Model, &asset.Company, &asset.CreatedAt, &asset.ArchivedAt)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// GetAllAssets retrieves all assets from the database
func (am *AssetModel) GetAllAssets() ([]*Asset, error) {
	// Write the SQL statement for retrieving all assets from the database
	stmt := `SELECT id, model, company, created_at, archive_at FROM asset`

	// Execute the SQL statement
	rows, err := am.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create a slice to hold the retrieved assets
	var assets []*Asset

	// Iterate through the rows and scan each row into an Asset struct
	for rows.Next() {
		var asset Asset
		err := rows.Scan(&asset.Id, &asset.Model, &asset.Company, &asset.CreatedAt, &asset.ArchivedAt)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return assets, nil
}
