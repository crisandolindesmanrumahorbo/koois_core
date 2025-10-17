package middleware

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"koois_core/internal/config"
	"net/http"
	"strings"
)

type Claims struct {
	Username string `json:"username"`
	Sub      string `json:"sub"`
	RoleId   int    `json:"role_id"`
	jwt.RegisteredClaims
}

type contextKey string

const (
	ClaimsContextKey contextKey = "claims"
)

func GetClaimsFromContext(r *http.Request) (*Claims, error) {
	claims, ok := r.Context().Value(ClaimsContextKey).(*Claims)
	if !ok {
		return nil, fmt.Errorf("claims not found in context")
	}
	return claims, nil
}

func JWT(next http.HandlerFunc, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractToken(r)
		if tokenString == "" {
			http.Error(w, `{"error":"missing token"}`, http.StatusUnauthorized)
			return
		}
		publicKeyStr := strings.ReplaceAll(cfg.JWT.Secret, "\\n", "\n")
		publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyStr))
		if err != nil {
			http.Error(w, `{"error":"invalid env"}`, http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return publicKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
			return
		}

		// Add claims to request context
		ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)

		next(w, r.WithContext(ctx))
	}
}
func extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}
