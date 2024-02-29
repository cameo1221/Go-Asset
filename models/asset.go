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

type AssetModel struct {
	DB *sql.DB
}

func (am *AssetModel) CreateAsset(asset *Asset) error {
	asset.Id = uuid.New()
	err := am.DB.QueryRow("INSERT INTO asset (id, model, company, created_at) VALUES ($1, $2, $3, $4) RETURNING id",asset.Id, asset.Model, asset.Company,time.Now()).Scan(&asset.Id)

	if err != nil {
		return fmt.Errorf("error creating asset: %w", err)
	}

	return nil
}

func (am *AssetModel) UpdateAsset(asset *Asset) error {
	stmt := `UPDATE asset SET model = $1, company = $2 WHERE id = $3`

	_, err := am.DB.Exec(stmt, asset.Model, asset.Company, asset.Id)
	if err != nil {
		return err
	}

	return nil
}

func (am *AssetModel) ArchiveAsset(id uuid.UUID) error {
	stmt := `UPDATE asset SET archive_at = $1 WHERE id = $2`

	_, err := am.DB.Exec(stmt, time.Now(), id)
	if err != nil {
		return err
	}

	return nil
}

func (am *AssetModel) GetAssetByID(id uuid.UUID) (*Asset, error) {
	stmt := `SELECT id, model, company, created_at, archive_at FROM asset WHERE id = $1`

	row := am.DB.QueryRow(stmt, id)

	var asset Asset
	err := row.Scan(&asset.Id, &asset.Model, &asset.Company, &asset.CreatedAt, &asset.ArchivedAt)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (am *AssetModel) GetAllAssets() ([]*Asset, error) {
	stmt := `SELECT id, model, company, created_at, archive_at FROM asset`
	rows, err := am.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []*Asset

	for rows.Next() {
		var asset Asset
		err := rows.Scan(&asset.Id, &asset.Model, &asset.Company, &asset.CreatedAt, &asset.ArchivedAt)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return assets, nil
}
