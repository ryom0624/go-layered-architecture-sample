package repository

import (
	"context"

	"layered-architecture-template/internal/domain/entity"
)

type ReadingListRepository interface {
	Create(ctx context.Context, readingList *entity.ReadingList) error
	Update(ctx context.Context, readingList *entity.ReadingList) error
	Delete(ctx context.Context, id uint) error
	GetByUserID(ctx context.Context, userID uint) ([]*entity.ReadingList, error)
	GetByID(ctx context.Context, id uint) (*entity.ReadingList, error)
	GetPublic(ctx context.Context) ([]*entity.ReadingList, error)
	AddItem(ctx context.Context, item *entity.ReadingListItem) error
	RemoveItem(ctx context.Context, readingListID, articleID uint) error
	GetItems(ctx context.Context, readingListID uint) ([]*entity.ReadingListItem, error)
	UpdateItem(ctx context.Context, item *entity.ReadingListItem) error
}