package usecase

import (
	"context"
	"fmt"
	"time"
	"todo-app/domain"

	"github.com/google/uuid"
)

type taskUsecase struct {
	repo domain.TaskRepository
}

func NewTaskUsecase(repo domain.TaskRepository) domain.TaskUsecase {
	return &taskUsecase{repo: repo}
}

func (u *taskUsecase) Create(ctx context.Context, userID uuid.UUID, req *domain.CreateTaskRequest) (*domain.Task, error) {
	now := time.Now()
	task := &domain.Task{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Status:      domain.StatusPending,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := u.repo.Create(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (u *taskUsecase) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Task, error) {
	return u.repo.FindByID(ctx, id, userID)
}

func (u *taskUsecase) GetAll(ctx context.Context, filter domain.TaskFilter) ([]*domain.Task, int64, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	return u.repo.FindAll(ctx, filter)
}

func (u *taskUsecase) Edit(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *domain.UpdateTaskRequest) (*domain.Task, error) {
	// LỚP KHOÁ 1: Hàm FindByID này đã được bạn nâng cấp SQL để chỉ lấy task của user
	task, err := u.repo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err // Lỗi "task not found" sẽ văng ra ngay nếu không đúng chủ
	}

	// Cập nhật các trường dữ liệu
	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Priority != nil {
		task.Priority = *req.Priority
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}
	task.UpdatedAt = time.Now()

	// Tiến hành lưu vào DB
	if err := u.repo.Update(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (u *taskUsecase) Remove(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	task, err := u.repo.FindByID(ctx, id, userID)
	if err != nil {
		return err
	}
	if task.UserID != userID {
		return fmt.Errorf("you do not have permission to delete this task")
	}
	return u.repo.Delete(ctx, id, userID)
}
