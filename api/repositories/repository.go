package repository

import (
	"time"

	"github.com/ngfenglong/ikou-backend/api/dto"
	"github.com/ngfenglong/ikou-backend/api/models"
)

type Repository interface {
	// auth
	GetUserByID(userID string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	InsertToken(userId string, refreshToken string, expiresAt time.Time) error
	RegisterUser(r dto.RegisterFormInputDTO) error
	CheckIfUserExists(r dto.RegisterFormInputDTO) (bool, bool, error)
	DeleteToken(refreshToken string) error
	FetchRefreshTokenFromDB(tokenStr string, userID string) (*models.RefreshToken, bool)

	// place
	GetAllPlaces(userID string) ([]*dto.PlaceDTO, error)
	GetPlaceById(id string, userID string) (*dto.PlaceDTO, error)
	GetPlacesByCategoryCode(category string, userID string) ([]*dto.PlaceDTO, error)
	GetPlacesBySubCategoryCode(code int, userID string) ([]*dto.PlaceDTO, error)
	SearchPlaceByKeyword(keyword string, userID string) ([]*dto.PlaceDTO, error)
	AddPlaceRequest(pr dto.PlaceRequestDto) error
	HasUserLikedPlace(userID string, placeID string) (bool, error)
	RemoveUserLikeFromPlace(userID string, placeID string) error
	AddUserLikeToPlace(userID string, placeID string) error

	// codestable
	GetAllCategory() ([]*models.CodeDecodeCategory, error)
	GetAllSubCategory() ([]*models.CodeDecodeSubCategory, error)
	GetAllSubCategoryByCategoryCode(categoryCode int) ([]*models.CodeDecodeSubCategory, error)
	GetAllAreas() ([]*models.CodeDecodeArea, error)

	// activity
	GetActivityByPlace(placeId string) ([]*models.Activity, error)
}
