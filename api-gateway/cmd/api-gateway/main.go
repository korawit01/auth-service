// @title Taskboard API Gateway
// @version 1.0
// @description API Gateway proxying auth-service and task-service.
// @BasePath /
// @schemes http
//
//go:generate go run github.com/swaggo/swag/cmd/swag@v1.16.6 init --parseInternal --parseDependency --dir .,../../internal/http,../../internal/client --generalInfo main.go --output ../../internal/http/docs
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/korawit01/auth-service/api-gateway/internal/client"
	myhttp "github.com/korawit01/auth-service/api-gateway/internal/http"
	"github.com/korawit01/auth-service/api-gateway/internal/http/docs"
)

func main() {
	addr := envOrDefault("GATEWAY_ADDR", ":8080")
	authURL := envOrDefault("AUTH_SERVICE_URL", "http://localhost:8081/auth")
	taskURL := envOrDefault("TASK_SERVICE_URL", "http://localhost:8082")

	authClient := client.NewAuthClient(authURL)
	taskClient := client.NewTaskClient(taskURL)
	handler := myhttp.NewHandler(authClient, taskClient)

	r := chi.NewRouter()
	r.Mount("/", handler.Routes())

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Host = swaggerHost(addr)

	go func() {
		log.Printf("api-gateway listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("shutting down api-gateway...")
	_ = srv.Shutdown(ctx)
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func swaggerHost(addr string) string {
	if host := os.Getenv("GATEWAY_SWAGGER_HOST"); host != "" {
		return host
	}

	trimmed := strings.TrimPrefix(addr, ":")
	if strings.Contains(addr, ":") && !strings.Contains(addr, "://") {
		if strings.HasPrefix(addr, ":") {
			return "localhost" + addr
		}
		return addr
	}

	if trimmed == "" {
		return "localhost:8080"
	}
	return "localhost:" + trimmed
}
