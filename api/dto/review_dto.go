package dto

import "time"

type ReviewDTO struct {
	ID                   string    `json:"id"`
	Rating               int       `json:"rating"`
	ReviewDescription    string    `json:"review_description"`
	ReviewerProfileImage string    `json:"reviewer_profile_image"`
	CreatedAt            time.Time `json:"-"`
	UpdatedAt            time.Time `json:"updated_at"`
	CreatedBy            string    `json:"reviewer_name"`
}
