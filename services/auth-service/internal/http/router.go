package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/korawit01/auth-service/services/auth-service/internal/service"
)

type Handler struct {
	userSvc service.UserService
}

func NewHandler(userSvc service.UserService) *Handler {
	return &Handler{userSvc: userSvc}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/register", h.handleRegister)
	r.Post("/login", h.handleLogin)

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
		log.Printf("register decode error: %v", err)
		http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
		return
	}

	user, err := h.userSvc.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		log.Printf("register service error: %v", err)
		http.Error(w, `{"error":"cannot create user"}`, http.StatusInternalServerError)
		return
	}

	// avoid sending hashed password back to clients
	user.Password = ""
	_ = json.NewEncoder(w).Encode(user)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("login decode error: %v", err)
		http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
		return
	}

	token, err := h.userSvc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		log.Printf("login service error: %v", err)
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{"token": token})
}
