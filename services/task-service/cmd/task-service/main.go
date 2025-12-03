package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/korawit01/auth-service/services/task-service/internal/config"
	myhttp "github.com/korawit01/auth-service/services/task-service/internal/http"
	"github.com/korawit01/auth-service/services/task-service/internal/http/handlers"
	"github.com/korawit01/auth-service/services/task-service/internal/repository"
	"github.com/korawit01/auth-service/services/task-service/internal/service"
)

func main() {
	cfg := config.FromEnv()
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	taskRepo := repository.NewTaskRepo(db)
	taskSvc := service.NewTaskService(taskRepo)
	taskHandler := handlers.NewTaskHandler(taskSvc)

	r := chi.NewRouter()
	r.Mount("/api", myhttp.NewRouter(taskHandler))

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: r,
	}

	go func() {
		log.Println("task-service listening on", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down task-service...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
