package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/korawit01/auth-service/services/task-service/internal/domain"
)

var ErrNotFound = errors.New("task not found")

type TaskRepository interface {
	ListByUser(ctx context.Context, userID int64) ([]*domain.Task, error)
	GetByID(ctx context.Context, userID, taskID int64) (*domain.Task, error)
	Create(ctx context.Context, t *domain.Task) error
	Update(ctx context.Context, t *domain.Task) error
	Delete(ctx context.Context, userID, taskID int64) error
}

type taskRepo struct {
	db *sql.DB
}

func NewTaskRepo(db *sql.DB) TaskRepository {
	return &taskRepo{db: db}
}

func (r *taskRepo) ListByUser(ctx context.Context, userID int64) ([]*domain.Task, error) {
	query := `
		SELECT id, user_id, title, description, status, due_date, created_at, updated_at
		FROM public.tasks
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		var t domain.Task
		err := rows.Scan(
			&t.ID, &t.UserID, &t.Title, &t.Description,
			&t.Status, &t.DueDate, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	return tasks, rows.Err()
}

func (r *taskRepo) GetByID(ctx context.Context, userID, taskID int64) (*domain.Task, error) {
	query := `
		SELECT id, user_id, title, description, status, due_date, created_at, updated_at
		FROM public.tasks
		WHERE id = $1 AND user_id = $2
	`
	var t domain.Task
	err := r.db.QueryRowContext(ctx, query, taskID, userID).Scan(
		&t.ID, &t.UserID, &t.Title, &t.Description,
		&t.Status, &t.DueDate, &t.CreatedAt, &t.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *taskRepo) Create(ctx context.Context, t *domain.Task) error {
	query := `
		INSERT INTO public.tasks (user_id, title, description, status, due_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(
		ctx, query,
		t.UserID, t.Title, t.Description, t.Status, t.DueDate,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *taskRepo) Update(ctx context.Context, t *domain.Task) error {
	t.UpdatedAt = time.Now()
	query := `
		UPDATE public.tasks
		SET title = $1,
		    description = $2,
		    status = $3,
		    due_date = $4,
		    updated_at = $5
		WHERE id = $6 AND user_id = $7
	`
	res, err := r.db.ExecContext(
		ctx, query,
		t.Title, t.Description, t.Status, t.DueDate, t.UpdatedAt,
		t.ID, t.UserID,
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *taskRepo) Delete(ctx context.Context, userID, taskID int64) error {
	query := `DELETE FROM public.tasks WHERE id = $1 AND user_id = $2`
	res, err := r.db.ExecContext(ctx, query, taskID, userID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

