package models

import "time"

type Review struct {
	ID                   string    `json:"id"`
	Rating               string    `json:"rating"`
	ReviewDescription    string    `json:"review_description"`
	ReviewerProfileImage string    `json:"reviewer_profile_image"`
	CreatedAt            time.Time `json:"-"`
	UpdatedAt            time.Time `json:"updated_at"`
	CreatedBy            string    `json:"reviewer_name"`
}
