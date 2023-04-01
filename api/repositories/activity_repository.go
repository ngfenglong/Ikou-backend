package repository

import (
	"context"
	"time"

	"github.com/ngfenglong/ikou-backend/api/models"
)

func (m *DBModel) GetActivityByPlace(placeId string) ([]*models.Activity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var activities []*models.Activity

	query := `
		SELECT a.id, a.activityname, a.hourDuration, a.created_at, a.updated_at, a.created_by, p.id, p.placeName, p.description, p.address, p.address, p.lat, p.lon, p.image_url, s.decode
		FROM Activities a 
		INNER JOIN Place p ON p.id = a.placeId
		Inner join CodedecodeSubcategories s on s.code = p.subCategoryCode
		WHERE placeId = ? 
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var a models.Activity
		err = rows.Scan(
			&a.ID,
			&a.ActvityName,
			&a.Duration,
			&a.Place,
			&a.CreatedAt,
			&a.UpdatedAt,
			&a.CreatedBy,
			&a.Place.ID,
			&a.Place.Name,
			&a.Place.Description,
			&a.Place.Address,
			&a.Place.Lat,
			&a.Place.Lon,
			&a.Place.ImageUrl,
			&a.Place.SubCategory,
		)
		if err != nil {
			return nil, err
		}

		activities = append(activities, &a)
	}
	return activities, nil

}
