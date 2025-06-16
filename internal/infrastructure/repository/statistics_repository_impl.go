package repository

import (
	"context"
	"layered-architecture-template/internal/domain/entity"
	"layered-architecture-template/internal/domain/repository"
	"time"

	"gorm.io/gorm"
)

type statisticsRepositoryImpl struct {
	db *gorm.DB
}

func NewStatisticsRepository(db *gorm.DB) repository.StatisticsRepository {
	return &statisticsRepositoryImpl{
		db: db,
	}
}

func (r *statisticsRepositoryImpl) CreateOrUpdateArticleStatistics(ctx context.Context, stats *entity.ArticleStatistics) error {
	var existing entity.ArticleStatistics
	err := r.db.WithContext(ctx).
		Where("article_id = ?", stats.ArticleID).
		First(&existing).Error
	
	if err == gorm.ErrRecordNotFound {
		return r.db.WithContext(ctx).Create(stats).Error
	} else if err != nil {
		return err
	}
	
	existing.TotalViews = stats.TotalViews
	existing.UniqueViews = stats.UniqueViews
	existing.AuthenticatedViews = stats.AuthenticatedViews
	existing.AnonymousViews = stats.AnonymousViews
	existing.AverageReadingTime = stats.AverageReadingTime
	existing.CompletionRate = stats.CompletionRate
	existing.TotalReadingTime = stats.TotalReadingTime
	existing.LastCalculatedAt = time.Now()
	
	return r.db.WithContext(ctx).Save(&existing).Error
}

func (r *statisticsRepositoryImpl) GetArticleStatistics(ctx context.Context, articleID uint) (*entity.ArticleStatistics, error) {
	var stats entity.ArticleStatistics
	err := r.db.WithContext(ctx).
		Where("article_id = ?", articleID).
		First(&stats).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (r *statisticsRepositoryImpl) GetAllArticleStatistics(ctx context.Context, limit, offset int) ([]*entity.ArticleStatistics, error) {
	var stats []*entity.ArticleStatistics
	err := r.db.WithContext(ctx).
		Order("total_views DESC").
		Limit(limit).
		Offset(offset).
		Preload("Article").
		Find(&stats).Error
	return stats, err
}

func (r *statisticsRepositoryImpl) RecalculateArticleStatistics(ctx context.Context, articleID uint) error {
	var stats entity.ArticleStatistics
	
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var totalViews, uniqueViews, authenticatedViews, anonymousViews int64
		var totalReadingTime int64
		var avgReadingTime float32
		var completionRate float32
		
		tx.Model(&entity.ArticleView{}).
			Where("article_id = ?", articleID).
			Count(&totalViews)
		
		tx.Model(&entity.ArticleView{}).
			Where("article_id = ?", articleID).
			Distinct("COALESCE(user_id, ip_address)").
			Count(&uniqueViews)
		
		tx.Model(&entity.ArticleView{}).
			Where("article_id = ? AND user_id IS NOT NULL", articleID).
			Count(&authenticatedViews)
		
		tx.Model(&entity.ArticleView{}).
			Where("article_id = ? AND user_id IS NULL", articleID).
			Count(&anonymousViews)
		
		tx.Model(&entity.ArticleView{}).
			Where("article_id = ?", articleID).
			Select("COALESCE(SUM(reading_time), 0)").
			Scan(&totalReadingTime)
		
		if totalViews > 0 {
			avgReadingTime = float32(totalReadingTime) / float32(totalViews)
		}
		
		var completedCount int64
		tx.Model(&entity.UserReadingHistory{}).
			Where("article_id = ? AND is_completed = ?", articleID, true).
			Count(&completedCount)
		
		var totalReadings int64
		tx.Model(&entity.UserReadingHistory{}).
			Where("article_id = ?", articleID).
			Count(&totalReadings)
		
		if totalReadings > 0 {
			completionRate = float32(completedCount) / float32(totalReadings)
		}
		
		stats = entity.ArticleStatistics{
			ArticleID:            articleID,
			TotalViews:           int(totalViews),
			UniqueViews:          int(uniqueViews),
			AuthenticatedViews:   int(authenticatedViews),
			AnonymousViews:       int(anonymousViews),
			AverageReadingTime:   avgReadingTime,
			CompletionRate:       completionRate,
			TotalReadingTime:     int(totalReadingTime),
			LastCalculatedAt:     time.Now(),
		}
		
		return r.CreateOrUpdateArticleStatistics(ctx, &stats)
	})
}

func (r *statisticsRepositoryImpl) CreateDailyStatistics(ctx context.Context, stats *entity.DailyStatistics) error {
	return r.db.WithContext(ctx).Create(stats).Error
}

func (r *statisticsRepositoryImpl) GetDailyStatistics(ctx context.Context, date time.Time) (*entity.DailyStatistics, error) {
	var stats entity.DailyStatistics
	err := r.db.WithContext(ctx).
		Where("date = ?", date.Format("2006-01-02")).
		First(&stats).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (r *statisticsRepositoryImpl) GetDailyStatisticsRange(ctx context.Context, startDate, endDate time.Time) ([]*entity.DailyStatistics, error) {
	var stats []*entity.DailyStatistics
	err := r.db.WithContext(ctx).
		Where("date >= ? AND date <= ?", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Order("date ASC").
		Find(&stats).Error
	return stats, err
}

func (r *statisticsRepositoryImpl) UpdateDailyStatistics(ctx context.Context, stats *entity.DailyStatistics) error {
	return r.db.WithContext(ctx).Save(stats).Error
}

func (r *statisticsRepositoryImpl) GetPlatformOverview(ctx context.Context) (*repository.PlatformOverview, error) {
	var overview repository.PlatformOverview
	today := time.Now().Format("2006-01-02")
	
	var totalUsers, totalArticles, totalViews int64
	r.db.WithContext(ctx).Model(&entity.User{}).Count(&totalUsers)
	r.db.WithContext(ctx).Model(&entity.Article{}).Count(&totalArticles)
	r.db.WithContext(ctx).Model(&entity.ArticleView{}).Count(&totalViews)
	
	overview.TotalUsers = int(totalUsers)
	overview.TotalArticles = int(totalArticles)
	overview.TotalViews = int(totalViews)
	
	var totalReadingTime int64
	r.db.WithContext(ctx).Model(&entity.ArticleView{}).
		Select("COALESCE(SUM(reading_time), 0)").
		Scan(&totalReadingTime)
	overview.TotalReadingTime = int(totalReadingTime)
	
	if overview.TotalViews > 0 {
		overview.AverageReadingTime = float32(totalReadingTime) / float32(overview.TotalViews)
	}
	
	var activeUsersToday, newUsersToday, publishedToday int64
	r.db.WithContext(ctx).Model(&entity.ArticleView{}).
		Where("DATE(viewed_at) = ?", today).
		Distinct("user_id").
		Count(&activeUsersToday)
	
	r.db.WithContext(ctx).Model(&entity.User{}).
		Where("DATE(created_at) = ?", today).
		Count(&newUsersToday)
	
	r.db.WithContext(ctx).Model(&entity.Article{}).
		Where("DATE(created_at) = ? AND status = ?", today, "published").
		Count(&publishedToday)
	
	overview.ActiveUsersToday = int(activeUsersToday)
	overview.NewUsersToday = int(newUsersToday)
	overview.PublishedToday = int(publishedToday)
	
	r.db.WithContext(ctx).
		Table("articles").
		Select("articles.*").
		Joins("JOIN article_views ON articles.id = article_views.article_id").
		Where("DATE(article_views.viewed_at) = ? AND articles.status = ?", today, "published").
		Group("articles.id").
		Order("COUNT(article_views.id) DESC").
		Limit(5).
		Find(&overview.PopularToday)
	
	return &overview, nil
}

func (r *statisticsRepositoryImpl) GetUserAnalytics(ctx context.Context, userID uint) (*repository.UserAnalytics, error) {
	var analytics repository.UserAnalytics
	analytics.UserID = userID
	
	var totalReadingTime int64
	r.db.WithContext(ctx).Model(&entity.UserReadingHistory{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(total_reading_time), 0)").
		Scan(&totalReadingTime)
	analytics.TotalReadingTime = int(totalReadingTime)
	
	var articlesRead int64
	r.db.WithContext(ctx).Model(&entity.UserReadingHistory{}).
		Where("user_id = ?", userID).
		Count(&articlesRead)
	analytics.ArticlesRead = int(articlesRead)
	
	var articlesCompleted int64
	r.db.WithContext(ctx).Model(&entity.UserReadingHistory{}).
		Where("user_id = ? AND is_completed = ?", userID, true).
		Count(&articlesCompleted)
	analytics.ArticlesCompleted = int(articlesCompleted)
	
	if analytics.ArticlesRead > 0 {
		analytics.AverageReadingTime = float32(totalReadingTime) / float32(analytics.ArticlesRead)
		analytics.CompletionRate = float32(articlesCompleted) / float32(analytics.ArticlesRead)
	}
	
	r.db.WithContext(ctx).Model(&entity.UserReadingHistory{}).
		Where("user_id = ?", userID).
		Select("last_viewed_at").
		Order("last_viewed_at DESC").
		Limit(1).
		Scan(&analytics.LastActiveDate)
	
	r.db.WithContext(ctx).Model(&entity.UserReadingHistory{}).
		Where("user_id = ?", userID).
		Order("last_viewed_at DESC").
		Limit(10).
		Preload("Article").
		Find(&analytics.RecentReadings)
	
	return &analytics, nil
}