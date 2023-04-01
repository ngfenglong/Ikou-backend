package controllers

import (
	"ikou/api/store"
	"ikou/internal/helper"
	"log"
	"net/http"
	"strconv"

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

	err = helper.WriteJSON(w, http.StatusOK, categories)
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

	err = helper.WriteJSON(w, http.StatusOK, subCategories)
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

	err = helper.WriteJSON(w, http.StatusOK, subCategories)
	if err != nil {
		log.Fatalf("Failed to convert to json: %v", err)
		return
	}
}
