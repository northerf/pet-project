package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"pet-project/internal/service"
	"pet-project/pkg/model"
	"strconv"
)

type NotificationHandler struct {
	NotificationService *service.NotificationService
}

type CreateNotificationRequest struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type MarkAsReadRequest struct {
	NotificationIDs []int `json:"notification_ids"`
}

func (h *NotificationHandler) CreateNotification(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		writeError(w, errors.New("Unauthorized"), http.StatusUnauthorized)
		return
	}

	var req CreateNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, errors.New("Invalid Request Body"), http.StatusBadRequest)
		return
	}

	if req.Message == "" {
		writeError(w, errors.New("Message is required"), http.StatusBadRequest)
		return
	}

	notif := &model.Notification{
		UserID:  userID,
		Type:    req.Type,
		Message: req.Message,
	}

	if err := h.NotificationService.Create(r.Context(), notif); err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{
		"id":      notif.ID,
		"message": "Notification created successfully",
	}, http.StatusCreated)
}

func (h *NotificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		writeError(w, errors.New("Unauthorized"), http.StatusUnauthorized)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10 // дефолтное значение
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	notifications, err := h.NotificationService.GetByUserID(r.Context(), userID, limit, offset)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, notifications)
}

func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		writeError(w, errors.New("Unauthorized"), http.StatusUnauthorized)
		return
	}

	var req MarkAsReadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, errors.New("Invalid Request Body"), http.StatusBadRequest)
		return
	}

	if len(req.NotificationIDs) == 0 {
		writeError(w, errors.New("No notification IDs provided"), http.StatusBadRequest)
		return
	}

	err := h.NotificationService.MarkAsRead(r.Context(), userID, req.NotificationIDs)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{
		"message": "Notifications marked as read",
	})
}

func (h *NotificationHandler) CountUnread(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		writeError(w, errors.New("Unauthorized"), http.StatusUnauthorized)
		return
	}

	count, err := h.NotificationService.CountUnread(r.Context(), userID)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{
		"unread_count": count,
	})
}
