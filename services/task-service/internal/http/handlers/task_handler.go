package handlers

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"

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

func (h *TaskHandler) RegisterRoutes(r fiber.Router) {
	r.Get("/tasks", h.listTasks)
	r.Get("/tasks/:id", h.getTask)
	r.Post("/tasks", h.createTask)
	r.Put("/tasks/:id", h.updateTask)
	r.Delete("/tasks/:id", h.deleteTask)
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

func (h *TaskHandler) listTasks(c *fiber.Ctx) error {
	userID, err := parseUserID(c)
	if err != nil {
		return handleError(c, err)
	}

	tasks, err := h.svc.ListTasks(c.UserContext(), userID)
	if err != nil {
		return handleError(c, err)
	}
	return writeJSON(c, fiber.StatusOK, tasks)
}

func (h *TaskHandler) getTask(c *fiber.Ctx) error {
	userID, err := parseUserID(c)
	if err != nil {
		return handleError(c, err)
	}

	idStr := c.Params("id")
	task, err := h.svc.GetTask(c.UserContext(), userID, idStr)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return writeError(c, fiber.StatusNotFound, "task not found")
		}
		return handleError(c, err)
	}
	return writeJSON(c, fiber.StatusOK, task)
}

func (h *TaskHandler) createTask(c *fiber.Ctx) error {
	userID, err := parseUserID(c)
	if err != nil {
		return handleError(c, err)
	}

	var req createTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "invalid JSON")
	}
	if req.Title == "" {
		return writeError(c, fiber.StatusBadRequest, "title is required")
	}

	var due *time.Time
	if req.DueDate != nil && *req.DueDate != "" {
		t, err := time.Parse(time.RFC3339, *req.DueDate)
		if err != nil {
			return writeError(c, fiber.StatusBadRequest, "invalid dueDate")
		}
		due = &t
	}

	input := service.CreateTaskInput{
		Title:       req.Title,
		Description: req.Description,
		Status:      domain.TaskStatus(req.Status),
		DueDate:     due,
	}

	task, err := h.svc.CreateTask(c.UserContext(), userID, input)
	if err != nil {
		if errors.Is(err, service.ErrInvalidStatus) {
			return writeError(c, fiber.StatusBadRequest, err.Error())
		}
		return handleError(c, err)
	}
	return writeJSON(c, fiber.StatusCreated, task)
}

func (h *TaskHandler) updateTask(c *fiber.Ctx) error {
	userID, err := parseUserID(c)
	if err != nil {
		return handleError(c, err)
	}

	idStr := c.Params("id")
	var req updateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "invalid JSON")
	}

	var due *time.Time
	if req.DueDate != nil {
		if *req.DueDate == "" {
			due = nil
		} else {
			t, err := time.Parse(time.RFC3339, *req.DueDate)
			if err != nil {
				return writeError(c, fiber.StatusBadRequest, "invalid dueDate")
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

	task, err := h.svc.UpdateTask(c.UserContext(), userID, idStr, input)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return writeError(c, fiber.StatusNotFound, "task not found")
		}
		if errors.Is(err, service.ErrInvalidStatus) {
			return writeError(c, fiber.StatusBadRequest, err.Error())
		}
		return handleError(c, err)
	}
	return writeJSON(c, fiber.StatusOK, task)
}

func (h *TaskHandler) deleteTask(c *fiber.Ctx) error {
	userID, err := parseUserID(c)
	if err != nil {
		return handleError(c, err)
	}

	idStr := c.Params("id")
	if err := h.svc.DeleteTask(c.UserContext(), userID, idStr); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return writeError(c, fiber.StatusNotFound, "task not found")
		}
		return handleError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
