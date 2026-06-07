package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ngfenglong/ikou-backend/api/controllers"
	"github.com/ngfenglong/ikou-backend/api/models"
	"github.com/ngfenglong/ikou-backend/api/store"
	"github.com/ngfenglong/ikou-backend/api/testutil"
)

func newCodestableStore(mock *testutil.MockDB) *store.Store {
	return &store.Store{DB: mock}
}

// NOTE: Error paths in codestable_controller.go call log.Fatalf which terminates
// the process, so only happy-path cases are testable here.

func TestGetAllCategories_Success_Returns200(t *testing.T) {
	mock := &testutil.MockDB{
		GetAllCategoryFunc: func() ([]*models.CodeDecodeCategory, error) {
			return []*models.CodeDecodeCategory{
				{ID: "c1", Code: 1, Decode: "Food"},
			}, nil
		},
	}
	cc := controllers.NewCodestableController(newCodestableStore(mock))

	r := httptest.NewRequest(http.MethodGet, "/categories", nil)
	w := httptest.NewRecorder()
	cc.GetAllCategories(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status: want 200, got %d", w.Code)
	}
	var cats []*models.CodeDecodeCategory
	json.Unmarshal(w.Body.Bytes(), &cats)
	if len(cats) != 1 || cats[0].Decode != "Food" {
		t.Errorf("unexpected categories: %v", cats)
	}
}

func TestGetAllSubCategories_Success_Returns200(t *testing.T) {
	mock := &testutil.MockDB{
		GetAllSubCategoryFunc: func() ([]*models.CodeDecodeSubCategory, error) {
			return []*models.CodeDecodeSubCategory{
				{ID: "s1", Code: 101, Decode: "Japanese"},
			}, nil
		},
	}
	cc := controllers.NewCodestableController(newCodestableStore(mock))

	r := httptest.NewRequest(http.MethodGet, "/subcategories", nil)
	w := httptest.NewRecorder()
	cc.GetAllSubCategories(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status: want 200, got %d", w.Code)
	}
}

func TestGetAllAreas_Success_Returns200(t *testing.T) {
	mock := &testutil.MockDB{
		GetAllAreasFunc: func() ([]*models.CodeDecodeArea, error) {
			return []*models.CodeDecodeArea{
				{ID: "a1", Code: 1, Decode: "Central"},
			}, nil
		},
	}
	cc := controllers.NewCodestableController(newCodestableStore(mock))

	r := httptest.NewRequest(http.MethodGet, "/areas", nil)
	w := httptest.NewRecorder()
	cc.GetAllAreas(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status: want 200, got %d", w.Code)
	}
}

func TestGetSubCategoriesByCategory_ValidCode_Returns200(t *testing.T) {
	mock := &testutil.MockDB{
		GetAllSubCategoryByCategoryCodeFunc: func(categoryCode int) ([]*models.CodeDecodeSubCategory, error) {
			return []*models.CodeDecodeSubCategory{
				{ID: "s2", Code: 102, Decode: "Korean"},
			}, nil
		},
	}
	cc := controllers.NewCodestableController(newCodestableStore(mock))

	r := httptest.NewRequest(http.MethodGet, "/subcategories/1", nil)
	r = withChiParam(r, "code", "1")
	w := httptest.NewRecorder()
	cc.GetSubCategoriesByCategory(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status: want 200, got %d", w.Code)
	}
}
