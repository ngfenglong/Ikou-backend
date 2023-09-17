package mapper

import (
	"github.com/ngfenglong/ikou-backend/api/dto"
	"github.com/ngfenglong/ikou-backend/api/models"
)

func MapToPlaceDTO(place *models.Place) *dto.PlaceDTO {
	avg_rating := 0
	if place.Reviews != nil && len(place.Reviews) > 0 {
		totalRating := 0
		for _, r := range place.Reviews {
			totalRating += r.Rating
		}
		avg_rating = totalRating / len(place.Reviews)
	}

	return &dto.PlaceDTO{
		ID:              place.ID,
		Name:            place.Name,
		Description:     place.Description,
		Address:         place.Address,
		Lat:             place.Lat,
		Lon:             place.Lon,
		AverageSpending: place.AverageSpending,
		AverageRating:   avg_rating,
		ImageUrl:        place.ImageUrl,
		SubCategory:     place.SubCategory,
		Category:        place.Category,
		Liked:           place.Liked,
		Area:            place.Area,
		Reviews:         MapToReviewsDTO(place.Reviews),
		CreatedAt:       place.CreatedAt,
		UpdatedAt:       place.UpdatedAt,
		CreatedBy:       place.CreatedBy,
	}
}
