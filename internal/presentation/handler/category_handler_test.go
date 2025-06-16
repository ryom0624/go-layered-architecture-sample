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

// MockCategoryUsecase implements usecase.CategoryUsecase for testing
type MockCategoryUsecase struct {
	categories map[uint]*entity.Category
	nextID     uint
}

func NewMockCategoryUsecase() *MockCategoryUsecase {
	return &MockCategoryUsecase{
		categories: make(map[uint]*entity.Category),
		nextID:     1,
	}
}

func (m *MockCategoryUsecase) CreateCategory(ctx context.Context, name, description string) (*entity.Category, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	// Check for duplicates
	for _, category := range m.categories {
		if category.Name == name {
			return nil, errors.New("category with this name already exists")
		}
	}

	category := &entity.Category{
		ID:          m.nextID,
		Name:        name,
		Description: description,
	}
	category.GenerateSlug()
	m.nextID++
	m.categories[category.ID] = category
	return category, nil
}

func (m *MockCategoryUsecase) GetCategory(ctx context.Context, slug string) (*entity.Category, error) {
	if slug == "" {
		return nil, errors.New("slug is required")
	}

	for _, category := range m.categories {
		if category.Slug == slug {
			return category, nil
		}
	}
	return nil, errors.New("Category not found")
}

func (m *MockCategoryUsecase) GetAllCategories(ctx context.Context) ([]*entity.Category, error) {
	var categories []*entity.Category
	for _, category := range m.categories {
		categories = append(categories, category)
	}
	return categories, nil
}

func (m *MockCategoryUsecase) UpdateCategory(ctx context.Context, id uint, name, description string) (*entity.Category, error) {
	if id == 0 {
		return nil, errors.New("category ID is required")
	}

	category, exists := m.categories[id]
	if !exists {
		return nil, errors.New("category not found")
	}

	// Check for duplicate names (excluding current category)
	if name != "" {
		for _, c := range m.categories {
			if c.ID != id && c.Name == name {
				return nil, errors.New("category with this name already exists")
			}
		}
		category.Name = name
		category.GenerateSlug()
	}
	
	if description != "" {
		category.Description = description
	}

	return category, nil
}

func (m *MockCategoryUsecase) DeleteCategory(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("category ID is required")
	}

	if _, exists := m.categories[id]; !exists {
		return errors.New("category not found")
	}

	delete(m.categories, id)
	return nil
}

func (m *MockCategoryUsecase) GetCategoriesWithCount(ctx context.Context) ([]*entity.CategoryWithCount, error) {
	var categories []*entity.CategoryWithCount
	for _, category := range m.categories {
		categories = append(categories, &entity.CategoryWithCount{
			Category:     category,
			ArticleCount: 0, // Mock implementation
		})
	}
	return categories, nil
}

func setupCategoryTestRouter(handler *CategoryHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	api := router.Group("/api/v1")
	{
		categories := api.Group("/categories")
		{
			categories.POST("", handler.CreateCategory)
			categories.GET("", handler.GetAllCategories)
			categories.GET("/with-count", handler.GetCategoriesWithCount)
			categories.GET("/:slug", handler.GetCategory)
			categories.PUT("/:id", handler.UpdateCategory)
			categories.DELETE("/:id", handler.DeleteCategory)
		}
	}

	return router
}

func TestCategoryHandler_CreateCategory(t *testing.T) {
	mockUsecase := NewMockCategoryUsecase()
	handler := NewCategoryHandler(mockUsecase)
	router := setupCategoryTestRouter(handler)

	t.Run("successful category creation", func(t *testing.T) {
		reqBody := CreateCategoryRequest{
			Name:        "Technology",
			Description: "Technology articles",
		}
		jsonData, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/categories", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		}

		var response entity.Category
		json.Unmarshal(w.Body.Bytes(), &response)
		if response.Name != reqBody.Name {
			t.Errorf("expected name %s, got %s", reqBody.Name, response.Name)
		}
	})

	t.Run("invalid JSON request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/categories", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("missing required name", func(t *testing.T) {
		reqBody := CreateCategoryRequest{
			Name:        "",
			Description: "Description without name",
		}
		jsonData, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/categories", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestCategoryHandler_GetCategory(t *testing.T) {
	mockUsecase := NewMockCategoryUsecase()
	handler := NewCategoryHandler(mockUsecase)
	router := setupCategoryTestRouter(handler)

	// Create a test category
	category, _ := mockUsecase.CreateCategory(context.Background(), "Science", "Science articles")

	t.Run("get existing category", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/categories/"+category.Slug, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response entity.Category
		json.Unmarshal(w.Body.Bytes(), &response)
		if response.Name != category.Name {
			t.Errorf("expected name %s, got %s", category.Name, response.Name)
		}
	})

	t.Run("get non-existent category", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/categories/non-existent", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestCategoryHandler_UpdateCategory(t *testing.T) {
	mockUsecase := NewMockCategoryUsecase()
	handler := NewCategoryHandler(mockUsecase)
	router := setupCategoryTestRouter(handler)

	// Create a test category
	_, _ = mockUsecase.CreateCategory(context.Background(), "Programming", "Programming articles")

	t.Run("successful category update", func(t *testing.T) {
		reqBody := UpdateCategoryRequest{
			Name:        "Software Development",
			Description: "Software development articles",
		}
		jsonData, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/categories/1", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response entity.Category
		json.Unmarshal(w.Body.Bytes(), &response)
		if response.Name != reqBody.Name {
			t.Errorf("expected name %s, got %s", reqBody.Name, response.Name)
		}
	})

	t.Run("update non-existent category", func(t *testing.T) {
		reqBody := UpdateCategoryRequest{
			Name:        "Updated Name",
			Description: "Updated Description",
		}
		jsonData, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/categories/999", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("invalid category ID", func(t *testing.T) {
		reqBody := UpdateCategoryRequest{
			Name:        "Updated Name",
			Description: "Updated Description",
		}
		jsonData, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/categories/invalid", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestCategoryHandler_DeleteCategory(t *testing.T) {
	mockUsecase := NewMockCategoryUsecase()
	handler := NewCategoryHandler(mockUsecase)
	router := setupCategoryTestRouter(handler)

	// Create a test category
	_, _ = mockUsecase.CreateCategory(context.Background(), "Web Development", "Web development articles")

	t.Run("successful category deletion", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/categories/1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("delete non-existent category", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/categories/999", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("invalid category ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/categories/invalid", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}