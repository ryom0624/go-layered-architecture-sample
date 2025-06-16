package repository

import (
	"context"
	"layered-architecture-template/internal/domain/entity"
	"time"
)

type ViewRepository interface {
	CreateView(ctx context.Context, view *entity.ArticleView) error
	GetViewsByArticle(ctx context.Context, articleID uint) ([]*entity.ArticleView, error)
	GetViewsByUser(ctx context.Context, userID uint) ([]*entity.ArticleView, error)
	GetViewsByIPAndArticle(ctx context.Context, ip string, articleID uint, since time.Time) ([]*entity.ArticleView, error)
	UpdateReadingTime(ctx context.Context, viewID uint, readingTime int) error
	
	CreateOrUpdateReadingHistory(ctx context.Context, history *entity.UserReadingHistory) error
	GetReadingHistory(ctx context.Context, userID uint, articleID uint) (*entity.UserReadingHistory, error)
	GetUserReadingHistories(ctx context.Context, userID uint, limit, offset int) ([]*entity.UserReadingHistory, error)
	UpdateReadingProgress(ctx context.Context, userID uint, articleID uint, progress float32, readingTime int, isCompleted bool) error
	
	GetPopularArticles(ctx context.Context, limit int, since time.Time) ([]*entity.Article, error)
	GetTrendingArticles(ctx context.Context, limit int, since time.Time) ([]*entity.Article, error)
}