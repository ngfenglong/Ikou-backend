package repository

import (
	"context"
	"time"

	"ikou/api/models"
)

func (m *DBModel) GetAllCategory() ([]*models.CodeDecodeCategory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var categories []*models.CodeDecodeCategory

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
		var c models.CodeDecodeCategory
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

func (m *DBModel) GetAllSubCategory() ([]*models.CodeDecodeSubCategory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var subCategories []*models.CodeDecodeSubCategory

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
		var c models.CodeDecodeSubCategory
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

func (m *DBModel) GetAllSubCategoryByCategoryCode(categoryCode int) ([]*models.CodeDecodeSubCategory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var subCategories []*models.CodeDecodeSubCategory

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
		var c models.CodeDecodeSubCategory
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
