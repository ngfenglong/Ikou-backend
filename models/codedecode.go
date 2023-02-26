package models

import (
	"context"
	"time"
)

type CodeDecodeCategory struct {
	ID        string    `json:"id"`
	Code      int       `json:"code"`
	Decode    string    `json:"decode"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type CodeDecodeSubCategory struct {
	ID           string    `json:"id"`
	Code         int       `json:"code"`
	Decode       string    `json:"decode"`
	IsActive     bool      `json:"is_active"`
	CategoryCode int       `json:"category_code"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

func (m *DBModel) GetAllCategory() ([]*CodeDecodeCategory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var categories []*CodeDecodeCategory

	query := `
		Select id, code, decode, isActive, created_at, updated_at 
		FROM CodeDecodeCategories 
		WHERE isActive = 1
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var c CodeDecodeCategory
		err = rows.Scan(
			&c.ID,
			&c.Code,
			&c.Decode,
			&c.IsActive,
			&c.CreatedAt,
			&c.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		categories = append(categories, &c)
	}

	return categories, nil
}

func (m *DBModel) GetAllSubCategory() ([]*CodeDecodeSubCategory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var subCategories []*CodeDecodeSubCategory

	query := `
		Select id, code, decode, isActive, categoryCode, created_at, updated_at 
		FROM CodeDecodeSubcategories 
		WHERE isActive = 1
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var c CodeDecodeSubCategory
		err = rows.Scan(
			&c.ID,
			&c.Code,
			&c.Decode,
			&c.IsActive,
			&c.CategoryCode,
			&c.CreatedAt,
			&c.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		subCategories = append(subCategories, &c)
	}

	return subCategories, nil
}

func (m *DBModel) GetAllSubCategoryByCategoryCode(categoryCode int) ([]*CodeDecodeSubCategory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var subCategories []*CodeDecodeSubCategory

	query := `
		Select id, code, decode, isActive, categoryCode, created_at, updated_at 
		FROM CodeDecodeSubcategories 
		WHERE isActive = 1 
		AND categoryCode = ?
	`
	rows, err := m.DB.QueryContext(ctx, query, categoryCode)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var c CodeDecodeSubCategory
		err = rows.Scan(
			&c.ID,
			&c.Code,
			&c.Decode,
			&c.IsActive,
			&c.CategoryCode,
			&c.CreatedAt,
			&c.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		subCategories = append(subCategories, &c)
	}

	return subCategories, nil
}
