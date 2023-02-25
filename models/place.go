package models

import (
	"context"
	"database/sql"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

type Place struct {
	ID              string    `json:"id"`
	Name            string    `json:"placeName"`
	Description     string    `json:"description"`
	Address         string    `json:"address"`
	Lat             string    `json:"lat"`
	Lon             string    `json:"lon"`
	ImageUrl        string    `json:"image_url"`
	SubCategoryCode int       `json:"subcategorycode"`
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedBy       string    `json:"created_by"`
}

// #region Place API
func (m *DBModel) GetAllPlaces() ([]*Place, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var places []*Place

	query := `
		SELECT id, placeName, description, address, lat, lon, image_url, subCategoryCode, created_at, updated_at, created_by
		FROM Places 
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var p Place
		err = rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Address,
			&p.Lat,
			&p.Lon,
			&p.ImageUrl,
			&p.SubCategoryCode,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.CreatedBy,
		)
		if err != nil {
			return nil, err
		}

		places = append(places, &p)
	}
	return places, nil
}

func (m *DBModel) GetPlaceById(id string) (Place, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var place Place

	row := m.DB.QueryRowContext(ctx, `
		SELECT 
			id, placename, description,address, lat, lon, image_url, 
			subcategorycode, created_at, updated_at, created_by
		FROM 
			places
		WHERE 
			id = ?`, id)

	err := row.Scan(
		&place.ID,
		&place.Name,
		&place.Description,
		&place.Address,
		&place.Lat,
		&place.Lon,
		&place.ImageUrl,
		&place.SubCategoryCode,
		&place.CreatedAt,
		&place.UpdatedAt,
		&place.CreatedBy,
	)

	if err != nil {
		// Can add som error handling here
		return place, err
	}

	return place, nil
}

func (m *DBModel) GetPlacesBySubCategoryCode(code int) ([]*Place, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var places []*Place

	query := `
		SELECT id, placeName, description, address, lat, lon, image_url, subCategoryCode, created_at, updated_at, created_by
		FROM Places 
		WHERE subCategoryCode = ?
	`

	rows, err := m.DB.QueryContext(ctx, query, code)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var p Place
		err = rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Address,
			&p.Lat,
			&p.Lon,
			&p.ImageUrl,
			&p.SubCategoryCode,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.CreatedBy,
		)
		if err != nil {
			return nil, err
		}

		places = append(places, &p)
	}
	return places, nil
}

//	#endregion
