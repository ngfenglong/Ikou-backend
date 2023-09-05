package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

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
	tokenStr := r.Header.Get("Authorization")

	var tokenClaims *helper.TokenDetail
	if tokenStr != "" {
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		claimsValid, tempClaims := helper.VerifyAccessToken(tokenStr)
		if claimsValid && tempClaims != nil {
			tokenClaims = tempClaims
		}
	}

	userID := ""
	if tokenClaims != nil {
		userID = tokenClaims.ID
	}

	places, err := pc.store.DB.GetAllPlaces(userID)
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
		helper.BadRequest(w, r, err)
		return
	}

	out, err := json.MarshalIndent(place, "", " ")

	if err != nil {
		helper.BadRequest(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (pc *PlaceController) GetPlacesBySubCategoryCode(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	subCategoryCode, err := strconv.Atoi(code)
	if err != nil {
		helper.BadRequest(w, r, err)
		return
	}

	places, err := pc.store.DB.GetPlacesBySubCategoryCode(subCategoryCode)
	if err != nil {
		helper.BadRequest(w, r, err)
		return
	}

	err = helper.WriteJSONResponse(w, http.StatusOK, places)
	if err != nil {
		helper.BadRequest(w, r, err)
		return
	}
}

func (pc *PlaceController) GetPlacesByCategory(w http.ResponseWriter, r *http.Request) {
	category := chi.URLParam(r, "category")

	if category == "" {
		helper.BadRequest(w, r, errors.New("category is invalid"))
		return
	}

	places, err := pc.store.DB.GetPlacesByCategoryCode(category)
	if err != nil {
		helper.BadRequest(w, r, err)
		return
	}

	err = helper.WriteJSONResponse(w, http.StatusOK, places)
	if err != nil {
		helper.BadRequest(w, r, err)
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
