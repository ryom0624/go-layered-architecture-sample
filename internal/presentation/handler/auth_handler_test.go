package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"layered-architecture-template/internal/domain/entity"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockAuthUsecase struct {
	users         map[string]*entity.User
	userIDCounter uint
}

func newMockAuthUsecase() *mockAuthUsecase {
	return &mockAuthUsecase{
		users:         make(map[string]*entity.User),
		userIDCounter: 1,
	}
}

func (m *mockAuthUsecase) Register(ctx context.Context, req *entity.RegisterRequest) (*entity.LoginResponse, error) {
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return nil, errors.New("missing required fields")
	}
	if len(req.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}
	if _, exists := m.users[req.Email]; exists {
		return nil, errors.New("email already exists")
	}

	user := &entity.User{
		ID:    m.userIDCounter,
		Name:  req.Name,
		Email: req.Email,
	}
	m.userIDCounter++
	m.users[req.Email] = user

	return &entity.LoginResponse{
		User:  *user,
		Token: "mock-access-token",
	}, nil
}

func (m *mockAuthUsecase) Login(ctx context.Context, req *entity.LoginRequest) (*entity.LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("missing required fields")
	}

	user, exists := m.users[req.Email]
	if !exists {
		return nil, errors.New("invalid email or password")
	}

	if req.Password != "password123" {
		return nil, errors.New("invalid email or password")
	}

	return &entity.LoginResponse{
		User:  *user,
		Token: "mock-access-token",
	}, nil
}

func (m *mockAuthUsecase) RefreshToken(ctx context.Context, refreshToken string) (*entity.LoginResponse, error) {
	if refreshToken == "" {
		return nil, errors.New("refresh token required")
	}
	if refreshToken != "valid-refresh-token" {
		return nil, errors.New("invalid refresh token")
	}

	// Return first user for testing
	for _, user := range m.users {
		return &entity.LoginResponse{
			User:  *user,
			Token: "new-access-token",
		}, nil
	}

	return nil, errors.New("no users found")
}

func (m *mockAuthUsecase) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return errors.New("refresh token required")
	}
	if refreshToken != "valid-refresh-token" {
		return errors.New("invalid refresh token")
	}
	return nil
}

func (m *mockAuthUsecase) LogoutAll(ctx context.Context, userID uint) error {
	if userID == 0 {
		return errors.New("user ID required")
	}
	return nil
}

func setupAuthTestRouter(authUsecase *mockAuthUsecase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	
	authHandler := NewAuthHandler(authUsecase)
	
	api := r.Group("/api/v1")
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authHandler.Logout)
		auth.POST("/logout-all", func(c *gin.Context) {
			c.Set("user_id", uint(1))
			authHandler.LogoutAll(c)
		})
	}
	
	return r
}

func TestAuthHandler_Register(t *testing.T) {
	mockUsecase := newMockAuthUsecase()
	router := setupAuthTestRouter(mockUsecase)

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedUser   bool
		expectedToken  bool
	}{
		{
			name: "valid registration",
			requestBody: `{
				"name": "John Doe",
				"email": "john@example.com",
				"password": "password123"
			}`,
			expectedStatus: http.StatusCreated,
			expectedUser:   true,
			expectedToken:  true,
		},
		{
			name:           "invalid JSON",
			requestBody:    `{"name": "John"`,
			expectedStatus: http.StatusBadRequest,
			expectedUser:   false,
			expectedToken:  false,
		},
		{
			name: "missing required fields",
			requestBody: `{
				"name": "",
				"email": "john@example.com",
				"password": "password123"
			}`,
			expectedStatus: http.StatusBadRequest,
			expectedUser:   false,
			expectedToken:  false,
		},
		{
			name: "short password",
			requestBody: `{
				"name": "John Doe",
				"email": "john@example.com",
				"password": "123"
			}`,
			expectedStatus: http.StatusBadRequest,
			expectedUser:   false,
			expectedToken:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Register() status = %d, want %d", w.Code, tt.expectedStatus)
			}

			if tt.expectedUser || tt.expectedToken {
				var response entity.LoginResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Register() failed to unmarshal response: %v", err)
					return
				}

				if tt.expectedUser && response.User.Email == "" {
					t.Error("Register() user not present in response")
				}
				if tt.expectedToken && response.Token == "" {
					t.Error("Register() token not present in response")
				}
			}
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	mockUsecase := newMockAuthUsecase()
	router := setupAuthTestRouter(mockUsecase)

	// Register a user first
	registerReq := `{
		"name": "John Doe",
		"email": "john@example.com",
		"password": "password123"
	}`
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(registerReq))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedUser   bool
		expectedToken  bool
	}{
		{
			name: "valid login",
			requestBody: `{
				"email": "john@example.com",
				"password": "password123"
			}`,
			expectedStatus: http.StatusOK,
			expectedUser:   true,
			expectedToken:  true,
		},
		{
			name:           "invalid JSON",
			requestBody:    `{"email": "john"`,
			expectedStatus: http.StatusBadRequest,
			expectedUser:   false,
			expectedToken:  false,
		},
		{
			name: "invalid credentials",
			requestBody: `{
				"email": "john@example.com",
				"password": "wrongpassword"
			}`,
			expectedStatus: http.StatusUnauthorized,
			expectedUser:   false,
			expectedToken:  false,
		},
		{
			name: "nonexistent user",
			requestBody: `{
				"email": "nonexistent@example.com",
				"password": "password123"
			}`,
			expectedStatus: http.StatusUnauthorized,
			expectedUser:   false,
			expectedToken:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Login() status = %d, want %d", w.Code, tt.expectedStatus)
			}

			if tt.expectedUser || tt.expectedToken {
				var response entity.LoginResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Login() failed to unmarshal response: %v", err)
					return
				}

				if tt.expectedUser && response.User.Email == "" {
					t.Error("Login() user not present in response")
				}
				if tt.expectedToken && response.Token == "" {
					t.Error("Login() token not present in response")
				}
			}
		})
	}
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	mockUsecase := newMockAuthUsecase()
	router := setupAuthTestRouter(mockUsecase)

	// Register a user first
	registerReq := `{
		"name": "John Doe",
		"email": "john@example.com",
		"password": "password123"
	}`
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(registerReq))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedToken  bool
	}{
		{
			name: "valid refresh token",
			requestBody: `{
				"refresh_token": "valid-refresh-token"
			}`,
			expectedStatus: http.StatusOK,
			expectedToken:  true,
		},
		{
			name:           "invalid JSON",
			requestBody:    `{"refresh_token"`,
			expectedStatus: http.StatusBadRequest,
			expectedToken:  false,
		},
		{
			name: "invalid refresh token",
			requestBody: `{
				"refresh_token": "invalid-token"
			}`,
			expectedStatus: http.StatusUnauthorized,
			expectedToken:  false,
		},
		{
			name: "missing refresh token",
			requestBody: `{
				"refresh_token": ""
			}`,
			expectedStatus: http.StatusBadRequest,
			expectedToken:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/api/v1/auth/refresh", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("RefreshToken() status = %d, want %d", w.Code, tt.expectedStatus)
			}

			if tt.expectedToken {
				var response entity.LoginResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("RefreshToken() failed to unmarshal response: %v", err)
					return
				}

				if response.Token == "" {
					t.Error("RefreshToken() token not present in response")
				}
			}
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	mockUsecase := newMockAuthUsecase()
	router := setupAuthTestRouter(mockUsecase)

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
	}{
		{
			name: "valid logout",
			requestBody: `{
				"refresh_token": "valid-refresh-token"
			}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid JSON",
			requestBody:    `{"refresh_token"`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid refresh token",
			requestBody: `{
				"refresh_token": "invalid-token"
			}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/api/v1/auth/logout", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Logout() status = %d, want %d", w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestAuthHandler_LogoutAll(t *testing.T) {
	mockUsecase := newMockAuthUsecase()
	router := setupAuthTestRouter(mockUsecase)

	tests := []struct {
		name           string
		expectedStatus int
	}{
		{
			name:           "valid logout all",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/api/v1/auth/logout-all", bytes.NewReader([]byte{}))
			req.Header.Set("Content-Type", "application/json")
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("LogoutAll() status = %d, want %d", w.Code, tt.expectedStatus)
			}
		})
	}
}