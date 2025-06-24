package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"pet-project/internal/service"
	"pet-project/pkg/model"
	"strconv"

	"github.com/go-chi/chi"
)

type CommentsHandler struct {
	CommentsService *service.CommentsService
}

type AddCommentRequest struct {
	TaskID int    `json:"taskID"`
	UserID int    `json:"userID"`
	Text   string `json:"text"`
}

type GetCommentsByTaskRequest struct {
	TaskID int `json:"taskID"`
}

type GetCommentsByUserRequest struct {
	UserID int `json:"userID"`
}

type UpdateCommentTextRequest struct {
	Text      string `json:"text"`
	CommentID int    `json:"comID"`
}

func (h *CommentsHandler) AddCommentRequest(w http.ResponseWriter, r *http.Request) {
	var req AddCommentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, errors.New("Invalid Request Body"), http.StatusBadRequest)
		return
	}

	if req.Text == "" {
		writeError(w, errors.New("Invalid Text"), http.StatusBadRequest)
		return
	}

	if req.TaskID <= 0 {
		writeError(w, errors.New("Invalid Task ID"), http.StatusBadRequest)
		return
	}

	if req.UserID <= 0 {
		writeError(w, errors.New("Invalid User ID"), http.StatusBadRequest)
		return
	}

	com := &model.Comments{
		TaskID: req.TaskID,
		UserID: req.UserID,
		Text:   req.Text,
	}

	err := h.CommentsService.AddComment(com)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	writeJSON(w, com, http.StatusCreated)
}

func (h *CommentsHandler) DeleteCommentRequest(w http.ResponseWriter, r *http.Request) {
	comID := chi.URLParam(r, "comID")

	intComID, err := strconv.Atoi(comID)
	if err != nil {
		writeError(w, errors.New("Invalid comment ID"), http.StatusBadRequest)
		return
	}

	err = h.CommentsService.DeleteComment(intComID)
	if err != nil {
		writeError(w, err, http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CommentsHandler) GetCommentsByTaskRequest(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")

	intTaskID, err := strconv.Atoi(taskID)
	if err != nil {
		writeError(w, errors.New("Invalid task ID"), http.StatusBadRequest)
		return
	}

	comments, err := h.CommentsService.GetCommentsByTask(intTaskID)
	if err != nil {
		writeError(w, err, http.StatusNotFound)
		return
	}
	writeJSON(w, comments)
}

func (h *CommentsHandler) GetCommentsByUserRequest(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	intUserID, err := strconv.Atoi(userID)
	if err != nil {
		writeError(w, errors.New("Invalid user ID"), http.StatusBadRequest)
		return
	}

	comments, err := h.CommentsService.GetCommentsByUser(intUserID)
	if err != nil {
		writeError(w, err, http.StatusNotFound)
		return
	}
	writeJSON(w, comments)
}

func (h *CommentsHandler) UpdateCommentTextRequest(w http.ResponseWriter, r *http.Request) {
	comID := chi.URLParam(r, "comID")

	var req UpdateCommentTextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, errors.New("Invalid Request Body"), http.StatusBadRequest)
		return
	}

	if req.Text == "" {
		writeError(w, errors.New("Comment text cannot be empty"), http.StatusBadRequest)
		return
	}

	intComID, err := strconv.Atoi(comID)
	if err != nil {
		writeError(w, errors.New("Invalid comment ID"), http.StatusBadRequest)
		return
	}

	err = h.CommentsService.UpdateCommentText(intComID, req.Text)
	if err != nil {
		writeError(w, err, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
