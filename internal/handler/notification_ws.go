package handler

import (
	"fmt"
	"net/http"
	"pet-project/internal/realtime"
	"pet-project/internal/service"

	"github.com/dgrijalva/jwt-go"
)

type NotificationWSHandler struct {
	ClientManager *realtime.ClientManager
	JwtSecret     []byte
	AuthService   *service.AuthService
}

func parseUserIDFromToken(tokenString string, secret []byte) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token: %w", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("user_id not found in token claims")
	}
	return int(userIDFloat), nil
}

func (h *NotificationWSHandler) WSNotifications(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Unauthorized: no token", http.StatusUnauthorized)
		return
	}

	userID, err := parseUserIDFromToken(token, h.JwtSecret)
	if err != nil || userID == 0 {
		http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
		return
	}

	h.ClientManager.ServeWS(w, r, userID)
}
