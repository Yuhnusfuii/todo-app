package response

import "github.com/gofiber/fiber/v2"

// Response is the standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// PaginatedResponse wraps a list with pagination metadata
type PaginatedResponse struct {
	Items      interface{} `json:"items"`
	TotalItems int64       `json:"total_items"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

// OK sends a successful response
func OK(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Data:    data,
	})
}

// Created sends a 201 response
func Created(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Success: true,
		Data:    data,
	})
}

// Error sends an error response
func Error(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(Response{
		Success: false,
		Error:   message,
	})
}

// Paginated sends a paginated list response
func Paginated(c *fiber.Ctx, items interface{}, total int64, page, limit int) error {
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Data: PaginatedResponse{
			Items:      items,
			TotalItems: total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
		},
	})
}
