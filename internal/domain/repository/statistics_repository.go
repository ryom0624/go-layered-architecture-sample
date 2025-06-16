package repository

import (
	"context"
	"layered-architecture-template/internal/domain/entity"
	"time"
)

type StatisticsRepository interface {
	CreateOrUpdateArticleStatistics(ctx context.Context, stats *entity.ArticleStatistics) error
	GetArticleStatistics(ctx context.Context, articleID uint) (*entity.ArticleStatistics, error)
	GetAllArticleStatistics(ctx context.Context, limit, offset int) ([]*entity.ArticleStatistics, error)
	RecalculateArticleStatistics(ctx context.Context, articleID uint) error
	
	CreateDailyStatistics(ctx context.Context, stats *entity.DailyStatistics) error
	GetDailyStatistics(ctx context.Context, date time.Time) (*entity.DailyStatistics, error)
	GetDailyStatisticsRange(ctx context.Context, startDate, endDate time.Time) ([]*entity.DailyStatistics, error)
	UpdateDailyStatistics(ctx context.Context, stats *entity.DailyStatistics) error
	
	GetPlatformOverview(ctx context.Context) (*PlatformOverview, error)
	GetUserAnalytics(ctx context.Context, userID uint) (*UserAnalytics, error)
}

type PlatformOverview struct {
	TotalUsers           int     `json:"total_users"`
	TotalArticles        int     `json:"total_articles"`
	TotalViews           int     `json:"total_views"`
	TotalReadingTime     int     `json:"total_reading_time"`
	AverageReadingTime   float32 `json:"average_reading_time"`
	ActiveUsersToday     int     `json:"active_users_today"`
	NewUsersToday        int     `json:"new_users_today"`
	PublishedToday       int     `json:"published_today"`
	PopularToday         []*entity.Article `json:"popular_today"`
}

type UserAnalytics struct {
	UserID               uint    `json:"user_id"`
	TotalReadingTime     int     `json:"total_reading_time"`
	ArticlesRead         int     `json:"articles_read"`
	ArticlesCompleted    int     `json:"articles_completed"`
	AverageReadingTime   float32 `json:"average_reading_time"`
	CompletionRate       float32 `json:"completion_rate"`
	FavoriteTopic        string  `json:"favorite_topic"`
	ReadingStreak        int     `json:"reading_streak"`
	LastActiveDate       time.Time `json:"last_active_date"`
	RecentReadings       []*entity.UserReadingHistory `json:"recent_readings"`
}