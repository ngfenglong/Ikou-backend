package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ngfenglong/ikou-backend/internal/helper"
)

type contextKey string

const UserIDKey contextKey = "userID"

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
		if tokenClaims != nil {
			userID = tokenClaims.ID
		}

		// Storing the userID in the request's context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
