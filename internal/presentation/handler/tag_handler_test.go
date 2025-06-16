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

// MockTagUsecase implements usecase.TagUsecase for testing
type MockTagUsecase struct {
	tags   map[uint]*entity.Tag
	nextID uint
}

func NewMockTagUsecase() *MockTagUsecase {
	return &MockTagUsecase{
		tags:   make(map[uint]*entity.Tag),
		nextID: 1,
	}
}

func (m *MockTagUsecase) CreateTag(ctx context.Context, name, color string) (*entity.Tag, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	// Check for duplicates
	for _, tag := range m.tags {
		if tag.Name == name {
			return nil, errors.New("tag with this name already exists")
		}
	}

	tag := &entity.Tag{
		ID:    m.nextID,
		Name:  name,
		Color: color,
	}
	if tag.Color == "" {
		tag.Color = "#3B82F6"
	}
	tag.GenerateSlug()
	m.nextID++
	m.tags[tag.ID] = tag
	return tag, nil
}

func (m *MockTagUsecase) GetTag(ctx context.Context, slug string) (*entity.Tag, error) {
	if slug == "" {
		return nil, errors.New("slug is required")
	}

	for _, tag := range m.tags {
		if tag.Slug == slug {
			return tag, nil
		}
	}
	return nil, errors.New("Tag not found")
}

func (m *MockTagUsecase) GetAllTags(ctx context.Context) ([]*entity.Tag, error) {
	var tags []*entity.Tag
	for _, tag := range m.tags {
		tags = append(tags, tag)
	}
	return tags, nil
}

func (m *MockTagUsecase) UpdateTag(ctx context.Context, id uint, name, color string) (*entity.Tag, error) {
	if id == 0 {
		return nil, errors.New("tag ID is required")
	}

	tag, exists := m.tags[id]
	if !exists {
		return nil, errors.New("tag not found")
	}

	// Check for duplicate names (excluding current tag)
	if name != "" {
		for _, t := range m.tags {
			if t.ID != id && t.Name == name {
				return nil, errors.New("tag with this name already exists")
			}
		}
		tag.Name = name
		tag.GenerateSlug()
	}
	
	if color != "" {
		tag.Color = color
	}

	return tag, nil
}

func (m *MockTagUsecase) DeleteTag(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("tag ID is required")
	}

	if _, exists := m.tags[id]; !exists {
		return errors.New("tag not found")
	}

	delete(m.tags, id)
	return nil
}

func (m *MockTagUsecase) GetOrCreateTags(ctx context.Context, tagNames []string) ([]*entity.Tag, error) {
	var tags []*entity.Tag
	for _, name := range tagNames {
		// Check if tag exists
		var existingTag *entity.Tag
		for _, tag := range m.tags {
			if tag.Name == name {
				existingTag = tag
				break
			}
		}
		
		if existingTag != nil {
			tags = append(tags, existingTag)
		} else {
			// Create new tag
			newTag, err := m.CreateTag(ctx, name, "#3B82F6")
			if err != nil {
				return nil, err
			}
			tags = append(tags, newTag)
		}
	}
	return tags, nil
}

func (m *MockTagUsecase) GetPopularTags(ctx context.Context, limit int) ([]*entity.TagWithCount, error) {
	// Mock implementation - return empty result
	return []*entity.TagWithCount{}, nil
}

func setupTagTestRouter(handler *TagHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	api := router.Group("/api/v1")
	{
		tags := api.Group("/tags")
		{
			tags.POST("", handler.CreateTag)
			tags.GET("", handler.GetAllTags)
			tags.GET("/popular", handler.GetPopularTags)
			tags.GET("/:slug", handler.GetTag)
			tags.PUT("/:id", handler.UpdateTag)
			tags.DELETE("/:id", handler.DeleteTag)
		}
	}

	return router
}

func TestTagHandler_CreateTag(t *testing.T) {
	mockUsecase := NewMockTagUsecase()
	handler := NewTagHandler(mockUsecase)
	router := setupTagTestRouter(handler)

	t.Run("successful tag creation", func(t *testing.T) {
		reqBody := CreateTagRequest{
			Name:  "Go",
			Color: "#00ADD8",
		}
		jsonData, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/tags", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		}

		var response entity.Tag
		json.Unmarshal(w.Body.Bytes(), &response)
		if response.Name != reqBody.Name {
			t.Errorf("expected name %s, got %s", reqBody.Name, response.Name)
		}
		if response.Color != reqBody.Color {
			t.Errorf("expected color %s, got %s", reqBody.Color, response.Color)
		}
	})

	t.Run("invalid JSON request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/tags", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("missing required name", func(t *testing.T) {
		reqBody := CreateTagRequest{
			Name:  "",
			Color: "#FF0000",
		}
		jsonData, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/tags", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("duplicate tag name", func(t *testing.T) {
		// Create first tag
		reqBody1 := CreateTagRequest{
			Name:  "JavaScript",
			Color: "#F7DF1E",
		}
		jsonData1, _ := json.Marshal(reqBody1)
		req1 := httptest.NewRequest(http.MethodPost, "/api/v1/tags", bytes.NewBuffer(jsonData1))
		req1.Header.Set("Content-Type", "application/json")
		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)

		// Try to create duplicate
		reqBody2 := CreateTagRequest{
			Name:  "JavaScript",
			Color: "#FF0000",
		}
		jsonData2, _ := json.Marshal(reqBody2)
		req2 := httptest.NewRequest(http.MethodPost, "/api/v1/tags", bytes.NewBuffer(jsonData2))
		req2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		if w2.Code != http.StatusConflict {
			t.Errorf("expected status %d, got %d", http.StatusConflict, w2.Code)
		}
	})
}

func TestTagHandler_GetTag(t *testing.T) {
	mockUsecase := NewMockTagUsecase()
	handler := NewTagHandler(mockUsecase)
	router := setupTagTestRouter(handler)

	// Create a test tag
	tag, _ := mockUsecase.CreateTag(context.Background(), "React", "#61DAFB")

	t.Run("get existing tag", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tags/"+tag.Slug, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response entity.Tag
		json.Unmarshal(w.Body.Bytes(), &response)
		if response.Name != tag.Name {
			t.Errorf("expected name %s, got %s", tag.Name, response.Name)
		}
	})

	t.Run("get non-existent tag", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tags/non-existent", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestTagHandler_UpdateTag(t *testing.T) {
	mockUsecase := NewMockTagUsecase()
	handler := NewTagHandler(mockUsecase)
	router := setupTagTestRouter(handler)

	// Create a test tag
	_, _ = mockUsecase.CreateTag(context.Background(), "Vue", "#4FC08D")

	t.Run("successful tag update", func(t *testing.T) {
		reqBody := UpdateTagRequest{
			Name:  "Vue.js",
			Color: "#00D4FF",
		}
		jsonData, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/tags/1", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response entity.Tag
		json.Unmarshal(w.Body.Bytes(), &response)
		if response.Name != reqBody.Name {
			t.Errorf("expected name %s, got %s", reqBody.Name, response.Name)
		}
		if response.Color != reqBody.Color {
			t.Errorf("expected color %s, got %s", reqBody.Color, response.Color)
		}
	})

	t.Run("update non-existent tag", func(t *testing.T) {
		reqBody := UpdateTagRequest{
			Name:  "Updated Name",
			Color: "#FF0000",
		}
		jsonData, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/tags/999", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("invalid tag ID", func(t *testing.T) {
		reqBody := UpdateTagRequest{
			Name:  "Updated Name",
			Color: "#FF0000",
		}
		jsonData, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/tags/invalid", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestTagHandler_DeleteTag(t *testing.T) {
	mockUsecase := NewMockTagUsecase()
	handler := NewTagHandler(mockUsecase)
	router := setupTagTestRouter(handler)

	// Create a test tag
	_, _ = mockUsecase.CreateTag(context.Background(), "Python", "#3776AB")

	t.Run("successful tag deletion", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/tags/1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("delete non-existent tag", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/tags/999", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("invalid tag ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/tags/invalid", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestTagHandler_GetPopularTags(t *testing.T) {
	mockUsecase := NewMockTagUsecase()
	handler := NewTagHandler(mockUsecase)
	router := setupTagTestRouter(handler)

	t.Run("get popular tags with default limit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tags/popular", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response []*entity.TagWithCount
		json.Unmarshal(w.Body.Bytes(), &response)
		// Mock implementation returns empty result
		if len(response) != 0 {
			t.Errorf("expected 0 tags from mock, got %d", len(response))
		}
	})

	t.Run("get popular tags with custom limit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tags/popular?limit=5", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})
}