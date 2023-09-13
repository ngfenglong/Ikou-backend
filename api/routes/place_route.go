package routes

import (
	"github.com/ngfenglong/ikou-backend/api/controllers"
	"github.com/ngfenglong/ikou-backend/api/middleware"
	"github.com/ngfenglong/ikou-backend/api/store"

	"github.com/go-chi/chi/v5"
)

func PlaceRoutes(store *store.Store) chi.Router {
	mux := chi.NewRouter()

	placeController := controllers.NewPlaceController(store)
	mux.Get("/", middleware.ExtractTokenMiddleware(placeController.GetAllPlaces))
	mux.Get("/getPlaceById/{id}", middleware.ExtractTokenMiddleware(placeController.GetPlaceById))
	mux.Get("/getPlacesBySubCategory/{code}", middleware.ExtractTokenMiddleware(placeController.GetPlacesBySubCategoryCode))
	mux.Get("/getPlacesByCategory/{category}", middleware.ExtractTokenMiddleware(placeController.GetPlacesByCategory))
	mux.Post("/searchPlaceByKeyword", middleware.ExtractTokenMiddleware(placeController.SearchPlacesByKeyword))
	mux.Post("/addPlaceRequest", middleware.ExtractTokenMiddleware(placeController.AddPlaceRequest))

	return mux
}
