package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"pet-project/internal/middleware"
	"pet-project/internal/service"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

type ProjectHandler struct {
	ProjectService *service.ProjectService
}

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	OwnerID     int    `json:"ownerID"`
}

type GetProjectInfoRequest struct {
	ProjectID int `json:"projectID"`
}

type UpdateProjectRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type DeleteProjectRequest struct {
	ProjectID string `json:"projectID"`
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, errors.New("Invalid Request Body"), http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		writeError(w, errors.New("Project required a name"), http.StatusBadRequest)
		return
	}

	project, err := h.ProjectService.CreateProject(req.Name, req.Description, req.OwnerID)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	writeJSON(w, project, http.StatusCreated)
}

func (h *ProjectHandler) GetProjectInfo(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	ownerID := getUserIDFromContext(r)

	intProjectID, err := strconv.Atoi(projectID)
	if err != nil {
		writeError(w, errors.New("Invalid project ID"), http.StatusBadRequest)
		return
	}

	project, err := h.ProjectService.GetByIDProject(intProjectID, ownerID)
	if err != nil {
		writeError(w, err, http.StatusNotFound)
		return
	}
	writeJSON(w, project)
}

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	ownerID := getUserIDFromContext(r)

	intProjectID, err := strconv.Atoi(projectID)
	if err != nil {
		writeError(w, errors.New("Invalid project ID"), http.StatusBadRequest)
		return
	}

	err = h.ProjectService.DeleteProject(intProjectID, ownerID)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	ownerID := getUserIDFromContext(r)

	if ownerID == 0 {
		writeError(w, errors.New("Unauthorized owner"), http.StatusUnauthorized)
		return
	}

	intProjectID, err := strconv.Atoi(projectID)
	if err != nil {
		writeError(w, errors.New("Invalid project ID"), http.StatusBadRequest)
		return
	}

	var req UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, errors.New("Invalid Request Body"), http.StatusBadRequest)
		return
	}

	project, err := h.ProjectService.GetByIDProject(intProjectID, ownerID)
	if err != nil {
		writeError(w, err, http.StatusNotFound)
		return
	}

	if req.Name != nil {
		project.Name = *req.Name
	}
	if req.Description != nil {
		project.Description = *req.Description
	}

	project.UpdatedAt = time.Now()

	err = h.ProjectService.UpdateProject(project)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, project)
}

func getUserIDFromContext(r *http.Request) int {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		return 0
	}
	return userID
}
