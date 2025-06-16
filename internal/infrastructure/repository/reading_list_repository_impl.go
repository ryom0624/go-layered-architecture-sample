package repository

import (
	"context"

	"gorm.io/gorm"
	"layered-architecture-template/internal/domain/entity"
	"layered-architecture-template/internal/domain/repository"
)

type readingListRepositoryImpl struct {
	db *gorm.DB
}

func NewReadingListRepository(db *gorm.DB) repository.ReadingListRepository {
	return &readingListRepositoryImpl{db: db}
}

func (r *readingListRepositoryImpl) Create(ctx context.Context, readingList *entity.ReadingList) error {
	return r.db.WithContext(ctx).Create(readingList).Error
}

func (r *readingListRepositoryImpl) Update(ctx context.Context, readingList *entity.ReadingList) error {
	return r.db.WithContext(ctx).Save(readingList).Error
}

func (r *readingListRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// First delete all reading list items
		if err := tx.Where("reading_list_id = ?", id).Delete(&entity.ReadingListItem{}).Error; err != nil {
			return err
		}
		// Then delete the reading list
		return tx.Delete(&entity.ReadingList{}, id).Error
	})
}

func (r *readingListRepositoryImpl) GetByUserID(ctx context.Context, userID uint) ([]*entity.ReadingList, error) {
	var readingLists []*entity.ReadingList
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Find(&readingLists).Error
	return readingLists, err
}

func (r *readingListRepositoryImpl) GetByID(ctx context.Context, id uint) (*entity.ReadingList, error) {
	var readingList entity.ReadingList
	err := r.db.WithContext(ctx).
		Preload("User").
		First(&readingList, id).Error
	if err != nil {
		return nil, err
	}
	return &readingList, nil
}

func (r *readingListRepositoryImpl) GetPublic(ctx context.Context) ([]*entity.ReadingList, error) {
	var readingLists []*entity.ReadingList
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("is_public = ?", true).
		Order("updated_at DESC").
		Find(&readingLists).Error
	return readingLists, err
}

func (r *readingListRepositoryImpl) AddItem(ctx context.Context, item *entity.ReadingListItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *readingListRepositoryImpl) RemoveItem(ctx context.Context, readingListID, articleID uint) error {
	return r.db.WithContext(ctx).
		Where("reading_list_id = ? AND article_id = ?", readingListID, articleID).
		Delete(&entity.ReadingListItem{}).Error
}

func (r *readingListRepositoryImpl) GetItems(ctx context.Context, readingListID uint) ([]*entity.ReadingListItem, error) {
	var items []*entity.ReadingListItem
	err := r.db.WithContext(ctx).
		Preload("Article").
		Preload("Article.Author").
		Where("reading_list_id = ?", readingListID).
		Order("added_at DESC").
		Find(&items).Error
	return items, err
}

func (r *readingListRepositoryImpl) UpdateItem(ctx context.Context, item *entity.ReadingListItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}