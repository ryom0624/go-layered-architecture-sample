package usecase

import (
	"context"
	"errors"
	"layered-architecture-template/internal/domain/entity"
	"layered-architecture-template/internal/domain/repository"
	"time"
	"layered-architecture-template/pkg/constants"
)

type ViewUsecase interface {
	TrackView(ctx context.Context, articleID uint, userID *uint, ipAddress, userAgent string) error
	UpdateReadingTime(ctx context.Context, viewID uint, readingTime int) error
	
	TrackReadingProgress(ctx context.Context, userID uint, articleID uint, progress float32, readingTime int) error
	GetReadingHistory(ctx context.Context, userID uint, articleID uint) (*entity.UserReadingHistory, error)
	GetUserReadingHistories(ctx context.Context, userID uint, page, limit int) ([]*entity.UserReadingHistory, error)
	
	GetPopularArticles(ctx context.Context, limit int, days int) ([]*entity.Article, error)
	GetTrendingArticles(ctx context.Context, limit int, hours int) ([]*entity.Article, error)
}

type viewUsecase struct {
	viewRepo    repository.ViewRepository
	articleRepo repository.ArticleRepository
}

func NewViewUsecase(viewRepo repository.ViewRepository, articleRepo repository.ArticleRepository) ViewUsecase {
	return &viewUsecase{
		viewRepo:    viewRepo,
		articleRepo: articleRepo,
	}
}

func (v *viewUsecase) TrackView(ctx context.Context, articleID uint, userID *uint, ipAddress, userAgent string) error {
	if articleID == 0 {
		return errors.New("article ID is required")
	}
	if ipAddress == "" {
		return errors.New("IP address is required")
	}
	
	article, err := v.articleRepo.GetByID(ctx, articleID)
	if err != nil {
		return errors.New("article not found")
	}
	if article.Status != "published" {
		return errors.New("article is not published")
	}
	
	since := time.Now().Add(-constants.DuplicateViewThresholdMinutes * time.Minute)
	existingViews, err := v.viewRepo.GetViewsByIPAndArticle(ctx, ipAddress, articleID, since)
	if err == nil && len(existingViews) > 0 {
		return nil
	}
	
	view := &entity.ArticleView{
		ArticleID: articleID,
		UserID:    userID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		ViewedAt:  time.Now(),
	}
	
	if err := v.viewRepo.CreateView(ctx, view); err != nil {
		return err
	}
	
	if userID != nil {
		history := &entity.UserReadingHistory{
			UserID:          *userID,
			ArticleID:       articleID,
			FirstViewedAt:   time.Now(),
			LastViewedAt:    time.Now(),
			ViewCount:       1,
		}
		return v.viewRepo.CreateOrUpdateReadingHistory(ctx, history)
	}
	
	return nil
}

func (v *viewUsecase) UpdateReadingTime(ctx context.Context, viewID uint, readingTime int) error {
	if viewID == 0 {
		return errors.New("view ID is required")
	}
	if readingTime < 0 {
		return errors.New("reading time must be non-negative")
	}
	
	return v.viewRepo.UpdateReadingTime(ctx, viewID, readingTime)
}

func (v *viewUsecase) TrackReadingProgress(ctx context.Context, userID uint, articleID uint, progress float32, readingTime int) error {
	if userID == 0 || articleID == 0 {
		return errors.New("user ID and article ID are required")
	}
	if progress < 0 || progress > 100 {
		return errors.New("progress must be between 0 and 100")
	}
	if readingTime < 0 {
		return errors.New("reading time must be non-negative")
	}
	
	article, err := v.articleRepo.GetByID(ctx, articleID)
	if err != nil {
		return errors.New("article not found")
	}
	if article.Status != "published" {
		return errors.New("article is not published")
	}
	
	isCompleted := progress >= constants.ReadCompletionThreshold
	
	return v.viewRepo.UpdateReadingProgress(ctx, userID, articleID, progress, readingTime, isCompleted)
}

func (v *viewUsecase) GetReadingHistory(ctx context.Context, userID uint, articleID uint) (*entity.UserReadingHistory, error) {
	if userID == 0 || articleID == 0 {
		return nil, errors.New("user ID and article ID are required")
	}
	
	return v.viewRepo.GetReadingHistory(ctx, userID, articleID)
}

func (v *viewUsecase) GetUserReadingHistories(ctx context.Context, userID uint, page, limit int) ([]*entity.UserReadingHistory, error) {
	if userID == 0 {
		return nil, errors.New("user ID is required")
	}
	if page < 1 {
		page = 1
	}
	if limit <= 0 || limit > constants.MaxPageSize {
		limit = constants.DefaultReadingHistoryPageSize
	}
	
	offset := (page - 1) * limit
	return v.viewRepo.GetUserReadingHistories(ctx, userID, limit, offset)
}

func (v *viewUsecase) GetPopularArticles(ctx context.Context, limit int, days int) ([]*entity.Article, error) {
	if limit <= 0 || limit > constants.MaxPageSize {
		limit = constants.DefaultPopularArticlesLimit
	}
	if days <= 0 {
		days = constants.PopularArticlesPeriodDays
	}
	
	since := time.Now().AddDate(0, 0, -days)
	return v.viewRepo.GetPopularArticles(ctx, limit, since)
}

func (v *viewUsecase) GetTrendingArticles(ctx context.Context, limit int, hours int) ([]*entity.Article, error) {
	if limit <= 0 || limit > constants.MaxPageSize {
		limit = constants.DefaultTrendingArticlesLimit
	}
	if hours <= 0 {
		hours = constants.TrendingArticlesPeriodHours
	}
	
	since := time.Now().Add(-time.Duration(hours) * time.Hour)
	return v.viewRepo.GetTrendingArticles(ctx, limit, since)
}