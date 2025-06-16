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

// MockArticleUsecase implements usecase.ArticleUsecase for testing
type MockArticleUsecase struct {
	articles map[uint]*entity.Article
	nextID   uint
}

func NewMockArticleUsecase() *MockArticleUsecase {
	return &MockArticleUsecase{
		articles: make(map[uint]*entity.Article),
		nextID:   1,
	}
}

func (m *MockArticleUsecase) CreateArticle(ctx context.Context, title, content string, authorID uint) (*entity.Article, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}
	if content == "" {
		return nil, errors.New("content is required")
	}
	if authorID == 0 {
		return nil, errors.New("author ID is required")
	}

	article := &entity.Article{
		ID:       m.nextID,
		Title:    title,
		Content:  content,
		AuthorID: authorID,
		Status:   "draft",
	}
	m.nextID++
	m.articles[article.ID] = article
	return article, nil
}

func (m *MockArticleUsecase) GetArticle(ctx context.Context, id uint) (*entity.Article, error) {
	if article, exists := m.articles[id]; exists {
		return article, nil
	}
	return nil, errors.New("article not found")
}

func (m *MockArticleUsecase) GetAllArticles(ctx context.Context) ([]*entity.Article, error) {
	var articles []*entity.Article
	for _, article := range m.articles {
		articles = append(articles, article)
	}
	return articles, nil
}

func (m *MockArticleUsecase) GetPublishedArticles(ctx context.Context) ([]*entity.Article, error) {
	var articles []*entity.Article
	for _, article := range m.articles {
		if article.Status == "published" {
			articles = append(articles, article)
		}
	}
	return articles, nil
}

func (m *MockArticleUsecase) UpdateArticle(ctx context.Context, id uint, title, content string) (*entity.Article, error) {
	if id == 0 {
		return nil, errors.New("article ID is required")
	}

	article, exists := m.articles[id]
	if !exists {
		return nil, errors.New("article not found")
	}

	if title != "" {
		article.Title = title
	}
	if content != "" {
		article.Content = content
	}

	return article, nil
}

func (m *MockArticleUsecase) DeleteArticle(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("article ID is required")
	}

	if _, exists := m.articles[id]; !exists {
		return errors.New("article not found")
	}

	delete(m.articles, id)
	return nil
}

func (m *MockArticleUsecase) PublishArticle(ctx context.Context, id uint) (*entity.Article, error) {
	article, exists := m.articles[id]
	if !exists {
		return nil, errors.New("article not found")
	}

	if article.Status == "published" {
		return nil, errors.New("article is already published")
	}

	article.Status = "published"
	return article, nil
}

func (m *MockArticleUsecase) UnpublishArticle(ctx context.Context, id uint) (*entity.Article, error) {
	article, exists := m.articles[id]
	if !exists {
		return nil, errors.New("article not found")
	}

	if article.Status == "draft" {
		return nil, errors.New("article is already unpublished")
	}

	article.Status = "draft"
	return article, nil
}

func (m *MockArticleUsecase) GetAllArticlesWithFilters(ctx context.Context, params *entity.ArticleSearchParams) (*entity.ArticleSearchResult, error) {
	var articles []*entity.Article
	for _, article := range m.articles {
		articles = append(articles, article)
	}
	
	return &entity.ArticleSearchResult{
		Articles:   articles,
		Total:      int64(len(articles)),
		Page:       1,
		Limit:      20,
		TotalPages: 1,
	}, nil
}

func (m *MockArticleUsecase) CreateArticleWithCategoryAndTags(ctx context.Context, title, content string, authorID uint, categoryID *uint, tagNames []string) (*entity.Article, error) {
	return m.CreateArticle(ctx, title, content, authorID)
}

func (m *MockArticleUsecase) GetArticlesByCategory(ctx context.Context, categorySlug string) ([]*entity.Article, error) {
	var articles []*entity.Article
	for _, article := range m.articles {
		articles = append(articles, article)
	}
	return articles, nil
}

func (m *MockArticleUsecase) GetArticlesByTags(ctx context.Context, tagSlugs []string) ([]*entity.Article, error) {
	var articles []*entity.Article
	for _, article := range m.articles {
		articles = append(articles, article)
	}
	return articles, nil
}

func (m *MockArticleUsecase) GetFilteredArticles(ctx context.Context, categorySlug string, tagSlugs []string) ([]*entity.Article, error) {
	var articles []*entity.Article
	for _, article := range m.articles {
		articles = append(articles, article)
	}
	return articles, nil
}

func (m *MockArticleUsecase) UpdateArticleWithCategoryAndTags(ctx context.Context, id uint, title, content string, categoryID *uint, tagNames []string) (*entity.Article, error) {
	return m.UpdateArticle(ctx, id, title, content)
}

func setupArticleTestRouter(handler *ArticleHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	api := router.Group("/api/v1")
	{
		articles := api.Group("/articles")
		{
			articles.POST("", handler.CreateArticle)
			articles.GET("", handler.GetAllArticles)
			articles.GET("/published", handler.GetPublishedArticles)
			articles.GET("/:id", handler.GetArticle)
			articles.PUT("/:id", handler.UpdateArticle)
			articles.DELETE("/:id", handler.DeleteArticle)
			articles.PUT("/:id/publish", handler.PublishArticle)
			articles.PUT("/:id/unpublish", handler.UnpublishArticle)
		}
	}

	return router
}

func TestArticleHandler_CreateArticle(t *testing.T) {
	mockUsecase := NewMockArticleUsecase()
	handler := NewArticleHandler(mockUsecase)
	router := setupArticleTestRouter(handler)

	t.Run("successful article creation", func(t *testing.T) {
		reqBody := CreateArticleRequest{
			Title:    "Test Article",
			Content:  "Test Content",
			AuthorID: 1,
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/articles", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		}

		var response entity.Article
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("failed to unmarshal response: %v", err)
		}

		if response.Title != "Test Article" {
			t.Errorf("expected title 'Test Article', got %s", response.Title)
		}
		if response.Content != "Test Content" {
			t.Errorf("expected content 'Test Content', got %s", response.Content)
		}
		if response.AuthorID != 1 {
			t.Errorf("expected author ID 1, got %d", response.AuthorID)
		}
		if response.Status != "draft" {
			t.Errorf("expected status 'draft', got %s", response.Status)
		}
	})

	t.Run("invalid JSON request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/articles", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("missing required fields", func(t *testing.T) {
		reqBody := CreateArticleRequest{
			Title: "Test Article",
			// Content and AuthorID missing
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/articles", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestArticleHandler_GetArticle(t *testing.T) {
	mockUsecase := NewMockArticleUsecase()
	handler := NewArticleHandler(mockUsecase)
	router := setupArticleTestRouter(handler)

	// Create a test article
	testArticle, _ := mockUsecase.CreateArticle(context.Background(), "Test Article", "Test Content", 1)

	t.Run("get existing article", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/articles/1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response entity.Article
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("failed to unmarshal response: %v", err)
		}

		if response.ID != testArticle.ID {
			t.Errorf("expected article ID %d, got %d", testArticle.ID, response.ID)
		}
	})

	t.Run("get non-existent article", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/articles/999", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("invalid article ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/articles/invalid", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestArticleHandler_GetAllArticles(t *testing.T) {
	mockUsecase := NewMockArticleUsecase()
	handler := NewArticleHandler(mockUsecase)
	router := setupArticleTestRouter(handler)

	t.Run("get all articles when empty", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/articles", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response []*entity.Article
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("failed to unmarshal response: %v", err)
		}

		if len(response) != 0 {
			t.Errorf("expected 0 articles, got %d", len(response))
		}
	})

	t.Run("get all articles with data", func(t *testing.T) {
		// Create test articles
		mockUsecase.CreateArticle(context.Background(), "Article 1", "Content 1", 1)
		mockUsecase.CreateArticle(context.Background(), "Article 2", "Content 2", 1)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/articles", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response []*entity.Article
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("failed to unmarshal response: %v", err)
		}

		if len(response) != 2 {
			t.Errorf("expected 2 articles, got %d", len(response))
		}
	})
}

func TestArticleHandler_GetPublishedArticles(t *testing.T) {
	mockUsecase := NewMockArticleUsecase()
	handler := NewArticleHandler(mockUsecase)
	router := setupArticleTestRouter(handler)

	// Create test articles
	_, _ = mockUsecase.CreateArticle(context.Background(), "Draft Article", "Content 1", 1)
	article2, _ := mockUsecase.CreateArticle(context.Background(), "Published Article", "Content 2", 1)

	// Publish one article
	mockUsecase.PublishArticle(context.Background(), article2.ID)

	t.Run("get only published articles", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/articles/published", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response []*entity.Article
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("failed to unmarshal response: %v", err)
		}

		if len(response) != 1 {
			t.Errorf("expected 1 published article, got %d", len(response))
		}

		if response[0].Status != "published" {
			t.Errorf("expected status 'published', got %s", response[0].Status)
		}
	})
}

func TestArticleHandler_UpdateArticle(t *testing.T) {
	mockUsecase := NewMockArticleUsecase()
	handler := NewArticleHandler(mockUsecase)
	router := setupArticleTestRouter(handler)

	// Create a test article
	mockUsecase.CreateArticle(context.Background(), "Original Title", "Original Content", 1)

	t.Run("successful article update", func(t *testing.T) {
		reqBody := UpdateArticleRequest{
			Title:   "Updated Title",
			Content: "Updated Content",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/articles/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response entity.Article
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("failed to unmarshal response: %v", err)
		}

		if response.Title != "Updated Title" {
			t.Errorf("expected title 'Updated Title', got %s", response.Title)
		}
		if response.Content != "Updated Content" {
			t.Errorf("expected content 'Updated Content', got %s", response.Content)
		}
	})

	t.Run("update non-existent article", func(t *testing.T) {
		reqBody := UpdateArticleRequest{
			Title: "Updated Title",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/articles/999", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("invalid article ID", func(t *testing.T) {
		reqBody := UpdateArticleRequest{
			Title: "Updated Title",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/articles/invalid", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestArticleHandler_PublishArticle(t *testing.T) {
	mockUsecase := NewMockArticleUsecase()
	handler := NewArticleHandler(mockUsecase)
	router := setupArticleTestRouter(handler)

	// Create a test article
	mockUsecase.CreateArticle(context.Background(), "Test Article", "Test Content", 1)

	t.Run("successful article publishing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/articles/1/publish", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response entity.Article
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("failed to unmarshal response: %v", err)
		}

		if response.Status != "published" {
			t.Errorf("expected status 'published', got %s", response.Status)
		}
	})

	t.Run("publish non-existent article", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/articles/999/publish", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}

func TestArticleHandler_DeleteArticle(t *testing.T) {
	mockUsecase := NewMockArticleUsecase()
	handler := NewArticleHandler(mockUsecase)
	router := setupArticleTestRouter(handler)

	// Create a test article
	mockUsecase.CreateArticle(context.Background(), "Test Article", "Test Content", 1)

	t.Run("successful article deletion", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/articles/1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
		}

		// Verify article is deleted by trying to get it
		getReq := httptest.NewRequest(http.MethodGet, "/api/v1/articles/1", nil)
		getW := httptest.NewRecorder()
		router.ServeHTTP(getW, getReq)

		if getW.Code != http.StatusNotFound {
			t.Errorf("expected article to be deleted, but got status %d", getW.Code)
		}
	})

	t.Run("delete non-existent article", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/articles/999", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("invalid article ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/articles/invalid", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}