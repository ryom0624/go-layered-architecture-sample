package repository

import (
	"context"
	"fmt"
	"layered-architecture-template/internal/domain/entity"
	"layered-architecture-template/internal/domain/repository"

	"gorm.io/gorm"
)

type tagRepositoryImpl struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) repository.TagRepository {
	return &tagRepositoryImpl{db: db}
}

func (r *tagRepositoryImpl) Create(ctx context.Context, tag *entity.Tag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

func (r *tagRepositoryImpl) GetByID(ctx context.Context, id uint) (*entity.Tag, error) {
	var tag entity.Tag
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepositoryImpl) GetBySlug(ctx context.Context, slug string) (*entity.Tag, error) {
	var tag entity.Tag
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepositoryImpl) GetAll(ctx context.Context) ([]*entity.Tag, error) {
	var tags []*entity.Tag
	err := r.db.WithContext(ctx).Order("name ASC").Find(&tags).Error
	return tags, err
}

func (r *tagRepositoryImpl) Update(ctx context.Context, tag *entity.Tag) error {
	return r.db.WithContext(ctx).Save(tag).Error
}

func (r *tagRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Tag{}, id).Error
}

func (r *tagRepositoryImpl) GetByNames(ctx context.Context, names []string) ([]*entity.Tag, error) {
	var tags []*entity.Tag
	err := r.db.WithContext(ctx).Where("name IN ?", names).Find(&tags).Error
	return tags, err
}

func (r *tagRepositoryImpl) GetPopular(ctx context.Context, limit int) ([]*entity.TagWithCount, error) {
	var results []struct {
		entity.Tag
		UsageCount int
	}

	err := r.db.WithContext(ctx).
		Table("tags").
		Select("tags.*, COUNT(article_tags.tag_id) as usage_count").
		Joins("LEFT JOIN article_tags ON tags.id = article_tags.tag_id").
		Group("tags.id").
		Order("usage_count DESC, tags.name ASC").
		Limit(limit).
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get popular tags: %w", err)
	}

	var tags []*entity.TagWithCount
	for _, result := range results {
		tags = append(tags, &entity.TagWithCount{
			Tag:        &result.Tag,
			UsageCount: result.UsageCount,
		})
	}

	return tags, nil
}

func (r *tagRepositoryImpl) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Tag{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

func (r *tagRepositoryImpl) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Tag{}).Where("slug = ?", slug).Count(&count).Error
	return count > 0, err
}

func (r *tagRepositoryImpl) CreateMultiple(ctx context.Context, tags []*entity.Tag) error {
	if len(tags) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&tags).Error
}