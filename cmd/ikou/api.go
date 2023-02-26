package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-TOKEN"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	mux.Get("/api/places", app.GetAllPlaces)
	mux.Get("/api/getPlaceById/{id}", app.GetPlaceById)
	mux.Get("/api/getPlacesBySubCategory/{code}", app.GetPlacesBySubCategoryCode)
	mux.Get("/api/codeDecodeCategories", app.GetAllCategories)
	mux.Get("/api/codeDecodeSubCategories", app.GetAllSubCategories)
	mux.Get("/api/codeDecodeSubCategories/{code}", app.GetSubCategoriesByCategory)

	return mux
}
