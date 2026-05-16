package http

import (
	"strconv"
	"todo-app/delivery/http/middleware"
	"todo-app/domain"
	"todo-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TaskHandler struct {
	taskUsecase domain.TaskUsecase
}

func NewTaskHandler(app *fiber.App, uc domain.TaskUsecase, jwtSecret string) {
	handler := &TaskHandler{taskUsecase: uc}

	api := app.Group("/api/v1")
	protected := api.Group("", middleware.JWTAuth(jwtSecret))

	protected.Post("/tasks", handler.Create)
	protected.Get("/tasks", handler.GetAll)
	protected.Get("/tasks/:id", handler.GetByID)
	protected.Patch("/tasks/:id", handler.Update)
	protected.Delete("/tasks/:id", handler.Delete)
}

func (h *TaskHandler) Create(c *fiber.Ctx) error {
	var req domain.CreateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userID := middleware.GetUserID(c)
	task, err := h.taskUsecase.Create(c.Context(), userID, &req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(task)
}

func (h *TaskHandler) GetAll(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	filter := domain.TaskFilter{
		UserID:  userID,
		Search:  c.Query("search"),
		SortBy:  c.Query("sort_by"),
		SortDir: c.Query("sort_dir"),
	}

	if s := c.Query("status"); s != "" {
		status := domain.TodoStatus(s)
		filter.Status = &status
	}
	if p := c.Query("priority"); p != "" {
		priority := domain.TodoPriority(p)
		filter.Priority = &priority
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	filter.Page = page
	filter.Limit = limit

	tasks, total, err := h.taskUsecase.GetAll(c.Context(), filter)
	if err != nil {
		return err
	}
	return response.Paginated(c, tasks, total, page, limit)
}

func (h *TaskHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}
	userID := middleware.GetUserID(c)
	task, err := h.taskUsecase.GetByID(c.Context(), id, userID)
	if err != nil {
		return err
	}
	if task.UserID != userID {
		return response.Error(c, fiber.StatusForbidden, "You do not have permission to access this task")
	}

	return response.OK(c, task)
}

func (h *TaskHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var req domain.UpdateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	userID := middleware.GetUserID(c)

	// Truyền thẳng userID xuống usecase
	task, err := h.taskUsecase.Edit(c.Context(), id, userID, &req)
	if err != nil {
		return err // Báo lỗi 404/500 trực tiếp từ tầng dưới
	}

	return response.OK(c, task)
}

func (h *TaskHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	// 1. Kiểm tra quyền sở hữu TRƯỚC khi xoá
	userID := middleware.GetUserID(c)
	existingTask, err := h.taskUsecase.GetByID(c.Context(), id, userID)
	if err != nil {
		return err
	}

	if existingTask.UserID != userID {
		return response.Error(c, fiber.StatusForbidden, "You do not have permission to delete this task")
	}

	// 2. Thực hiện xoá nếu đúng chủ sở hữu
	if err := h.taskUsecase.Remove(c.Context(), id, userID); err != nil {
		return response.Error(c, fiber.StatusNotFound, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
