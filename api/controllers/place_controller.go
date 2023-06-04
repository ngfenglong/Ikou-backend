package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/ngfenglong/ikou-backend/api/store"
	"github.com/ngfenglong/ikou-backend/internal/helper"

	"github.com/go-chi/chi/v5"
)

type PlaceController struct {
	store *store.Store
}

func NewPlaceController(store *store.Store) *PlaceController {
	return &PlaceController{store: store}
}

func (pc *PlaceController) GetAllPlaces(w http.ResponseWriter, r *http.Request) {
	places, err := pc.store.DB.GetAllPlaces()
	if err != nil {
		// ToDo: do some proper error handling here
		log.Fatalf("Failed to execute queries: %v", err)
		return
	}

	err = helper.WriteJSONResponse(w, http.StatusOK, places)
	if err != nil {
		log.Fatalf("Failed to convert to json: %v", err)
		return
	}
}

// Get place with details such as comments, liked, etc...
func (pc *PlaceController) GetPlaceById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	place, err := pc.store.DB.GetPlaceById(id)

	if err != nil {
		log.Fatalf("Failed to execute queries: %v", err)
		return
	}

	out, err := json.MarshalIndent(place, "", " ")

	if err != nil {
		log.Fatalf("Failed to convert to json: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (pc *PlaceController) GetPlacesBySubCategoryCode(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	subCategoryCode, err := strconv.Atoi(code)
	if err != nil {
		log.Fatalf("Failed to convert parameter into int: %v", err)
		return
	}

	places, err := pc.store.DB.GetPlacesBySubCategoryCode(subCategoryCode)
	if err != nil {
		log.Fatalf("Failed to execute queries: %v", err)
		return
	}

	err = helper.WriteJSONResponse(w, http.StatusOK, places)
	if err != nil {
		log.Fatalf("Failed to convert to json: %v", err)
		return
	}
}

func (pc *PlaceController) SearchPlacesByKeyword(w http.ResponseWriter, r *http.Request) {
	var searchPlaceRequestDto struct {
		Keyword string `json:"keyword"`
	}

	err := helper.ReadJSON(w, r, &searchPlaceRequestDto)
	if err != nil {
		helper.BadRequest(w, r, err)
		return
	}

	if len(searchPlaceRequestDto.Keyword) < 3 {
		helper.BadRequest(w, r, errors.New("keyword length must be at least 3"))
		return
	}

	places, err := pc.store.DB.SearchPlaceByKeyword(searchPlaceRequestDto.Keyword)
	if err != nil {
		helper.BadRequest(w, r, err)
		return
	}

	out, err := json.MarshalIndent(places, "", "")
	if err != nil {
		helper.BadRequest(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
