package helper_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ngfenglong/ikou-backend/internal/helper"
	"golang.org/x/crypto/bcrypt"
)

// --- WriteJSONResponse ---

func TestWriteJSONResponse_SetsStatusAndContentType(t *testing.T) {
	w := httptest.NewRecorder()
	payload := map[string]string{"key": "value"}

	err := helper.WriteJSONResponse(w, http.StatusCreated, payload)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.Code != http.StatusCreated {
		t.Errorf("status: want 201, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type: want application/json, got %q", ct)
	}

	var got map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("body not valid JSON: %v", err)
	}
	if got["key"] != "value" {
		t.Errorf("body key: want %q, got %q", "value", got["key"])
	}
}

// --- ReadJSON ---

func TestReadJSON_ValidBody_DecodesIntoTarget(t *testing.T) {
	body := `{"username":"alice","password":"secret"}`
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	w := httptest.NewRecorder()

	var got struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := helper.ReadJSON(w, r, &got)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Username != "alice" {
		t.Errorf("username: want alice, got %q", got.Username)
	}
}

func TestReadJSON_DoubleJSON_ReturnsError(t *testing.T) {
	body := `{"x":1}{"y":2}`
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	w := httptest.NewRecorder()

	var got map[string]int
	err := helper.ReadJSON(w, r, &got)

	if err == nil {
		t.Fatal("expected error for double JSON value, got nil")
	}
}

func TestReadJSON_OversizedBody_ReturnsError(t *testing.T) {
	body := strings.Repeat("a", 1048577) // 1 byte over the 1 MB limit
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"x":"`+body+`"}`))
	w := httptest.NewRecorder()

	var got map[string]string
	err := helper.ReadJSON(w, r, &got)

	if err == nil {
		t.Fatal("expected error for oversized body, got nil")
	}
}

func TestReadJSON_InvalidJSON_ReturnsError(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`not json`))
	w := httptest.NewRecorder()

	var got map[string]string
	err := helper.ReadJSON(w, r, &got)

	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

// --- BadRequest ---

func TestBadRequest_Returns400WithErrorBody(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	helper.BadRequest(w, r, errors.New("something went wrong"))

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: want 400, got %d", w.Code)
	}
	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["error"] != true {
		t.Error("expected error=true in response body")
	}
	if body["message"] != "something went wrong" {
		t.Errorf("message: got %q", body["message"])
	}
}

// --- Unauthorized ---

func TestUnauthorized_Returns401(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	helper.Unauthorized(w, r, errors.New("not allowed"))

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status: want 401, got %d", w.Code)
	}
}

// --- InternalServerError ---

func TestInternalServerError_Returns500(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	helper.InternalServerError(w, r, errors.New("boom"))

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status: want 500, got %d", w.Code)
	}
}

// --- InvalidCredential ---

func TestInvalidCredential_Returns401WithFixedMessage(t *testing.T) {
	w := httptest.NewRecorder()

	helper.InvalidCredential(w)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status: want 401, got %d", w.Code)
	}
	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["message"] != "Invalid authentication credential" {
		t.Errorf("message: got %q", body["message"])
	}
}

// --- ConflictErrorResponse ---

func TestConflictErrorResponse_Returns409(t *testing.T) {
	w := httptest.NewRecorder()

	helper.ConflictErrorResponse("Username already exists", w)

	if w.Code != http.StatusConflict {
		t.Errorf("status: want 409, got %d", w.Code)
	}
	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["message"] != "Username already exists" {
		t.Errorf("message: got %q", body["message"])
	}
}

// --- PasswordMatches ---

func TestPasswordMatches_CorrectPassword_ReturnsTrue(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret123"), 4)

	ok, err := helper.PasswordMatches(string(hash), "secret123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected true for matching password")
	}
}

func TestPasswordMatches_WrongPassword_ReturnsFalse(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret123"), 4)

	ok, err := helper.PasswordMatches(string(hash), "wrong")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected false for wrong password")
	}
}

func TestPasswordMatches_InvalidHash_ReturnsError(t *testing.T) {
	_, err := helper.PasswordMatches("not-a-hash", "password")

	if err == nil {
		t.Fatal("expected error for invalid hash, got nil")
	}
}

// --- WriteJSONResponse with custom header ---

func TestWriteJSONResponse_CustomHeader_Applied(t *testing.T) {
	w := httptest.NewRecorder()
	extra := http.Header{"X-Custom": []string{"yes"}}

	helper.WriteJSONResponse(w, http.StatusOK, struct{}{}, extra)

	if w.Header().Get("X-Custom") != "yes" {
		t.Errorf("custom header not applied, got %q", w.Header().Get("X-Custom"))
	}
}

// ensure ReadJSON does not accept an empty body without a JSON value
func TestReadJSON_EmptyBody_ReturnsError(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	var got map[string]string
	err := helper.ReadJSON(w, r, &got)

	if err == nil {
		t.Fatal("expected error for empty body, got nil")
	}
}
