package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ngfenglong/ikou-backend/api/dto"
	"github.com/ngfenglong/ikou-backend/api/store"
	"github.com/ngfenglong/ikou-backend/internal/helper"
)

type AuthController struct {
	store *store.Store
}

func NewAuthController(store *store.Store) *AuthController {
	return &AuthController{store: store}
}

func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var loginCredentialInput dto.LoginCredentialInput

	err := helper.ReadJSON(w, r, &loginCredentialInput)
	if err != nil {
		helper.BadRequest(w, r, err)
		return
	}

	user, err := ac.store.DB.GetUserByUsername(loginCredentialInput.Username)
	if err != nil {
		helper.InvalidCredential(w)
		fmt.Print("Username")
		return
	}

	validPassword, err := helper.PasswordMatches(user.Password, loginCredentialInput.Password)
	if !validPassword || err != nil {
		helper.InvalidCredential(w)
		fmt.Print("Password mismatch")
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

	var payload struct {
		Error        bool      `json:"error"`
		Message      string    `json:"message"`
		AccessToken  string    `json:"access_token"`
		RefreshToken string    `json:"refresh_token"`
		Expiry       time.Time `json:"expiry"`
	}

	payload.Error = false
	payload.AccessToken = accessToken
	payload.RefreshToken = refreshToken
	payload.Expiry = accessExpiry

	err = helper.WriteJSONResponse(w, http.StatusOK, payload)
	if err != nil {
		helper.BadRequest(w, r, err)
	}
}

func (ac *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {

}
