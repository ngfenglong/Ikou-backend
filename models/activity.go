package models

import (
	"context"
	"time"
)

type Activity struct {
	ID          string    `json:"id"`
	ActvityName string    `json:"activityName"`
	Duration    int       `json:"duration"`
	Place       Place     `json:"place"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   string    `json:"created_by"`
}

func (m *DBModel) GetActivityByPlace(placeId string) ([]*Activity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var activities []*Activity

	query := `
		SELECT a.id, a.activityname, a.hourDuration, a.created_at, a.updated_at, a.created_by, p.id, p.placeName, p.description, p.address, p.address, p.lat, p.lon, p.image_url, p.subCategoryCode
		FROM Activities a 
		INNER JOIN Place p
		ON p.id = a.placeId
		WHERE placeId = ? 
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var a Activity
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
			&a.Place.SubCategoryCode,
		)
		if err != nil {
			return nil, err
		}

		activities = append(activities, &a)
	}
	return activities, nil

}
