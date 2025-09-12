package db

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"todo-wails-go/internal/domain/models"
	"todo-wails-go/internal/domain/ports"
)

// MemoryRepository implements TaskRepository interface using in-memory storage
type MemoryRepository struct {
	tasks map[string]*models.Task
	mutex sync.RWMutex
}

// NewMemoryRepository creates a new in-memory repository
func NewMemoryRepository() ports.TaskRepository {
	return &MemoryRepository{
		tasks: make(map[string]*models.Task),
	}
}

// Create creates a new task
func (r *MemoryRepository) Create(ctx context.Context, task *models.Task) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.tasks[task.ID] = task
	return nil
}

// GetByID retrieves a task by ID
func (r *MemoryRepository) GetByID(ctx context.Context, id string) (*models.Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, fmt.Errorf("task not found")
	}

	// Return a copy to avoid race conditions
	taskCopy := *task
	return &taskCopy, nil
}

// GetAll retrieves all tasks with optional filtering and sorting
func (r *MemoryRepository) GetAll(ctx context.Context, filter *models.FilterOptions) ([]*models.Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var tasks []*models.Task

	for _, task := range r.tasks {
		// Apply filters
		if filter != nil {
			if filter.Status != nil && task.Status != *filter.Status {
				continue
			}
			if filter.Priority != nil && task.Priority != *filter.Priority {
				continue
			}
			if filter.DateFrom != nil && task.CreatedAt.Before(*filter.DateFrom) {
				continue
			}
			if filter.DateTo != nil && task.CreatedAt.After(*filter.DateTo) {
				continue
			}
		}

		// Create a copy to avoid race conditions
		taskCopy := *task
		tasks = append(tasks, &taskCopy)
	}

	// Apply sorting
	if filter != nil && filter.SortBy != "" {
		sort.Slice(tasks, func(i, j int) bool {
			switch filter.SortBy {
			case "title":
				if filter.SortOrder == "desc" {
					return tasks[i].Title > tasks[j].Title
				}
				return tasks[i].Title < tasks[j].Title
			case "priority":
				if filter.SortOrder == "desc" {
					return tasks[i].Priority > tasks[j].Priority
				}
				return tasks[i].Priority < tasks[j].Priority
			case "due_date":
				if tasks[i].DueDate == nil && tasks[j].DueDate == nil {
					return false
				}
				if tasks[i].DueDate == nil {
					return false
				}
				if tasks[j].DueDate == nil {
					return true
				}
				if filter.SortOrder == "desc" {
					return tasks[i].DueDate.After(*tasks[j].DueDate)
				}
				return tasks[i].DueDate.Before(*tasks[j].DueDate)
			case "created_at":
				fallthrough
			default:
				if filter.SortOrder == "desc" {
					return tasks[i].CreatedAt.After(tasks[j].CreatedAt)
				}
				return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
			}
		})
	} else {
		// Default sort by created_at desc
		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].CreatedAt.After(tasks[j].CreatedAt)
		})
	}

	return tasks, nil
}

// Update updates an existing task
func (r *MemoryRepository) Update(ctx context.Context, task *models.Task) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.tasks[task.ID]; !exists {
		return fmt.Errorf("task not found")
	}

	r.tasks[task.ID] = task
	return nil
}

// Delete deletes a task by ID
func (r *MemoryRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.tasks[id]; !exists {
		return fmt.Errorf("task not found")
	}

	delete(r.tasks, id)
	return nil
}

// Close closes the repository (no-op for memory repository)
func (r *MemoryRepository) Close() error {
	return nil
}
