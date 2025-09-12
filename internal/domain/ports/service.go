package ports

import (
	"context"
	"todo-wails-go/internal/domain/models"
)

// TaskService defines the interface for task business logic
type TaskService interface {
	CreateTask(ctx context.Context, req *models.CreateTaskRequest) (*models.Task, error)
	GetTask(ctx context.Context, id string) (*models.Task, error)
	GetTasks(ctx context.Context, filter *models.FilterOptions) ([]*models.Task, error)
	UpdateTask(ctx context.Context, req *models.UpdateTaskRequest) (*models.Task, error)
	DeleteTask(ctx context.Context, id string) error
	ToggleTaskStatus(ctx context.Context, id string) (*models.Task, error)
}
