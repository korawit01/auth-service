package http

import (
	"errors"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrMissingAuthorization = errors.New("missing Authorization header")
	ErrInvalidAuthorization = errors.New("invalid Authorization header")
)

func (h *Handler) requireAuth(c *fiber.Ctx) error {
	token, err := readBearerToken(c)
	if err != nil {
		return writeError(c, fiber.StatusUnauthorized, err.Error())
	}

	claims, err := h.authClient.VerifyToken(c.UserContext(), token)
	if err != nil {
		log.Printf("verify token error: %v", err)
		return writeError(c, fiber.StatusUnauthorized, "unauthorized")
	}

	if strings.TrimSpace(claims.UserID) == "" {
		return writeError(c, fiber.StatusUnauthorized, "unauthorized")
	}

	c.Locals("userID", claims.UserID)
	return c.Next()
}

func readBearerToken(c *fiber.Ctx) (string, error) {
	authHeader := strings.TrimSpace(c.Get("Authorization"))
	if authHeader == "" {
		return "", ErrMissingAuthorization
	}

	lower := strings.ToLower(authHeader)
	const prefix = "bearer "
	if !strings.HasPrefix(lower, prefix) {
		return "", ErrInvalidAuthorization
	}

	token := strings.TrimSpace(authHeader[len(prefix):])
	if token == "" {
		return "", ErrInvalidAuthorization
	}
	return token, nil
}
