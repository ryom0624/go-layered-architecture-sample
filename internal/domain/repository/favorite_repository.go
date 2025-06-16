package repository

import (
	"context"

	"layered-architecture-template/internal/domain/entity"
)

type FavoriteRepository interface {
	Create(ctx context.Context, favorite *entity.Favorite) error
	CreateWithTx(ctx context.Context, tx Transaction, favorite *entity.Favorite) error
	Delete(ctx context.Context, userID, articleID uint) error
	DeleteWithTx(ctx context.Context, tx Transaction, userID, articleID uint) error
	GetByUserID(ctx context.Context, userID uint) ([]*entity.Favorite, error)
	GetByArticleID(ctx context.Context, articleID uint) ([]*entity.Favorite, error)
	IsFavorited(ctx context.Context, userID, articleID uint) (bool, error)
	GetFavoriteCount(ctx context.Context, articleID uint) (int, error)
}