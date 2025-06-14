package usecase

import (
	"context"
	"errors"
	"testing"

	"layered-architecture-template/internal/domain/entity"
)

// MockUserRepository implements repository.UserRepository for testing
type MockUserRepository struct {
	users map[uint]*entity.User
	nextID uint
	emailIndex map[string]*entity.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[uint]*entity.User),
		nextID: 1,
		emailIndex: make(map[string]*entity.User),
	}
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}
	user.ID = m.nextID
	m.nextID++
	m.users[user.ID] = user
	m.emailIndex[user.Email] = user
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	if user, exists := m.emailIndex[email]; exists {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]*entity.User, error) {
	var users []*entity.User
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	if user == nil || user.ID == 0 {
		return errors.New("invalid user")
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

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	if user, exists := m.users[id]; exists {
		delete(m.users, id)
		delete(m.emailIndex, user.Email)
		return nil
	}
	return errors.New("user not found")
}

func TestUserUsecase_CreateUser(t *testing.T) {
	mockRepo := NewMockUserRepository()
	usecase := NewUserUsecase(mockRepo)
	ctx := context.Background()

	t.Run("successful user creation", func(t *testing.T) {
		user, err := usecase.CreateUser(ctx, "John Doe", "john@example.com")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if user == nil {
			t.Error("expected user to be created, got nil")
		}
		if user.Name != "John Doe" {
			t.Errorf("expected name 'John Doe', got %s", user.Name)
		}
		if user.Email != "john@example.com" {
			t.Errorf("expected email 'john@example.com', got %s", user.Email)
		}
		if user.ID == 0 {
			t.Error("expected user ID to be set")
		}
	})

	t.Run("empty name validation", func(t *testing.T) {
		user, err := usecase.CreateUser(ctx, "", "test@example.com")
		if err == nil {
			t.Error("expected error for empty name")
		}
		if user != nil {
			t.Error("expected nil user for invalid input")
		}
		expectedError := "name and email are required"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("empty email validation", func(t *testing.T) {
		user, err := usecase.CreateUser(ctx, "John Doe", "")
		if err == nil {
			t.Error("expected error for empty email")
		}
		if user != nil {
			t.Error("expected nil user for invalid input")
		}
		expectedError := "name and email are required"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("duplicate email validation", func(t *testing.T) {
		// First user creation should succeed
		_, err := usecase.CreateUser(ctx, "First User", "duplicate@example.com")
		if err != nil {
			t.Errorf("expected no error for first user, got %v", err)
		}

		// Second user with same email should fail
		user, err := usecase.CreateUser(ctx, "Second User", "duplicate@example.com")
		if err == nil {
			t.Error("expected error for duplicate email")
		}
		if user != nil {
			t.Error("expected nil user for duplicate email")
		}
		expectedError := "user with this email already exists"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})
}

func TestUserUsecase_GetUser(t *testing.T) {
	mockRepo := NewMockUserRepository()
	usecase := NewUserUsecase(mockRepo)
	ctx := context.Background()

	// Create a test user
	createdUser, _ := usecase.CreateUser(ctx, "Test User", "test@example.com")

	t.Run("get existing user", func(t *testing.T) {
		user, err := usecase.GetUser(ctx, createdUser.ID)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if user == nil {
			t.Error("expected user to be found, got nil")
		}
		if user.ID != createdUser.ID {
			t.Errorf("expected user ID %d, got %d", createdUser.ID, user.ID)
		}
	})

	t.Run("get non-existent user", func(t *testing.T) {
		user, err := usecase.GetUser(ctx, 999)
		if err == nil {
			t.Error("expected error for non-existent user")
		}
		if user != nil {
			t.Error("expected nil user for non-existent ID")
		}
	})
}

func TestUserUsecase_UpdateUser(t *testing.T) {
	mockRepo := NewMockUserRepository()
	usecase := NewUserUsecase(mockRepo)
	ctx := context.Background()

	// Create a test user
	createdUser, _ := usecase.CreateUser(ctx, "Original Name", "original@example.com")

	t.Run("successful user update", func(t *testing.T) {
		createdUser.Name = "Updated Name"
		createdUser.Email = "updated@example.com"
		
		err := usecase.UpdateUser(ctx, createdUser)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Verify the update
		updatedUser, _ := usecase.GetUser(ctx, createdUser.ID)
		if updatedUser.Name != "Updated Name" {
			t.Errorf("expected name 'Updated Name', got %s", updatedUser.Name)
		}
		if updatedUser.Email != "updated@example.com" {
			t.Errorf("expected email 'updated@example.com', got %s", updatedUser.Email)
		}
	})

	t.Run("update user with zero ID", func(t *testing.T) {
		user := &entity.User{
			ID: 0,
			Name: "Test",
			Email: "test@example.com",
		}
		
		err := usecase.UpdateUser(ctx, user)
		if err == nil {
			t.Error("expected error for zero ID")
		}
		expectedError := "user ID is required"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})
}

func TestUserUsecase_DeleteUser(t *testing.T) {
	mockRepo := NewMockUserRepository()
	usecase := NewUserUsecase(mockRepo)
	ctx := context.Background()

	// Create a test user
	createdUser, _ := usecase.CreateUser(ctx, "Test User", "test@example.com")

	t.Run("successful user deletion", func(t *testing.T) {
		err := usecase.DeleteUser(ctx, createdUser.ID)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Verify the user is deleted
		_, err = usecase.GetUser(ctx, createdUser.ID)
		if err == nil {
			t.Error("expected error when getting deleted user")
		}
	})

	t.Run("delete user with zero ID", func(t *testing.T) {
		err := usecase.DeleteUser(ctx, 0)
		if err == nil {
			t.Error("expected error for zero ID")
		}
		expectedError := "user ID is required"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("delete non-existent user", func(t *testing.T) {
		err := usecase.DeleteUser(ctx, 999)
		if err == nil {
			t.Error("expected error for non-existent user")
		}
	})
}

func TestUserUsecase_GetAllUsers(t *testing.T) {
	mockRepo := NewMockUserRepository()
	usecase := NewUserUsecase(mockRepo)
	ctx := context.Background()

	t.Run("get all users when empty", func(t *testing.T) {
		users, err := usecase.GetAllUsers(ctx)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(users) != 0 {
			t.Errorf("expected 0 users, got %d", len(users))
		}
	})

	t.Run("get all users with data", func(t *testing.T) {
		// Create test users
		usecase.CreateUser(ctx, "User 1", "user1@example.com")
		usecase.CreateUser(ctx, "User 2", "user2@example.com")
		usecase.CreateUser(ctx, "User 3", "user3@example.com")

		users, err := usecase.GetAllUsers(ctx)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(users) != 3 {
			t.Errorf("expected 3 users, got %d", len(users))
		}
	})
}

func TestUserUsecase_GetUserByEmail(t *testing.T) {
	mockRepo := NewMockUserRepository()
	usecase := NewUserUsecase(mockRepo)
	ctx := context.Background()

	// Create a test user
	createdUser, _ := usecase.CreateUser(ctx, "Test User", "findme@example.com")

	t.Run("get user by existing email", func(t *testing.T) {
		user, err := usecase.GetUserByEmail(ctx, "findme@example.com")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if user == nil {
			t.Error("expected user to be found, got nil")
		}
		if user.ID != createdUser.ID {
			t.Errorf("expected user ID %d, got %d", createdUser.ID, user.ID)
		}
	})

	t.Run("get user by non-existent email", func(t *testing.T) {
		user, err := usecase.GetUserByEmail(ctx, "notfound@example.com")
		if err == nil {
			t.Error("expected error for non-existent email")
		}
		if user != nil {
			t.Error("expected nil user for non-existent email")
		}
	})
}