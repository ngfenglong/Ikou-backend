package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/ngfenglong/ikou-backend/api/dto"
	"github.com/ngfenglong/ikou-backend/api/middleware"
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
	userID := r.Context().Value(middleware.UserIDKey).(string)

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
	userID := r.Context().Value(middleware.UserIDKey).(string)
	place, err := pc.store.DB.GetPlaceById(id, userID)

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
	userID := r.Context().Value(middleware.UserIDKey).(string)
	code := chi.URLParam(r, "code")
	subCategoryCode, err := strconv.Atoi(code)
	if err != nil {
		helper.BadRequest(w, r, err)
		return
	}

	places, err := pc.store.DB.GetPlacesBySubCategoryCode(subCategoryCode, userID)
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
	userID := r.Context().Value(middleware.UserIDKey).(string)
	category := chi.URLParam(r, "category")

	if category == "" {
		helper.BadRequest(w, r, errors.New("category is invalid"))
		return
	}

	places, err := pc.store.DB.GetPlacesByCategoryCode(category, userID)
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
	userID := r.Context().Value(middleware.UserIDKey).(string)
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

	places, err := pc.store.DB.SearchPlaceByKeyword(searchPlaceRequestDto.Keyword, userID)
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

func (pc *PlaceController) AddPlaceRequest(w http.ResponseWriter, r *http.Request) {
	userName := r.Context().Value(middleware.UserNameKey).(string)
	if userName == "" {
		helper.InvalidCredential(w)
	}

	var npr dto.PlaceRequestDto
	err := helper.ReadJSON(w, r, &npr)
	if err != nil {
		helper.BadRequest(w, r, err)
		return
	}

	npr.CreatedBy = userName

	err = pc.store.DB.AddPlaceRequest(npr)
	if err != nil {
		helper.BadRequest(w, r, err)
		return
	}

	var payload dto.SuccessResponseDto
	payload.Error = false
	payload.Message = "Place Request submitted successfully"

	err = helper.WriteJSONResponse(w, http.StatusCreated, payload)
	if err != nil {
		helper.BadRequest(w, r, err)
	}
}

func (pc *PlaceController) ToggleLike(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middleware.UserIDKey).(string)
	if userId == "" {
		helper.InvalidCredential(w)
	}

	placeId := chi.URLParam(r, "placeId")

	// Check if liked_places have records
	liked, err := pc.store.DB.HasUserLikedPlace(userId, placeId)
	if err != nil {
		helper.BadRequest(w, r, err)
		return
	}

	payloadMessage := ""
	if liked {
		err = pc.store.DB.RemoveUserLikeFromPlace(userId, placeId)
		payloadMessage = "Removed like Successfully"
	} else {
		err = pc.store.DB.AddUserLikeToPlace(userId, placeId)
		payloadMessage = "Liked successfully"
	}

	if err != nil {
		helper.BadRequest(w, r, err)
		return
	}

	var payload dto.SuccessResponseDto
	payload.Error = false
	payload.Message = payloadMessage

	err = helper.WriteJSONResponse(w, http.StatusOK, payload)
	if err != nil {
		helper.BadRequest(w, r, err)
	}

}
