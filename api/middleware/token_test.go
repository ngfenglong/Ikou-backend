package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ngfenglong/ikou-backend/api/middleware"
	"github.com/ngfenglong/ikou-backend/internal/helper"
)

func handlerCapture(userID, userName *string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		*userID = r.Context().Value(middleware.UserIDKey).(string)
		*userName = r.Context().Value(middleware.UserNameKey).(string)
		w.WriteHeader(http.StatusOK)
	}
}

func TestExtractTokenMiddleware_NoHeader_EmptyContext(t *testing.T) {
	var gotID, gotName string
	handler := middleware.ExtractTokenMiddleware(handlerCapture(&gotID, &gotName))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if gotID != "" {
		t.Errorf("userID: want empty, got %q", gotID)
	}
	if gotName != "" {
		t.Errorf("userName: want empty, got %q", gotName)
	}
}

func TestExtractTokenMiddleware_ValidToken_PopulatesContext(t *testing.T) {
	td := &helper.TokenDetail{
		ID:       "user-42",
		Username: "alice",
		Email:    "alice@example.com",
	}
	tokenStr, _, err := helper.GenerateAccessToken(td)
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	var gotID, gotName string
	handler := middleware.ExtractTokenMiddleware(handlerCapture(&gotID, &gotName))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Authorization", "Bearer "+tokenStr)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if gotID != "user-42" {
		t.Errorf("userID: want user-42, got %q", gotID)
	}
	if gotName != "alice" {
		t.Errorf("userName: want alice, got %q", gotName)
	}
}

func TestExtractTokenMiddleware_InvalidToken_EmptyContext(t *testing.T) {
	var gotID, gotName string
	handler := middleware.ExtractTokenMiddleware(handlerCapture(&gotID, &gotName))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Authorization", "Bearer garbage.not.valid")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if gotID != "" {
		t.Errorf("userID: want empty for invalid token, got %q", gotID)
	}
	if gotName != "" {
		t.Errorf("userName: want empty for invalid token, got %q", gotName)
	}
}

func TestExtractTokenMiddleware_AlwaysCallsNext(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	// Wrap next in a handler that also sets context values (required by handlerCapture)
	wrapped := middleware.ExtractTokenMiddleware(func(w http.ResponseWriter, r *http.Request) {
		called = true
		// Check that context keys exist and are strings (even if empty)
		_ = r.Context().Value(middleware.UserIDKey).(string)
		_ = r.Context().Value(middleware.UserNameKey).(string)
	})
	_ = next
	_ = context.Background()

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, r)

	if !called {
		t.Error("expected next handler to be called")
	}
}
