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
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (m *DBModel) GetPlaceById(id int) (Place, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var place Place

	row := m.DB.QueryRowContext(ctx, `
		SELECT 
			id, name, description 
		FROM 
			places
		WHERE 
			id = ?`, id)

	err := row.Scan(
		&place.ID,
		&place.Name,
		&place.Description,
	)

	if err != nil {
		// Can add som error handling here
		return place, err
	}

	return place, nil
}
