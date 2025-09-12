package usecase

import (
	"context"
	"time"

	"todo-wails-go/internal/domain/models"
	"todo-wails-go/internal/domain/ports"
)

// TaskUseCase implements the application use cases
type TaskUseCase struct {
	service ports.TaskService
}

// NewTaskUseCase creates a new task use case
func NewTaskUseCase(service ports.TaskService) *TaskUseCase {
	return &TaskUseCase{service: service}
}

// CreateTask creates a new task
func (uc *TaskUseCase) CreateTask(ctx context.Context, req *models.CreateTaskRequest) (*models.Task, error) {
	return uc.service.CreateTask(ctx, req)
}

// GetTask retrieves a task by ID
func (uc *TaskUseCase) GetTask(ctx context.Context, id string) (*models.Task, error) {
	return uc.service.GetTask(ctx, id)
}

// GetTasks retrieves all tasks with optional filtering
func (uc *TaskUseCase) GetTasks(ctx context.Context, filter *models.FilterOptions) ([]*models.Task, error) {
	return uc.service.GetTasks(ctx, filter)
}

// UpdateTask updates an existing task
func (uc *TaskUseCase) UpdateTask(ctx context.Context, req *models.UpdateTaskRequest) (*models.Task, error) {
	return uc.service.UpdateTask(ctx, req)
}

// DeleteTask deletes a task by ID
func (uc *TaskUseCase) DeleteTask(ctx context.Context, id string) error {
	return uc.service.DeleteTask(ctx, id)
}

// ToggleTaskStatus toggles the completion status of a task
func (uc *TaskUseCase) ToggleTaskStatus(ctx context.Context, id string) (*models.Task, error) {
	return uc.service.ToggleTaskStatus(ctx, id)
}

// GetTasksByStatus retrieves tasks filtered by status
func (uc *TaskUseCase) GetTasksByStatus(ctx context.Context, status models.Status) ([]*models.Task, error) {
	filter := &models.FilterOptions{
		Status:    &status,
		SortBy:    "created_at",
		SortOrder: "desc",
	}
	return uc.service.GetTasks(ctx, filter)
}

// GetTasksByPriority retrieves tasks filtered by priority
func (uc *TaskUseCase) GetTasksByPriority(ctx context.Context, priority models.Priority) ([]*models.Task, error) {
	filter := &models.FilterOptions{
		Priority:  &priority,
		SortBy:    "created_at",
		SortOrder: "desc",
	}
	return uc.service.GetTasks(ctx, filter)
}

// GetTasksByDateRange retrieves tasks within a date range
func (uc *TaskUseCase) GetTasksByDateRange(ctx context.Context, from, to time.Time) ([]*models.Task, error) {
	filter := &models.FilterOptions{
		DateFrom:  &from,
		DateTo:    &to,
		SortBy:    "created_at",
		SortOrder: "desc",
	}
	return uc.service.GetTasks(ctx, filter)
}

// GetOverdueTasks retrieves overdue tasks
func (uc *TaskUseCase) GetOverdueTasks(ctx context.Context) ([]*models.Task, error) {
	now := time.Now()
	filter := &models.FilterOptions{
		Status:    &[]models.Status{models.StatusActive}[0],
		DateTo:    &now,
		SortBy:    "due_date",
		SortOrder: "asc",
	}
	return uc.service.GetTasks(ctx, filter)
}
