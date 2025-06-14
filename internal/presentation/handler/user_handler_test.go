package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"layered-architecture-template/internal/domain/entity"

	"github.com/gin-gonic/gin"
)

// MockUserUsecase implements usecase.UserUsecase for testing
type MockUserUsecase struct {
	users map[uint]*entity.User
	nextID uint
	emailIndex map[string]*entity.User
}

func NewMockUserUsecase() *MockUserUsecase {
	return &MockUserUsecase{
		users: make(map[uint]*entity.User),
		nextID: 1,
		emailIndex: make(map[string]*entity.User),
	}
}

func (m *MockUserUsecase) CreateUser(ctx context.Context, name, email string) (*entity.User, error) {
	if name == "" || email == "" {
		return nil, errors.New("name and email are required")
	}
	
	if _, exists := m.emailIndex[email]; exists {
		return nil, errors.New("user with this email already exists")
	}

	user := &entity.User{
		ID: m.nextID,
		Name: name,
		Email: email,
	}
	m.nextID++
	m.users[user.ID] = user
	m.emailIndex[email] = user
	return user, nil
}

func (m *MockUserUsecase) GetUser(ctx context.Context, id uint) (*entity.User, error) {
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func (m *MockUserUsecase) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	if user, exists := m.emailIndex[email]; exists {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func (m *MockUserUsecase) GetAllUsers(ctx context.Context) ([]*entity.User, error) {
	var users []*entity.User
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func (m *MockUserUsecase) UpdateUser(ctx context.Context, user *entity.User) error {
	if user == nil || user.ID == 0 {
		return errors.New("user ID is required")
	}
	if _, exists := m.users[user.ID]; !exists {
		return errors.New("user not found")
	}
	// Remove old email index
	for email, u := range m.emailIndex {
		if u.ID == user.ID {
			delete(m.emailIndex, email)
			break
		}
	}
	m.users[user.ID] = user
	m.emailIndex[user.Email] = user
	return nil
}

func (m *MockUserUsecase) DeleteUser(ctx context.Context, id uint) error {
	if user, exists := m.users[id]; exists {
		delete(m.users, id)
		delete(m.emailIndex, user.Email)
		return nil
	}
	return errors.New("user not found")
}

func setupTestRouter(handler *UserHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	api := router.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.POST("", handler.CreateUser)
			users.GET("", handler.GetAllUsers)
			users.GET("/:id", handler.GetUser)
			users.PUT("/:id", handler.UpdateUser)
			users.DELETE("/:id", handler.DeleteUser)
		}
	}
	
	return router
}

func TestUserHandler_CreateUser(t *testing.T) {
	mockUsecase := NewMockUserUsecase()
	handler := NewUserHandler(mockUsecase)
	router := setupTestRouter(handler)

	t.Run("successful user creation", func(t *testing.T) {
		reqBody := CreateUserRequest{
			Name:  "John Doe",
			Email: "john@example.com",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		}

		var response entity.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("failed to unmarshal response: %v", err)
		}

		if response.Name != "John Doe" {
			t.Errorf("expected name 'John Doe', got %s", response.Name)
		}
		if response.Email != "john@example.com" {
			t.Errorf("expected email 'john@example.com', got %s", response.Email)
		}
		if response.ID == 0 {
			t.Error("expected user ID to be set")
		}
	})

	t.Run("invalid JSON request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("missing required fields", func(t *testing.T) {
		reqBody := CreateUserRequest{
			Name: "John Doe",
			// Email is missing
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("duplicate email", func(t *testing.T) {
		// Create first user
		reqBody := CreateUserRequest{
			Name:  "First User",
			Email: "duplicate@example.com",
		}
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Try to create second user with same email
		reqBody2 := CreateUserRequest{
			Name:  "Second User",
			Email: "duplicate@example.com",
		}
		jsonBody2, _ := json.Marshal(reqBody2)
		req2 := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(jsonBody2))
		req2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		if w2.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w2.Code)
		}
	})
}

func TestUserHandler_GetUser(t *testing.T) {
	mockUsecase := NewMockUserUsecase()
	handler := NewUserHandler(mockUsecase)
	router := setupTestRouter(handler)

	// Create a test user
	testUser, _ := mockUsecase.CreateUser(context.Background(), "Test User", "test@example.com")

	t.Run("get existing user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response entity.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("failed to unmarshal response: %v", err)
		}

		if response.ID != testUser.ID {
			t.Errorf("expected user ID %d, got %d", testUser.ID, response.ID)
		}
	})

	t.Run("get non-existent user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/999", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("invalid user ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/invalid", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestUserHandler_GetAllUsers(t *testing.T) {
	mockUsecase := NewMockUserUsecase()
	handler := NewUserHandler(mockUsecase)
	router := setupTestRouter(handler)

	t.Run("get all users when empty", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response []*entity.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("failed to unmarshal response: %v", err)
		}

		if len(response) != 0 {
			t.Errorf("expected 0 users, got %d", len(response))
		}
	})

	t.Run("get all users with data", func(t *testing.T) {
		// Create test users
		mockUsecase.CreateUser(context.Background(), "User 1", "user1@example.com")
		mockUsecase.CreateUser(context.Background(), "User 2", "user2@example.com")

		req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response []*entity.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("failed to unmarshal response: %v", err)
		}

		if len(response) != 2 {
			t.Errorf("expected 2 users, got %d", len(response))
		}
	})
}

func TestUserHandler_UpdateUser(t *testing.T) {
	mockUsecase := NewMockUserUsecase()
	handler := NewUserHandler(mockUsecase)
	router := setupTestRouter(handler)

	// Create a test user
	_, _ = mockUsecase.CreateUser(context.Background(), "Original Name", "original@example.com")

	t.Run("successful user update", func(t *testing.T) {
		reqBody := UpdateUserRequest{
			Name:  "Updated Name",
			Email: "updated@example.com",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/users/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response entity.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("failed to unmarshal response: %v", err)
		}

		if response.Name != "Updated Name" {
			t.Errorf("expected name 'Updated Name', got %s", response.Name)
		}
		if response.Email != "updated@example.com" {
			t.Errorf("expected email 'updated@example.com', got %s", response.Email)
		}
	})

	t.Run("partial update", func(t *testing.T) {
		reqBody := UpdateUserRequest{
			Name: "Partially Updated Name",
			// Email not provided, should remain unchanged
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/users/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response entity.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("failed to unmarshal response: %v", err)
		}

		if response.Name != "Partially Updated Name" {
			t.Errorf("expected name 'Partially Updated Name', got %s", response.Name)
		}
		// Email should remain from previous test
		if response.Email != "updated@example.com" {
			t.Errorf("expected email to remain 'updated@example.com', got %s", response.Email)
		}
	})

	t.Run("update non-existent user", func(t *testing.T) {
		reqBody := UpdateUserRequest{
			Name: "Updated Name",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/users/999", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("invalid user ID", func(t *testing.T) {
		reqBody := UpdateUserRequest{
			Name: "Updated Name",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/users/invalid", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestUserHandler_DeleteUser(t *testing.T) {
	mockUsecase := NewMockUserUsecase()
	handler := NewUserHandler(mockUsecase)
	router := setupTestRouter(handler)

	// Create a test user
	mockUsecase.CreateUser(context.Background(), "Test User", "test@example.com")

	t.Run("successful user deletion", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
		}

		// Verify user is deleted by trying to get it
		getReq := httptest.NewRequest(http.MethodGet, "/api/v1/users/1", nil)
		getW := httptest.NewRecorder()
		router.ServeHTTP(getW, getReq)

		if getW.Code != http.StatusNotFound {
			t.Errorf("expected user to be deleted, but got status %d", getW.Code)
		}
	})

	t.Run("delete non-existent user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/999", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("invalid user ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/invalid", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}