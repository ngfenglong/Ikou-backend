package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ngfenglong/ikou-backend/api/dto"
	"github.com/ngfenglong/ikou-backend/api/mapper"
	"github.com/ngfenglong/ikou-backend/api/models"
	"github.com/ngfenglong/ikou-backend/internal/util"
)

// #region Place API
func (m *DBModel) GetAllPlaces(userID string) ([]*dto.PlaceDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	placesMap := make(map[string]*models.Place)

	query := `
	SELECT 
			p.id, p.placename, p.description, p.address, p.lat, p.lon, p.imageUrl, p.averageSpending, 
			s.decode, c.decode, a.decode, p.created_at, p.updated_at, p.created_by,
			r.id, r.rating, r.reviewDescription, u.profileImage, r.created_at, r.updated_at, u.username, 
			CASE WHEN lp.id IS NOT NULL THEN true ELSE false END AS liked
		FROM 
			Places p
			Inner Join CodeDecodeSubcategories s on s.code = p.subCategoryCode
			Inner Join CodeDecodeCategories c on c.code = s.categorycode
			Inner Join CodeDecodeAreas a on a.code = p.areaCode
			Left Join Reviews r on p.id = r.place_id
			Left Join Users u on u.id = r.created_by
			LEFT JOIN Liked_Places lp on p.id = lp.placeId AND lp.userId = ? 
		ORDER BY p.id
	`

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var p models.Place
		var r models.Review
		var liked sql.NullBool

		var rID sql.NullString
		var rRating sql.NullInt32
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
			&p.Area,
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
			&liked,
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
			r.Rating = util.CoalesceNullInt(rRating)
			r.ReviewDescription = util.CoalesceNullString(rReviewDescription)
			r.ReviewerProfileImage = util.CoalesceNullString(rReviewerProfileImage)
			r.CreatedAt = rCreatedAt.Time
			r.UpdatedAt = rUpdatedAt.Time
			r.CreatedBy = rCreatedBy.String
			place.Reviews = append(place.Reviews, &r)
		}

		place.Liked = liked.Valid && liked.Bool
	}

	places := make([]*dto.PlaceDTO, 0, len(placesMap))
	for _, place := range placesMap {
		places = append(places, mapper.MapToPlaceDTO(place))
	}

	return places, nil
}

func (m *DBModel) GetPlaceById(id string, userID string) (*dto.PlaceDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var place models.Place
	var liked sql.NullBool

	row := m.DB.QueryRowContext(ctx, `
		SELECT 
			p.id, p.placename, p.description, p.address, p.lat, p.lon, p.imageUrl, p.averageSpending, 
			s.decode, c.decode, a.decode, p.created_at, p.updated_at, p.created_by, 
			CASE WHEN lp.id IS NOT NULL THEN true ELSE false END AS liked
		FROM 
			Places p
			Inner Join CodeDecodeSubcategories s on s.code = p.subCategoryCode
			Inner Join CodeDecodeCategories c on c.code = s.categorycode
			Inner Join CodeDecodeAreas a on a.code = p.areaCode
			LEFT JOIN Liked_Places lp on p.id = lp.placeId AND lp.userId = ? 
		WHERE 
			p.id = ?`, userID, id)

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
		&place.Area,
		&place.Category,
		&place.CreatedAt,
		&place.UpdatedAt,
		&place.CreatedBy,
		&liked,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("record not found")
		} else {
			return nil, err
		}
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
		return nil, err
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
			return nil, err
		}
		reviews = append(reviews, &r)
	}

	if len(reviews) == 0 {
		place.Reviews = []*models.Review{}
	} else {
		place.Reviews = reviews
	}

	place.Liked = liked.Valid && liked.Bool

	return mapper.MapToPlaceDTO(&place), nil
}

func (m *DBModel) GetPlacesByCategoryCode(category string, userID string) ([]*dto.PlaceDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	placesMap := make(map[string]*models.Place)

	query := `
		SELECT 
			p.id, p.placename, p.description, p.address, p.lat, p.lon, p.imageUrl, p.averageSpending, 
			s.decode, c.decode, a.decode, p.created_at, p.updated_at, p.created_by,
			r.id, r.rating, r.reviewDescription, u.profileImage, r.created_at, r.updated_at, u.username, 
			CASE WHEN lp.id IS NOT NULL THEN true ELSE false END AS liked
		FROM Places p
		Inner Join CodeDecodeSubcategories s on s.code = p.subCategoryCode 
		INNER JOIN CodeDecodeCategories c on c.code = s.categoryCode
		Inner Join CodeDecodeAreas a on a.code = p.areaCode
		Left Join Reviews r on p.id = r.place_id
		Left Join Users u on u.id = r.created_by
		LEFT JOIN Liked_Places lp on p.id = lp.placeId AND lp.userId = ? 
		WHERE c.decode = ?
		Order by p.id
	`

	rows, err := m.DB.QueryContext(ctx, query, userID, category)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var p models.Place
		var r models.Review
		var liked sql.NullBool

		var rID sql.NullString
		var rRating sql.NullInt32
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
			&p.Area,
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
			&liked,
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
			r.Rating = util.CoalesceNullInt(rRating)
			r.ReviewDescription = util.CoalesceNullString(rReviewDescription)
			r.ReviewerProfileImage = util.CoalesceNullString(rReviewerProfileImage)
			r.CreatedAt = rCreatedAt.Time
			r.UpdatedAt = rUpdatedAt.Time
			r.CreatedBy = rCreatedBy.String
			place.Reviews = append(place.Reviews, &r)
		}

		place.Liked = liked.Valid && liked.Bool
	}

	places := make([]*dto.PlaceDTO, 0, len(placesMap))
	for _, place := range placesMap {
		places = append(places, mapper.MapToPlaceDTO(place))
	}

	return places, nil
}

func (m *DBModel) GetPlacesBySubCategoryCode(code int, userID string) ([]*dto.PlaceDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	placesMap := make(map[string]*models.Place)

	query := `
		SELECT 
		p.id, p.placename, p.description, p.address, p.lat, p.lon, p.imageUrl, p.averageSpending, 
		s.decode, c.decode, a.decode, p.created_at, p.updated_at, p.created_by,
		r.id, r.rating, r.reviewDescription, u.profileImage, r.created_at, r.updated_at, u.username, 
		CASE WHEN lp.id IS NOT NULL THEN true ELSE false END AS liked
	FROM 
		Places p
		Inner Join CodeDecodeSubcategories s on s.code = p.subCategoryCode
		Inner Join CodeDecodeCategories c on c.code = s.categorycode
		Inner Join CodeDecodeAreas a on a.code = p.areaCode
		Left Join Reviews r on p.id = r.place_id
		Left Join Users u on u.id = r.created_by
		LEFT JOIN Liked_Places lp on p.id = lp.placeId AND lp.userId = ? 
	WHERE subCategoryCode = ?
	ORDER BY p.id
	`

	rows, err := m.DB.QueryContext(ctx, query, userID, code)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var p models.Place
		var r models.Review
		var liked sql.NullBool

		var rID sql.NullString
		var rRating sql.NullInt32
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
			&p.Area,
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
			&liked,
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
			r.Rating = util.CoalesceNullInt(rRating)
			r.ReviewDescription = util.CoalesceNullString(rReviewDescription)
			r.ReviewerProfileImage = util.CoalesceNullString(rReviewerProfileImage)
			r.CreatedAt = rCreatedAt.Time
			r.UpdatedAt = rUpdatedAt.Time
			r.CreatedBy = rCreatedBy.String
			place.Reviews = append(place.Reviews, &r)
		}
		place.Liked = liked.Valid && liked.Bool
	}

	places := make([]*dto.PlaceDTO, 0, len(placesMap))
	for _, place := range placesMap {
		places = append(places, mapper.MapToPlaceDTO(place))
	}

	return places, nil
}

func (m *DBModel) SearchPlaceByKeyword(keyword string, userID string) ([]*dto.PlaceDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	placesMap := make(map[string]*models.Place)

	addressKeyword := "%" + keyword + "%"
	placenameKeyword := "%" + keyword + "%"

	query := `
		SELECT 
			p.id, p.placename, p.description, p.address, p.lat, p.lon, p.imageUrl, p.averageSpending, 
			s.decode, c.decode, a.decode, p.created_at, p.updated_at, p.created_by,
			r.id, r.rating, r.reviewDescription, u.profileImage, r.created_at,  r.updated_at, u.username, 
			CASE WHEN lp.id IS NOT NULL THEN true ELSE false END AS liked
		FROM 
			Places p
			Inner Join CodeDecodeSubcategories s on s.code = p.subCategoryCode
			Inner Join CodeDecodeCategories c on c.code = s.categorycode
			Inner Join CodeDecodeAreas a on a.code = p.areaCode
			Left Join Reviews r on p.id = r.place_id
			Left Join Users u on u.id = r.created_by
			LEFT JOIN Liked_Places lp on p.id = lp.placeId AND lp.userId = ? 
		WHERE p.address like ? OR p.placename like ?
	`

	rows, err := m.DB.QueryContext(ctx, query, userID, addressKeyword, placenameKeyword)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var p models.Place
		var r models.Review
		var liked sql.NullBool

		var rID sql.NullString
		var rRating sql.NullInt32
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
			&p.Area,
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
			&liked,
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
			r.Rating = util.CoalesceNullInt(rRating)
			r.ReviewDescription = util.CoalesceNullString(rReviewDescription)
			r.ReviewerProfileImage = util.CoalesceNullString(rReviewerProfileImage)
			r.CreatedAt = rCreatedAt.Time
			r.UpdatedAt = rUpdatedAt.Time
			r.CreatedBy = rCreatedBy.String
			place.Reviews = append(place.Reviews, &r)
		}

		place.Liked = liked.Valid && liked.Bool
	}

	places := make([]*dto.PlaceDTO, 0, len(placesMap))
	for _, place := range placesMap {
		places = append(places, mapper.MapToPlaceDTO(place))
	}

	return places, nil
}

func (m *DBModel) AddPlaceRequest(pr dto.PlaceRequestDto) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	stmt := `
		Insert into Place_Requests (placeName, description, address, imageUrl, subCategoryCode, areaCode, created_by) 
		VALUE (?,?,?,?,?,?,?)
	`

	_, err := m.DB.ExecContext(ctx, stmt, pr.Name, pr.Description, pr.Address, pr.ImageUrl, pr.SubCategory, pr.Area, pr.CreatedBy)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) HasUserLikedPlace(userID string, placeID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	checkLikedStmt := `
		Select Count(*) 
		From Liked_Places 
		WHERE userId = ? AND placeId = ?  
	`

	row := m.DB.QueryRowContext(ctx, checkLikedStmt, userID, placeID)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (m *DBModel) RemoveUserLikeFromPlace(userID string, placeID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `Delete From Liked_Places WHERE userId = ? AND placeId = ?`

	_, err := m.DB.ExecContext(ctx, stmt, userID, placeID)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) AddUserLikeToPlace(userID string, placeID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT into Liked_Places (userId, placeId, created_by) VALUES (?,?,?)`

	_, err := m.DB.ExecContext(ctx, stmt, userID, placeID, userID)
	if err != nil {
		return err
	}

	return nil
}

//	#endregion
