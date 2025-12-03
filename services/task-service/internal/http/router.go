package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/korawit01/auth-service/services/task-service/internal/http/handlers"
)

func NewRouter(taskHandler *handlers.TaskHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(api chi.Router) {
		taskHandler.RegisterRoutes(api)
	})

	return r
}
