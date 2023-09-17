package repository

import (
	"context"
	"time"

	"github.com/ngfenglong/ikou-backend/api/models"
)

func (m *DBModel) GetAllCategory() ([]*models.CodeDecodeCategory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	categoriesMap := make(map[string]*models.CodeDecodeCategory)

	query := `
		Select c.id, c.code, c.decode, c.isActive, c.created_at, c.updated_at, 
		s.id, s.code, s.decode,s.isActive, s.created_at, s.updated_at  
		FROM CodeDecodeCategories c
		INNER JOIN CodeDecodeSubCategories s ON c.code = s.categoryCode
		WHERE c.isActive = 1 AND s.isActive = 1
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var c models.CodeDecodeCategory
		var s models.CodeDecodeSubCategory

		err = rows.Scan(
			&c.ID,
			&c.Code,
			&c.Decode,
			&c.IsActive,
			&c.CreatedAt,
			&c.UpdatedAt,
			&s.ID,
			&s.Code,
			&s.Decode,
			&s.IsActive,
			&s.CreatedAt,
			&s.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		category, exists := categoriesMap[c.ID]
		if !exists {
			category = &c
			category.SubCategories = []*models.CodeDecodeSubCategory{}
			categoriesMap[c.ID] = category
		}
		category.SubCategories = append(category.SubCategories, &s)
	}
	categories := make([]*models.CodeDecodeCategory, 0, len(categoriesMap))
	for _, category := range categoriesMap {
		categories = append(categories, category)
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

func (m *DBModel) GetAllAreas() ([]*models.CodeDecodeArea, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var areas []*models.CodeDecodeArea

	query := `
				Select id, code, decode, isActive, created_at, updated_at
    			From CodeDecodeAreas
			`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var a models.CodeDecodeArea

		err = rows.Scan(
			&a.ID,
			&a.Code,
			&a.Decode,
			&a.IsActive,
			&a.CreatedAt,
			&a.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		areas = append(areas, &a)
	}

	return areas, nil
}
