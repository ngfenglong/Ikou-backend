package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ngfenglong/ikou-backend/api/store"
	"github.com/ngfenglong/ikou-backend/internal/helper"

	"github.com/go-chi/chi/v5"
)

type CodestableController struct {
	store *store.Store
}

func NewCodestableController(store *store.Store) *CodestableController {
	return &CodestableController{store: store}
}

func (cc *CodestableController) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := cc.store.DB.GetAllCategory()
	if err != nil {
		log.Fatalf("Failed to execute queries: %v", err)
		return
	}

	err = helper.WriteJSONResponse(w, http.StatusOK, categories)
	if err != nil {
		log.Fatalf("Failed to convert to json: %v", err)
		return
	}
}

func (cc *CodestableController) GetAllSubCategories(w http.ResponseWriter, r *http.Request) {
	subCategories, err := cc.store.DB.GetAllSubCategory()
	if err != nil {
		log.Fatalf("Failed to execute queries: %v", err)
		return
	}

	err = helper.WriteJSONResponse(w, http.StatusOK, subCategories)
	if err != nil {
		log.Fatalf("Failed to convert to json: %v", err)
		return
	}
}

func (cc *CodestableController) GetSubCategoriesByCategory(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	categoryCode, err := strconv.Atoi(code)
	if err != nil {
		log.Fatalf("Failed to convert to parameter to int: %v", err)
		return
	}

	subCategories, err := cc.store.DB.GetAllSubCategoryByCategoryCode(categoryCode)
	if err != nil {
		log.Fatalf("Failed to execute queries: %v", err)
		return
	}

	err = helper.WriteJSONResponse(w, http.StatusOK, subCategories)
	if err != nil {
		log.Fatalf("Failed to convert to json: %v", err)
		return
	}
}

func (cc *CodestableController) GetAllAreas(w http.ResponseWriter, r *http.Request) {
	areas, err := cc.store.DB.GetAllAreas()
	if err != nil {
		log.Fatalf("Failed to execute queries: %v", err)
		return
	}

	err = helper.WriteJSONResponse(w, http.StatusOK, areas)
	if err != nil {
		log.Fatalf("Failed to convert to json: %v", err)
		return
	}
}
