package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/ngfenglong/ikou-backend/api/controllers"
	"github.com/ngfenglong/ikou-backend/api/store"
)

func AuthRoutes(store *store.Store) chi.Router {
	mux := chi.NewRouter()

	authController := controllers.NewAuthController(store)

	mux.Post("/login", authController.Login)
	mux.Post("/Register", authController.Register)
	mux.Post("/refresh", authController.Refresh)

	return mux
}
