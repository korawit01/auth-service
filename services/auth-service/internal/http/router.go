package http

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/korawit01/auth-service/services/auth-service/internal/service"
)

type Handler struct {
	userSvc service.UserService
}

func NewHandler(userSvc service.UserService) *Handler {
	return &Handler{userSvc: userSvc}
}

func (h *Handler) RegisterRoutes(r fiber.Router) {
	r.Use(requestid.New())
	r.Use(logger.New())
	r.Use(recover.New())

	r.Post("/register", h.handleRegister)
	r.Post("/login", h.handleLogin)
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) handleRegister(c *fiber.Ctx) error {
	var req registerRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("register decode error: %v", err)
		return writeError(c, fiber.StatusBadRequest, "bad request")
	}

	user, err := h.userSvc.Register(c.UserContext(), req.Email, req.Password)
	if err != nil {
		log.Printf("register service error: %v", err)
		return writeError(c, fiber.StatusInternalServerError, "cannot create user")
	}

	// avoid sending hashed password back to clients
	user.Password = ""
	return writeJSON(c, fiber.StatusCreated, user)
}

func (h *Handler) handleLogin(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("login decode error: %v", err)
		return writeError(c, fiber.StatusBadRequest, "bad request")
	}

	token, err := h.userSvc.Login(c.UserContext(), req.Email, req.Password)
	if err != nil {
		log.Printf("login service error: %v", err)
		return writeError(c, fiber.StatusUnauthorized, "invalid credentials")
	}

	return writeJSON(c, fiber.StatusOK, map[string]string{"token": token})
}
