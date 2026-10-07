package testutil

import (
	"time"

	"github.com/ngfenglong/ikou-backend/api/dto"
	"github.com/ngfenglong/ikou-backend/api/models"
)

// MockDB is a test double for repository.Repository. Set each Func field to
// control the return value for that method in a given test.
type MockDB struct {
	GetUserByIDFunc                     func(userID string) (*models.User, error)
	GetUserByUsernameFunc               func(username string) (*models.User, error)
	InsertTokenFunc                     func(userId string, refreshToken string, expiresAt time.Time) error
	RegisterUserFunc                    func(r dto.RegisterFormInputDTO) error
	CheckIfUserExistsFunc               func(r dto.RegisterFormInputDTO) (bool, bool, error)
	DeleteTokenFunc                     func(refreshToken string) error
	FetchRefreshTokenFromDBFunc         func(tokenStr string, userID string) (*models.RefreshToken, bool)
	GetAllPlacesFunc                    func(userID string) ([]*dto.PlaceDTO, error)
	GetPlaceByIdFunc                    func(id string, userID string) (*dto.PlaceDTO, error)
	GetPlacesByCategoryCodeFunc         func(category string, userID string) ([]*dto.PlaceDTO, error)
	GetPlacesBySubCategoryCodeFunc      func(code int, userID string) ([]*dto.PlaceDTO, error)
	SearchPlaceByKeywordFunc            func(keyword string, userID string) ([]*dto.PlaceDTO, error)
	AddPlaceRequestFunc                 func(pr dto.PlaceRequestDto) error
	HasUserLikedPlaceFunc               func(userID string, placeID string) (bool, error)
	RemoveUserLikeFromPlaceFunc         func(userID string, placeID string) error
	AddUserLikeToPlaceFunc              func(userID string, placeID string) error
	GetAllCategoryFunc                  func() ([]*models.CodeDecodeCategory, error)
	GetAllSubCategoryFunc               func() ([]*models.CodeDecodeSubCategory, error)
	GetAllSubCategoryByCategoryCodeFunc func(categoryCode int) ([]*models.CodeDecodeSubCategory, error)
	GetAllAreasFunc                     func() ([]*models.CodeDecodeArea, error)
	GetActivityByPlaceFunc              func(placeId string) ([]*models.Activity, error)
}

func (m *MockDB) GetUserByID(userID string) (*models.User, error) {
	return m.GetUserByIDFunc(userID)
}
func (m *MockDB) GetUserByUsername(username string) (*models.User, error) {
	return m.GetUserByUsernameFunc(username)
}
func (m *MockDB) InsertToken(userId string, refreshToken string, expiresAt time.Time) error {
	return m.InsertTokenFunc(userId, refreshToken, expiresAt)
}
func (m *MockDB) RegisterUser(r dto.RegisterFormInputDTO) error {
	return m.RegisterUserFunc(r)
}
func (m *MockDB) CheckIfUserExists(r dto.RegisterFormInputDTO) (bool, bool, error) {
	return m.CheckIfUserExistsFunc(r)
}
func (m *MockDB) DeleteToken(refreshToken string) error {
	return m.DeleteTokenFunc(refreshToken)
}
func (m *MockDB) FetchRefreshTokenFromDB(tokenStr string, userID string) (*models.RefreshToken, bool) {
	return m.FetchRefreshTokenFromDBFunc(tokenStr, userID)
}
func (m *MockDB) GetAllPlaces(userID string) ([]*dto.PlaceDTO, error) {
	return m.GetAllPlacesFunc(userID)
}
func (m *MockDB) GetPlaceById(id string, userID string) (*dto.PlaceDTO, error) {
	return m.GetPlaceByIdFunc(id, userID)
}
func (m *MockDB) GetPlacesByCategoryCode(category string, userID string) ([]*dto.PlaceDTO, error) {
	return m.GetPlacesByCategoryCodeFunc(category, userID)
}
func (m *MockDB) GetPlacesBySubCategoryCode(code int, userID string) ([]*dto.PlaceDTO, error) {
	return m.GetPlacesBySubCategoryCodeFunc(code, userID)
}
func (m *MockDB) SearchPlaceByKeyword(keyword string, userID string) ([]*dto.PlaceDTO, error) {
	return m.SearchPlaceByKeywordFunc(keyword, userID)
}
func (m *MockDB) AddPlaceRequest(pr dto.PlaceRequestDto) error {
	return m.AddPlaceRequestFunc(pr)
}
func (m *MockDB) HasUserLikedPlace(userID string, placeID string) (bool, error) {
	return m.HasUserLikedPlaceFunc(userID, placeID)
}
func (m *MockDB) RemoveUserLikeFromPlace(userID string, placeID string) error {
	return m.RemoveUserLikeFromPlaceFunc(userID, placeID)
}
func (m *MockDB) AddUserLikeToPlace(userID string, placeID string) error {
	return m.AddUserLikeToPlaceFunc(userID, placeID)
}
func (m *MockDB) GetAllCategory() ([]*models.CodeDecodeCategory, error) {
	return m.GetAllCategoryFunc()
}
func (m *MockDB) GetAllSubCategory() ([]*models.CodeDecodeSubCategory, error) {
	return m.GetAllSubCategoryFunc()
}
func (m *MockDB) GetAllSubCategoryByCategoryCode(categoryCode int) ([]*models.CodeDecodeSubCategory, error) {
	return m.GetAllSubCategoryByCategoryCodeFunc(categoryCode)
}
func (m *MockDB) GetAllAreas() ([]*models.CodeDecodeArea, error) {
	return m.GetAllAreasFunc()
}
func (m *MockDB) GetActivityByPlace(placeId string) ([]*models.Activity, error) {
	return m.GetActivityByPlaceFunc(placeId)
}
