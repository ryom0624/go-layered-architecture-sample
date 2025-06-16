package repository

import (
	"context"
	"layered-architecture-template/internal/domain/entity"
)

type TagRepository interface {
	Create(ctx context.Context, tag *entity.Tag) error
	GetByID(ctx context.Context, id uint) (*entity.Tag, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Tag, error)
	GetAll(ctx context.Context) ([]*entity.Tag, error)
	Update(ctx context.Context, tag *entity.Tag) error
	Delete(ctx context.Context, id uint) error
	GetByNames(ctx context.Context, names []string) ([]*entity.Tag, error)
	GetPopular(ctx context.Context, limit int) ([]*entity.TagWithCount, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
	ExistsBySlug(ctx context.Context, slug string) (bool, error)
	CreateMultiple(ctx context.Context, tags []*entity.Tag) error
}