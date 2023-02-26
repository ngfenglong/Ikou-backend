package main

import (
	"encoding/json"
	"ikou/internal/helper"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *application) GetAllPlaces(w http.ResponseWriter, r *http.Request) {
	places, err := app.DB.GetAllPlaces()
	if err != nil {
		// ToDo: do some proper error handling here
		app.errorLog.Println(err)
		return
	}

	err = helper.WriteJSON(w, http.StatusOK, places)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
}

// Get place with details such as comments, liked, etc...
func (app *application) GetPlaceById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	place, err := app.DB.GetPlaceById(id)

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	out, err := json.MarshalIndent(place, "", " ")

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (app *application) GetPlacesBySubCategoryCode(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	subCategoryCode, err := strconv.Atoi(code)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	places, err := app.DB.GetPlacesBySubCategoryCode(subCategoryCode)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	err = helper.WriteJSON(w, http.StatusOK, places)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
}

func (app *application) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := app.DB.GetAllCategory()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	err = helper.WriteJSON(w, http.StatusOK, categories)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
}

func (app *application) GetAllSubCategories(w http.ResponseWriter, r *http.Request) {
	subCategories, err := app.DB.GetAllSubCategory()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	err = helper.WriteJSON(w, http.StatusOK, subCategories)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
}

func (app *application) GetSubCategoriesByCategory(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	categoryCode, err := strconv.Atoi(code)
	subCategories, err := app.DB.GetAllSubCategoryByCategoryCode(categoryCode)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	err = helper.WriteJSON(w, http.StatusOK, subCategories)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
}
