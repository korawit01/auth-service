package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/korawit01/auth-service/api-gateway/internal/client"
	myhttp "github.com/korawit01/auth-service/api-gateway/internal/http"
)

func main() {
	addr := envOrDefault("GATEWAY_ADDR", ":8080")
	authURL := envOrDefault("AUTH_SERVICE_URL", "http://localhost:8081/auth")

	authClient := client.NewAuthClient(authURL)
	handler := myhttp.NewHandler(authClient)

	r := chi.NewRouter()
	r.Mount("/", handler.Routes())

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

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
