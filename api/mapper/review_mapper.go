package mapper

import (
	"github.com/ngfenglong/ikou-backend/api/dto"
	"github.com/ngfenglong/ikou-backend/api/models"
)

func MapToReviewsDTO(reviews []*models.Review) []*dto.ReviewDTO {
	var reviewsDTO []*dto.ReviewDTO

	for _, review := range reviews {
		reviewDTO := dto.ReviewDTO{
			ID:                   review.ID,
			Rating:               review.Rating,
			ReviewDescription:    review.ReviewDescription,
			ReviewerProfileImage: review.ReviewerProfileImage,
			CreatedAt:            review.CreatedAt,
			UpdatedAt:            review.UpdatedAt,
			CreatedBy:            review.CreatedBy,
		}
		reviewsDTO = append(reviewsDTO, &reviewDTO)
	}

	return reviewsDTO
}
