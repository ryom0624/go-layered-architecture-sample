package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"layered-architecture-template/internal/domain/entity"
	"layered-architecture-template/internal/domain/repository"
)

// MockArticleRepository implements repository.ArticleRepository for testing
type MockArticleRepository struct {
	articles map[uint]*entity.Article
	nextID   uint
}

func NewMockArticleRepository() *MockArticleRepository {
	return &MockArticleRepository{
		articles: make(map[uint]*entity.Article),
		nextID:   1,
	}
}

func (m *MockArticleRepository) Create(ctx context.Context, article *entity.Article) error {
	if article == nil {
		return errors.New("article cannot be nil")
	}
	article.ID = m.nextID
	m.nextID++
	m.articles[article.ID] = article
	return nil
}

func (m *MockArticleRepository) CreateWithTx(ctx context.Context, tx repository.Transaction, article *entity.Article) error {
	return m.Create(ctx, article)
}

func (m *MockArticleRepository) GetByID(ctx context.Context, id uint) (*entity.Article, error) {
	if article, exists := m.articles[id]; exists {
		return article, nil
	}
	return nil, errors.New("article not found")
}

func (m *MockArticleRepository) GetAll(ctx context.Context) ([]*entity.Article, error) {
	var articles []*entity.Article
	for _, article := range m.articles {
		articles = append(articles, article)
	}
	return articles, nil
}

func (m *MockArticleRepository) GetByStatus(ctx context.Context, status string) ([]*entity.Article, error) {
	var articles []*entity.Article
	for _, article := range m.articles {
		if article.Status == status {
			articles = append(articles, article)
		}
	}
	return articles, nil
}

func (m *MockArticleRepository) Update(ctx context.Context, article *entity.Article) error {
	if article == nil || article.ID == 0 {
		return errors.New("invalid article")
	}
	if _, exists := m.articles[article.ID]; !exists {
		return errors.New("article not found")
	}
	m.articles[article.ID] = article
	return nil
}

func (m *MockArticleRepository) UpdateWithTx(ctx context.Context, tx repository.Transaction, article *entity.Article) error {
	return m.Update(ctx, article)
}

func (m *MockArticleRepository) Delete(ctx context.Context, id uint) error {
	if _, exists := m.articles[id]; exists {
		delete(m.articles, id)
		return nil
	}
	return errors.New("article not found")
}

func (m *MockArticleRepository) DeleteWithTx(ctx context.Context, tx repository.Transaction, id uint) error {
	return m.Delete(ctx, id)
}

// Search methods for search functionality
func (m *MockArticleRepository) Search(ctx context.Context, params *entity.ArticleSearchParams) (*entity.ArticleSearchResult, error) {
	var articles []*entity.Article
	for _, article := range m.articles {
		articles = append(articles, article)
	}
	
	// Simple implementation for testing
	totalPages := 1
	if len(articles) > params.Limit {
		totalPages = (len(articles) + params.Limit - 1) / params.Limit
	}
	
	return &entity.ArticleSearchResult{
		Articles:   articles,
		Total:      int64(len(articles)),
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

func (m *MockArticleRepository) SearchWithFilters(ctx context.Context, query string, filters map[string]interface{}) ([]*entity.Article, error) {
	var articles []*entity.Article
	for _, article := range m.articles {
		articles = append(articles, article)
	}
	return articles, nil
}

func (m *MockArticleRepository) GetByDateRange(ctx context.Context, from, to time.Time) ([]*entity.Article, error) {
	var articles []*entity.Article
	for _, article := range m.articles {
		articles = append(articles, article)
	}
	return articles, nil
}

func (m *MockArticleRepository) GetPopular(ctx context.Context, limit int) ([]*entity.Article, error) {
	var articles []*entity.Article
	for _, article := range m.articles {
		if article.Status == "published" {
			articles = append(articles, article)
		}
	}
	return articles, nil
}

func (m *MockArticleRepository) CountByFilters(ctx context.Context, filters map[string]interface{}) (int64, error) {
	return int64(len(m.articles)), nil
}

// MockTransactionManager implements repository.TransactionManager for testing
type MockTransactionManager struct{}

func (m *MockTransactionManager) WithTransaction(ctx context.Context, fn func(tx repository.Transaction) error) error {
	tx := &MockTransaction{}
	return fn(tx)
}

// MockTransaction implements repository.Transaction for testing
type MockTransaction struct{}

func (m *MockTransaction) Commit() error   { return nil }
func (m *MockTransaction) Rollback() error { return nil }
func (m *MockTransaction) GetDB() interface{} { return nil }

func TestArticleUsecase_CreateArticle(t *testing.T) {
	mockArticleRepo := NewMockArticleRepository()
	mockUserRepo := NewMockUserRepository()
	mockCategoryRepo := NewMockCategoryRepository()
	mockTagRepo := NewMockTagRepository()
	mockTxManager := &MockTransactionManager{}
	usecase := NewArticleUsecase(mockArticleRepo, mockUserRepo, mockCategoryRepo, mockTagRepo, mockTxManager)
	ctx := context.Background()

	// Create a test user first
	testUser := &entity.User{Name: "Test Author", Email: "author@example.com"}
	_ = mockUserRepo.Create(ctx, testUser)
	testUser.ID = 1

	t.Run("successful article creation", func(t *testing.T) {
		article, err := usecase.CreateArticle(ctx, "Test Title", "Test Content", 1)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if article == nil {
			t.Error("expected article to be created, got nil")
		}
		if article.Title != "Test Title" {
			t.Errorf("expected title 'Test Title', got %s", article.Title)
		}
		if article.Content != "Test Content" {
			t.Errorf("expected content 'Test Content', got %s", article.Content)
		}
		if article.AuthorID != 1 {
			t.Errorf("expected author ID 1, got %d", article.AuthorID)
		}
		if article.Status != "draft" {
			t.Errorf("expected status 'draft', got %s", article.Status)
		}
	})

	t.Run("empty title validation", func(t *testing.T) {
		article, err := usecase.CreateArticle(ctx, "", "Test Content", 1)
		if err == nil {
			t.Error("expected error for empty title")
		}
		if article != nil {
			t.Error("expected nil article for invalid input")
		}
		expectedError := "title is required"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("empty content validation", func(t *testing.T) {
		article, err := usecase.CreateArticle(ctx, "Test Title", "", 1)
		if err == nil {
			t.Error("expected error for empty content")
		}
		if article != nil {
			t.Error("expected nil article for invalid input")
		}
		expectedError := "content is required"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("zero author ID validation", func(t *testing.T) {
		article, err := usecase.CreateArticle(ctx, "Test Title", "Test Content", 0)
		if err == nil {
			t.Error("expected error for zero author ID")
		}
		if article != nil {
			t.Error("expected nil article for invalid input")
		}
		expectedError := "author ID is required"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("non-existent author validation", func(t *testing.T) {
		article, err := usecase.CreateArticle(ctx, "Test Title", "Test Content", 999)
		if err == nil {
			t.Error("expected error for non-existent author")
		}
		if article != nil {
			t.Error("expected nil article for invalid author")
		}
		expectedError := "author not found"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})
}

func TestArticleUsecase_GetArticle(t *testing.T) {
	mockArticleRepo := NewMockArticleRepository()
	mockUserRepo := NewMockUserRepository()
	mockCategoryRepo := NewMockCategoryRepository()
	mockTagRepo := NewMockTagRepository()
	mockTxManager := &MockTransactionManager{}
	usecase := NewArticleUsecase(mockArticleRepo, mockUserRepo, mockCategoryRepo, mockTagRepo, mockTxManager)
	ctx := context.Background()

	// Create a test article
	testUser := &entity.User{Name: "Test Author", Email: "author@example.com"}
	_ = mockUserRepo.Create(ctx, testUser)
	testUser.ID = 1
	createdArticle, _ := usecase.CreateArticle(ctx, "Test Article", "Test Content", 1)

	t.Run("get existing article", func(t *testing.T) {
		article, err := usecase.GetArticle(ctx, createdArticle.ID)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if article == nil {
			t.Error("expected article to be found, got nil")
		}
		if article.ID != createdArticle.ID {
			t.Errorf("expected article ID %d, got %d", createdArticle.ID, article.ID)
		}
	})

	t.Run("get non-existent article", func(t *testing.T) {
		article, err := usecase.GetArticle(ctx, 999)
		if err == nil {
			t.Error("expected error for non-existent article")
		}
		if article != nil {
			t.Error("expected nil article for non-existent ID")
		}
	})
}

func TestArticleUsecase_UpdateArticle(t *testing.T) {
	mockArticleRepo := NewMockArticleRepository()
	mockUserRepo := NewMockUserRepository()
	mockTxManager := &MockTransactionManager{}
	usecase := NewArticleUsecase(mockArticleRepo, mockUserRepo, NewMockCategoryRepository(), NewMockTagRepository(), mockTxManager)
	ctx := context.Background()

	// Create a test article
	testUser := &entity.User{Name: "Test Author", Email: "author@example.com"}
	_ = mockUserRepo.Create(ctx, testUser)
	testUser.ID = 1
	createdArticle, _ := usecase.CreateArticle(ctx, "Original Title", "Original Content", 1)

	t.Run("successful article update", func(t *testing.T) {
		updatedArticle, err := usecase.UpdateArticle(ctx, createdArticle.ID, "Updated Title", "Updated Content")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if updatedArticle.Title != "Updated Title" {
			t.Errorf("expected title 'Updated Title', got %s", updatedArticle.Title)
		}
		if updatedArticle.Content != "Updated Content" {
			t.Errorf("expected content 'Updated Content', got %s", updatedArticle.Content)
		}
	})

	t.Run("partial update with empty title", func(t *testing.T) {
		originalTitle := createdArticle.Title
		updatedArticle, err := usecase.UpdateArticle(ctx, createdArticle.ID, "", "New Content Only")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if updatedArticle.Title != originalTitle {
			t.Errorf("expected title to remain '%s', got %s", originalTitle, updatedArticle.Title)
		}
		if updatedArticle.Content != "New Content Only" {
			t.Errorf("expected content 'New Content Only', got %s", updatedArticle.Content)
		}
	})

	t.Run("update non-existent article", func(t *testing.T) {
		article, err := usecase.UpdateArticle(ctx, 999, "Title", "Content")
		if err == nil {
			t.Error("expected error for non-existent article")
		}
		if article != nil {
			t.Error("expected nil article for non-existent ID")
		}
	})

	t.Run("update with zero ID", func(t *testing.T) {
		_, err := usecase.UpdateArticle(ctx, 0, "Title", "Content")
		if err == nil {
			t.Error("expected error for zero ID")
		}
		expectedError := "article ID is required"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})
}

func TestArticleUsecase_PublishArticle(t *testing.T) {
	mockArticleRepo := NewMockArticleRepository()
	mockUserRepo := NewMockUserRepository()
	mockTxManager := &MockTransactionManager{}
	usecase := NewArticleUsecase(mockArticleRepo, mockUserRepo, NewMockCategoryRepository(), NewMockTagRepository(), mockTxManager)
	ctx := context.Background()

	// Create a test article
	testUser := &entity.User{Name: "Test Author", Email: "author@example.com"}
	_ = mockUserRepo.Create(ctx, testUser)
	testUser.ID = 1
	createdArticle, _ := usecase.CreateArticle(ctx, "Test Article", "Test Content", 1)

	t.Run("successful article publishing", func(t *testing.T) {
		publishedArticle, err := usecase.PublishArticle(ctx, createdArticle.ID)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if publishedArticle.Status != "published" {
			t.Errorf("expected status 'published', got %s", publishedArticle.Status)
		}
	})

	t.Run("publish already published article", func(t *testing.T) {
		// First publish
		usecase.PublishArticle(ctx, createdArticle.ID)

		// Try to publish again
		_, err := usecase.PublishArticle(ctx, createdArticle.ID)
		if err == nil {
			t.Error("expected error for already published article")
		}
		expectedError := "article is already published"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("publish non-existent article", func(t *testing.T) {
		_, err := usecase.PublishArticle(ctx, 999)
		if err == nil {
			t.Error("expected error for non-existent article")
		}
	})
}

func TestArticleUsecase_GetPublishedArticles(t *testing.T) {
	mockArticleRepo := NewMockArticleRepository()
	mockUserRepo := NewMockUserRepository()
	mockTxManager := &MockTransactionManager{}
	usecase := NewArticleUsecase(mockArticleRepo, mockUserRepo, NewMockCategoryRepository(), NewMockTagRepository(), mockTxManager)
	ctx := context.Background()

	// Create test user
	testUser := &entity.User{Name: "Test Author", Email: "author@example.com"}
	_ = mockUserRepo.Create(ctx, testUser)
	testUser.ID = 1

	// Create test articles
	_, _ = usecase.CreateArticle(ctx, "Draft Article", "Content 1", 1)
	article2, _ := usecase.CreateArticle(ctx, "Published Article 1", "Content 2", 1)
	article3, _ := usecase.CreateArticle(ctx, "Published Article 2", "Content 3", 1)

	// Publish some articles
	usecase.PublishArticle(ctx, article2.ID)
	usecase.PublishArticle(ctx, article3.ID)

	t.Run("get only published articles", func(t *testing.T) {
		publishedArticles, err := usecase.GetPublishedArticles(ctx)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(publishedArticles) != 2 {
			t.Errorf("expected 2 published articles, got %d", len(publishedArticles))
		}
		for _, article := range publishedArticles {
			if article.Status != "published" {
				t.Errorf("expected all articles to be published, found status %s", article.Status)
			}
		}
	})
}

func TestArticleUsecase_DeleteArticle(t *testing.T) {
	mockArticleRepo := NewMockArticleRepository()
	mockUserRepo := NewMockUserRepository()
	mockTxManager := &MockTransactionManager{}
	usecase := NewArticleUsecase(mockArticleRepo, mockUserRepo, NewMockCategoryRepository(), NewMockTagRepository(), mockTxManager)
	ctx := context.Background()

	// Create a test article
	testUser := &entity.User{Name: "Test Author", Email: "author@example.com"}
	_ = mockUserRepo.Create(ctx, testUser)
	testUser.ID = 1
	createdArticle, _ := usecase.CreateArticle(ctx, "Test Article", "Test Content", 1)

	t.Run("successful article deletion", func(t *testing.T) {
		err := usecase.DeleteArticle(ctx, createdArticle.ID)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Verify the article is deleted
		_, err = usecase.GetArticle(ctx, createdArticle.ID)
		if err == nil {
			t.Error("expected error when getting deleted article")
		}
	})

	t.Run("delete with zero ID", func(t *testing.T) {
		err := usecase.DeleteArticle(ctx, 0)
		if err == nil {
			t.Error("expected error for zero ID")
		}
		expectedError := "article ID is required"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("delete non-existent article", func(t *testing.T) {
		err := usecase.DeleteArticle(ctx, 999)
		if err == nil {
			t.Error("expected error for non-existent article")
		}
	})
}
