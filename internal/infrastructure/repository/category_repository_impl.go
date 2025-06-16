package repository

import (
	"context"
	"fmt"
	"layered-architecture-template/internal/domain/entity"
	"layered-architecture-template/internal/domain/repository"

	"gorm.io/gorm"
)

type categoryRepositoryImpl struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) repository.CategoryRepository {
	return &categoryRepositoryImpl{db: db}
}

func (r *categoryRepositoryImpl) Create(ctx context.Context, category *entity.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *categoryRepositoryImpl) GetByID(ctx context.Context, id uint) (*entity.Category, error) {
	var category entity.Category
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepositoryImpl) GetBySlug(ctx context.Context, slug string) (*entity.Category, error) {
	var category entity.Category
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepositoryImpl) GetAll(ctx context.Context) ([]*entity.Category, error) {
	var categories []*entity.Category
	err := r.db.WithContext(ctx).Order("name ASC").Find(&categories).Error
	return categories, err
}

func (r *categoryRepositoryImpl) Update(ctx context.Context, category *entity.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *categoryRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Category{}, id).Error
}

func (r *categoryRepositoryImpl) GetWithArticleCount(ctx context.Context) ([]*entity.CategoryWithCount, error) {
	var results []struct {
		entity.Category
		ArticleCount int
	}

	err := r.db.WithContext(ctx).
		Table("categories").
		Select("categories.*, COUNT(articles.id) as article_count").
		Joins("LEFT JOIN articles ON categories.id = articles.category_id").
		Group("categories.id").
		Order("categories.name ASC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get categories with article count: %w", err)
	}

	var categories []*entity.CategoryWithCount
	for _, result := range results {
		categories = append(categories, &entity.CategoryWithCount{
			Category:     &result.Category,
			ArticleCount: result.ArticleCount,
		})
	}

	return categories, nil
}

func (r *categoryRepositoryImpl) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Category{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

func (r *categoryRepositoryImpl) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Category{}).Where("slug = ?", slug).Count(&count).Error
	return count > 0, err
}