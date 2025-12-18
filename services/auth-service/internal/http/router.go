package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

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

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

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
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("register decode error: %v", err)
		writeError(w, http.StatusBadRequest, "bad request")
		return
	}

	user, err := h.userSvc.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		log.Printf("register service error: %v", err)
		writeError(w, http.StatusInternalServerError, "cannot create user")
		return
	}

	// avoid sending hashed password back to clients
	user.Password = ""
	writeJSON(w, http.StatusCreated, user)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("login decode error: %v", err)
		writeError(w, http.StatusBadRequest, "bad request")
		return
	}

	token, err := h.userSvc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		log.Printf("login service error: %v", err)
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}
