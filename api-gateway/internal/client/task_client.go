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

func NewTaskClient(baseURL string) *TaskClient {
	return &TaskClient{
		baseURL:    baseURL,
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

func (c *TaskClient) CreateTask(ctx context.Context, title, description string) (string, error) {
	payload := map[string]string{
		"title":       title,
		"description": description,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/tasks", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("create task failed: %s: %s", resp.Status, strings.TrimSpace(string(bodyBytes)))
	}	
	// Implementation for creating a task would go here.
	return "", nil
}

func (c *TaskClient) GetTask(ctx context.Context, id int64) (string, error) {
	// Implementation for getting a task would go here.
	return "", nil
}
