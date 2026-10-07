package controllers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/ngfenglong/ikou-backend/api/controllers"
	"github.com/ngfenglong/ikou-backend/api/dto"
	"github.com/ngfenglong/ikou-backend/api/middleware"
	"github.com/ngfenglong/ikou-backend/api/store"
	"github.com/ngfenglong/ikou-backend/api/testutil"
)

func newPlaceStore(mock *testutil.MockDB) *store.Store {
	return &store.Store{DB: mock}
}

func withUserContext(r *http.Request, userID, userName string) *http.Request {
	ctx := context.WithValue(r.Context(), middleware.UserIDKey, userID)
	ctx = context.WithValue(ctx, middleware.UserNameKey, userName)
	return r.WithContext(ctx)
}

func withChiParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

// --- GetAllPlaces ---

func TestGetAllPlaces_Success_Returns200(t *testing.T) {
	mock := &testutil.MockDB{
		GetAllPlacesFunc: func(userID string) ([]*dto.PlaceDTO, error) {
			return []*dto.PlaceDTO{{ID: "p1", Name: "Ramen House"}}, nil
		},
	}
	pc := controllers.NewPlaceController(newPlaceStore(mock))

	r := httptest.NewRequest(http.MethodGet, "/places", nil)
	r = withUserContext(r, "user-1", "alice")
	w := httptest.NewRecorder()
	pc.GetAllPlaces(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status: want 200, got %d", w.Code)
	}
	var places []*dto.PlaceDTO
	json.Unmarshal(w.Body.Bytes(), &places)
	if len(places) != 1 || places[0].ID != "p1" {
		t.Errorf("unexpected places response: %v", places)
	}
}

// --- GetPlaceById ---

func TestGetPlaceById_Found_Returns200(t *testing.T) {
	mock := &testutil.MockDB{
		GetPlaceByIdFunc: func(id, userID string) (*dto.PlaceDTO, error) {
			return &dto.PlaceDTO{ID: id, Name: "Sushi Bar"}, nil
		},
	}
	pc := controllers.NewPlaceController(newPlaceStore(mock))

	r := httptest.NewRequest(http.MethodGet, "/places/p1", nil)
	r = withUserContext(r, "user-1", "alice")
	r = withChiParam(r, "id", "p1")
	w := httptest.NewRecorder()
	pc.GetPlaceById(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status: want 200, got %d", w.Code)
	}
}

func TestGetPlaceById_NotFound_Returns400(t *testing.T) {
	mock := &testutil.MockDB{
		GetPlaceByIdFunc: func(id, userID string) (*dto.PlaceDTO, error) {
			return nil, errors.New("record not found")
		},
	}
	pc := controllers.NewPlaceController(newPlaceStore(mock))

	r := httptest.NewRequest(http.MethodGet, "/places/unknown", nil)
	r = withUserContext(r, "user-1", "alice")
	r = withChiParam(r, "id", "unknown")
	w := httptest.NewRecorder()
	pc.GetPlaceById(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: want 400, got %d", w.Code)
	}
}

// --- GetPlacesByCategory ---

func TestGetPlacesByCategory_EmptyCategory_Returns400(t *testing.T) {
	mock := &testutil.MockDB{}
	pc := controllers.NewPlaceController(newPlaceStore(mock))

	r := httptest.NewRequest(http.MethodGet, "/places/category/", nil)
	r = withUserContext(r, "user-1", "alice")
	// chi param is empty string
	r = withChiParam(r, "category", "")
	w := httptest.NewRecorder()
	pc.GetPlacesByCategory(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: want 400, got %d", w.Code)
	}
}

func TestGetPlacesByCategory_ValidCategory_Returns200(t *testing.T) {
	mock := &testutil.MockDB{
		GetPlacesByCategoryCodeFunc: func(category, userID string) ([]*dto.PlaceDTO, error) {
			return []*dto.PlaceDTO{}, nil
		},
	}
	pc := controllers.NewPlaceController(newPlaceStore(mock))

	r := httptest.NewRequest(http.MethodGet, "/places/category/Food", nil)
	r = withUserContext(r, "user-1", "alice")
	r = withChiParam(r, "category", "Food")
	w := httptest.NewRecorder()
	pc.GetPlacesByCategory(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status: want 200, got %d", w.Code)
	}
}

// --- SearchPlacesByKeyword ---

func TestSearchPlacesByKeyword_ShortKeyword_Returns400(t *testing.T) {
	mock := &testutil.MockDB{}
	pc := controllers.NewPlaceController(newPlaceStore(mock))

	body := `{"keyword":"ab"}`
	r := httptest.NewRequest(http.MethodPost, "/places/search", strings.NewReader(body))
	r = withUserContext(r, "user-1", "alice")
	w := httptest.NewRecorder()
	pc.SearchPlacesByKeyword(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: want 400, got %d", w.Code)
	}
}

func TestSearchPlacesByKeyword_ValidKeyword_Returns200(t *testing.T) {
	mock := &testutil.MockDB{
		SearchPlaceByKeywordFunc: func(keyword, userID string) ([]*dto.PlaceDTO, error) {
			return []*dto.PlaceDTO{}, nil
		},
	}
	pc := controllers.NewPlaceController(newPlaceStore(mock))

	body := `{"keyword":"ramen"}`
	r := httptest.NewRequest(http.MethodPost, "/places/search", strings.NewReader(body))
	r = withUserContext(r, "user-1", "alice")
	w := httptest.NewRecorder()
	pc.SearchPlacesByKeyword(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status: want 200, got %d", w.Code)
	}
}

// --- ToggleLike ---

func TestToggleLike_NotLiked_LikesPlace(t *testing.T) {
	addCalled := false
	mock := &testutil.MockDB{
		HasUserLikedPlaceFunc: func(userID, placeID string) (bool, error) {
			return false, nil
		},
		AddUserLikeToPlaceFunc: func(userID, placeID string) error {
			addCalled = true
			return nil
		},
	}
	pc := controllers.NewPlaceController(newPlaceStore(mock))

	r := httptest.NewRequest(http.MethodPost, "/places/p1/like", nil)
	r = withUserContext(r, "user-1", "alice")
	r = withChiParam(r, "placeId", "p1")
	w := httptest.NewRecorder()
	pc.ToggleLike(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status: want 200, got %d", w.Code)
	}
	if !addCalled {
		t.Error("expected AddUserLikeToPlace to be called")
	}
}

func TestToggleLike_AlreadyLiked_UnlikesPlace(t *testing.T) {
	removeCalled := false
	mock := &testutil.MockDB{
		HasUserLikedPlaceFunc: func(userID, placeID string) (bool, error) {
			return true, nil
		},
		RemoveUserLikeFromPlaceFunc: func(userID, placeID string) error {
			removeCalled = true
			return nil
		},
	}
	pc := controllers.NewPlaceController(newPlaceStore(mock))

	r := httptest.NewRequest(http.MethodPost, "/places/p1/like", nil)
	r = withUserContext(r, "user-1", "alice")
	r = withChiParam(r, "placeId", "p1")
	w := httptest.NewRecorder()
	pc.ToggleLike(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status: want 200, got %d", w.Code)
	}
	if !removeCalled {
		t.Error("expected RemoveUserLikeFromPlace to be called")
	}
}
