package models

import "time"

type Review struct {
	ID                string    `json:"id"`
	Rating            string    `json:"placeName"`
	ReviewDescription string    `json:"description"`
	CreatedAt         time.Time `json:"-"`
	UpdatedAt         time.Time `json:"updated_at"`
	CreatedBy         string    `json:"created_by"`
}
