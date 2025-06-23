package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"pet-project/internal/service"
	"pet-project/pkg/model"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type updateUserRequest struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type deleteUserRequest struct {
	Email string `json:"email"`
}

func (h *AuthHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	if req.Name == nil {
		writeError(w, errors.New("Name is required"), http.StatusBadRequest)
	}
	if req.Email == nil {
		writeError(w, errors.New("Email is required"), http.StatusBadRequest)
	}
	if req.Password == nil {
		writeError(w, errors.New("Password is required"), http.StatusBadRequest)
	}

	user := &model.User{
		Name:     *req.Name,
		Email:    *req.Email,
		Password: *req.Password,
	}

	err := h.AuthService.UpdateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{
	"message":"User updated successfully"
	}`))
}

func (h *AuthHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var req deleteUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}
	err := h.AuthService.DeleteUser(req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{
	"message":"User deleted successfully"
	}`))
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}
	err := h.AuthService.Register(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{
	"message":"User registred successfully"
	}`))
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}
	token, err := h.AuthService.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
