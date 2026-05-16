package middleware

import (
	"strings"

	jwtpkg "todo-app/pkg/jwt"
	"todo-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const userIDKey = "userID"

// JWTAuth returns a Fiber middleware that validates JWT tokens
func JWTAuth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return response.Error(c, fiber.StatusUnauthorized, "missing authorization header")
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return response.Error(c, fiber.StatusUnauthorized, "invalid authorization header format")
		}

		claims, err := jwtpkg.ParseToken(parts[1], secret)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, "invalid or expired token")
		}
		if claims.TokenType != jwtpkg.TokenTypeAccess {
			return response.Error(c, fiber.StatusUnauthorized, "invalid token type")
		}

		c.Locals(userIDKey, claims.UserID)
		return c.Next()
	}
}

// GetUserID extracts the user UUID from Fiber context locals
func GetUserID(c *fiber.Ctx) uuid.UUID {
	id, _ := c.Locals(userIDKey).(uuid.UUID)
	return id
}
