package usecase

import (
	"context"
	"errors"
	"layered-architecture-template/internal/domain/entity"
	"layered-architecture-template/internal/domain/repository"
	"time"
)

type StatisticsUsecase interface {
	GetArticleStatistics(ctx context.Context, articleID uint) (*entity.ArticleStatistics, error)
	RecalculateArticleStatistics(ctx context.Context, articleID uint) error
	GetAllArticleStatistics(ctx context.Context, page, limit int) ([]*entity.ArticleStatistics, error)
	
	GetDailyStatistics(ctx context.Context, date time.Time) (*entity.DailyStatistics, error)
	GetDailyStatisticsRange(ctx context.Context, startDate, endDate time.Time) ([]*entity.DailyStatistics, error)
	CalculateDailyStatistics(ctx context.Context, date time.Time) error
	
	GetPlatformOverview(ctx context.Context) (*repository.PlatformOverview, error)
	GetUserAnalytics(ctx context.Context, userID uint) (*repository.UserAnalytics, error)
}

type statisticsUsecase struct {
	statsRepo   repository.StatisticsRepository
	viewRepo    repository.ViewRepository
	userRepo    repository.UserRepository
	articleRepo repository.ArticleRepository
}

func NewStatisticsUsecase(
	statsRepo repository.StatisticsRepository,
	viewRepo repository.ViewRepository,
	userRepo repository.UserRepository,
	articleRepo repository.ArticleRepository,
) StatisticsUsecase {
	return &statisticsUsecase{
		statsRepo:   statsRepo,
		viewRepo:    viewRepo,
		userRepo:    userRepo,
		articleRepo: articleRepo,
	}
}

func (s *statisticsUsecase) GetArticleStatistics(ctx context.Context, articleID uint) (*entity.ArticleStatistics, error) {
	if articleID == 0 {
		return nil, errors.New("article ID is required")
	}
	
	article, err := s.articleRepo.GetByID(ctx, articleID)
	if err != nil {
		return nil, errors.New("article not found")
	}
	if article.Status != "published" {
		return nil, errors.New("statistics only available for published articles")
	}
	
	stats, err := s.statsRepo.GetArticleStatistics(ctx, articleID)
	if err != nil {
		if err := s.statsRepo.RecalculateArticleStatistics(ctx, articleID); err != nil {
			return nil, err
		}
		return s.statsRepo.GetArticleStatistics(ctx, articleID)
	}
	
	return stats, nil
}

func (s *statisticsUsecase) RecalculateArticleStatistics(ctx context.Context, articleID uint) error {
	if articleID == 0 {
		return errors.New("article ID is required")
	}
	
	_, err := s.articleRepo.GetByID(ctx, articleID)
	if err != nil {
		return errors.New("article not found")
	}
	
	return s.statsRepo.RecalculateArticleStatistics(ctx, articleID)
}

func (s *statisticsUsecase) GetAllArticleStatistics(ctx context.Context, page, limit int) ([]*entity.ArticleStatistics, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	offset := (page - 1) * limit
	return s.statsRepo.GetAllArticleStatistics(ctx, limit, offset)
}

func (s *statisticsUsecase) GetDailyStatistics(ctx context.Context, date time.Time) (*entity.DailyStatistics, error) {
	return s.statsRepo.GetDailyStatistics(ctx, date)
}

func (s *statisticsUsecase) GetDailyStatisticsRange(ctx context.Context, startDate, endDate time.Time) ([]*entity.DailyStatistics, error) {
	if startDate.After(endDate) {
		return nil, errors.New("start date must be before end date")
	}
	
	if endDate.Sub(startDate).Hours() > 24*365 {
		return nil, errors.New("date range cannot exceed 365 days")
	}
	
	return s.statsRepo.GetDailyStatisticsRange(ctx, startDate, endDate)
}

func (s *statisticsUsecase) CalculateDailyStatistics(ctx context.Context, date time.Time) error {
	stats := &entity.DailyStatistics{
		Date: date,
	}
	
	return s.statsRepo.CreateDailyStatistics(ctx, stats)
}

func (s *statisticsUsecase) GetPlatformOverview(ctx context.Context) (*repository.PlatformOverview, error) {
	return s.statsRepo.GetPlatformOverview(ctx)
}

func (s *statisticsUsecase) GetUserAnalytics(ctx context.Context, userID uint) (*repository.UserAnalytics, error) {
	if userID == 0 {
		return nil, errors.New("user ID is required")
	}
	
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	
	return s.statsRepo.GetUserAnalytics(ctx, userID)
}