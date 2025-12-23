package handlers

import (
	"errors"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrMissingHeader = errors.New("missing X-User-ID")
	ErrUnauthorized  = errors.New("unauthorized")
)

func parseUserID(c *fiber.Ctx) (string, error) {
	id := strings.TrimSpace(c.Get("X-User-ID"))
	if id == "" {
		return "", ErrUnauthorized
	}
	return id, nil
}

func writeJSON(c *fiber.Ctx, status int, v interface{}) error {
	return c.Status(status).JSON(v)
}

func writeError(c *fiber.Ctx, status int, msg string) error {
	return writeJSON(c, status, map[string]string{"error": msg})
}

func handleError(c *fiber.Ctx, err error) error {
	if errors.Is(err, ErrUnauthorized) {
		return writeError(c, fiber.StatusUnauthorized, err.Error())
	}
	log.Printf("handler error: %v", err)
	return writeError(c, fiber.StatusInternalServerError, err.Error())
}
