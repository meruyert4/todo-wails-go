package service

import (
	"context"
	"fmt"
	"time"

	"todo-wails-go/internal/domain/models"
	"todo-wails-go/internal/domain/ports"

	"github.com/google/uuid"
)

// TaskService implements the task business logic
type TaskService struct {
	repo ports.TaskRepository
}

// NewTaskService creates a new task service
func NewTaskService(repo ports.TaskRepository) ports.TaskService {
	return &TaskService{repo: repo}
}

// CreateTask creates a new task
func (s *TaskService) CreateTask(ctx context.Context, req *models.CreateTaskRequest) (*models.Task, error) {
	// Validate input
	if req.Title == "" {
		return nil, fmt.Errorf("title is required")
	}

	// Generate ID and timestamps
	now := time.Now()
	task := &models.Task{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      models.StatusActive,
		DueDate:     req.DueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Save to repository
	if err := s.repo.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}

// GetTask retrieves a task by ID
func (s *TaskService) GetTask(ctx context.Context, id string) (*models.Task, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	return s.repo.GetByID(ctx, id)
}

// GetTasks retrieves all tasks with optional filtering
func (s *TaskService) GetTasks(ctx context.Context, filter *models.FilterOptions) ([]*models.Task, error) {
	return s.repo.GetAll(ctx, filter)
}

// UpdateTask updates an existing task
func (s *TaskService) UpdateTask(ctx context.Context, req *models.UpdateTaskRequest) (*models.Task, error) {
	// Validate input
	if req.ID == "" {
		return nil, fmt.Errorf("id is required")
	}

	if req.Title == "" {
		return nil, fmt.Errorf("title is required")
	}

	// Get existing task
	task, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	// Update fields
	task.Title = req.Title
	task.Description = req.Description
	task.Priority = req.Priority
	task.Status = req.Status
	task.DueDate = req.DueDate
	task.UpdatedAt = time.Now()

	// Save changes
	if err := s.repo.Update(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return task, nil
}

// DeleteTask deletes a task by ID
func (s *TaskService) DeleteTask(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}

	// Check if task exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	return s.repo.Delete(ctx, id)
}

// ToggleTaskStatus toggles the completion status of a task
func (s *TaskService) ToggleTaskStatus(ctx context.Context, id string) (*models.Task, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	// Get existing task
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	// Toggle status
	if task.Status == models.StatusActive {
		task.Status = models.StatusCompleted
	} else {
		task.Status = models.StatusActive
	}

	task.UpdatedAt = time.Now()

	// Save changes
	if err := s.repo.Update(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return task, nil
}
