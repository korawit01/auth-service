package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

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

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Post("/register", h.handleRegister)
	r.Post("/login", h.handleLogin)

	r.Post("/tasks", h.handleCreateTask)
	r.Get("/tasks", h.handleGetTasks)
	r.Get("/tasks/{id}", h.handleGetTask)
	r.Put("/tasks/{id}", h.handleUpdateTask)
	r.Delete("/tasks/{id}", h.handleDeleteTask)

	// Swagger docs (live generated from annotations)
	r.Handle("/swagger/*", SwaggerRoutes())

	return r
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
func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad request")
		return
	}

	user, err := h.authClient.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		log.Printf("register error: %v", err)
		writeError(w, http.StatusBadGateway, "register failed")
		return
	}

	writeJSON(w, http.StatusOK, user)
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
func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad request")
		return
	}

	token, err := h.authClient.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		log.Printf("login error: %v", err)
		writeError(w, http.StatusBadGateway, "login failed")
		return
	}

	writeJSON(w, http.StatusOK, TokenResponse{Token: token})
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
func (h *Handler) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUserID(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad request")
		return
	}

	task, err := h.taskClient.CreateTask(r.Context(), userID, req.Title, req.Description)
	if err != nil {
		log.Printf("create task error: %v", err)
		if te, ok := err.(*client.TaskError); ok {
			writeError(w, te.StatusCode, te.Message)
			return
		}
		writeError(w, http.StatusBadGateway, fmt.Sprintf("create task failed: %s", err.Error()))
		return
	}

	writeJSON(w, http.StatusCreated, TaskResponse(*task))
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
func (h *Handler) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUserID(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	task, err := h.taskClient.GetTasks(r.Context(), userID)
	if err != nil {
		log.Printf("get tasks error: %v", err)
		if te, ok := err.(*client.TaskError); ok {
			writeError(w, te.StatusCode, te.Message)
			return
		}
		writeError(w, http.StatusBadGateway, fmt.Sprintf("get tasks failed: %s", err.Error()))
		return
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

	writeJSON(w, http.StatusOK, res)
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
func (h *Handler) handleGetTask(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUserID(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	id := chi.URLParam(r, "id")

	task, err := h.taskClient.GetTask(r.Context(), userID, id)
	if err != nil {
		log.Printf("get task error: %v", err)
		if te, ok := err.(*client.TaskError); ok {
			writeError(w, te.StatusCode, te.Message)
			return
		}
		writeError(w, http.StatusBadGateway, fmt.Sprintf("get task failed: %s", err.Error()))
		return
	}

	writeJSON(w, http.StatusOK, TaskResponse(*task))
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
func (h *Handler) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUserID(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	id := chi.URLParam(r, "id")

	var req UpdateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad request")
		return
	}

	task, err := h.taskClient.UpdateTask(r.Context(), userID, id, client.UpdateTaskInput(req))
	if err != nil {
		log.Printf("update task error: %v", err)
		if te, ok := err.(*client.TaskError); ok {
			writeError(w, te.StatusCode, te.Message)
			return
		}
		writeError(w, http.StatusBadGateway, fmt.Sprintf("update task failed: %s", err.Error()))
		return
	}

	writeJSON(w, http.StatusOK, TaskResponse(*task))
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
func (h *Handler) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	userID, err := parseUserID(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	id := chi.URLParam(r, "id")

	if err := h.taskClient.DeleteTask(r.Context(), userID, id); err != nil {
		log.Printf("delete task error: %v", err)
		if te, ok := err.(*client.TaskError); ok {
			writeError(w, te.StatusCode, te.Message)
			return
		}
		writeError(w, http.StatusBadGateway, fmt.Sprintf("delete task failed: %s", err.Error()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
