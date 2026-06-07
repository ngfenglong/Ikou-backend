package controllers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ngfenglong/ikou-backend/api/controllers"
	"github.com/ngfenglong/ikou-backend/api/dto"
	"github.com/ngfenglong/ikou-backend/api/models"
	"github.com/ngfenglong/ikou-backend/api/store"
	"github.com/ngfenglong/ikou-backend/api/testutil"
	"github.com/ngfenglong/ikou-backend/internal/helper"
	"golang.org/x/crypto/bcrypt"
)

func newAuthStore(mock *testutil.MockDB) *store.Store {
	return &store.Store{DB: mock}
}

var validHash []byte

func init() {
	h, err := bcrypt.GenerateFromPassword([]byte("correct-password"), 4)
	if err != nil {
		panic(err)
	}
	validHash = h
}

func validUser() *models.User {
	return &models.User{
		ID:           "user-1",
		Username:     "alice",
		Email:        "alice@example.com",
		Password:     string(validHash),
		FirstName:    "Alice",
		LastName:     "Smith",
		ProfileImage: "/img.jpg",
	}
}

// --- Login ---

func TestLogin_ValidCredentials_Returns200WithTokens(t *testing.T) {
	mock := &testutil.MockDB{
		GetUserByUsernameFunc: func(username string) (*models.User, error) {
			return validUser(), nil
		},
		InsertTokenFunc: func(userId, refreshToken string, expiresAt time.Time) error {
			return nil
		},
	}
	ac := controllers.NewAuthController(newAuthStore(mock))

	body := `{"username":"alice","password":"correct-password"}`
	r := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	w := httptest.NewRecorder()
	ac.Login(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status: want 200, got %d — body: %s", w.Code, w.Body.String())
	}
	var resp dto.LoginResponseDTO
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.AccessToken == "" {
		t.Error("expected non-empty access_token")
	}
	if resp.RefreshToken == "" {
		t.Error("expected non-empty refresh_token")
	}
	if resp.User.UserName != "alice" {
		t.Errorf("user.username: want alice, got %q", resp.User.UserName)
	}
}

func TestLogin_InvalidJSON_Returns400(t *testing.T) {
	mock := &testutil.MockDB{}
	ac := controllers.NewAuthController(newAuthStore(mock))

	r := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("not json"))
	w := httptest.NewRecorder()
	ac.Login(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: want 400, got %d", w.Code)
	}
}

func TestLogin_UnknownUser_Returns401(t *testing.T) {
	mock := &testutil.MockDB{
		GetUserByUsernameFunc: func(username string) (*models.User, error) {
			return nil, errors.New("not found")
		},
	}
	ac := controllers.NewAuthController(newAuthStore(mock))

	body := `{"username":"ghost","password":"pass"}`
	r := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	w := httptest.NewRecorder()
	ac.Login(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status: want 401, got %d", w.Code)
	}
}

func TestLogin_WrongPassword_Returns401(t *testing.T) {
	mock := &testutil.MockDB{
		GetUserByUsernameFunc: func(username string) (*models.User, error) {
			return validUser(), nil
		},
	}
	ac := controllers.NewAuthController(newAuthStore(mock))

	body := `{"username":"alice","password":"wrong-password"}`
	r := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	w := httptest.NewRecorder()
	ac.Login(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status: want 401, got %d", w.Code)
	}
}

// --- Register ---

func TestRegister_Success_Returns201(t *testing.T) {
	mock := &testutil.MockDB{
		CheckIfUserExistsFunc: func(r dto.RegisterFormInputDTO) (bool, bool, error) {
			return false, false, nil
		},
		RegisterUserFunc: func(r dto.RegisterFormInputDTO) error {
			return nil
		},
	}
	ac := controllers.NewAuthController(newAuthStore(mock))

	body := `{"username":"bob","email":"bob@example.com","password":"pass123","first_name":"Bob","last_name":"Jones","country":"SG"}`
	r := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	w := httptest.NewRecorder()
	ac.Register(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("status: want 201, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestRegister_UsernameConflict_Returns409(t *testing.T) {
	mock := &testutil.MockDB{
		CheckIfUserExistsFunc: func(r dto.RegisterFormInputDTO) (bool, bool, error) {
			return true, false, nil // username exists
		},
	}
	ac := controllers.NewAuthController(newAuthStore(mock))

	body := `{"username":"alice","email":"new@example.com","password":"pass"}`
	r := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	w := httptest.NewRecorder()
	ac.Register(w, r)

	if w.Code != http.StatusConflict {
		t.Errorf("status: want 409, got %d", w.Code)
	}
}

func TestRegister_EmailConflict_Returns409(t *testing.T) {
	mock := &testutil.MockDB{
		CheckIfUserExistsFunc: func(r dto.RegisterFormInputDTO) (bool, bool, error) {
			return false, true, nil // email exists
		},
	}
	ac := controllers.NewAuthController(newAuthStore(mock))

	body := `{"username":"newuser","email":"alice@example.com","password":"pass"}`
	r := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	w := httptest.NewRecorder()
	ac.Register(w, r)

	if w.Code != http.StatusConflict {
		t.Errorf("status: want 409, got %d", w.Code)
	}
}

// --- RefreshToken ---

func TestRefreshToken_EmptyToken_Returns400(t *testing.T) {
	mock := &testutil.MockDB{}
	ac := controllers.NewAuthController(newAuthStore(mock))

	body := `{"refreshToken":""}`
	r := httptest.NewRequest(http.MethodPost, "/refresh", strings.NewReader(body))
	w := httptest.NewRecorder()
	ac.RefreshToken(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: want 400, got %d", w.Code)
	}
}

func TestRefreshToken_InvalidToken_Returns400(t *testing.T) {
	mock := &testutil.MockDB{}
	ac := controllers.NewAuthController(newAuthStore(mock))

	body := `{"refreshToken":"garbage.token.here"}`
	r := httptest.NewRequest(http.MethodPost, "/refresh", strings.NewReader(body))
	w := httptest.NewRecorder()
	ac.RefreshToken(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: want 400, got %d", w.Code)
	}
}

func TestRefreshToken_ValidToken_Returns200WithNewAccessToken(t *testing.T) {
	refreshToken, _, err := helper.GenerateRefreshToken("user-1")
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	mock := &testutil.MockDB{
		FetchRefreshTokenFromDBFunc: func(tokenStr string, userID string) (*models.RefreshToken, bool) {
			return &models.RefreshToken{
				UserID:    "user-1",
				Token:     tokenStr,
				ExpiresAt: time.Now().Add(time.Hour * 24),
			}, true
		},
		GetUserByIDFunc: func(userID string) (*models.User, error) {
			return validUser(), nil
		},
	}
	ac := controllers.NewAuthController(newAuthStore(mock))

	body := `{"refreshToken":"` + refreshToken + `"}`
	r := httptest.NewRequest(http.MethodPost, "/refresh", strings.NewReader(body))
	w := httptest.NewRecorder()
	ac.RefreshToken(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status: want 200, got %d — body: %s", w.Code, w.Body.String())
	}
	var resp dto.RefreshTokenResponseDTO
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.AccessToken == "" {
		t.Error("expected non-empty access_token in refresh response")
	}
}

// --- Logout ---

func TestLogout_ValidRequest_Returns200(t *testing.T) {
	mock := &testutil.MockDB{
		DeleteTokenFunc: func(refreshToken string) error {
			return nil
		},
	}
	ac := controllers.NewAuthController(newAuthStore(mock))

	body := `{"refresh_token":"some-token"}`
	r := httptest.NewRequest(http.MethodPost, "/logout", strings.NewReader(body))
	w := httptest.NewRecorder()
	ac.Logout(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status: want 200, got %d", w.Code)
	}
}
