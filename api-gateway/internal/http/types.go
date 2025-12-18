package http

import "time"

type RegisterRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"secret123"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"secret123"`
}

type CreateTaskRequest struct {
	Title       string `json:"title" example:"Write design doc"`
	Description string `json:"description" example:"Outline API contract and data model"`
}

type UpdateTaskInput struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Status      *string    `json:"status,omitempty"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
}
type TokenResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type TaskResponse struct {
	RowID       string     `json:"row_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	Status      string     `json:"status"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}
