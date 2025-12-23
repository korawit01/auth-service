package http

import "github.com/gofiber/fiber/v2"

func writeJSON(c *fiber.Ctx, status int, v interface{}) error {
	return c.Status(status).JSON(v)
}

func writeError(c *fiber.Ctx, status int, msg string) error {
	return writeJSON(c, status, map[string]string{"error": msg})
}
