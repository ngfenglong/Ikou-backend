package mapper_test

import (
	"testing"

	"github.com/ngfenglong/ikou-backend/api/mapper"
	"github.com/ngfenglong/ikou-backend/api/models"
)

func TestMapToPlaceDTO_NoReviews_AverageRatingZero(t *testing.T) {
	place := &models.Place{
		ID:      "p1",
		Name:    "Ramen House",
		Reviews: nil,
	}

	dto := mapper.MapToPlaceDTO(place)

	if dto.AverageRating != 0 {
		t.Errorf("AverageRating: want 0, got %d", dto.AverageRating)
	}
	if dto.ID != "p1" {
		t.Errorf("ID: want p1, got %q", dto.ID)
	}
}

func TestMapToPlaceDTO_EmptyReviews_AverageRatingZero(t *testing.T) {
	place := &models.Place{Reviews: []*models.Review{}}

	dto := mapper.MapToPlaceDTO(place)

	if dto.AverageRating != 0 {
		t.Errorf("AverageRating: want 0 for empty reviews, got %d", dto.AverageRating)
	}
}

func TestMapToPlaceDTO_SingleReview_AverageRatingCorrect(t *testing.T) {
	place := &models.Place{
		Reviews: []*models.Review{{Rating: 4}},
	}

	dto := mapper.MapToPlaceDTO(place)

	if dto.AverageRating != 4 {
		t.Errorf("AverageRating: want 4, got %d", dto.AverageRating)
	}
}

func TestMapToPlaceDTO_MultipleReviews_AverageRatingRoundsDown(t *testing.T) {
	// 3 + 4 + 4 = 11 / 3 = 3 (integer division, truncates)
	place := &models.Place{
		Reviews: []*models.Review{
			{Rating: 3},
			{Rating: 4},
			{Rating: 4},
		},
	}

	dto := mapper.MapToPlaceDTO(place)

	if dto.AverageRating != 3 {
		t.Errorf("AverageRating: want 3 (integer division of 11/3), got %d", dto.AverageRating)
	}
}

func TestMapToPlaceDTO_FieldsAreMapped(t *testing.T) {
	place := &models.Place{
		ID:              "abc",
		Name:            "Cafe Latte",
		Description:     "cozy",
		Address:         "1 Main St",
		Lat:             "1.23",
		Lon:             "4.56",
		AverageSpending: 20,
		ImageUrl:        "http://img.example.com/img.jpg",
		SubCategory:     "Coffee",
		Category:        "Food",
		Area:            "Central",
		Liked:           true,
		CreatedBy:       "admin",
	}

	dto := mapper.MapToPlaceDTO(place)

	checks := map[string]bool{
		"ID":          dto.ID == "abc",
		"Name":        dto.Name == "Cafe Latte",
		"Description": dto.Description == "cozy",
		"Address":     dto.Address == "1 Main St",
		"Liked":       dto.Liked == true,
		"Category":    dto.Category == "Food",
		"CreatedBy":   dto.CreatedBy == "admin",
	}
	for field, ok := range checks {
		if !ok {
			t.Errorf("field %s not mapped correctly", field)
		}
	}
}
