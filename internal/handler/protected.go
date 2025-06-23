package handler

import (
	"encoding/json"
	"net/http"
	"pet-project/internal/middleware"
)

type ProtectedHandler struct{}

func (h *ProtectedHandler) ProtectedEndpoint(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "This is a protected endpoint",
		"user_id": userID,
	})
}
