package repository

import (
	"context"
	"layered-architecture-template/internal/domain/entity"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *entity.Category) error
	GetByID(ctx context.Context, id uint) (*entity.Category, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Category, error)
	GetAll(ctx context.Context) ([]*entity.Category, error)
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id uint) error
	GetWithArticleCount(ctx context.Context) ([]*entity.CategoryWithCount, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
	ExistsBySlug(ctx context.Context, slug string) (bool, error)
}