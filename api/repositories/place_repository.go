package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/ngfenglong/ikou-backend/api/models"
	"github.com/ngfenglong/ikou-backend/internal/util"
)

// #region Place API
func (m *DBModel) GetAllPlaces() ([]*models.Place, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	placesMap := make(map[string]*models.Place)

	query := `
	SELECT 
			p.id, p.placename, p.description, p.address, p.lat, p.lon, p.imageUrl, p.averageSpending, 
			s.decode, c.decode, p.created_at, p.updated_at, p.created_by,
			r.id, r.rating, r.reviewDescription, u.profileImage, r.created_at, r.updated_at, u.username
		FROM 
			Places p
			Inner Join CodeDecodeSubcategories s on s.code = p.subCategoryCode
			Inner Join CodeDecodeCategories c on c.code = s.categorycode
			Left Join Reviews r on p.id = r.place_id
			Left Join Users u on u.id = r.created_by
		ORDER BY p.id
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var p models.Place
		var r models.Review

		var rID sql.NullString
		var rRating sql.NullString
		var rReviewDescription sql.NullString
		var rReviewerProfileImage sql.NullString
		var rCreatedAt sql.NullTime
		var rUpdatedAt sql.NullTime
		var rCreatedBy sql.NullString

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
			&rID,
			&rRating,
			&rReviewDescription,
			&rReviewerProfileImage,
			&rCreatedAt,
			&rUpdatedAt,
			&rCreatedBy,
		)
		if err != nil {
			return nil, err
		}

		place, exists := placesMap[p.ID]
		if !exists {
			place = &p
			place.Reviews = []*models.Review{}
			placesMap[p.ID] = place
		}

		if rID.Valid && rRating.Valid && rCreatedBy.Valid {
			r.ID = rID.String
			r.Rating = rRating.String
			r.ReviewDescription = util.CoalesceNullString(rReviewDescription)
			r.ReviewerProfileImage = util.CoalesceNullString(rReviewerProfileImage)
			r.CreatedAt = rCreatedAt.Time
			r.UpdatedAt = rUpdatedAt.Time
			r.CreatedBy = rCreatedBy.String
			place.Reviews = append(place.Reviews, &r)
		}
	}

	places := make([]*models.Place, 0, len(placesMap))
	for _, place := range placesMap {
		places = append(places, place)
	}

	return places, nil
}

func (m *DBModel) GetPlaceById(id string) (models.Place, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var place models.Place

	row := m.DB.QueryRowContext(ctx, `
		SELECT 
			p.id, p.placename, p.description, p.address, p.lat, p.lon, p.imageUrl, p.averageSpending, 
			s.decode, c.decode, p.created_at, p.updated_at, p.created_by
		FROM 
			Places p
			Inner Join CodeDecodeSubcategories s on s.code = p.subCategoryCode
			Inner Join CodeDecodeCategories c on c.code = s.categorycode
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
			r.id, r.rating, r.reviewDescription, u.profileImage, r.created_at, r.updated_at, u.username
		from 
			Reviews r
		inner join 
			Users u on u.id = r.created_by
		where
			place_id = ?
	`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return place, err
	}

	defer rows.Close()
	var reviews []*models.Review

	for rows.Next() {
		var r models.Review
		err = rows.Scan(
			&r.ID,
			&r.Rating,
			&r.ReviewDescription,
			&r.ReviewerProfileImage,
			&r.CreatedAt,
			&r.UpdatedAt,
			&r.CreatedBy,
		)
		if err != nil {
			return place, err
		}
		reviews = append(reviews, &r)
	}

	if len(reviews) == 0 {
		place.Reviews = []*models.Review{}
	} else {
		place.Reviews = reviews
	}

	return place, nil
}

func (m *DBModel) GetPlacesByCategoryCode(category string) ([]*models.Place, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	placesMap := make(map[string]*models.Place)

	query := `
		SELECT 
			p.id, p.placename, p.description, p.address, p.lat, p.lon, p.imageUrl, p.averageSpending, 
			s.decode, c.decode, p.created_at, p.updated_at, p.created_by,
			r.id, r.rating, r.reviewDescription, u.profileImage, r.created_at, r.updated_at, u.username
		FROM Places p
		Inner Join CodeDecodeSubcategories s on s.code = p.subCategoryCode 
		INNER JOIN CodeDecodeCategories c on c.code = s.categoryCode
		Left Join Reviews r on p.id = r.place_id
		Left Join Users u on u.id = r.created_by
		WHERE c.decode = ?
		Order by p.id
	`

	rows, err := m.DB.QueryContext(ctx, query, category)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var p models.Place
		var r models.Review

		var rID sql.NullString
		var rRating sql.NullString
		var rReviewDescription sql.NullString
		var rReviewerProfileImage sql.NullString
		var rCreatedAt sql.NullTime
		var rUpdatedAt sql.NullTime
		var rCreatedBy sql.NullString

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
			&rID,
			&rRating,
			&rReviewDescription,
			&rReviewerProfileImage,
			&rCreatedAt,
			&rUpdatedAt,
			&rCreatedBy,
		)

		if err != nil {
			return nil, err
		}

		place, exists := placesMap[p.ID]
		if !exists {
			place = &p
			place.Reviews = []*models.Review{}
			placesMap[p.ID] = place
		}

		if rID.Valid && rRating.Valid && rCreatedBy.Valid {
			r.ID = rID.String
			r.Rating = rRating.String
			r.ReviewDescription = util.CoalesceNullString(rReviewDescription)
			r.ReviewerProfileImage = util.CoalesceNullString(rReviewerProfileImage)
			r.CreatedAt = rCreatedAt.Time
			r.UpdatedAt = rUpdatedAt.Time
			r.CreatedBy = rCreatedBy.String
			place.Reviews = append(place.Reviews, &r)
		}
	}

	places := make([]*models.Place, 0, len(placesMap))
	for _, place := range placesMap {
		places = append(places, place)
	}

	return places, nil
}

func (m *DBModel) GetPlacesBySubCategoryCode(code int) ([]*models.Place, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var places []*models.Place

	query := `
		SELECT 
			p.id, p.placename, p.description, p.address, p.lat, p.lon, p.imageUrl, p.averageSpending, 
			s.decode, p.created_at, p.updated_at, p.created_by
		FROM Places p
		Inner Join CodeDecodeSubcategories s on s.code = p.subCategoryCode
		WHERE subCategoryCode = ?
	`

	rows, err := m.DB.QueryContext(ctx, query, code)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var p models.Place
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

func (m *DBModel) SearchPlaceByKeyword(keyword string) ([]*models.Place, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var places []*models.Place

	addressKeyword := "%" + keyword + "%"
	placenameKeyword := "%" + keyword + "%"

	query := `
		SELECT 
			p.id, p.placename, p.description, p.address, p.lat, p.lon, p.imageUrl, p.averageSpending, 
			s.decode, p.created_at, p.updated_at, p.created_by
		FROM Places p
		Inner Join CodeDecodeSubcategories s on s.code = p.subCategoryCode
		WHERE p.address like ? OR p.placename like ?
	`

	rows, err := m.DB.QueryContext(ctx, query, addressKeyword, placenameKeyword)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var p models.Place
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
