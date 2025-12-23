package http

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var ErrUnauthorized = errors.New("unauthorized")

func parseUserID(c *fiber.Ctx) (string, error) {
	if v := c.Locals("userID"); v != nil {
		if userID, ok := v.(string); ok {
			userID = strings.TrimSpace(userID)
			if userID != "" {
				return userID, nil
			}
		}
	}

	userID := strings.TrimSpace(c.Get("X-User-ID"))
	if userID == "" {
		return "", ErrUnauthorized
	}
	return userID, nil
}

func writeJSON(c *fiber.Ctx, status int, v interface{}) error {
	return c.Status(status).JSON(v)
}

func writeError(c *fiber.Ctx, status int, msg string) error {
	return writeJSON(c, status, map[string]string{"error": msg})
}
