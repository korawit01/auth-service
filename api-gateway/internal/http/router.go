package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/korawit01/auth-service/api-gateway/internal/client"
)

type Handler struct {
	authClient *client.AuthClient
	taskClient *client.TaskClient
}

func NewHandler(authClient *client.AuthClient) *Handler {
	return &Handler{authClient: authClient}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	r.Post("/register", h.handleRegister)
	r.Post("/login", h.handleLogin)
	r.Post("/tasks", h.handleCreateTask)

	// Swagger docs
	r.Mount("/swagger", http.StripPrefix("/swagger", SwaggerRoutes()))

	return r
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req registerRequest
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

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req loginRequest
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

	_ = json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *Handler) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
		return
	}

	task, err := h.taskClient.CreateTask(r.Context(), req.Title, req.Description)
	if err != nil {
		log.Printf("create task error: %v", err)
		http.Error(w, `{"error":"create task failed"}`, http.StatusBadGateway)
		return
	}

	_ = json.NewEncoder(w).Encode(task)
}
