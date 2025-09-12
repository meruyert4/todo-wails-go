package models

import (
	"time"
)

// Priority represents task priority levels
type Priority int

const (
	PriorityLow Priority = iota
	PriorityMedium
	PriorityHigh
)

// Status represents task completion status
type Status int

const (
	StatusActive Status = iota
	StatusCompleted
)

// Task represents a todo task
type Task struct {
	ID          string     `json:"id" db:"id"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description" db:"description"`
	Priority    Priority   `json:"priority" db:"priority"`
	Status      Status     `json:"status" db:"status"`
	DueDate     *time.Time `json:"dueDate,omitempty" db:"due_date"`
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time  `json:"updatedAt" db:"updated_at"`
}

// CreateTaskRequest represents request to create a new task
type CreateTaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    Priority   `json:"priority"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
}

// UpdateTaskRequest represents request to update a task
type UpdateTaskRequest struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    Priority   `json:"priority"`
	Status      Status     `json:"status"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
}

// FilterOptions represents filtering and sorting options
type FilterOptions struct {
	Status    *Status    `json:"status,omitempty"`
	Priority  *Priority  `json:"priority,omitempty"`
	DateFrom  *time.Time `json:"dateFrom,omitempty"`
	DateTo    *time.Time `json:"dateTo,omitempty"`
	SortBy    string     `json:"sortBy"`    // "created_at", "due_date", "priority", "title"
	SortOrder string     `json:"sortOrder"` // "asc", "desc"
}
