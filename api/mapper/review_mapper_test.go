package mapper_test

import (
	"testing"

	"github.com/ngfenglong/ikou-backend/api/mapper"
	"github.com/ngfenglong/ikou-backend/api/models"
)

func TestMapToReviewsDTO_NilSlice_ReturnsNil(t *testing.T) {
	got := mapper.MapToReviewsDTO(nil)

	if got != nil {
		t.Errorf("expected nil for nil input, got %v", got)
	}
}

func TestMapToReviewsDTO_SingleReview_MapsAllFields(t *testing.T) {
	reviews := []*models.Review{
		{
			ID:                   "r1",
			Rating:               5,
			ReviewDescription:    "excellent",
			ReviewerProfileImage: "http://example.com/img.png",
			CreatedBy:            "alice",
		},
	}

	dtos := mapper.MapToReviewsDTO(reviews)

	if len(dtos) != 1 {
		t.Fatalf("expected 1 dto, got %d", len(dtos))
	}
	got := dtos[0]
	if got.ID != "r1" {
		t.Errorf("ID: want r1, got %q", got.ID)
	}
	if got.Rating != 5 {
		t.Errorf("Rating: want 5, got %d", got.Rating)
	}
	if got.ReviewDescription != "excellent" {
		t.Errorf("ReviewDescription: want excellent, got %q", got.ReviewDescription)
	}
	if got.CreatedBy != "alice" {
		t.Errorf("CreatedBy: want alice, got %q", got.CreatedBy)
	}
}

func TestMapToReviewsDTO_MultipleReviews_AllMapped(t *testing.T) {
	reviews := []*models.Review{
		{ID: "r1", Rating: 3},
		{ID: "r2", Rating: 5},
	}

	dtos := mapper.MapToReviewsDTO(reviews)

	if len(dtos) != 2 {
		t.Fatalf("expected 2 dtos, got %d", len(dtos))
	}
	if dtos[0].ID != "r1" || dtos[1].ID != "r2" {
		t.Errorf("order or IDs incorrect: got %q %q", dtos[0].ID, dtos[1].ID)
	}
}
