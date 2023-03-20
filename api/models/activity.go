package models

import (
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
