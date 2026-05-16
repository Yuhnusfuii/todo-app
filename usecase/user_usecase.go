package usecase

import (
	"context"
	"errors"
	"time"

	"todo-app/config"
	"todo-app/domain"
	apperrors "todo-app/pkg/errors"
	jwtpkg "todo-app/pkg/jwt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo domain.UserRepository
	jwtCfg   config.JWTConfig
}

// NewUserUsecase creates a new UserUsecase
func NewUserUsecase(userRepo domain.UserRepository, jwtCfg config.JWTConfig) domain.UserUsecase {
	return &userUsecase{userRepo: userRepo, jwtCfg: jwtCfg}
}

func (u *userUsecase) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	existing, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, apperrors.ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}

	user := &domain.User{
		ID:           uuid.New(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	token, err := jwtpkg.GenerateToken(user.ID, user.Email, jwtpkg.TokenTypeAccess, u.jwtCfg.Secret, u.jwtCfg.ExpiresIn)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}
	refreshToken, err := jwtpkg.GenerateToken(user.ID, user.Email, jwtpkg.TokenTypeRefresh, u.jwtCfg.Secret, u.jwtCfg.RefreshExpiresIn)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}

	return &domain.AuthResponse{Token: token, RefreshToken: refreshToken, User: user}, nil
}

func (u *userUsecase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error) {
	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrInvalidPassword
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, apperrors.ErrInvalidPassword
	}

	token, err := jwtpkg.GenerateToken(user.ID, user.Email, jwtpkg.TokenTypeAccess, u.jwtCfg.Secret, u.jwtCfg.ExpiresIn)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}
	refreshToken, err := jwtpkg.GenerateToken(user.ID, user.Email, jwtpkg.TokenTypeRefresh, u.jwtCfg.Secret, u.jwtCfg.RefreshExpiresIn)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}

	return &domain.AuthResponse{Token: token, RefreshToken: refreshToken, User: user}, nil
}

func (u *userUsecase) RefreshToken(ctx context.Context, req *domain.RefreshTokenRequest) (*domain.AuthResponse, error) {
	claims, err := jwtpkg.ParseToken(req.RefreshToken, u.jwtCfg.Secret)
	if err != nil {
		return nil, apperrors.ErrUnauthorized
	}
	if claims.TokenType != jwtpkg.TokenTypeRefresh {
		return nil, apperrors.ErrUnauthorized
	}

	user, err := u.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrUnauthorized
		}
		return nil, err
	}

	accessToken, err := jwtpkg.GenerateToken(user.ID, user.Email, jwtpkg.TokenTypeAccess, u.jwtCfg.Secret, u.jwtCfg.ExpiresIn)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}
	newRefreshToken, err := jwtpkg.GenerateToken(user.ID, user.Email, jwtpkg.TokenTypeRefresh, u.jwtCfg.Secret, u.jwtCfg.RefreshExpiresIn)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}

	return &domain.AuthResponse{Token: accessToken, RefreshToken: newRefreshToken}, nil
}

func (u *userUsecase) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	return u.userRepo.FindByID(ctx, userID)
}
