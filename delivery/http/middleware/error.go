package middleware

import (
	apperrors "todo-app/pkg/errors"
	"todo-app/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler is a Fiber error handler that maps AppError to HTTP responses
func ErrorHandler(c *fiber.Ctx, err error) error {
	if appErr, ok := err.(*apperrors.AppError); ok {
		return response.Error(c, appErr.Code, appErr.Message)
	}
	if fiberErr, ok := err.(*fiber.Error); ok {
		return response.Error(c, fiberErr.Code, fiberErr.Message)
	}
	return response.Error(c, fiber.StatusInternalServerError, "internal server error")
}

// Recovery is a panic recovery middleware
func Recovery() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fiber.NewError(fiber.StatusInternalServerError, "internal server error")
			}
		}()
		return c.Next()
	}
}
