package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"todo-wails-go/internal/adapter/db"
	"todo-wails-go/internal/adapter/handler"
	"todo-wails-go/internal/adapter/service"
	"todo-wails-go/internal/usecase"
)

// App struct
type App struct {
	ctx     context.Context
	handler *handler.TaskHandler
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize database connection
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		// Default connection string for local development
		connStr = "host=localhost port=5432 user=postgres password=postgres dbname=todo_app sslmode=disable"
	}

	// Create repository
	repo, err := db.NewPostgresRepository(connStr)
	if err != nil {
		log.Printf("Warning: Failed to connect to PostgreSQL: %v", err)
		log.Println("Using in-memory storage instead")
		// For development, we'll use a simple in-memory storage
		// In production, you should always use PostgreSQL
		repo = db.NewMemoryRepository()
	}
	defer repo.Close()

	// Create service
	taskService := service.NewTaskService(repo)

	// Create use case
	taskUseCase := usecase.NewTaskUseCase(taskService)

	// Create handler
	a.handler = handler.NewTaskHandler(taskUseCase)
}

// CreateTask creates a new task
func (a *App) CreateTask(reqJSON string) (string, error) {
	if a.handler == nil {
		return "", fmt.Errorf("database not initialized")
	}
	return a.handler.CreateTask(a.ctx, reqJSON)
}

// GetTask retrieves a task by ID
func (a *App) GetTask(id string) (string, error) {
	if a.handler == nil {
		return "", fmt.Errorf("database not initialized")
	}
	return a.handler.GetTask(a.ctx, id)
}

// GetTasks retrieves all tasks with optional filtering
func (a *App) GetTasks(filterJSON string) (string, error) {
	if a.handler == nil {
		return "", fmt.Errorf("database not initialized")
	}
	return a.handler.GetTasks(a.ctx, filterJSON)
}

// UpdateTask updates an existing task
func (a *App) UpdateTask(reqJSON string) (string, error) {
	if a.handler == nil {
		return "", fmt.Errorf("database not initialized")
	}
	return a.handler.UpdateTask(a.ctx, reqJSON)
}

// DeleteTask deletes a task by ID
func (a *App) DeleteTask(id string) error {
	if a.handler == nil {
		return fmt.Errorf("database not initialized")
	}
	return a.handler.DeleteTask(a.ctx, id)
}

// ToggleTaskStatus toggles the completion status of a task
func (a *App) ToggleTaskStatus(id string) (string, error) {
	if a.handler == nil {
		return "", fmt.Errorf("database not initialized")
	}
	return a.handler.ToggleTaskStatus(a.ctx, id)
}

// GetTasksByStatus retrieves tasks filtered by status
func (a *App) GetTasksByStatus(status int) (string, error) {
	if a.handler == nil {
		return "", fmt.Errorf("database not initialized")
	}
	return a.handler.GetTasksByStatus(a.ctx, status)
}

// GetTasksByPriority retrieves tasks filtered by priority
func (a *App) GetTasksByPriority(priority int) (string, error) {
	if a.handler == nil {
		return "", fmt.Errorf("database not initialized")
	}
	return a.handler.GetTasksByPriority(a.ctx, priority)
}

// GetOverdueTasks retrieves overdue tasks
func (a *App) GetOverdueTasks() (string, error) {
	if a.handler == nil {
		return "", fmt.Errorf("database not initialized")
	}
	return a.handler.GetOverdueTasks(a.ctx)
}
