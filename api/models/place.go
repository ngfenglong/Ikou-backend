package models

import (
	"time"
)

type Place struct {
	ID              string    `json:"id"`
	Name            string    `json:"placeName"`
	Description     string    `json:"description"`
	Address         string    `json:"address"`
	Lat             string    `json:"lat"`
	Lon             string    `json:"lon"`
	ImageUrl        string    `json:"image_url"`
	SubCategoryCode int       `json:"subcategorycode"`
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedBy       string    `json:"created_by"`
}
