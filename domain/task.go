package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TodoStatus string
type TodoPriority string

const (
	StatusPending    TodoStatus = "pending"
	StatusInProgress TodoStatus = "in_progress"
	StatusCompleted  TodoStatus = "completed"

	PriorityLow    TodoPriority = "low"
	PriorityMedium TodoPriority = "medium"
	PriorityHigh   TodoPriority = "high"
)

type Task struct {
	ID          uuid.UUID    `json:"id"`
	UserID      uuid.UUID    `json:"user_id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      TodoStatus   `json:"status"`
	Priority    TodoPriority `json:"priority"`
	DueDate     *time.Time   `json:"due_date,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletedAt   *time.Time   `json:"deleted_at,omitempty"`
}

type TaskRepository interface {
	Create(ctx context.Context, task *Task) error
	FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Task, error)
	FindAll(ctx context.Context, filter TaskFilter) ([]*Task, int64, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type TaskUsecase interface {
	Create(ctx context.Context, userID uuid.UUID, req *CreateTaskRequest) (*Task, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Task, error)
	GetAll(ctx context.Context, filter TaskFilter) ([]*Task, int64, error)
	Edit(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *UpdateTaskRequest) (*Task, error)
	Remove(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type TaskFilter struct {
	UserID   uuid.UUID
	Status   *TodoStatus
	Priority *TodoPriority
	Search   string
	SortBy   string
	SortDir  string
	Page     int
	Limit    int
}

type CreateTaskRequest struct {
	Title       string       `json:"title" validate:"required,min=1,max=255"`
	Description string       `json:"description" validate:"max=1000"`
	Priority    TodoPriority `json:"priority" validate:"required,oneof=low medium high"`
	DueDate     *time.Time   `json:"due_date,omitempty"`
}

func (r *CreateTaskRequest) Validate() error {
	if r.Title == "" {
		return fmt.Errorf("title is required")
	}
	if len(r.Title) > 255 {
		return fmt.Errorf("title cannot exceed 255 characters")
	}
	if len(r.Description) > 1000 {
		return fmt.Errorf("description cannot exceed 1000 characters")
	}
	if r.Priority != PriorityLow && r.Priority != PriorityMedium && r.Priority != PriorityHigh {
		return fmt.Errorf("priority must be one of: low, medium, high")
	}
	return nil
}

type UpdateTaskRequest struct {
	Title       *string       `json:"title" validate:"omitempty,min=1,max=255"`
	Description *string       `json:"description" validate:"omitempty,max=1000"`
	Priority    *TodoPriority `json:"priority" validate:"omitempty,oneof=low medium high"`
	DueDate     *time.Time    `json:"due_date,omitempty"`
}
