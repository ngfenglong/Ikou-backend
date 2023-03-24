package repository

import (
	"context"
	"database/sql"
	. "ikou/api/models"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

// #region Place API
func (m *DBModel) GetAllPlaces() ([]*Place, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var places []*Place

	query := `
	SELECT 
		p.id, p.placename, p.description, p.address, p.lat, p.lon, p.imageUrl, p.averageSpending, 
		s.decode, c.decode, p.created_at, p.updated_at, p.created_by
	FROM 
		places p
		Inner Join CodedecodeSubcategories s on s.code = p.subCategoryCode
		Inner Join codedecodeCategories c on c.code = s.categorycode
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
			&p.AverageSpending,
			&p.SubCategory,
			&p.Category,
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
			p.id, p.placename, p.description, p.address, p.lat, p.lon, p.imageUrl, p.averageSpending, 
			s.decode, c.decode, p.created_at, p.updated_at, p.created_by
		FROM 
			places p
			Inner Join CodedecodeSubcategories s on s.code = p.subCategoryCode
			Inner Join codedecodeCategories c on c.code = s.categorycode
		WHERE 
			p.id = ?`, id)

	err := row.Scan(
		&place.ID,
		&place.Name,
		&place.Description,
		&place.Address,
		&place.Lat,
		&place.Lon,
		&place.ImageUrl,
		&place.AverageSpending,
		&place.SubCategory,
		&place.Category,
		&place.CreatedAt,
		&place.UpdatedAt,
		&place.CreatedBy,
	)

	if err != nil {
		// Can add som error handling here
		return place, err
	}

	query := `
		select 
			id, rating, reviewDescription, created_at, updated_at, created_by
		from 
			reviews
		where
			place_id = ?
	`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return place, err
	}

	defer rows.Close()
	var reviews []*Review

	for rows.Next() {
		var r Review
		err = rows.Scan(
			&r.ID,
			&r.Rating,
			&r.ReviewDescription,
			&r.CreatedAt,
			&r.UpdatedAt,
			&r.CreatedBy,
		)
		if err != nil {
			return place, err
		}
		reviews = append(reviews, &r)
	}

	place.Reviews = reviews

	return place, nil
}

func (m *DBModel) GetPlacesBySubCategoryCode(code int) ([]*Place, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var places []*Place

	query := `
		SELECT 
			p.id, p.placename, p.description, p.address, p.lat, p.lon, p.imageUrl, p.averageSpending, 
			s.decode, p.created_at, p.updated_at, p.created_by
		FROM Places 
		Inner Join CodedecodeSubcategories s on s.code = p.subCategoryCode
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
			&p.AverageSpending,
			&p.SubCategory,
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
