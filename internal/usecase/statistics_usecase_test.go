package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"layered-architecture-template/internal/domain/entity"
	"layered-architecture-template/internal/domain/repository"
)

type MockStatisticsRepository struct {
	articleStats  map[uint]*entity.ArticleStatistics
	dailyStats    map[string]*entity.DailyStatistics // key: date string
	nextStatsID   uint
}

func NewMockStatisticsRepository() *MockStatisticsRepository {
	return &MockStatisticsRepository{
		articleStats: make(map[uint]*entity.ArticleStatistics),
		dailyStats:   make(map[string]*entity.DailyStatistics),
		nextStatsID:  1,
	}
}

func (m *MockStatisticsRepository) CreateOrUpdateArticleStatistics(ctx context.Context, stats *entity.ArticleStatistics) error {
	if stats == nil {
		return errors.New("stats cannot be nil")
	}
	if existing, exists := m.articleStats[stats.ArticleID]; exists {
		existing.TotalViews = stats.TotalViews
		existing.UniqueViews = stats.UniqueViews
		existing.AuthenticatedViews = stats.AuthenticatedViews
		existing.AnonymousViews = stats.AnonymousViews
		existing.AverageReadingTime = stats.AverageReadingTime
		existing.CompletionRate = stats.CompletionRate
		existing.TotalReadingTime = stats.TotalReadingTime
		existing.LastCalculatedAt = time.Now()
	} else {
		if stats.ID == 0 {
			stats.ID = m.nextStatsID
			m.nextStatsID++
		}
		m.articleStats[stats.ArticleID] = stats
	}
	return nil
}

func (m *MockStatisticsRepository) GetArticleStatistics(ctx context.Context, articleID uint) (*entity.ArticleStatistics, error) {
	if stats, exists := m.articleStats[articleID]; exists {
		return stats, nil
	}
	return nil, errors.New("article statistics not found")
}

func (m *MockStatisticsRepository) GetAllArticleStatistics(ctx context.Context, limit, offset int) ([]*entity.ArticleStatistics, error) {
	var stats []*entity.ArticleStatistics
	for _, stat := range m.articleStats {
		stats = append(stats, stat)
	}
	
	// Simple pagination simulation
	start := offset
	end := offset + limit
	if start >= len(stats) {
		return []*entity.ArticleStatistics{}, nil
	}
	if end > len(stats) {
		end = len(stats)
	}
	return stats[start:end], nil
}

func (m *MockStatisticsRepository) RecalculateArticleStatistics(ctx context.Context, articleID uint) error {
	// Simulate recalculation by creating/updating stats
	stats := &entity.ArticleStatistics{
		ArticleID:            articleID,
		TotalViews:           10,
		UniqueViews:          8,
		AuthenticatedViews:   6,
		AnonymousViews:       4,
		AverageReadingTime:   120.5,
		CompletionRate:       0.75,
		TotalReadingTime:     1205,
		LastCalculatedAt:     time.Now(),
	}
	return m.CreateOrUpdateArticleStatistics(ctx, stats)
}

func (m *MockStatisticsRepository) CreateDailyStatistics(ctx context.Context, stats *entity.DailyStatistics) error {
	if stats == nil {
		return errors.New("stats cannot be nil")
	}
	if stats.ID == 0 {
		stats.ID = m.nextStatsID
		m.nextStatsID++
	}
	key := stats.Date.Format("2006-01-02")
	m.dailyStats[key] = stats
	return nil
}

func (m *MockStatisticsRepository) GetDailyStatistics(ctx context.Context, date time.Time) (*entity.DailyStatistics, error) {
	key := date.Format("2006-01-02")
	if stats, exists := m.dailyStats[key]; exists {
		return stats, nil
	}
	return nil, errors.New("daily statistics not found")
}

func (m *MockStatisticsRepository) GetDailyStatisticsRange(ctx context.Context, startDate, endDate time.Time) ([]*entity.DailyStatistics, error) {
	var stats []*entity.DailyStatistics
	current := startDate
	for !current.After(endDate) {
		key := current.Format("2006-01-02")
		if stat, exists := m.dailyStats[key]; exists {
			stats = append(stats, stat)
		}
		current = current.AddDate(0, 0, 1)
	}
	return stats, nil
}

func (m *MockStatisticsRepository) UpdateDailyStatistics(ctx context.Context, stats *entity.DailyStatistics) error {
	if stats == nil || stats.ID == 0 {
		return errors.New("invalid daily statistics")
	}
	key := stats.Date.Format("2006-01-02")
	m.dailyStats[key] = stats
	return nil
}

func (m *MockStatisticsRepository) GetPlatformOverview(ctx context.Context) (*repository.PlatformOverview, error) {
	return &repository.PlatformOverview{
		TotalUsers:         100,
		TotalArticles:      50,
		TotalViews:         1000,
		TotalReadingTime:   50000,
		AverageReadingTime: 120.5,
		ActiveUsersToday:   25,
		NewUsersToday:      5,
		PublishedToday:     3,
		PopularToday:       []*entity.Article{},
	}, nil
}

func (m *MockStatisticsRepository) GetUserAnalytics(ctx context.Context, userID uint) (*repository.UserAnalytics, error) {
	return &repository.UserAnalytics{
		UserID:             userID,
		TotalReadingTime:   5000,
		ArticlesRead:       15,
		ArticlesCompleted:  12,
		AverageReadingTime: 333.3,
		CompletionRate:     0.8,
		FavoriteTopic:      "Technology",
		ReadingStreak:      7,
		LastActiveDate:     time.Now(),
		RecentReadings:     []*entity.UserReadingHistory{},
	}, nil
}

func TestStatisticsUsecase_GetArticleStatistics(t *testing.T) {
	statsRepo := NewMockStatisticsRepository()
	viewRepo := NewMockViewRepository()
	userRepo := NewMockUserRepository()
	articleRepo := NewMockArticleRepository()
	usecase := NewStatisticsUsecase(statsRepo, viewRepo, userRepo, articleRepo)
	ctx := context.Background()

	// Create a test article
	article := &entity.Article{
		Title:    "Test Article",
		Content:  "Test content",
		AuthorID: 1,
		Status:   "published",
	}
	articleRepo.Create(ctx, article)

	t.Run("get existing article statistics", func(t *testing.T) {
		// Create test statistics
		stats := &entity.ArticleStatistics{
			ArticleID:            article.ID,
			TotalViews:           100,
			UniqueViews:          80,
			AuthenticatedViews:   60,
			AnonymousViews:       40,
			AverageReadingTime:   150.0,
			CompletionRate:       0.8,
			TotalReadingTime:     15000,
			LastCalculatedAt:     time.Now(),
		}
		statsRepo.CreateOrUpdateArticleStatistics(ctx, stats)

		result, err := usecase.GetArticleStatistics(ctx, article.ID)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result.TotalViews != 100 {
			t.Errorf("expected total views 100, got %d", result.TotalViews)
		}
		if result.CompletionRate != 0.8 {
			t.Errorf("expected completion rate 0.8, got %f", result.CompletionRate)
		}
	})

	t.Run("get statistics for article without existing stats triggers recalculation", func(t *testing.T) {
		// Create another test article
		article2 := &entity.Article{
			Title:    "Test Article 2",
			Content:  "Test content 2",
			AuthorID: 1,
			Status:   "published",
		}
		articleRepo.Create(ctx, article2)

		result, err := usecase.GetArticleStatistics(ctx, article2.ID)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		// Should have been calculated by recalculation
		if result.TotalViews != 10 { // Mock returns 10
			t.Errorf("expected total views 10, got %d", result.TotalViews)
		}
	})

	t.Run("invalid article ID", func(t *testing.T) {
		_, err := usecase.GetArticleStatistics(ctx, 0)
		if err == nil {
			t.Error("expected error for invalid article ID")
		}
		expectedError := "article ID is required"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("article not found", func(t *testing.T) {
		_, err := usecase.GetArticleStatistics(ctx, 999)
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

		_, err := usecase.GetArticleStatistics(ctx, unpublished.ID)
		if err == nil {
			t.Error("expected error for unpublished article")
		}
		expectedError := "statistics only available for published articles"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})
}

func TestStatisticsUsecase_GetDailyStatisticsRange(t *testing.T) {
	statsRepo := NewMockStatisticsRepository()
	viewRepo := NewMockViewRepository()
	userRepo := NewMockUserRepository()
	articleRepo := NewMockArticleRepository()
	usecase := NewStatisticsUsecase(statsRepo, viewRepo, userRepo, articleRepo)
	ctx := context.Background()

	// Create test daily statistics
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 5; i++ {
		date := startDate.AddDate(0, 0, i)
		stats := &entity.DailyStatistics{
			Date:               date,
			TotalViews:         100 + i*10,
			UniqueUsers:        50 + i*5,
			NewUsers:           5 + i,
			TotalReadingTime:   5000 + i*500,
			AverageReadingTime: float32(100 + i*10),
		}
		statsRepo.CreateDailyStatistics(ctx, stats)
	}

	t.Run("get valid date range", func(t *testing.T) {
		endDate := startDate.AddDate(0, 0, 3)
		stats, err := usecase.GetDailyStatisticsRange(ctx, startDate, endDate)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(stats) != 4 { // 4 days inclusive
			t.Errorf("expected 4 days of stats, got %d", len(stats))
		}
	})

	t.Run("invalid date range", func(t *testing.T) {
		endDate := startDate.AddDate(0, 0, -1) // End before start
		_, err := usecase.GetDailyStatisticsRange(ctx, startDate, endDate)
		if err == nil {
			t.Error("expected error for invalid date range")
		}
		expectedError := "start date must be before end date"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("date range too large", func(t *testing.T) {
		endDate := startDate.AddDate(2, 0, 0) // 2 years later
		_, err := usecase.GetDailyStatisticsRange(ctx, startDate, endDate)
		if err == nil {
			t.Error("expected error for large date range")
		}
		expectedError := "date range cannot exceed 365 days"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})
}

func TestStatisticsUsecase_GetPlatformOverview(t *testing.T) {
	statsRepo := NewMockStatisticsRepository()
	viewRepo := NewMockViewRepository()
	userRepo := NewMockUserRepository()
	articleRepo := NewMockArticleRepository()
	usecase := NewStatisticsUsecase(statsRepo, viewRepo, userRepo, articleRepo)
	ctx := context.Background()

	t.Run("get platform overview", func(t *testing.T) {
		overview, err := usecase.GetPlatformOverview(ctx)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if overview.TotalUsers != 100 {
			t.Errorf("expected total users 100, got %d", overview.TotalUsers)
		}
		if overview.TotalArticles != 50 {
			t.Errorf("expected total articles 50, got %d", overview.TotalArticles)
		}
		if overview.TotalViews != 1000 {
			t.Errorf("expected total views 1000, got %d", overview.TotalViews)
		}
	})
}

func TestStatisticsUsecase_GetUserAnalytics(t *testing.T) {
	statsRepo := NewMockStatisticsRepository()
	viewRepo := NewMockViewRepository()
	userRepo := NewMockUserRepository()
	articleRepo := NewMockArticleRepository()
	usecase := NewStatisticsUsecase(statsRepo, viewRepo, userRepo, articleRepo)
	ctx := context.Background()

	// Create a test user
	user := &entity.User{
		Name:  "Test User",
		Email: "test@example.com",
	}
	userRepo.Create(ctx, user)

	t.Run("get user analytics", func(t *testing.T) {
		analytics, err := usecase.GetUserAnalytics(ctx, user.ID)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if analytics.UserID != user.ID {
			t.Errorf("expected user ID %d, got %d", user.ID, analytics.UserID)
		}
		if analytics.TotalReadingTime != 5000 {
			t.Errorf("expected total reading time 5000, got %d", analytics.TotalReadingTime)
		}
		if analytics.CompletionRate != 0.8 {
			t.Errorf("expected completion rate 0.8, got %f", analytics.CompletionRate)
		}
	})

	t.Run("invalid user ID", func(t *testing.T) {
		_, err := usecase.GetUserAnalytics(ctx, 0)
		if err == nil {
			t.Error("expected error for invalid user ID")
		}
		expectedError := "user ID is required"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("user not found", func(t *testing.T) {
		_, err := usecase.GetUserAnalytics(ctx, 999)
		if err == nil {
			t.Error("expected error for non-existent user")
		}
		expectedError := "user not found"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})
}