package models

import (
	"time"
)

type CodeDecodeCategory struct {
	ID            string                   `json:"id"`
	Code          int                      `json:"code"`
	Decode        string                   `json:"decode"`
	IsActive      bool                     `json:"is_active"`
	SubCategories []*CodeDecodeSubCategory `json:"sub_cateogries"`
	CreatedAt     time.Time                `json:"-"`
	UpdatedAt     time.Time                `json:"-"`
}

type CodeDecodeSubCategory struct {
	ID           string    `json:"id"`
	Code         int       `json:"code"`
	Decode       string    `json:"decode"`
	IsActive     bool      `json:"is_active"`
	CategoryCode int       `json:"category_code"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}
