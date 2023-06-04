package routes

import (
	"github.com/ngfenglong/ikou-backend/api/controllers"
	"github.com/ngfenglong/ikou-backend/api/store"

	"github.com/go-chi/chi/v5"
)

func PlaceRoutes(store *store.Store) chi.Router {
	mux := chi.NewRouter()

	placeController := controllers.NewPlaceController(store)
	mux.Get("/", placeController.GetAllPlaces)
	mux.Get("/getPlaceById/{id}", placeController.GetPlaceById)
	mux.Get("/getPlacesBySubCategory/{code}", placeController.GetPlacesBySubCategoryCode)
	mux.Post("/searchPlaceByKeyword", placeController.SearchPlacesByKeyword)

	return mux
}
