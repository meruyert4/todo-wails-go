package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"todo-wails-go/internal/domain/models"
	"todo-wails-go/internal/usecase"
)

// TaskHandler handles task-related HTTP requests
type TaskHandler struct {
	useCase *usecase.TaskUseCase
}

// NewTaskHandler creates a new task handler
func NewTaskHandler(useCase *usecase.TaskUseCase) *TaskHandler {
	return &TaskHandler{useCase: useCase}
}

// CreateTask creates a new task
func (h *TaskHandler) CreateTask(ctx context.Context, reqJSON string) (string, error) {
	var req models.CreateTaskRequest
	if err := json.Unmarshal([]byte(reqJSON), &req); err != nil {
		return "", fmt.Errorf("invalid request format: %w", err)
	}

	task, err := h.useCase.CreateTask(ctx, &req)
	if err != nil {
		return "", err
	}

	result, err := json.Marshal(task)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}

	return string(result), nil
}

// GetTask retrieves a task by ID
func (h *TaskHandler) GetTask(ctx context.Context, id string) (string, error) {
	task, err := h.useCase.GetTask(ctx, id)
	if err != nil {
		return "", err
	}

	result, err := json.Marshal(task)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}

	return string(result), nil
}

// GetTasks retrieves all tasks with optional filtering
func (h *TaskHandler) GetTasks(ctx context.Context, filterJSON string) (string, error) {
	var filter *models.FilterOptions
	if filterJSON != "" {
		filter = &models.FilterOptions{}
		if err := json.Unmarshal([]byte(filterJSON), filter); err != nil {
			return "", fmt.Errorf("invalid filter format: %w", err)
		}
	}

	tasks, err := h.useCase.GetTasks(ctx, filter)
	if err != nil {
		return "", err
	}

	result, err := json.Marshal(tasks)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}

	return string(result), nil
}

// UpdateTask updates an existing task
func (h *TaskHandler) UpdateTask(ctx context.Context, reqJSON string) (string, error) {
	var req models.UpdateTaskRequest
	if err := json.Unmarshal([]byte(reqJSON), &req); err != nil {
		return "", fmt.Errorf("invalid request format: %w", err)
	}

	task, err := h.useCase.UpdateTask(ctx, &req)
	if err != nil {
		return "", err
	}

	result, err := json.Marshal(task)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}

	return string(result), nil
}

// DeleteTask deletes a task by ID
func (h *TaskHandler) DeleteTask(ctx context.Context, id string) error {
	return h.useCase.DeleteTask(ctx, id)
}

// ToggleTaskStatus toggles the completion status of a task
func (h *TaskHandler) ToggleTaskStatus(ctx context.Context, id string) (string, error) {
	task, err := h.useCase.ToggleTaskStatus(ctx, id)
	if err != nil {
		return "", err
	}

	result, err := json.Marshal(task)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}

	return string(result), nil
}

// GetTasksByStatus retrieves tasks filtered by status
func (h *TaskHandler) GetTasksByStatus(ctx context.Context, status int) (string, error) {
	tasks, err := h.useCase.GetTasksByStatus(ctx, models.Status(status))
	if err != nil {
		return "", err
	}

	result, err := json.Marshal(tasks)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}

	return string(result), nil
}

// GetTasksByPriority retrieves tasks filtered by priority
func (h *TaskHandler) GetTasksByPriority(ctx context.Context, priority int) (string, error) {
	tasks, err := h.useCase.GetTasksByPriority(ctx, models.Priority(priority))
	if err != nil {
		return "", err
	}

	result, err := json.Marshal(tasks)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}

	return string(result), nil
}

// GetOverdueTasks retrieves overdue tasks
func (h *TaskHandler) GetOverdueTasks(ctx context.Context) (string, error) {
	tasks, err := h.useCase.GetOverdueTasks(ctx)
	if err != nil {
		return "", err
	}

	result, err := json.Marshal(tasks)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}

	return string(result), nil
}
