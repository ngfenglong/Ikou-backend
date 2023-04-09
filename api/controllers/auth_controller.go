package controllers

import (
	"net/http"

	"github.com/ngfenglong/ikou-backend/api/dto"
	"github.com/ngfenglong/ikou-backend/api/store"
	"github.com/ngfenglong/ikou-backend/internal/helper"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	store *store.Store
}

func NewAuthController(store *store.Store) *AuthController {
	return &AuthController{store: store}
}

func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var loginCredentialInput dto.LoginCredentialInputDto

	err := helper.ReadJSON(w, r, &loginCredentialInput)
	if err != nil {
		helper.BadRequest(w, r, err)
		return
	}

	user, err := ac.store.DB.GetUserByUsername(loginCredentialInput.Username)
	if err != nil {
		helper.InvalidCredential(w)
		return
	}

	validPassword, err := helper.PasswordMatches(user.Password, loginCredentialInput.Password)
	if !validPassword || err != nil {
		helper.InvalidCredential(w)
		return
	}

	tokenDetail := &helper.TokenDetail{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}

	accessToken, accessExpiry, err := helper.GenerateAccessToken(tokenDetail)
	if err != nil {
		helper.BadRequest(w, r, err)
		return
	}

	refreshToken, refreshExpiry, err := helper.GenerateRefreshToken(tokenDetail)
	if err != nil {
		helper.BadRequest(w, r, err)
		return
	}

	err = ac.store.DB.InsertToken(user.ID, refreshToken, refreshExpiry)
	if err != nil {
		helper.BadRequest(w, r, err)
		return
	}

	var payload dto.LoginResponseDto
	// var userDto
	payload.Error = false
	payload.AccessToken = accessToken
	payload.RefreshToken = refreshToken
	payload.Expiry = accessExpiry
	payload.User = dto.UserDto{
		UserName:     user.Username,
		FirstName:    user.FirstName,
		Email:        user.Email,
		LastName:     user.LastName,
		Country:      user.Country,
		ProfileImage: user.ProfileImage,
	}

	err = helper.WriteJSONResponse(w, http.StatusOK, payload)
	if err != nil {
		helper.BadRequest(w, r, err)
	}
}

func (ac *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {

}

func (ac *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var rfi dto.RegisterFormInputDto

	err := helper.ReadJSON(w, r, &rfi)
	if err != nil {
		helper.BadRequest(w, r, err)
		return
	}

	// Todo: Add Validation Handling
	usernameExists, emailExists, err := ac.store.DB.CheckIfUserExists(rfi)
	if err != nil {
		helper.BadRequest(w, r, err)
	}

	if usernameExists {
		helper.ConflictErrorResponse("Username already exists", w)
		return
	}

	if emailExists {
		helper.ConflictErrorResponse("Email already exists", w)
		return
	}

	if rfi.ProfileImage == "" {
		rfi.ProfileImage = "/images/no_profile.jpeg"
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rfi.Password), 10)
	if err != nil {
		helper.BadRequest(w, r, err)
		return
	}

	rfi.Password = string(hashedPassword)

	err = ac.store.DB.RegisterUser(rfi)
	if err != nil {
		helper.BadRequest(w, r, err)
		return
	}
}
