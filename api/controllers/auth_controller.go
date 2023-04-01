package controllers

import "github.com/ngfenglong/ikou-backend/api/store"

type AuthController struct {
	store *store.Store
}

func NewAuthController(store *store.Store) *AuthController {
	return &AuthController{store: store}
}
