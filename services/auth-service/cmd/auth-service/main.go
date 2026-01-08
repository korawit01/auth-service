package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/korawit01/auth-service/services/auth-service/internal/config"
	myhttp "github.com/korawit01/auth-service/services/auth-service/internal/http"
	"github.com/korawit01/auth-service/services/auth-service/internal/repository"
	"github.com/korawit01/auth-service/services/auth-service/internal/service"
	"github.com/korawit01/auth-service/services/auth-service/internal/token"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.SetPrefix("auth-service: ")
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	log.Printf("connecting to database (%s)", describeDSN(cfg.DB.URL))
	db, err := sql.Open("pgx", cfg.DB.URL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("database ping failed: %v", err)
	}

	userRepo := repository.NewUserRepo(db)
	tokenProv := token.NewJWTProvider(cfg.JWT.Secret)
	userSvc := service.NewUserService(userRepo, tokenProv)
	h := myhttp.NewHandler(userSvc)

	// Ensure schema exists (no-op if already created)
	if err := ensureSchema(context.Background(), db); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()
	h.RegisterRoutes(app.Group("/auth"))

	// graceful shutdown
	go func() {
		log.Printf("auth-service listening on %s", cfg.HTTP.Addr)
		if err := app.Listen(cfg.HTTP.Addr); err != nil && !isServerClosed(err) {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	if err := app.ShutdownWithContext(ctxShutdown); err != nil {
		log.Printf("auth-service shutdown error: %v", err)
	}
}

func ensureSchema(ctx context.Context, db *sql.DB) error {
	const createUsers = `
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TABLE IF NOT EXISTS public.users (
    row_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);`
	_, err := db.ExecContext(ctx, createUsers)
	return err
}

func describeDSN(dsn string) string {
	u, err := url.Parse(dsn)
	if err != nil {
		return "invalid dsn"
	}
	user := ""
	if u.User != nil {
		user = u.User.Username()
	}
	return "scheme=" + u.Scheme + " host=" + u.Host + " db=" + strings.TrimPrefix(u.Path, "/") + " user=" + user
}

func isServerClosed(err error) bool {
	return errors.Is(err, net.ErrClosed) || strings.Contains(err.Error(), "Server closed")
}
