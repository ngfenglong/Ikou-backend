package controllers

import "ikou/api/store"

type AuthController struct {
	store *store.Store
}

func NewAuthController(store *store.Store) *AuthController {
	return &AuthController{store: store}
}
