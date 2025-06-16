package usecase

import (
	"context"
	"errors"
	"layered-architecture-template/internal/domain/entity"
	"layered-architecture-template/pkg/jwt"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type mockAuthRepository struct {
	users         map[string]*entity.User
	refreshTokens map[string]*entity.RefreshToken
	userIDCounter uint
}

func newMockAuthRepository() *mockAuthRepository {
	return &mockAuthRepository{
		users:         make(map[string]*entity.User),
		refreshTokens: make(map[string]*entity.RefreshToken),
		userIDCounter: 1,
	}
}

// Helper function to create a hashed password for testing
func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func (m *mockAuthRepository) Register(ctx context.Context, user *entity.User) error {
	if _, exists := m.users[user.Email]; exists {
		return errors.New("email already exists")
	}
	user.ID = m.userIDCounter
	m.userIDCounter++
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	// Store a copy of the user with the hashed password
	userCopy := &entity.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password, // This will be the hashed password from the usecase
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	m.users[user.Email] = userCopy
	return nil
}

func (m *mockAuthRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, exists := m.users[email]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (m *mockAuthRepository) StoreRefreshToken(ctx context.Context, token *entity.RefreshToken) error {
	token.ID = uint(len(m.refreshTokens) + 1)
	token.CreatedAt = time.Now()
	m.refreshTokens[token.Token] = token
	return nil
}

func (m *mockAuthRepository) GetRefreshToken(ctx context.Context, token string) (*entity.RefreshToken, error) {
	refreshToken, exists := m.refreshTokens[token]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}
	
	// Find user by UserID
	for _, user := range m.users {
		if user.ID == refreshToken.UserID {
			refreshToken.User = *user
			break
		}
	}
	
	return refreshToken, nil
}

func (m *mockAuthRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	delete(m.refreshTokens, token)
	return nil
}

func (m *mockAuthRepository) DeleteUserRefreshTokens(ctx context.Context, userID uint) error {
	for token, refreshToken := range m.refreshTokens {
		if refreshToken.UserID == userID {
			delete(m.refreshTokens, token)
		}
	}
	return nil
}

func TestAuthUsecase_Register(t *testing.T) {
	repo := newMockAuthRepository()
	jwtManager := jwt.NewJWTManager("test-secret", 15*time.Minute)
	usecase := NewAuthUsecase(repo, jwtManager, 7*24*time.Hour)

	tests := []struct {
		name    string
		req     *entity.RegisterRequest
		wantErr bool
	}{
		{
			name: "valid registration",
			req: &entity.RegisterRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "duplicate email",
			req: &entity.RegisterRequest{
				Name:     "Jane Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "invalid name",
			req: &entity.RegisterRequest{
				Name:     "",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "invalid password",
			req: &entity.RegisterRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "123",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := usecase.Register(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("Register() response is nil")
					return
				}
				if resp.User.Email != tt.req.Email {
					t.Errorf("Register() email = %v, want %v", resp.User.Email, tt.req.Email)
				}
				if resp.Token == "" {
					t.Error("Register() token is empty")
				}
			}
		})
	}
}

func TestAuthUsecase_Login(t *testing.T) {
	repo := newMockAuthRepository()
	jwtManager := jwt.NewJWTManager("test-secret", 15*time.Minute)
	usecase := NewAuthUsecase(repo, jwtManager, 7*24*time.Hour)

	// Manually create a user with hashed password in the mock repository
	hashedPassword := hashPassword("password123")
	testUser := &entity.User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: hashedPassword,
	}
	repo.users["john@example.com"] = testUser
	repo.userIDCounter = 2

	tests := []struct {
		name    string
		req     *entity.LoginRequest
		wantErr bool
	}{
		{
			name: "valid login",
			req: &entity.LoginRequest{
				Email:    "john@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			req: &entity.LoginRequest{
				Email:    "notexist@example.com",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "invalid password",
			req: &entity.LoginRequest{
				Email:    "john@example.com",
				Password: "wrongpassword",
			},
			wantErr: true,
		},
		{
			name: "empty email",
			req: &entity.LoginRequest{
				Email:    "",
				Password: "password123",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := usecase.Login(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("Login() response is nil")
					return
				}
				if resp.User.Email != tt.req.Email {
					t.Errorf("Login() email = %v, want %v", resp.User.Email, tt.req.Email)
				}
				if resp.Token == "" {
					t.Error("Login() token is empty")
				}
			}
		})
	}
}

func TestAuthUsecase_RefreshToken(t *testing.T) {
	repo := newMockAuthRepository()
	jwtManager := jwt.NewJWTManager("test-secret", 15*time.Minute)
	usecase := NewAuthUsecase(repo, jwtManager, 7*24*time.Hour)

	// Manually create a user and refresh token
	testUser := &entity.User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: hashPassword("password123"),
	}
	repo.users["john@example.com"] = testUser
	
	refreshTokenStr := "test-refresh-token"
	refreshToken := &entity.RefreshToken{
		ID:        1,
		UserID:    testUser.ID,
		Token:     refreshTokenStr,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		User:      *testUser,
	}
	repo.refreshTokens[refreshTokenStr] = refreshToken

	tests := []struct {
		name         string
		refreshToken string
		wantErr      bool
	}{
		{
			name:         "valid refresh token",
			refreshToken: refreshTokenStr,
			wantErr:      false,
		},
		{
			name:         "invalid refresh token",
			refreshToken: "invalid-token",
			wantErr:      true,
		},
		{
			name:         "empty refresh token",
			refreshToken: "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := usecase.RefreshToken(context.Background(), tt.refreshToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("RefreshToken() response is nil")
					return
				}
				if resp.User.Email != testUser.Email {
					t.Errorf("RefreshToken() email = %v, want %v", resp.User.Email, testUser.Email)
				}
				if resp.Token == "" {
					t.Error("RefreshToken() token is empty")
				}
			}
		})
	}
}

func TestAuthUsecase_Logout(t *testing.T) {
	repo := newMockAuthRepository()
	jwtManager := jwt.NewJWTManager("test-secret", 15*time.Minute)
	usecase := NewAuthUsecase(repo, jwtManager, 7*24*time.Hour)

	// Create a refresh token manually
	refreshTokenStr := "test-refresh-token"
	refreshToken := &entity.RefreshToken{
		ID:        1,
		UserID:    1,
		Token:     refreshTokenStr,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	repo.refreshTokens[refreshTokenStr] = refreshToken

	err := usecase.Logout(context.Background(), refreshTokenStr)
	if err != nil {
		t.Errorf("Logout() error = %v", err)
	}

	// Verify token is deleted
	_, exists := repo.refreshTokens[refreshTokenStr]
	if exists {
		t.Error("Logout() did not delete refresh token")
	}
}

func TestAuthUsecase_LogoutAll(t *testing.T) {
	repo := newMockAuthRepository()
	jwtManager := jwt.NewJWTManager("test-secret", 15*time.Minute)
	usecase := NewAuthUsecase(repo, jwtManager, 7*24*time.Hour)

	userID := uint(1)
	
	// Create multiple refresh tokens for the same user
	token1 := &entity.RefreshToken{
		ID:        1,
		UserID:    userID,
		Token:     "token1",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	token2 := &entity.RefreshToken{
		ID:        2,
		UserID:    userID,
		Token:     "token2",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	repo.refreshTokens["token1"] = token1
	repo.refreshTokens["token2"] = token2

	err := usecase.LogoutAll(context.Background(), userID)
	if err != nil {
		t.Errorf("LogoutAll() error = %v", err)
	}

	// Verify all tokens for user are deleted
	for _, refreshToken := range repo.refreshTokens {
		if refreshToken.UserID == userID {
			t.Error("LogoutAll() did not delete all refresh tokens for user")
		}
	}
}