package dto

import "time"

type PlaceDTO struct {
	ID              string       `json:"id"`
	Name            string       `json:"placeName"`
	Description     string       `json:"description"`
	Address         string       `json:"address"`
	Lat             string       `json:"lat"`
	Lon             string       `json:"lon"`
	AverageSpending int          `json:"average_spending"`
	AverageRating   int          `json:"average_rating"`
	ImageUrl        string       `json:"image_url"`
	SubCategory     string       `json:"sub_category"`
	Category        string       `json:"category"`
	Area            string       `json:"area"`
	Reviews         []*ReviewDTO `json:"reviews"`
	CreatedAt       time.Time    `json:"-"`
	UpdatedAt       time.Time    `json:"updated_at"`
	CreatedBy       string       `json:"created_by"`
}
