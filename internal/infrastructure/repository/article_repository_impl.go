package repository

import (
	"context"
	"layered-architecture-template/internal/domain/entity"
	"layered-architecture-template/internal/domain/repository"

	"gorm.io/gorm"
)

type articleRepositoryImpl struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) repository.ArticleRepository {
	return &articleRepositoryImpl{db: db}
}

func (r *articleRepositoryImpl) Create(ctx context.Context, article *entity.Article) error {
	return r.db.WithContext(ctx).Create(article).Error
}

func (r *articleRepositoryImpl) CreateWithTx(ctx context.Context, tx repository.Transaction, article *entity.Article) error {
	db := tx.GetDB().(*gorm.DB)
	return db.WithContext(ctx).Create(article).Error
}

func (r *articleRepositoryImpl) GetByID(ctx context.Context, id uint) (*entity.Article, error) {
	var article entity.Article
	err := r.db.WithContext(ctx).Preload("Author").First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepositoryImpl) GetAll(ctx context.Context) ([]*entity.Article, error) {
	var articles []*entity.Article
	err := r.db.WithContext(ctx).Preload("Author").Find(&articles).Error
	return articles, err
}

func (r *articleRepositoryImpl) GetByStatus(ctx context.Context, status string) ([]*entity.Article, error) {
	var articles []*entity.Article
	err := r.db.WithContext(ctx).Preload("Author").Where("status = ?", status).Find(&articles).Error
	return articles, err
}

func (r *articleRepositoryImpl) Update(ctx context.Context, article *entity.Article) error {
	return r.db.WithContext(ctx).Save(article).Error
}

func (r *articleRepositoryImpl) UpdateWithTx(ctx context.Context, tx repository.Transaction, article *entity.Article) error {
	db := tx.GetDB().(*gorm.DB)
	return db.WithContext(ctx).Save(article).Error
}

func (r *articleRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Article{}, id).Error
}

func (r *articleRepositoryImpl) DeleteWithTx(ctx context.Context, tx repository.Transaction, id uint) error {
	db := tx.GetDB().(*gorm.DB)
	return db.WithContext(ctx).Delete(&entity.Article{}, id).Error
}