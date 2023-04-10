package helper

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var accessTokenSecret = os.Getenv("JTW_ACCESS_SECRET")
var refreshTokenSecret = os.Getenv("JTW_REFRESH_SECRET")

type TokenDetail struct {
	ID       string
	Email    string
	Username string
}

func GenerateAccessToken(td *TokenDetail) (string, time.Time, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	expiry := time.Now().Add(time.Hour * 24 * 3)
	claims := token.Claims.(jwt.MapClaims)
	claims["ID"] = td.ID
	claims["Email"] = td.Email
	claims["Username"] = td.Username
	claims["exp"] = expiry

	signedToken, err := token.SignedString([]byte(accessTokenSecret))
	if err != nil {
		return "", expiry, err
	}

	return signedToken, expiry, nil
}

func GenerateRefreshToken(td *TokenDetail) (string, time.Time, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	expiry := time.Now().Add(time.Hour * 24 * 7)

	claims := token.Claims.(jwt.MapClaims)
	claims["ID"] = td.ID
	claims["exp"] = expiry

	signedToken, err := token.SignedString([]byte(accessTokenSecret))
	if err != nil {
		return "", expiry, err
	}

	return signedToken, expiry, nil
}
