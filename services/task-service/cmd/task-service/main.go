package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/korawit01/auth-service/services/task-service/internal/config"
	myhttp "github.com/korawit01/auth-service/services/task-service/internal/http"
	"github.com/korawit01/auth-service/services/task-service/internal/http/handlers"
	"github.com/korawit01/auth-service/services/task-service/internal/repository"
	"github.com/korawit01/auth-service/services/task-service/internal/service"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.SetPrefix("task-service: ")
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

	app := fiber.New()
	myhttp.RegisterRoutes(app.Group("/api"), taskHandler)

	go func() {
		log.Println("task-service listening on", cfg.Addr)
		if err := app.Listen(cfg.Addr); err != nil && !isServerClosed(err) {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down task-service...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("task-service shutdown error: %v", err)
	}
}

func isServerClosed(err error) bool {
	return errors.Is(err, net.ErrClosed) || strings.Contains(err.Error(), "Server closed")
}
