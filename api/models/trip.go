package models

import "time"

type Trip struct {
	ID       string `json:""`
	Tripname string `json:""`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy string    `json:"created_by"`
}
