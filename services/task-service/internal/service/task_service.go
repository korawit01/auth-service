package service

import (
	"context"
	"errors"
	"time"

	"github.com/korawit01/auth-service/services/task-service/internal/domain"
	"github.com/korawit01/auth-service/services/task-service/internal/repository"
)

var (
	ErrInvalidStatus = errors.New("invalid status")
)

type TaskService interface {
	ListTasks(ctx context.Context, userID int64) ([]*domain.Task, error)
	GetTask(ctx context.Context, userID, taskID int64) (*domain.Task, error)
	CreateTask(ctx context.Context, userID int64, input CreateTaskInput) (*domain.Task, error)
	UpdateTask(ctx context.Context, userID, taskID int64, input UpdateTaskInput) (*domain.Task, error)
	DeleteTask(ctx context.Context, userID, taskID int64) error
}

type taskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

// DTOs
type CreateTaskInput struct {
	Title       string
	Description *string
	Status      domain.TaskStatus
	DueDate     *time.Time
}

type UpdateTaskInput struct {
	Title       *string
	Description *string
	Status      *domain.TaskStatus
	DueDate     *time.Time
}

func validateStatus(s domain.TaskStatus) bool {
	switch s {
	case domain.StatusTodo, domain.StatusInProgress, domain.StatusDone:
		return true
	default:
		return false
	}
}

func (s *taskService) ListTasks(ctx context.Context, userID int64) ([]*domain.Task, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *taskService) GetTask(ctx context.Context, userID, taskID int64) (*domain.Task, error) {
	return s.repo.GetByID(ctx, userID, taskID)
}

func (s *taskService) CreateTask(ctx context.Context, userID int64, input CreateTaskInput) (*domain.Task, error) {
	status := input.Status
	if status == "" {
		status = domain.StatusTodo
	}
	if !validateStatus(status) {
		return nil, ErrInvalidStatus
	}

	t := &domain.Task{
		UserID:      userID,
		Title:       input.Title,
		Description: input.Description,
		Status:      status,
		DueDate:     input.DueDate,
	}
	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *taskService) UpdateTask(ctx context.Context, userID, taskID int64, input UpdateTaskInput) (*domain.Task, error) {
	t, err := s.repo.GetByID(ctx, userID, taskID)
	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		t.Title = *input.Title
	}
	if input.Description != nil {
		t.Description = input.Description
	}
	if input.Status != nil {
		if !validateStatus(*input.Status) {
			return nil, ErrInvalidStatus
		}
		t.Status = *input.Status
	}
	if input.DueDate != nil {
		t.DueDate = input.DueDate
	}

	if err := s.repo.Update(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *taskService) DeleteTask(ctx context.Context, userID, taskID int64) error {
	return s.repo.Delete(ctx, userID, taskID)
}
