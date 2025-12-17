package repository

import (
	"context"
	"database/sql"

	"github.com/korawit01/auth-service/services/auth-service/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, u *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, u *domain.User) error {
	query := `
		INSERT INTO public.users (email, password_hash)
		VALUES ($1, $2)
		RETURNING row_id, created_at
	`
	return r.db.QueryRowContext(ctx, query, u.Email, u.Password).
		Scan(&u.RowId, &u.CreatedAt)
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	query := `
		SELECT row_id, email, password_hash, created_at
		FROM public.users
		WHERE email = $1
	`
	err := r.db.QueryRowContext(ctx, query, email).
		Scan(&u.RowId, &u.Email, &u.Password, &u.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &u, nil
}
