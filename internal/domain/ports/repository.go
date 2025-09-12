package ports

import (
	"context"
	"todo-wails-go/internal/domain/models"
)

// TaskRepository defines the interface for task data operations
type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	GetByID(ctx context.Context, id string) (*models.Task, error)
	GetAll(ctx context.Context, filter *models.FilterOptions) ([]*models.Task, error)
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id string) error
	Close() error
}
