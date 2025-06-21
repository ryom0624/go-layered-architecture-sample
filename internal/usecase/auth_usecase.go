package usecase

import (
	"context"
	"errors"
	"layered-architecture-template/internal/domain/entity"
	"layered-architecture-template/internal/domain/repository"
	"layered-architecture-template/pkg/jwt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUsecase interface {
	Register(ctx context.Context, req *entity.RegisterRequest) (*entity.LoginResponse, error)
	Login(ctx context.Context, req *entity.LoginRequest) (*entity.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*entity.LoginResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, userID uint) error
}

type authUsecaseImpl struct {
	authRepo           repository.AuthRepository
	jwtManager         *jwt.JWTManager
	refreshTokenExpiry time.Duration
}

func NewAuthUsecase(authRepo repository.AuthRepository, jwtManager *jwt.JWTManager, refreshTokenExpiry time.Duration) AuthUsecase {
	return &authUsecaseImpl{
		authRepo:           authRepo,
		jwtManager:         jwtManager,
		refreshTokenExpiry: refreshTokenExpiry,
	}
}

func (u *authUsecaseImpl) Register(ctx context.Context, req *entity.RegisterRequest) (*entity.LoginResponse, error) {
	if err := u.validateRegisterRequest(req); err != nil {
		return nil, err
	}

	existingUser, err := u.authRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already exists")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	user := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := u.authRepo.Register(ctx, user); err != nil {
		return nil, err
	}

	accessToken, err := u.jwtManager.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	log.Println("accessToken", accessToken)

	// Generate refresh token with retry logic
	var refreshTokenStr string
	var refreshToken *entity.RefreshToken
	maxRetries := 3
	
	for i := 0; i < maxRetries; i++ {
		refreshTokenStr, err = jwt.GenerateRefreshToken()
		if err != nil {
			return nil, err
		}

		refreshToken = &entity.RefreshToken{
			UserID:    user.ID,
			Token:     refreshTokenStr,
			ExpiresAt: time.Now().Add(u.refreshTokenExpiry),
		}

		if err := u.authRepo.StoreRefreshToken(ctx, refreshToken); err != nil {
			// If it's a unique constraint violation, retry with a new token
			if i < maxRetries-1 {
				continue
			}
			return nil, err
		}
		break
	}

	user.Password = ""

	return &entity.LoginResponse{
		User:  *user,
		Token: accessToken,
	}, nil
}

func (u *authUsecaseImpl) Login(ctx context.Context, req *entity.LoginRequest) (*entity.LoginResponse, error) {
	if err := u.validateLoginRequest(req); err != nil {
		return nil, err
	}

	user, err := u.authRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	accessToken, err := u.jwtManager.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	// Generate refresh token with retry logic
	var refreshTokenStr string
	var refreshToken *entity.RefreshToken
	maxRetries := 3
	
	for i := 0; i < maxRetries; i++ {
		refreshTokenStr, err = jwt.GenerateRefreshToken()
		if err != nil {
			return nil, err
		}

		refreshToken = &entity.RefreshToken{
			UserID:    user.ID,
			Token:     refreshTokenStr,
			ExpiresAt: time.Now().Add(u.refreshTokenExpiry),
		}

		if err := u.authRepo.StoreRefreshToken(ctx, refreshToken); err != nil {
			// If it's a unique constraint violation, retry with a new token
			if i < maxRetries-1 {
				continue
			}
			return nil, err
		}
		break
	}

	user.Password = ""

	return &entity.LoginResponse{
		User:  *user,
		Token: accessToken,
	}, nil
}

func (u *authUsecaseImpl) RefreshToken(ctx context.Context, refreshToken string) (*entity.LoginResponse, error) {
	storedToken, err := u.authRepo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid refresh token")
		}
		return nil, err
	}

	if time.Now().After(storedToken.ExpiresAt) {
		u.authRepo.DeleteRefreshToken(ctx, refreshToken)
		return nil, errors.New("refresh token expired")
	}

	accessToken, err := u.jwtManager.GenerateToken(&storedToken.User)
	if err != nil {
		return nil, err
	}

	// Delete old token first to avoid constraint violations
	if err := u.authRepo.DeleteRefreshToken(ctx, refreshToken); err != nil {
		return nil, err
	}

	// Generate new refresh token with retry logic
	var newRefreshTokenStr string
	var newRefreshToken *entity.RefreshToken
	maxRetries := 3
	
	for i := 0; i < maxRetries; i++ {
		newRefreshTokenStr, err = jwt.GenerateRefreshToken()
		if err != nil {
			return nil, err
		}

		newRefreshToken = &entity.RefreshToken{
			UserID:    storedToken.UserID,
			Token:     newRefreshTokenStr,
			ExpiresAt: time.Now().Add(u.refreshTokenExpiry),
		}

		if err := u.authRepo.StoreRefreshToken(ctx, newRefreshToken); err != nil {
			// If it's a unique constraint violation, retry with a new token
			if i < maxRetries-1 {
				continue
			}
			return nil, err
		}
		break
	}

	storedToken.User.Password = ""

	return &entity.LoginResponse{
		User:  storedToken.User,
		Token: accessToken,
	}, nil
}

func (u *authUsecaseImpl) Logout(ctx context.Context, refreshToken string) error {
	return u.authRepo.DeleteRefreshToken(ctx, refreshToken)
}

func (u *authUsecaseImpl) LogoutAll(ctx context.Context, userID uint) error {
	return u.authRepo.DeleteUserRefreshTokens(ctx, userID)
}

func (u *authUsecaseImpl) validateRegisterRequest(req *entity.RegisterRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if len(req.Name) < 2 || len(req.Name) > 100 {
		return errors.New("name must be between 2 and 100 characters")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	return nil
}

func (u *authUsecaseImpl) validateLoginRequest(req *entity.LoginRequest) error {
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	return nil
}
