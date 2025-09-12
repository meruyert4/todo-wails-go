package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"todo-wails-go/internal/domain/models"
	"todo-wails-go/internal/domain/ports"

	_ "github.com/lib/pq"
)

// PostgresRepository implements TaskRepository interface
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(connStr string) (ports.TaskRepository, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	repo := &PostgresRepository{db: db}

	// Create table if not exists
	if err := repo.createTable(); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return repo, nil
}

// createTable creates the tasks table if it doesn't exist
func (r *PostgresRepository) createTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id VARCHAR(36) PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		priority INTEGER NOT NULL DEFAULT 0,
		status INTEGER NOT NULL DEFAULT 0,
		due_date TIMESTAMP,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);
	
	CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
	CREATE INDEX IF NOT EXISTS idx_tasks_priority ON tasks(priority);
	CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at);
	CREATE INDEX IF NOT EXISTS idx_tasks_due_date ON tasks(due_date);
	`

	_, err := r.db.Exec(query)
	return err
}

// Create creates a new task
func (r *PostgresRepository) Create(ctx context.Context, task *models.Task) error {
	query := `
		INSERT INTO tasks (id, title, description, priority, status, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		task.ID, task.Title, task.Description, task.Priority, task.Status,
		task.DueDate, task.CreatedAt, task.UpdatedAt)

	return err
}

// GetByID retrieves a task by ID
func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*models.Task, error) {
	query := `
		SELECT id, title, description, priority, status, due_date, created_at, updated_at
		FROM tasks WHERE id = $1
	`

	task := &models.Task{}
	var dueDate sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID, &task.Title, &task.Description, &task.Priority, &task.Status,
		&dueDate, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, err
	}

	if dueDate.Valid {
		task.DueDate = &dueDate.Time
	}

	return task, nil
}

// GetAll retrieves all tasks with optional filtering and sorting
func (r *PostgresRepository) GetAll(ctx context.Context, filter *models.FilterOptions) ([]*models.Task, error) {
	query := "SELECT id, title, description, priority, status, due_date, created_at, updated_at FROM tasks"
	args := []interface{}{}
	argIndex := 1

	// Build WHERE clause
	whereClauses := []string{}

	if filter != nil {
		if filter.Status != nil {
			whereClauses = append(whereClauses, fmt.Sprintf("status = $%d", argIndex))
			args = append(args, *filter.Status)
			argIndex++
		}

		if filter.Priority != nil {
			whereClauses = append(whereClauses, fmt.Sprintf("priority = $%d", argIndex))
			args = append(args, *filter.Priority)
			argIndex++
		}

		if filter.DateFrom != nil {
			whereClauses = append(whereClauses, fmt.Sprintf("created_at >= $%d", argIndex))
			args = append(args, *filter.DateFrom)
			argIndex++
		}

		if filter.DateTo != nil {
			whereClauses = append(whereClauses, fmt.Sprintf("created_at <= $%d", argIndex))
			args = append(args, *filter.DateTo)
			argIndex++
		}
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Build ORDER BY clause
	if filter != nil && filter.SortBy != "" {
		orderBy := filter.SortBy
		if orderBy == "due_date" {
			orderBy = "due_date NULLS LAST"
		}
		query += " ORDER BY " + orderBy

		if filter.SortOrder == "desc" {
			query += " DESC"
		} else {
			query += " ASC"
		}
	} else {
		query += " ORDER BY created_at DESC"
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task := &models.Task{}
		var dueDate sql.NullTime

		err := rows.Scan(
			&task.ID, &task.Title, &task.Description, &task.Priority, &task.Status,
			&dueDate, &task.CreatedAt, &task.UpdatedAt)

		if err != nil {
			return nil, err
		}

		if dueDate.Valid {
			task.DueDate = &dueDate.Time
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Update updates an existing task
func (r *PostgresRepository) Update(ctx context.Context, task *models.Task) error {
	query := `
		UPDATE tasks 
		SET title = $2, description = $3, priority = $4, status = $5, due_date = $6, updated_at = $7
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query,
		task.ID, task.Title, task.Description, task.Priority, task.Status,
		task.DueDate, task.UpdatedAt)

	return err
}

// Delete deletes a task by ID
func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM tasks WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// Close closes the database connection
func (r *PostgresRepository) Close() error {
	return r.db.Close()
}
