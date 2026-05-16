package middleware

import (
	"time"
	"todo-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// RateLimiter khởi tạo middleware giới hạn số lượng request
func RateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        5,                // Tối đa 5 requests
		Expiration: 10 * time.Second, // Trong vòng 10 giây (sửa lại theo ý muốn)

		// Hàm này sẽ chạy khi user vượt quá giới hạn
		LimitReached: func(c *fiber.Ctx) error {
			return response.Error(c, fiber.StatusTooManyRequests, "Too many requests. Please try again later.")
		},

		// Mặc định Fiber sẽ chặn theo địa chỉ IP (c.IP())
		// Nếu bạn muốn cấu hình phức tạp hơn (VD: bỏ qua cho localhost), có thể dùng Next
	})
}
