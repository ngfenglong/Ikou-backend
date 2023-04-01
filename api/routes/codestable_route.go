package routes

import (
	"ikou/api/controllers"
	"ikou/api/store"

	"github.com/go-chi/chi/v5"
)

func CodestableRoutes(store *store.Store) chi.Router {
	mux := chi.NewRouter()

	codestableController := controllers.NewCodestableController(store)

	mux.Get("/codeDecodeCategories", codestableController.GetAllCategories)
	mux.Get("/codeDecodeSubCategories", codestableController.GetAllSubCategories)
	mux.Get("/codeDecodeSubCategories/{code}", codestableController.GetSubCategoriesByCategory)

	return mux
}
