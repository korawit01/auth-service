package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// TaskClient is a placeholder for communicating with the task service.
type TaskClient struct {
	baseURL    string
	httpClient *http.Client
}

// TaskError captures upstream status and message for clearer error reporting.
type TaskError struct {
	StatusCode int
	Message    string
}

func (e *TaskError) Error() string {
	return e.Message
}

type UpdateTaskInput struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Status      *string    `json:"status,omitempty"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
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

func NewTaskClient(baseURL string) *TaskClient {
	return &TaskClient{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *TaskClient) Health(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/health", nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("task service unhealthy: %s", resp.Status)
	}
	return nil
}

func (c *TaskClient) CreateTask(ctx context.Context, userID, title, description string) (*TaskResponse, error) {
	payload := map[string]string{
		"title":       title,
		"description": description,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/task", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", strings.TrimSpace(userID))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		msg := strings.TrimSpace(string(bodyBytes))
		if msg == "" {
			msg = "task service returned empty body"
		}
		return nil, &TaskError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("task service error (%s): %s", resp.Status, msg),
		}
	}

	var task TaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}

	return &task, nil
}

func (c *TaskClient) GetTasks(ctx context.Context, userID string) (*[]TaskResponse, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/tasks", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", strings.TrimSpace(userID))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		msg := strings.TrimSpace(string(bodyBytes))
		if msg == "" {
			msg = "task service returned empty body"
		}
		return nil, &TaskError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("task service error (%s): %s", resp.Status, msg),
		}
	}

	var task []TaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}

	return &task, nil
}

func (c *TaskClient) GetTask(ctx context.Context, userID, taskId string) (*TaskResponse, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/task/"+taskId, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", strings.TrimSpace(userID))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		msg := strings.TrimSpace(string(bodyBytes))
		if msg == "" {
			msg = "task service returned empty body"
		}
		return nil, &TaskError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("task service error (%s): %s", resp.Status, msg),
		}
	}

	var task TaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}

	return &task, nil
}

func (c *TaskClient) UpdateTask(ctx context.Context, userID, taskId string, req UpdateTaskInput) (*TaskResponse, error) {
	payload := make(map[string]interface{})
	if req.Title != nil {
		payload["title"] = *req.Title
	}
	if req.Description != nil {
		payload["description"] = *req.Description
	}
	if req.Status != nil {
		payload["status"] = string(*req.Status)
	}
	if req.DueDate != nil {
		payload["dueDate"] = req.DueDate.Format(time.RFC3339)
	}

	body, _ := json.Marshal(payload)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, c.baseURL+"/task/"+taskId, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-User-ID", strings.TrimSpace(userID))
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		msg := strings.TrimSpace(string(bodyBytes))
		if msg == "" {
			msg = "task service returned empty body"
		}
		return nil, &TaskError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("task service error (%s): %s", resp.Status, msg),
		}
	}

	var task TaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}

	return &task, nil
}

func (c *TaskClient) DeleteTask(ctx context.Context, userID, taskId string) (*TaskResponse, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.baseURL+"/task/"+taskId, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", strings.TrimSpace(userID))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		msg := strings.TrimSpace(string(bodyBytes))
		if msg == "" {
			msg = "task service returned empty body"
		}
		return nil, &TaskError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("task service error (%s): %s", resp.Status, msg),
		}
	}

	var task TaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}

	return &task, nil
}
