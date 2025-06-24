package handler

import (
	"encoding/json"
	"net/http"
)

func writeError(w http.ResponseWriter, err error, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

func writeJSON(w http.ResponseWriter, data interface{}, statusCode ...int) {
	if len(statusCode) == 0 {
		statusCode = append(statusCode, http.StatusOK)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode[0])
	json.NewEncoder(w).Encode(data)
}
