package repository

import (
	"context"

	"gorm.io/gorm"
	"layered-architecture-template/internal/domain/entity"
	"layered-architecture-template/internal/domain/repository"
)

type favoriteRepositoryImpl struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) repository.FavoriteRepository {
	return &favoriteRepositoryImpl{db: db}
}

func (r *favoriteRepositoryImpl) Create(ctx context.Context, favorite *entity.Favorite) error {
	return r.db.WithContext(ctx).Create(favorite).Error
}

func (r *favoriteRepositoryImpl) CreateWithTx(ctx context.Context, tx repository.Transaction, favorite *entity.Favorite) error {
	gormTx := tx.GetDB().(*gorm.DB)
	return gormTx.WithContext(ctx).Create(favorite).Error
}

func (r *favoriteRepositoryImpl) Delete(ctx context.Context, userID, articleID uint) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND article_id = ?", userID, articleID).Delete(&entity.Favorite{}).Error
}

func (r *favoriteRepositoryImpl) DeleteWithTx(ctx context.Context, tx repository.Transaction, userID, articleID uint) error {
	gormTx := tx.GetDB().(*gorm.DB)
	return gormTx.WithContext(ctx).Where("user_id = ? AND article_id = ?", userID, articleID).Delete(&entity.Favorite{}).Error
}

func (r *favoriteRepositoryImpl) GetByUserID(ctx context.Context, userID uint) ([]*entity.Favorite, error) {
	var favorites []*entity.Favorite
	err := r.db.WithContext(ctx).
		Preload("Article").
		Preload("Article.Author").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&favorites).Error
	return favorites, err
}

func (r *favoriteRepositoryImpl) GetByArticleID(ctx context.Context, articleID uint) ([]*entity.Favorite, error) {
	var favorites []*entity.Favorite
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("article_id = ?", articleID).
		Order("created_at DESC").
		Find(&favorites).Error
	return favorites, err
}

func (r *favoriteRepositoryImpl) IsFavorited(ctx context.Context, userID, articleID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Favorite{}).
		Where("user_id = ? AND article_id = ?", userID, articleID).
		Count(&count).Error
	return count > 0, err
}

func (r *favoriteRepositoryImpl) GetFavoriteCount(ctx context.Context, articleID uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Favorite{}).
		Where("article_id = ?", articleID).
		Count(&count).Error
	return int(count), err
}