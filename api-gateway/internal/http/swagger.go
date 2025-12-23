package http

import (
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	_ "github.com/korawit01/auth-service/api-gateway/internal/http/docs"
)

// SwaggerRoutes serves Swagger UI for the generated spec.
func SwaggerRoutes() fiber.Handler {
	return fiberSwagger.FiberWrapHandler(fiberSwagger.URL("/swagger/doc.json"))
}
