package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type contextKey string

const userIDKey contextKey = "userID"

func AuthMiddleware(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return secret, nil
			})

			if err != nil {
				log.Printf("Token parsing error: %v", err)
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				log.Printf("Token is invalid")
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				log.Printf("Failed to parse claims")
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			userIDFloat, ok := claims["user_id"].(float64)
			if !ok {
				log.Printf("Failed to get user_id from claims: %v", claims)
				http.Error(w, "User ID not found in token", http.StatusUnauthorized)
				return
			}

			userID := int(userIDFloat)
			log.Printf("Successfully authenticated user ID: %d", userID)

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func GetUserID(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(userIDKey).(int)
	return userID, ok
}
