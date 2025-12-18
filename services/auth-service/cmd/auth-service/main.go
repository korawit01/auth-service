package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
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
	cfg := config.FromEnv()
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	log.Printf("connecting to database (%s)", describeDSN(cfg.DatabaseURL))
	db, err := sql.Open("pgx", cfg.DatabaseURL)
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
	tokenProv := token.NewJWTProvider(cfg.JWTSecret)
	userSvc := service.NewUserService(userRepo, tokenProv)
	h := myhttp.NewHandler(userSvc)

	// Ensure schema exists (no-op if already created)
	if err := ensureSchema(context.Background(), db); err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Mount("/auth", h.Routes())

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: r,
	}

	// graceful shutdown
	go func() {
		log.Printf("auth-service listening on %s", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	_ = srv.Shutdown(ctxShutdown)
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
