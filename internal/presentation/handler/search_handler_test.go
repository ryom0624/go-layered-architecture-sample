package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"layered-architecture-template/internal/domain/entity"

	"github.com/gin-gonic/gin"
)

// MockSearchUsecase implements usecase.SearchUsecase for testing
type MockSearchUsecase struct {
	articles []*entity.Article
}

func NewMockSearchUsecase() *MockSearchUsecase {
	return &MockSearchUsecase{
		articles: []*entity.Article{
			{
				ID:      1,
				Title:   "Go Programming",
				Content: "Learn Go programming language",
				Status:  "published",
			},
			{
				ID:      2,
				Title:   "Python Guide",
				Content: "Python programming tutorial",
				Status:  "published",
			},
		},
	}
}

func (m *MockSearchUsecase) SearchArticles(ctx context.Context, params *entity.ArticleSearchParams) (*entity.ArticleSearchResult, error) {
	if params.Query != "" && len(params.Query) > 100 {
		return nil, &entity.ValidationError{Field: "query", Message: "query must be 100 characters or less"}
	}
	
	return &entity.ArticleSearchResult{
		Articles:   m.articles,
		Total:      int64(len(m.articles)),
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: 1,
	}, nil
}

func (m *MockSearchUsecase) GetPopularArticles(ctx context.Context, limit int) ([]*entity.Article, error) {
	if limit > len(m.articles) {
		return m.articles, nil
	}
	return m.articles[:limit], nil
}

func (m *MockSearchUsecase) GetRecentArticles(ctx context.Context, limit int) ([]*entity.Article, error) {
	if limit > len(m.articles) {
		return m.articles, nil
	}
	return m.articles[:limit], nil
}

func (m *MockSearchUsecase) ValidateSearchParams(params *entity.ArticleSearchParams) error {
	if params == nil {
		return &entity.ValidationError{Field: "params", Message: "search parameters are required"}
	}
	return params.Validate()
}

func setupSearchTestRouter(handler *SearchHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	api := router.Group("/api/v1")
	{
		api.GET("/search", handler.SearchArticles)
		
		articles := api.Group("/articles")
		{
			articles.GET("/popular", handler.GetPopularArticles)
			articles.GET("/recent", handler.GetRecentArticles)
		}
	}

	return router
}

func TestSearchHandler_SearchArticles(t *testing.T) {
	mockUsecase := NewMockSearchUsecase()
	handler := NewSearchHandler(mockUsecase)
	router := setupSearchTestRouter(handler)

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		checkResponse  bool
	}{
		{
			name:           "search without parameters",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			checkResponse:  true,
		},
		{
			name:           "search with query parameter",
			queryParams:    "?query=Go",
			expectedStatus: http.StatusOK,
			checkResponse:  true,
		},
		{
			name:           "search with pagination",
			queryParams:    "?page=1&limit=10",
			expectedStatus: http.StatusOK,
			checkResponse:  true,
		},
		{
			name:           "search with filters",
			queryParams:    "?query=Go&status=published&author_id=1",
			expectedStatus: http.StatusOK,
			checkResponse:  true,
		},
		{
			name:           "search with date range",
			queryParams:    "?date_from=2023-01-01&date_to=2023-12-31",
			expectedStatus: http.StatusOK,
			checkResponse:  true,
		},
		{
			name:           "search with sorting",
			queryParams:    "?sort_by=title&sort_order=asc",
			expectedStatus: http.StatusOK,
			checkResponse:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/v1/search"+tt.queryParams, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResponse && w.Code == http.StatusOK {
				var result entity.ArticleSearchResult
				if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				
				if len(result.Articles) == 0 {
					t.Errorf("Expected articles in response, got none")
				}
			}
		})
	}
}

func TestSearchHandler_GetPopularArticles(t *testing.T) {
	mockUsecase := NewMockSearchUsecase()
	handler := NewSearchHandler(mockUsecase)
	router := setupSearchTestRouter(handler)

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		checkResponse  bool
	}{
		{
			name:           "get popular articles without limit",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			checkResponse:  true,
		},
		{
			name:           "get popular articles with limit",
			queryParams:    "?limit=5",
			expectedStatus: http.StatusOK,
			checkResponse:  true,
		},
		{
			name:           "get popular articles with invalid limit",
			queryParams:    "?limit=invalid",
			expectedStatus: http.StatusOK,
			checkResponse:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/v1/articles/popular"+tt.queryParams, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResponse && w.Code == http.StatusOK {
				var articles []*entity.Article
				if err := json.Unmarshal(w.Body.Bytes(), &articles); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
			}
		})
	}
}

func TestSearchHandler_GetRecentArticles(t *testing.T) {
	mockUsecase := NewMockSearchUsecase()
	handler := NewSearchHandler(mockUsecase)
	router := setupSearchTestRouter(handler)

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		checkResponse  bool
	}{
		{
			name:           "get recent articles without limit",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			checkResponse:  true,
		},
		{
			name:           "get recent articles with limit",
			queryParams:    "?limit=5",
			expectedStatus: http.StatusOK,
			checkResponse:  true,
		},
		{
			name:           "get recent articles with invalid limit",
			queryParams:    "?limit=invalid",
			expectedStatus: http.StatusOK,
			checkResponse:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/v1/articles/recent"+tt.queryParams, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResponse && w.Code == http.StatusOK {
				var articles []*entity.Article
				if err := json.Unmarshal(w.Body.Bytes(), &articles); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
			}
		})
	}
}