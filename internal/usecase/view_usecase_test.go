package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"layered-architecture-template/internal/domain/entity"
)

type MockViewRepository struct {
	views            map[uint]*entity.ArticleView
	readingHistories map[string]*entity.UserReadingHistory // key: "userID-articleID"
	nextViewID       uint
}

func NewMockViewRepository() *MockViewRepository {
	return &MockViewRepository{
		views:            make(map[uint]*entity.ArticleView),
		readingHistories: make(map[string]*entity.UserReadingHistory),
		nextViewID:       1,
	}
}

func (m *MockViewRepository) CreateView(ctx context.Context, view *entity.ArticleView) error {
	if view == nil {
		return errors.New("view cannot be nil")
	}
	view.ID = m.nextViewID
	m.nextViewID++
	m.views[view.ID] = view
	return nil
}

func (m *MockViewRepository) GetViewsByArticle(ctx context.Context, articleID uint) ([]*entity.ArticleView, error) {
	var views []*entity.ArticleView
	for _, view := range m.views {
		if view.ArticleID == articleID {
			views = append(views, view)
		}
	}
	return views, nil
}

func (m *MockViewRepository) GetViewsByUser(ctx context.Context, userID uint) ([]*entity.ArticleView, error) {
	var views []*entity.ArticleView
	for _, view := range m.views {
		if view.UserID != nil && *view.UserID == userID {
			views = append(views, view)
		}
	}
	return views, nil
}

func (m *MockViewRepository) GetViewsByIPAndArticle(ctx context.Context, ip string, articleID uint, since time.Time) ([]*entity.ArticleView, error) {
	var views []*entity.ArticleView
	for _, view := range m.views {
		if view.IPAddress == ip && view.ArticleID == articleID && view.ViewedAt.After(since) {
			views = append(views, view)
		}
	}
	return views, nil
}

func (m *MockViewRepository) UpdateReadingTime(ctx context.Context, viewID uint, readingTime int) error {
	if view, exists := m.views[viewID]; exists {
		view.ReadingTime = readingTime
		return nil
	}
	return errors.New("view not found")
}

func (m *MockViewRepository) CreateOrUpdateReadingHistory(ctx context.Context, history *entity.UserReadingHistory) error {
	if history == nil {
		return errors.New("history cannot be nil")
	}
	key := getHistoryKey(history.UserID, history.ArticleID)
	if existing, exists := m.readingHistories[key]; exists {
		existing.TotalReadingTime += history.TotalReadingTime
		existing.ReadingProgress = history.ReadingProgress
		existing.IsCompleted = history.IsCompleted
		existing.LastViewedAt = history.LastViewedAt
		existing.ViewCount++
		if history.IsCompleted && existing.CompletedAt == nil {
			existing.CompletedAt = history.CompletedAt
		}
	} else {
		if history.ID == 0 {
			history.ID = m.nextViewID
			m.nextViewID++
		}
		m.readingHistories[key] = history
	}
	return nil
}

func (m *MockViewRepository) GetReadingHistory(ctx context.Context, userID uint, articleID uint) (*entity.UserReadingHistory, error) {
	key := getHistoryKey(userID, articleID)
	if history, exists := m.readingHistories[key]; exists {
		return history, nil
	}
	return nil, errors.New("reading history not found")
}

func (m *MockViewRepository) GetUserReadingHistories(ctx context.Context, userID uint, limit, offset int) ([]*entity.UserReadingHistory, error) {
	var histories []*entity.UserReadingHistory
	for _, history := range m.readingHistories {
		if history.UserID == userID {
			histories = append(histories, history)
		}
	}
	
	// Simple pagination simulation
	start := offset
	end := offset + limit
	if start >= len(histories) {
		return []*entity.UserReadingHistory{}, nil
	}
	if end > len(histories) {
		end = len(histories)
	}
	return histories[start:end], nil
}

func (m *MockViewRepository) UpdateReadingProgress(ctx context.Context, userID uint, articleID uint, progress float32, readingTime int, isCompleted bool) error {
	key := getHistoryKey(userID, articleID)
	if history, exists := m.readingHistories[key]; exists {
		history.ReadingProgress = progress
		history.TotalReadingTime += readingTime
		history.IsCompleted = isCompleted
		history.LastViewedAt = time.Now()
		history.ViewCount++
		if isCompleted && history.CompletedAt == nil {
			now := time.Now()
			history.CompletedAt = &now
		}
		return nil
	}
	return errors.New("reading history not found")
}

func (m *MockViewRepository) GetPopularArticles(ctx context.Context, limit int, since time.Time) ([]*entity.Article, error) {
	// Mock implementation - return empty for simplicity
	return []*entity.Article{}, nil
}

func (m *MockViewRepository) GetTrendingArticles(ctx context.Context, limit int, since time.Time) ([]*entity.Article, error) {
	// Mock implementation - return empty for simplicity
	return []*entity.Article{}, nil
}

func getHistoryKey(userID, articleID uint) string {
	return string(rune(userID)) + "-" + string(rune(articleID))
}


func TestViewUsecase_TrackView(t *testing.T) {
	viewRepo := NewMockViewRepository()
	articleRepo := NewMockArticleRepository()
	usecase := NewViewUsecase(viewRepo, articleRepo)
	ctx := context.Background()

	// Create a test article
	article := &entity.Article{
		Title:    "Test Article",
		Content:  "Test content",
		AuthorID: 1,
		Status:   "published",
	}
	articleRepo.Create(ctx, article)

	t.Run("successful view tracking for authenticated user", func(t *testing.T) {
		userID := uint(1)
		err := usecase.TrackView(ctx, article.ID, &userID, "192.168.1.1", "test-agent")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Verify view was created
		views, _ := viewRepo.GetViewsByArticle(ctx, article.ID)
		if len(views) != 1 {
			t.Errorf("expected 1 view, got %d", len(views))
		}

		// Verify reading history was created
		history, err := viewRepo.GetReadingHistory(ctx, userID, article.ID)
		if err != nil {
			t.Errorf("expected reading history to be created, got error: %v", err)
		}
		if history.UserID != userID {
			t.Errorf("expected user ID %d, got %d", userID, history.UserID)
		}
	})

	t.Run("successful view tracking for anonymous user", func(t *testing.T) {
		err := usecase.TrackView(ctx, article.ID, nil, "192.168.1.2", "test-agent")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Verify view was created
		views, _ := viewRepo.GetViewsByArticle(ctx, article.ID)
		if len(views) != 2 { // Previous test + this one
			t.Errorf("expected 2 views, got %d", len(views))
		}
	})

	t.Run("invalid article ID", func(t *testing.T) {
		err := usecase.TrackView(ctx, 0, nil, "192.168.1.1", "test-agent")
		if err == nil {
			t.Error("expected error for invalid article ID")
		}
		expectedError := "article ID is required"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("missing IP address", func(t *testing.T) {
		err := usecase.TrackView(ctx, article.ID, nil, "", "test-agent")
		if err == nil {
			t.Error("expected error for missing IP address")
		}
		expectedError := "IP address is required"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("article not found", func(t *testing.T) {
		err := usecase.TrackView(ctx, 999, nil, "192.168.1.1", "test-agent")
		if err == nil {
			t.Error("expected error for non-existent article")
		}
		expectedError := "article not found"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("unpublished article", func(t *testing.T) {
		// Create unpublished article
		unpublished := &entity.Article{
			Title:    "Draft Article",
			Content:  "Draft content",
			AuthorID: 1,
			Status:   "draft",
		}
		articleRepo.Create(ctx, unpublished)

		err := usecase.TrackView(ctx, unpublished.ID, nil, "192.168.1.1", "test-agent")
		if err == nil {
			t.Error("expected error for unpublished article")
		}
		expectedError := "article is not published"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})
}

func TestViewUsecase_TrackReadingProgress(t *testing.T) {
	viewRepo := NewMockViewRepository()
	articleRepo := NewMockArticleRepository()
	usecase := NewViewUsecase(viewRepo, articleRepo)
	ctx := context.Background()

	// Create a test article
	article := &entity.Article{
		Title:    "Test Article",
		Content:  "Test content",
		AuthorID: 1,
		Status:   "published",
	}
	articleRepo.Create(ctx, article)

	// Create initial reading history
	history := &entity.UserReadingHistory{
		UserID:        1,
		ArticleID:     article.ID,
		FirstViewedAt: time.Now(),
		LastViewedAt:  time.Now(),
	}
	viewRepo.CreateOrUpdateReadingHistory(ctx, history)

	t.Run("successful progress tracking", func(t *testing.T) {
		err := usecase.TrackReadingProgress(ctx, 1, article.ID, 50.0, 300)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Verify progress was updated
		updatedHistory, _ := viewRepo.GetReadingHistory(ctx, 1, article.ID)
		if updatedHistory.ReadingProgress != 50.0 {
			t.Errorf("expected progress 50.0, got %f", updatedHistory.ReadingProgress)
		}
	})

	t.Run("completion tracking", func(t *testing.T) {
		err := usecase.TrackReadingProgress(ctx, 1, article.ID, 95.0, 600)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Verify completion was detected
		updatedHistory, _ := viewRepo.GetReadingHistory(ctx, 1, article.ID)
		if !updatedHistory.IsCompleted {
			t.Error("expected article to be marked as completed")
		}
	})

	t.Run("invalid progress range", func(t *testing.T) {
		err := usecase.TrackReadingProgress(ctx, 1, article.ID, 150.0, 300)
		if err == nil {
			t.Error("expected error for invalid progress")
		}
		expectedError := "progress must be between 0 and 100"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("negative reading time", func(t *testing.T) {
		err := usecase.TrackReadingProgress(ctx, 1, article.ID, 50.0, -100)
		if err == nil {
			t.Error("expected error for negative reading time")
		}
		expectedError := "reading time must be non-negative"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})
}

func TestViewUsecase_GetUserReadingHistories(t *testing.T) {
	viewRepo := NewMockViewRepository()
	articleRepo := NewMockArticleRepository()
	usecase := NewViewUsecase(viewRepo, articleRepo)
	ctx := context.Background()

	// Create test reading histories
	for i := 1; i <= 3; i++ {
		history := &entity.UserReadingHistory{
			UserID:        1,
			ArticleID:     uint(i),
			FirstViewedAt: time.Now(),
			LastViewedAt:  time.Now(),
		}
		viewRepo.CreateOrUpdateReadingHistory(ctx, history)
	}

	t.Run("get reading histories with default pagination", func(t *testing.T) {
		histories, err := usecase.GetUserReadingHistories(ctx, 1, 1, 20)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(histories) != 3 {
			t.Errorf("expected 3 histories, got %d", len(histories))
		}
	})

	t.Run("invalid user ID", func(t *testing.T) {
		_, err := usecase.GetUserReadingHistories(ctx, 0, 1, 20)
		if err == nil {
			t.Error("expected error for invalid user ID")
		}
		expectedError := "user ID is required"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("pagination defaults", func(t *testing.T) {
		histories, err := usecase.GetUserReadingHistories(ctx, 1, 0, 0)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		// Should apply default page=1, limit=20
		if len(histories) != 3 {
			t.Errorf("expected 3 histories, got %d", len(histories))
		}
	})
}