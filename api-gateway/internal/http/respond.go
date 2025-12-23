package http

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var ErrMissingUserID = errors.New("missing X-User-ID")

func parseUserID(c *fiber.Ctx) (string, error) {
	userID := strings.TrimSpace(c.Get("X-User-ID"))
	if userID == "" {
		return "", ErrMissingUserID
	}
	return userID, nil
}

func writeJSON(c *fiber.Ctx, status int, v interface{}) error {
	return c.Status(status).JSON(v)
}

func writeError(c *fiber.Ctx, status int, msg string) error {
	return writeJSON(c, status, map[string]string{"error": msg})
}
