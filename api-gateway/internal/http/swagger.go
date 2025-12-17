package http

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/korawit01/auth-service/api-gateway/internal/http/docs"
)

// SwaggerRoutes serves Swagger UI for the generated spec.
func SwaggerRoutes() http.Handler {
	return httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	)
}
