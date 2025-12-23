package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/korawit01/auth-service/services/task-service/internal/http/handlers"
)

func RegisterRoutes(r fiber.Router, taskHandler *handlers.TaskHandler) {
	r.Use(requestid.New())
	r.Use(logger.New())
	r.Use(recover.New())

	taskHandler.RegisterRoutes(r)
}
