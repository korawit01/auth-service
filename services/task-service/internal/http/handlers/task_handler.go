package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/korawit01/auth-service/services/task-service/internal/domain"
	"github.com/korawit01/auth-service/services/task-service/internal/repository"
	"github.com/korawit01/auth-service/services/task-service/internal/service"
)

type TaskHandler struct {
	svc service.TaskService
}

func NewTaskHandler(svc service.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

func (h *TaskHandler) RegisterRoutes(r chi.Router) {
	r.Get("/tasks", h.listTasks)
	r.Get("/tasks/{id}", h.getTask)
	r.Post("/tasks", h.createTask)
	r.Put("/tasks/{id}", h.updateTask)
	r.Delete("/tasks/{id}", h.deleteTask)
}

type createTaskRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Status      string  `json:"status"`
	DueDate     *string `json:"dueDate"` // ISO8601 string
}

type updateTaskRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
	DueDate     *string `json:"dueDate"`
}

func (h *TaskHandler) listTasks(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUserID(r)
	if err != nil {
		handleError(w, err)
		return
	}

	tasks, err := h.svc.ListTasks(r.Context(), userID)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) getTask(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUserID(r)
	if err != nil {
		handleError(w, err)
		return
	}

	idStr := chi.URLParam(r, "id")
	task, err := h.svc.GetTask(r.Context(), userID, idStr)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUserID(r)
	if err != nil {
		handleError(w, err)
		return
	}

	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	var due *time.Time
	if req.DueDate != nil && *req.DueDate != "" {
		t, err := time.Parse(time.RFC3339, *req.DueDate)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid dueDate")
			return
		}
		due = &t
	}

	input := service.CreateTaskInput{
		Title:       req.Title,
		Description: req.Description,
		Status:      domain.TaskStatus(req.Status),
		DueDate:     due,
	}

	task, err := h.svc.CreateTask(r.Context(), userID, input)
	if err != nil {
		if errors.Is(err, service.ErrInvalidStatus) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) updateTask(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUserID(r)
	if err != nil {
		handleError(w, err)
		return
	}

	idStr := chi.URLParam(r, "id")
	var req updateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	var due *time.Time
	if req.DueDate != nil {
		if *req.DueDate == "" {
			due = nil
		} else {
			t, err := time.Parse(time.RFC3339, *req.DueDate)
			if err != nil {
				writeError(w, http.StatusBadRequest, "invalid dueDate")
				return
			}
			due = &t
		}
	}

	var statusPtr *domain.TaskStatus
	if req.Status != nil {
		s := domain.TaskStatus(*req.Status)
		statusPtr = &s
	}

	input := service.UpdateTaskInput{
		Title:       req.Title,
		Description: req.Description,
		Status:      statusPtr,
		DueDate:     due,
	}

	task, err := h.svc.UpdateTask(r.Context(), userID, idStr, input)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		if errors.Is(err, service.ErrInvalidStatus) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) deleteTask(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUserID(r)
	if err != nil {
		handleError(w, err)
		return
	}

	idStr := chi.URLParam(r, "id")
	if err := h.svc.DeleteTask(r.Context(), userID, idStr); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
