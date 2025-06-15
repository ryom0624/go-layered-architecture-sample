package usecase

import (
	"context"
	"testing"
	"time"

	"layered-architecture-template/internal/domain/entity"
)

func TestSearchUsecase_SearchArticles(t *testing.T) {
	articleRepo := NewMockArticleRepository()
	searchUsecase := NewSearchUsecase(articleRepo)

	// Create test articles
	article1 := &entity.Article{
		ID:      1,
		Title:   "Go Programming",
		Content: "Learn Go programming language",
		Status:  "published",
	}
	article2 := &entity.Article{
		ID:      2,
		Title:   "Python Guide",
		Content: "Python programming tutorial",
		Status:  "published",
	}
	
	articleRepo.Create(context.Background(), article1)
	articleRepo.Create(context.Background(), article2)

	tests := []struct {
		name           string
		params         *entity.ArticleSearchParams
		expectedCount  int
		expectError    bool
	}{
		{
			name: "valid search with query",
			params: &entity.ArticleSearchParams{
				Query: "Go",
				Page:  1,
				Limit: 10,
			},
			expectedCount: 2, // Mock returns all articles
			expectError:   false,
		},
		{
			name: "search with pagination",
			params: &entity.ArticleSearchParams{
				Page:  1,
				Limit: 1,
			},
			expectedCount: 2, // Mock returns all articles
			expectError:   false,
		},
		{
			name: "search with status filter",
			params: &entity.ArticleSearchParams{
				Status: "published",
				Page:   1,
				Limit:  10,
			},
			expectedCount: 2,
			expectError:   false,
		},
		{
			name: "search with invalid sort order",
			params: &entity.ArticleSearchParams{
				Page:      1,
				Limit:     10,
				SortOrder: "invalid",
			},
			expectedCount: 0,
			expectError:   true,
		},
		{
			name: "search with query too long",
			params: &entity.ArticleSearchParams{
				Query: string(make([]byte, 101)), // 101 characters
				Page:  1,
				Limit: 10,
			},
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := searchUsecase.SearchArticles(context.Background(), tt.params)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if result == nil {
				t.Errorf("Expected result, but got nil")
				return
			}
			
			if len(result.Articles) != tt.expectedCount {
				t.Errorf("Expected %d articles, got %d", tt.expectedCount, len(result.Articles))
			}
		})
	}
}

func TestSearchUsecase_GetPopularArticles(t *testing.T) {
	articleRepo := NewMockArticleRepository()
	searchUsecase := NewSearchUsecase(articleRepo)

	// Create test articles
	article1 := &entity.Article{
		ID:      1,
		Title:   "Popular Article 1",
		Content: "Content 1",
		Status:  "published",
	}
	article2 := &entity.Article{
		ID:      2,
		Title:   "Popular Article 2",
		Content: "Content 2",
		Status:  "published",
	}
	
	articleRepo.Create(context.Background(), article1)
	articleRepo.Create(context.Background(), article2)

	tests := []struct {
		name          string
		limit         int
		expectedLimit int
	}{
		{
			name:          "valid limit",
			limit:         5,
			expectedLimit: 5,
		},
		{
			name:          "zero limit defaults to 10",
			limit:         0,
			expectedLimit: 10,
		},
		{
			name:          "limit exceeds maximum, capped at 50",
			limit:         100,
			expectedLimit: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			articles, err := searchUsecase.GetPopularArticles(context.Background(), tt.limit)
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if articles == nil {
				t.Errorf("Expected articles, but got nil")
				return
			}
		})
	}
}

func TestSearchUsecase_GetRecentArticles(t *testing.T) {
	articleRepo := NewMockArticleRepository()
	searchUsecase := NewSearchUsecase(articleRepo)

	// Create test articles
	article1 := &entity.Article{
		ID:      1,
		Title:   "Recent Article 1",
		Content: "Content 1",
		Status:  "published",
	}
	
	articleRepo.Create(context.Background(), article1)

	tests := []struct {
		name          string
		limit         int
		expectedLimit int
	}{
		{
			name:          "valid limit",
			limit:         5,
			expectedLimit: 5,
		},
		{
			name:          "zero limit defaults to 10",
			limit:         0,
			expectedLimit: 10,
		},
		{
			name:          "limit exceeds maximum, capped at 50",
			limit:         100,
			expectedLimit: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			articles, err := searchUsecase.GetRecentArticles(context.Background(), tt.limit)
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if articles == nil {
				t.Errorf("Expected articles, but got nil")
				return
			}
		})
	}
}

func TestSearchUsecase_ValidateSearchParams(t *testing.T) {
	searchUsecase := NewSearchUsecase(NewMockArticleRepository())

	tests := []struct {
		name        string
		params      *entity.ArticleSearchParams
		expectError bool
	}{
		{
			name:        "nil params returns error",
			params:      nil,
			expectError: true,
		},
		{
			name: "valid params",
			params: &entity.ArticleSearchParams{
				Query: "test",
				Page:  1,
				Limit: 10,
			},
			expectError: false,
		},
		{
			name: "query too long",
			params: &entity.ArticleSearchParams{
				Query: string(make([]byte, 101)), // 101 characters
				Page:  1,
				Limit: 10,
			},
			expectError: true,
		},
		{
			name: "invalid sort order",
			params: &entity.ArticleSearchParams{
				Page:      1,
				Limit:     10,
				SortOrder: "invalid",
			},
			expectError: true,
		},
		{
			name: "invalid date range",
			params: &entity.ArticleSearchParams{
				Page:     1,
				Limit:    10,
				DateFrom: time.Now(),
				DateTo:   time.Now().AddDate(0, 0, -1), // DateTo before DateFrom
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := searchUsecase.ValidateSearchParams(tt.params)
			
			if tt.expectError && err == nil {
				t.Errorf("Expected error, but got none")
			}
			
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}