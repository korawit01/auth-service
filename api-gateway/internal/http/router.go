package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

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

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	r.Post("/register", h.handleRegister)
	r.Post("/login", h.handleLogin)

	r.Post("/task", h.handleCreateTask)
	r.Get("/tasks", h.handleGetTasks)
	r.Get("/task/{id}", h.handleGetTask)
	r.Put("/task/{id}", h.handleUpdateTask)
	r.Delete("/task/{id}", h.handleDeleteTask)

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
	w.Header().Set("Content-Type", "application/json")
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
		return
	}

	user, err := h.authClient.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		log.Printf("register error: %v", err)
		http.Error(w, `{"error":"register failed"}`, http.StatusBadGateway)
		return
	}

	_ = json.NewEncoder(w).Encode(user)
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
	w.Header().Set("Content-Type", "application/json")
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
		return
	}

	token, err := h.authClient.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		log.Printf("login error: %v", err)
		http.Error(w, `{"error":"login failed"}`, http.StatusBadGateway)
		return
	}

	_ = json.NewEncoder(w).Encode(TokenResponse{Token: token})
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
// @Router /task [post]
func (h *Handler) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.Header.Get("X-User-ID")
	if strings.TrimSpace(userID) == "" {
		http.Error(w, `{"error":"missing X-User-ID"}`, http.StatusUnauthorized)
		return
	}

	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
		return
	}

	task, err := h.taskClient.CreateTask(r.Context(), userID, req.Title, req.Description)
	if err != nil {
		log.Printf("create task error: %v", err)
		if te, ok := err.(*client.TaskError); ok {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, te.Message), te.StatusCode)
			return
		}
		http.Error(w, fmt.Sprintf(`{"error":"create task failed: %s"}`, err.Error()), http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(TaskResponse(*task))
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
	w.Header().Set("Content-Type", "application/json")
	userID := r.Header.Get("X-User-ID")
	if strings.TrimSpace(userID) == "" {
		http.Error(w, `{"error":"missing X-User-ID"}`, http.StatusUnauthorized)
		return
	}

	task, err := h.taskClient.GetTasks(r.Context(), userID)
	if err != nil {
		log.Printf("get tasks error: %v", err)
		if te, ok := err.(*client.TaskError); ok {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, te.Message), te.StatusCode)
			return
		}
		http.Error(w, fmt.Sprintf(`{"error":"get tasks failed: %s"}`, err.Error()), http.StatusBadGateway)
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

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

// handleGetTask godoc
// @Summary Get tasks
// @Description Get tasks belonging to the provided user and task id.
// @Tags tasks
// @Produce json
// @Param X-User-ID header string true "User ID"
// @Param id path string true "Task ID"
// @Success 200 {array} TaskResponse
// @Failure 401 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /task/{id} [get]
func (h *Handler) handleGetTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.Header.Get("X-User-ID")
	taskId := r.PathValue("id")
	if strings.TrimSpace(userID) == "" {
		http.Error(w, `{"error":"missing X-User-ID"}`, http.StatusUnauthorized)
		return
	}

	task, err := h.taskClient.GetTask(r.Context(), userID, taskId)
	if err != nil {
		log.Printf("get task error: %v", err)
		if te, ok := err.(*client.TaskError); ok {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, te.Message), te.StatusCode)
			return
		}
		http.Error(w, fmt.Sprintf(`{"error":"get task failed: %s"}`, err.Error()), http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(TaskResponse(*task))
}

// handleUpdateTask godoc
// @Summary Update tasks
// @Description Update tasks belonging to the provided user and task id.
// @Tags tasks
// @Produce json
// @Param X-User-ID header string true "User ID"
// @Param id path string true "Task ID"
// @Param request body UpdateTaskInput true "Task payload"
// @Success 200 {array} TaskResponse
// @Failure 401 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /task/{id} [put]
func (h *Handler) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.Header.Get("X-User-ID")
	taskId := r.PathValue("id")
	if strings.TrimSpace(userID) == "" {
		http.Error(w, `{"error":"missing X-User-ID"}`, http.StatusUnauthorized)
		return
	}

	var req UpdateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
		return
	}
	task, err := h.taskClient.UpdateTask(r.Context(), userID, taskId, client.UpdateTaskInput(req))
	if err != nil {
		log.Printf("get task error: %v", err)
		if te, ok := err.(*client.TaskError); ok {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, te.Message), te.StatusCode)
			return
		}
		http.Error(w, fmt.Sprintf(`{"error":"get task failed: %s"}`, err.Error()), http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(TaskResponse(*task))
}

// handleDeleteTask godoc
// @Summary Delete tasks
// @Description Delete tasks belonging to the provided user and task id.
// @Tags tasks
// @Produce json
// @Param X-User-ID header string true "User ID"
// @Param id path string true "Task ID"
// @Success 200 {array} TaskResponse
// @Failure 401 {object} ErrorResponse
// @Failure 502 {object} ErrorResponse
// @Router /task/{id} [delete]
func (h *Handler) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.Header.Get("X-User-ID")
	taskId := r.PathValue("id")
	if strings.TrimSpace(userID) == "" {
		http.Error(w, `{"error":"missing X-User-ID"}`, http.StatusUnauthorized)
		return
	}

	task, err := h.taskClient.DeleteTask(r.Context(), userID, taskId)
	if err != nil {
		log.Printf("delete task error: %v", err)
		if te, ok := err.(*client.TaskError); ok {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, te.Message), te.StatusCode)
			return
		}
		http.Error(w, fmt.Sprintf(`{"error":"delete task failed: %s"}`, err.Error()), http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(TaskResponse(*task))
}