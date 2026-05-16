package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"
	"todo-app/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskModel struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null;index"`
	Title       string     `gorm:"type:varchar(255);not null"`
	Description string     `gorm:"type:text"`
	Status      string     `gorm:"type:varchar(20);not null;default:'pending';index"`
	Priority    string     `gorm:"type:varchar(10);not null;default:'medium';index"`
	DueDate     *time.Time `gorm:"index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"index"`
}

func (TaskModel) TableName() string { return "tasks" }

func toEntity(t *TaskModel) *domain.Task {
	return &domain.Task{
		ID:          t.ID,
		UserID:      t.UserID,
		Title:       t.Title,
		Description: t.Description,
		Status:      domain.TodoStatus(t.Status),
		Priority:    domain.TodoPriority(t.Priority),
		DueDate:     t.DueDate,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
		DeletedAt:   t.DeletedAt,
	}
}

func toModel(t *domain.Task) *TaskModel {
	return &TaskModel{
		ID:          t.ID,
		UserID:      t.UserID,
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		Priority:    string(t.Priority),
		DueDate:     t.DueDate,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
		DeletedAt:   t.DeletedAt,
	}
}

type taskPostgresRepo struct {
	db *gorm.DB
}

func NewTaskPostgresRepo(db *gorm.DB) domain.TaskRepository {
	db.AutoMigrate(&TaskModel{})

	return &taskPostgresRepo{db: db}
}

func (r *taskPostgresRepo) Create(ctx context.Context, task *domain.Task) error {
	model := toModel(task)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	task.ID = model.ID
	task.CreatedAt = model.CreatedAt
	task.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *taskPostgresRepo) FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Task, error) {
	var model TaskModel
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ? AND deleted_at IS NULL", id, userID).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("task not found")
		}
		return nil, err
	}
	return toEntity(&model), nil
}

func (r *taskPostgresRepo) FindAll(ctx context.Context, filter domain.TaskFilter) ([]*domain.Task, int64, error) {
	query := r.db.WithContext(ctx).Model(&TaskModel{}).
		Where("deleted_at IS NULL").
		Where("user_id = ?", filter.UserID)

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Priority != nil {
		query = query.Where("priority = ?", *filter.Priority)
	}
	if filter.Search != "" {
		query = query.Where("title ILIKE ?", fmt.Sprintf("%%%s%%", filter.Search))
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sorting
	sortDir := "DESC"
	if filter.SortDir == "asc" {
		sortDir = "ASC"
	}
	switch filter.SortBy {
	case "due_date":
		query = query.Order(fmt.Sprintf("due_date %s", sortDir))
	case "priority":
		// Explicit ordering: high(3) > medium(2) > low(1)
		if sortDir == "DESC" {
			query = query.Order("CASE priority WHEN 'high' THEN 3 WHEN 'medium' THEN 2 ELSE 1 END DESC")
		} else {
			query = query.Order("CASE priority WHEN 'high' THEN 3 WHEN 'medium' THEN 2 ELSE 1 END ASC")
		}
	default:
		query = query.Order(fmt.Sprintf("created_at %s", sortDir))
	}

	// Pagination
	offset := (filter.Page - 1) * filter.Limit
	query = query.Offset(offset).Limit(filter.Limit)

	var models []TaskModel
	if err := query.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	tasks := make([]*domain.Task, len(models))
	for i, m := range models {
		m := m
		tasks[i] = toEntity(&m)
	}
	return tasks, total, nil
}

func (r *taskPostgresRepo) Update(ctx context.Context, task *domain.Task) error {
	model := toModel(task)
	result := r.db.WithContext(ctx).
		Model(&TaskModel{}).
		Where("id = ? AND user_id = ? AND deleted_at IS NULL", task.ID, task.UserID).
		Updates(model)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func (r *taskPostgresRepo) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Model(&TaskModel{}).
		Where("id = ? AND user_id = ? AND deleted_at IS NULL", id, userID).
		Update("deleted_at", time.Now())
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}
