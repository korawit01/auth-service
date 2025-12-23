package http

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/korawit01/auth-service/api-gateway/internal/client"
)

type Handler struct {
	authClient *client.AuthClient
	taskClient *client.TaskClient
}

func NewHandler(authClient *client.AuthClient, taskClient *client.TaskClient) *Handler {
	return &Handler{
		authClient: authClient,
		taskClient: taskClient,
	}
}

func (h *Handler) RegisterRoutes(r fiber.Router) {
	r.Use(requestid.New())
	r.Use(logger.New())
	r.Use(recover.New())
	r.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders: "Content-Type, Authorization, X-User-ID",
	}))

	r.Get("/health", func(c *fiber.Ctx) error {
		return writeJSON(c, fiber.StatusOK, map[string]string{"status": "ok"})
	})

	r.Post("/register", h.handleRegister)
	r.Post("/login", h.handleLogin)

	r.Post("/tasks", h.handleCreateTask)
	r.Get("/tasks", h.handleGetTasks)
	r.Get("/tasks/:id", h.handleGetTask)
	r.Put("/tasks/:id", h.handleUpdateTask)
	r.Delete("/tasks/:id", h.handleDeleteTask)

	// Swagger docs (live generated from annotations)
	r.Get("/swagger/*", SwaggerRoutes())
}

// handleRegister godoc
// @Summary Register user
// @Description Create a new user via auth-service.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "User credentials"
// @Success 200 {object} client.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /register [post]
func (h *Handler) handleRegister(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "bad request")
	}

	user, err := h.authClient.Register(c.UserContext(), req.Email, req.Password)
	if err != nil {
		log.Printf("register error: %v", err)
		return writeError(c, fiber.StatusBadGateway, "register failed")
	}

	return writeJSON(c, fiber.StatusOK, user)
}

// handleLogin godoc
// @Summary Login
// @Description Login via auth-service and receive JWT.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "User credentials"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /login [post]
func (h *Handler) handleLogin(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "bad request")
	}

	token, err := h.authClient.Login(c.UserContext(), req.Email, req.Password)
	if err != nil {
		log.Printf("login error: %v", err)
		return writeError(c, fiber.StatusBadGateway, "login failed")
	}

	return writeJSON(c, fiber.StatusOK, TokenResponse{Token: token})
}

// handleCreateTask godoc
// @Summary Create task
// @Description Create a task via task-service.
// @Tags tasks
// @Accept json
// @Produce json
// @Param X-User-ID header string true "User ID"
// @Param request body CreateTaskRequest true "Task payload"
// @Success 201 {object} TaskResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /tasks [post]
func (h *Handler) handleCreateTask(c *fiber.Ctx) error {
	userID, err := parseUserID(c)
	if err != nil {
		return writeError(c, fiber.StatusUnauthorized, err.Error())
	}

	var req CreateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "bad request")
	}

	task, err := h.taskClient.CreateTask(c.UserContext(), userID, req.Title, req.Description, req.Status)
	if err != nil {
		log.Printf("create task error: %v", err)
		if te, ok := err.(*client.TaskError); ok {
			return writeError(c, te.StatusCode, te.Message)
		}
		return writeError(c, fiber.StatusBadGateway, fmt.Sprintf("create task failed: %s", err.Error()))
	}

	return writeJSON(c, fiber.StatusCreated, TaskResponse(*task))
}

// handleGetTasks godoc
// @Summary List tasks
// @Description List tasks belonging to the provided user.
// @Tags tasks
// @Produce json
// @Param X-User-ID header string true "User ID"
// @Success 200 {array} TaskResponse
// @Failure 401 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /tasks [get]
func (h *Handler) handleGetTasks(c *fiber.Ctx) error {
	userID, err := parseUserID(c)
	if err != nil {
		return writeError(c, fiber.StatusUnauthorized, err.Error())
	}

	task, err := h.taskClient.GetTasks(c.UserContext(), userID)
	if err != nil {
		log.Printf("get tasks error: %v", err)
		if te, ok := err.(*client.TaskError); ok {
			return writeError(c, te.StatusCode, te.Message)
		}
		return writeError(c, fiber.StatusBadGateway, fmt.Sprintf("get tasks failed: %s", err.Error()))
	}

	var res []TaskResponse
	if task != nil {
		res = make([]TaskResponse, len(*task))
		for i, t := range *task {
			res[i] = TaskResponse(t)
		}
	} else {
		res = []TaskResponse{}
	}

	return writeJSON(c, fiber.StatusOK, res)
}

// handleGetTask godoc
// @Summary Get task
// @Description Get a task belonging to the provided user.
// @Tags tasks
// @Produce json
// @Param X-User-ID header string true "User ID"
// @Param id path string true "Task ID"
// @Success 200 {object} TaskResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /tasks/{id} [get]
func (h *Handler) handleGetTask(c *fiber.Ctx) error {
	userID, err := parseUserID(c)
	if err != nil {
		return writeError(c, fiber.StatusUnauthorized, err.Error())
	}
	id := c.Params("id")

	task, err := h.taskClient.GetTask(c.UserContext(), userID, id)
	if err != nil {
		log.Printf("get task error: %v", err)
		if te, ok := err.(*client.TaskError); ok {
			return writeError(c, te.StatusCode, te.Message)
		}
		return writeError(c, fiber.StatusBadGateway, fmt.Sprintf("get task failed: %s", err.Error()))
	}

	return writeJSON(c, fiber.StatusOK, TaskResponse(*task))
}

// handleUpdateTask godoc
// @Summary Update task
// @Description Update a task belonging to the provided user.
// @Tags tasks
// @Accept json
// @Produce json
// @Param X-User-ID header string true "User ID"
// @Param id path string true "Task ID"
// @Param request body UpdateTaskInput true "Task payload"
// @Success 200 {object} TaskResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /tasks/{id} [put]
func (h *Handler) handleUpdateTask(c *fiber.Ctx) error {
	userID, err := parseUserID(c)
	if err != nil {
		return writeError(c, fiber.StatusUnauthorized, err.Error())
	}
	id := c.Params("id")

	var req UpdateTaskInput
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "bad request")
	}

	task, err := h.taskClient.UpdateTask(c.UserContext(), userID, id, client.UpdateTaskInput(req))
	if err != nil {
		log.Printf("update task error: %v", err)
		if te, ok := err.(*client.TaskError); ok {
			return writeError(c, te.StatusCode, te.Message)
		}
		return writeError(c, fiber.StatusBadGateway, fmt.Sprintf("update task failed: %s", err.Error()))
	}

	return writeJSON(c, fiber.StatusOK, TaskResponse(*task))
}

// handleDeleteTask godoc
// @Summary Delete task
// @Description Delete a task belonging to the provided user.
// @Tags tasks
// @Produce json
// @Param X-User-ID header string true "User ID"
// @Param id path string true "Task ID"
// @Success 204 {string} string "no content"
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /tasks/{id} [delete]
func (h *Handler) handleDeleteTask(c *fiber.Ctx) error {
	userID, err := parseUserID(c)
	if err != nil {
		return writeError(c, fiber.StatusUnauthorized, err.Error())
	}
	id := c.Params("id")

	if err := h.taskClient.DeleteTask(c.UserContext(), userID, id); err != nil {
		log.Printf("delete task error: %v", err)
		if te, ok := err.(*client.TaskError); ok {
			return writeError(c, te.StatusCode, te.Message)
		}
		return writeError(c, fiber.StatusBadGateway, fmt.Sprintf("delete task failed: %s", err.Error()))
	}

	return c.SendStatus(fiber.StatusNoContent)
}
