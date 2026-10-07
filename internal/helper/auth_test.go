package helper_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/ngfenglong/ikou-backend/internal/helper"
)

var testTokenDetail = &helper.TokenDetail{
	ID:           "user-abc123",
	Email:        "tester@example.com",
	Username:     "tester",
	ProfileImage: "https://example.com/avatar.png",
	FirstName:    "Test",
	LastName:     "User",
}

// --- GenerateAccessToken ---

func TestGenerateAccessToken_ReturnsSignedToken(t *testing.T) {
	tokenStr, expiry, err := helper.GenerateAccessToken(testTokenDetail)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tokenStr == "" {
		t.Fatal("expected non-empty token string")
	}
	if time.Until(expiry) < 70*time.Hour || time.Until(expiry) > 73*time.Hour {
		t.Errorf("expected expiry ~72h from now, got %v", time.Until(expiry))
	}
}

// --- VerifyAccessToken ---

func TestVerifyAccessToken_ValidToken_ReturnsClaimsRoundTrip(t *testing.T) {
	tokenStr, _, err := helper.GenerateAccessToken(testTokenDetail)
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	valid, got := helper.VerifyAccessToken(tokenStr)

	if !valid {
		t.Fatal("expected token to be valid")
	}
	if got == nil {
		t.Fatal("expected non-nil TokenDetail")
	}
	if got.ID != testTokenDetail.ID {
		t.Errorf("ID: want %q, got %q", testTokenDetail.ID, got.ID)
	}
	if got.Email != testTokenDetail.Email {
		t.Errorf("Email: want %q, got %q", testTokenDetail.Email, got.Email)
	}
	if got.Username != testTokenDetail.Username {
		t.Errorf("Username: want %q, got %q", testTokenDetail.Username, got.Username)
	}
	if got.FirstName != testTokenDetail.FirstName {
		t.Errorf("FirstName: want %q, got %q", testTokenDetail.FirstName, got.FirstName)
	}
	if got.LastName != testTokenDetail.LastName {
		t.Errorf("LastName: want %q, got %q", testTokenDetail.LastName, got.LastName)
	}
}

func TestVerifyAccessToken_GarbageString_ReturnsFalse(t *testing.T) {
	valid, claims := helper.VerifyAccessToken("not.a.token")

	if valid {
		t.Error("expected invalid token to return false")
	}
	if claims != nil {
		t.Error("expected nil claims for invalid token")
	}
}

func TestVerifyAccessToken_WrongSecret_ReturnsFalse(t *testing.T) {
	// Craft a token signed with a non-empty secret the server doesn't know.
	tok := jwt.New(jwt.SigningMethodHS256)
	claims := tok.Claims.(jwt.MapClaims)
	claims["id"] = "x"
	claims["email"] = "x@x.com"
	claims["username"] = "x"
	claims["profile_image"] = ""
	claims["first_name"] = "x"
	claims["last_name"] = "x"
	claims["exp"] = time.Now().Add(time.Hour).Unix()

	tokenStr, err := tok.SignedString([]byte("wrong-secret"))
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	valid, detail := helper.VerifyAccessToken(tokenStr)

	if valid {
		t.Error("expected token signed with wrong secret to be rejected")
	}
	if detail != nil {
		t.Error("expected nil TokenDetail for rejected token")
	}
}

// --- GenerateRefreshToken ---

func TestGenerateRefreshToken_ReturnsSignedToken(t *testing.T) {
	tokenStr, expiry, err := helper.GenerateRefreshToken("user-abc123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tokenStr == "" {
		t.Fatal("expected non-empty refresh token string")
	}
	if time.Until(expiry) < 166*time.Hour || time.Until(expiry) > 169*time.Hour {
		t.Errorf("expected expiry ~168h from now, got %v", time.Until(expiry))
	}
}

// --- VerifyRefreshToken ---

func TestVerifyRefreshToken_ValidToken_ReturnsClaimsRoundTrip(t *testing.T) {
	userID := "user-abc123"
	tokenStr, _, err := helper.GenerateRefreshToken(userID)
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	valid, got := helper.VerifyRefreshToken(tokenStr)

	if !valid {
		t.Fatal("expected refresh token to be valid")
	}
	if got == nil {
		t.Fatal("expected non-nil RefreshTokenClaims")
	}
	if got.ID != userID {
		t.Errorf("ID: want %q, got %q", userID, got.ID)
	}
	if got.Exp <= 0 {
		t.Errorf("expected positive Exp, got %d", got.Exp)
	}
}

func TestVerifyRefreshToken_GarbageString_ReturnsFalse(t *testing.T) {
	valid, claims := helper.VerifyRefreshToken("garbage")

	if valid {
		t.Error("expected invalid refresh token to return false")
	}
	if claims != nil {
		t.Error("expected nil claims for invalid refresh token")
	}
}

func TestVerifyRefreshToken_WrongSecret_ReturnsFalse(t *testing.T) {
	tok := jwt.New(jwt.SigningMethodHS256)
	c := tok.Claims.(jwt.MapClaims)
	c["id"] = "user-x"
	c["exp"] = time.Now().Add(time.Hour).Unix()

	tokenStr, err := tok.SignedString([]byte("wrong-secret"))
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	valid, claims := helper.VerifyRefreshToken(tokenStr)

	if valid {
		t.Error("expected refresh token with wrong secret to be rejected")
	}
	if claims != nil {
		t.Error("expected nil claims for rejected refresh token")
	}
}

// --- IsTokenExpiryValid ---

func TestIsTokenExpiryValid_FutureTime_ReturnsTrue(t *testing.T) {
	if !helper.IsTokenExpiryValid(time.Now().Add(time.Hour)) {
		t.Error("expected future expiry to be valid")
	}
}

func TestIsTokenExpiryValid_PastTime_ReturnsFalse(t *testing.T) {
	if helper.IsTokenExpiryValid(time.Now().Add(-time.Second)) {
		t.Error("expected past expiry to be invalid")
	}
}
