package helper

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var accessTokenSecret = os.Getenv("JTW_ACCESS_SECRET")
var refreshTokenSecret = os.Getenv("JTW_REFRESH_SECRET")

type TokenDetail struct {
	ID           string
	Email        string
	Username     string
	ProfileImage string
	FirstName    string
	LastName     string
}

type RefreshTokenClaims struct {
	ID  string
	Exp int64
}

func GenerateAccessToken(td *TokenDetail) (string, time.Time, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	expiry := time.Now().Add(time.Hour * 24 * 3)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = td.ID
	claims["email"] = td.Email
	claims["username"] = td.Username
	claims["profile_image"] = td.ProfileImage
	claims["first_name"] = td.FirstName
	claims["last_name"] = td.LastName
	claims["exp"] = expiry.Unix()

	signedToken, err := token.SignedString([]byte(accessTokenSecret))
	if err != nil {
		return "", expiry, err
	}

	return signedToken, expiry, nil
}

func GenerateRefreshToken(userID string) (string, time.Time, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	expiry := time.Now().Add(time.Hour * 24 * 7)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = userID
	claims["exp"] = expiry.Unix()

	signedToken, err := token.SignedString([]byte(refreshTokenSecret))
	if err != nil {
		return "", expiry, err
	}

	return signedToken, expiry, nil
}

func VerifyAccessToken(tokenStr string) (bool, *TokenDetail) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(accessTokenSecret), nil
	})

	if err != nil {
		return false, nil
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Convert claims to your TokenDetail struct and return
		return true, &TokenDetail{
			ID:           claims["id"].(string),
			Email:        claims["email"].(string),
			Username:     claims["username"].(string),
			ProfileImage: claims["profile_image"].(string),
			FirstName:    claims["first_name"].(string),
			LastName:     claims["last_name"].(string),
		}
	}
	return false, nil
}

func VerifyJWTToken(tokenStr string) (bool, *RefreshTokenClaims) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(refreshTokenSecret), nil
	})

	if err != nil {
		return false, nil
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expValue, ok := claims["exp"].(float64) // JWT standard represents time in float64
		if !ok {
			return false, nil
		}

		// Convert claims to your RefreshTokenClaims struct and return
		return true, &RefreshTokenClaims{
			ID:  claims["id"].(string),
			Exp: int64(expValue),
		}
	}

	return false, nil
}

func IsTokenExpiryValid(expTime time.Time) bool {
	currentTime := time.Now()
	return currentTime.Before(expTime)
}
