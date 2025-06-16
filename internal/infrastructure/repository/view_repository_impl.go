package repository

import (
	"context"
	"layered-architecture-template/internal/domain/entity"
	"layered-architecture-template/internal/domain/repository"
	"time"

	"gorm.io/gorm"
)

type viewRepositoryImpl struct {
	db *gorm.DB
}

func NewViewRepository(db *gorm.DB) repository.ViewRepository {
	return &viewRepositoryImpl{
		db: db,
	}
}

func (r *viewRepositoryImpl) CreateView(ctx context.Context, view *entity.ArticleView) error {
	return r.db.WithContext(ctx).Create(view).Error
}

func (r *viewRepositoryImpl) GetViewsByArticle(ctx context.Context, articleID uint) ([]*entity.ArticleView, error) {
	var views []*entity.ArticleView
	err := r.db.WithContext(ctx).
		Where("article_id = ?", articleID).
		Order("viewed_at DESC").
		Find(&views).Error
	return views, err
}

func (r *viewRepositoryImpl) GetViewsByUser(ctx context.Context, userID uint) ([]*entity.ArticleView, error) {
	var views []*entity.ArticleView
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("viewed_at DESC").
		Find(&views).Error
	return views, err
}

func (r *viewRepositoryImpl) GetViewsByIPAndArticle(ctx context.Context, ip string, articleID uint, since time.Time) ([]*entity.ArticleView, error) {
	var views []*entity.ArticleView
	err := r.db.WithContext(ctx).
		Where("ip_address = ? AND article_id = ? AND viewed_at > ?", ip, articleID, since).
		Order("viewed_at DESC").
		Find(&views).Error
	return views, err
}

func (r *viewRepositoryImpl) UpdateReadingTime(ctx context.Context, viewID uint, readingTime int) error {
	return r.db.WithContext(ctx).
		Model(&entity.ArticleView{}).
		Where("id = ?", viewID).
		Update("reading_time", readingTime).Error
}

func (r *viewRepositoryImpl) CreateOrUpdateReadingHistory(ctx context.Context, history *entity.UserReadingHistory) error {
	var existing entity.UserReadingHistory
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND article_id = ?", history.UserID, history.ArticleID).
		First(&existing).Error
	
	if err == gorm.ErrRecordNotFound {
		return r.db.WithContext(ctx).Create(history).Error
	} else if err != nil {
		return err
	}
	
	existing.TotalReadingTime += history.TotalReadingTime
	existing.ReadingProgress = history.ReadingProgress
	existing.IsCompleted = history.IsCompleted
	existing.LastViewedAt = history.LastViewedAt
	existing.ViewCount++
	if history.IsCompleted && existing.CompletedAt == nil {
		existing.CompletedAt = history.CompletedAt
	}
	
	return r.db.WithContext(ctx).Save(&existing).Error
}

func (r *viewRepositoryImpl) GetReadingHistory(ctx context.Context, userID uint, articleID uint) (*entity.UserReadingHistory, error) {
	var history entity.UserReadingHistory
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND article_id = ?", userID, articleID).
		First(&history).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}

func (r *viewRepositoryImpl) GetUserReadingHistories(ctx context.Context, userID uint, limit, offset int) ([]*entity.UserReadingHistory, error) {
	var histories []*entity.UserReadingHistory
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("last_viewed_at DESC").
		Limit(limit).
		Offset(offset).
		Preload("Article").
		Find(&histories).Error
	return histories, err
}

func (r *viewRepositoryImpl) UpdateReadingProgress(ctx context.Context, userID uint, articleID uint, progress float32, readingTime int, isCompleted bool) error {
	updateData := map[string]interface{}{
		"reading_progress":     progress,
		"total_reading_time":   gorm.Expr("total_reading_time + ?", readingTime),
		"last_viewed_at":       time.Now(),
		"is_completed":         isCompleted,
		"view_count":           gorm.Expr("view_count + 1"),
	}
	
	if isCompleted {
		updateData["completed_at"] = time.Now()
	}
	
	return r.db.WithContext(ctx).
		Model(&entity.UserReadingHistory{}).
		Where("user_id = ? AND article_id = ?", userID, articleID).
		Updates(updateData).Error
}

func (r *viewRepositoryImpl) GetPopularArticles(ctx context.Context, limit int, since time.Time) ([]*entity.Article, error) {
	var articles []*entity.Article
	err := r.db.WithContext(ctx).
		Table("articles").
		Select("articles.*, COUNT(article_views.id) as view_count").
		Joins("LEFT JOIN article_views ON articles.id = article_views.article_id").
		Where("articles.status = ? AND article_views.viewed_at > ?", "published", since).
		Group("articles.id").
		Order("view_count DESC").
		Limit(limit).
		Find(&articles).Error
	return articles, err
}

func (r *viewRepositoryImpl) GetTrendingArticles(ctx context.Context, limit int, since time.Time) ([]*entity.Article, error) {
	var articles []*entity.Article
	err := r.db.WithContext(ctx).
		Table("articles").
		Select("articles.*, COUNT(article_views.id) as recent_views").
		Joins("LEFT JOIN article_views ON articles.id = article_views.article_id").
		Where("articles.status = ? AND article_views.viewed_at > ?", "published", since).
		Group("articles.id").
		Having("recent_views > 0").
		Order("recent_views DESC, articles.created_at DESC").
		Limit(limit).
		Find(&articles).Error
	return articles, err
}