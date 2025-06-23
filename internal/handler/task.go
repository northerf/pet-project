package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"pet-project/internal/service"
	"pet-project/pkg/model"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

type TaskHandler struct {
	TaskService *service.TaskService
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	AssignedTo  int    `json:"assignedTo"`
	ProjectID   int    `json:"projectID"`
	DueDate     string `json:"dueDate"`
}

type UpdateTaskRequest struct {
	TaskID      int     `json:"taskID"`
	AssignedTo  int     `json:"assignedTo"`
	Title       *string `json:"title"`
	Status      *string `json:"status"`
	Priority    *string `json:"priority"`
	Description *string `json:"description"`
}

type GetByIDTaskRequest struct {
	TaskID int `json:"taskID"`
}

type ListByProjectTaskRequest struct {
	ProjectID int `json:"projectID"`
}

type DeleteTaskRequest struct {
	TaskID     int `json:"taskID"`
	AssignedTo int `json:"assignedTo"`
}

func (h *TaskHandler) CreateTaskRequest(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, errors.New("Invalid Request Body"), http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		writeError(w, errors.New("Task required a name"), http.StatusBadRequest)
		return
	}

	if req.Status == "" {
		writeError(w, errors.New("Task required status"), http.StatusBadRequest)
		return
	}

	if req.Priority == "" {
		writeError(w, errors.New("Task required priority"), http.StatusBadRequest)
		return
	}

	if req.ProjectID == 0 {
		writeError(w, errors.New("Task required projectID"), http.StatusBadRequest)
		return
	}

	var dueDate *time.Time
	if req.DueDate != "" {
		t, err := time.Parse(time.RFC3339, req.DueDate)
		if err != nil {
			writeError(w, errors.New("Invalid dueDate format, use RFC3339"), http.StatusBadRequest)
			return
		}
		dueDate = &t
	}

	task := &model.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		AssignedTo:  req.AssignedTo,
		ProjectID:   req.ProjectID,
		DueDate:     dueDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.TaskService.CreateTask(task); err != nil {
		writeError(w, errors.New("Failed to create task"), http.StatusBadRequest)
		log.Println("Failed to create task", err)
		return
	}
}

func (h *TaskHandler) UpdateProjectRequest(w http.ResponseWriter, r *http.Request) {
	var req UpdateTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, errors.New("Invalid Request Body"), http.StatusBadRequest)
		return
	}

	if req.AssignedTo <= 0 {
		writeError(w, errors.New("Invalid assigned to"), http.StatusBadRequest)
		return
	}

	if *req.Title == "" {
		writeError(w, errors.New("Task required a name"), http.StatusBadRequest)
		return
	}

	if req.Status != nil && *req.Status != "pending" && *req.Status != "in_progress" && *req.Status != "done" {
		writeError(w, errors.New("Invalid status"), http.StatusBadRequest)
		return
	}
	if req.Priority != nil && *req.Priority != "low" && *req.Priority != "medium" && *req.Priority != "high" {
		writeError(w, errors.New("Invalid priority"), http.StatusBadRequest)
		return
	}

	task := &model.Task{
		ID:          req.TaskID,
		AssignedTo:  req.AssignedTo,
		Title:       *req.Title,
		Status:      *req.Status,
		Priority:    *req.Priority,
		Description: *req.Description,
	}
	task.UpdatedAt = time.Now()
	if req.Description == nil {
		existingTask, err := h.TaskService.GetByIDTask(req.TaskID, req.AssignedTo)
		if err != nil {
			writeError(w, errors.New("Failed to get existing task for description"), http.StatusBadRequest)
			return
		}
		task.Description = existingTask.Description
	}

	if err := h.TaskService.UpdateTask(task, req.AssignedTo); err != nil {
		writeError(w, errors.New("Failed to update task"), http.StatusBadRequest)
		return
	}
}

func (h *TaskHandler) GetByIDTaskRequest(w http.ResponseWriter, r *http.Request) {

	taskID := chi.URLParam(r, "taskID")
	AssignedToID := getUserIDFromContext(r)

	intTaskID, err := strconv.Atoi(taskID)
	if err != nil {
		writeError(w, errors.New("Invalid task ID"), http.StatusBadRequest)
	}

	task, err := h.TaskService.GetByIDTask(intTaskID, AssignedToID)
	if err != nil {
		writeError(w, err, http.StatusNotFound)
		return
	}
	writeJSON(w, task)
}

func (h *TaskHandler) ListByProjectTaskRequest(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	intProjectID, err := strconv.Atoi(projectID)
	if err != nil {
		writeError(w, errors.New("Invalid project ID"), http.StatusBadRequest)
	}

	task, err := h.TaskService.ListByProjectTask(intProjectID)
	if err != nil {
		writeError(w, err, http.StatusNotFound)
		return
	}
	writeJSON(w, task)
}

func (h *TaskHandler) DeleteTaskRequest(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	AssignedToID := getUserIDFromContext(r)

	intTaskID, err := strconv.Atoi(taskID)
	if err != nil {
		writeError(w, errors.New("Invalid task ID"), http.StatusBadRequest)
		return
	}

	task := h.TaskService.DeleteTask(intTaskID, AssignedToID)
	if err != nil {
		writeError(w, errors.New("Failed delete task"), http.StatusBadRequest)
		return
	}
	writeJSON(w, task)
}
