package http

import (
	"todo-app/delivery/http/middleware"
	"todo-app/domain"
	apperrors "todo-app/pkg/errors"
	"todo-app/pkg/response"
	"todo-app/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// UserHandler handles user HTTP requests
type UserHandler struct {
	userUsecase domain.UserUsecase
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(app *fiber.App, userUsecase domain.UserUsecase) *UserHandler {
	handler := &UserHandler{userUsecase: userUsecase}

	app.Post("/api/v1/auth/register", handler.Register)
	app.Post("/api/v1/auth/login", handler.Login)
	app.Post("/api/v1/auth/refresh", handler.RefreshToken)
	app.Get("/api/v1/users/me", handler.GetProfile)

	return handler
}

// Register godoc
// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param body body domain.RegisterRequest true "Register request"
// @Success 201 {object} response.Response
// @Router /api/v1/auth/register [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req domain.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.ErrBadRequest
	}
	if errs := validator.Validate(&req); errs != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"errors":  errs,
		})
	}
	result, err := h.userUsecase.Register(c.Context(), &req)
	if err != nil {
		return err
	}
	return response.Created(c, result)
}

// Login godoc
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param body body domain.LoginRequest true "Login request"
// @Success 200 {object} response.Response
// @Router /api/v1/auth/login [post]
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.ErrBadRequest
	}
	if errs := validator.Validate(&req); errs != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"errors":  errs,
		})
	}
	result, err := h.userUsecase.Login(c.Context(), &req)
	if err != nil {
		return err
	}
	return response.OK(c, result)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Tags auth
// @Accept json
// @Produce json
// @Param body body domain.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} response.Response
// @Router /api/v1/auth/refresh [post]
func (h *UserHandler) RefreshToken(c *fiber.Ctx) error {
	var req domain.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return apperrors.ErrBadRequest
	}
	if errs := validator.Validate(&req); errs != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"errors":  errs,
		})
	}
	result, err := h.userUsecase.RefreshToken(c.Context(), &req)
	if err != nil {
		return err
	}
	return response.OK(c, result)
}

// GetProfile godoc
// @Summary Get current user profile
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/users/me [get]
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	user, err := h.userUsecase.GetProfile(c.Context(), userID)
	if err != nil {
		return err
	}
	return response.OK(c, user)
}
