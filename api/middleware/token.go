package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ngfenglong/ikou-backend/internal/helper"
)

type contextKey string

const UserIDKey contextKey = "userID"
const UserNameKey contextKey = "userName"

func ExtractTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		userName := ""
		if tokenClaims != nil {
			userID = tokenClaims.ID
			userName = tokenClaims.Username
		}

		// Storing the userID in the request's context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, UserNameKey, userName)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
