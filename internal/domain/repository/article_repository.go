package repository

import (
	"context"
	"layered-architecture-template/internal/domain/entity"
)

type ArticleRepository interface {
	Create(ctx context.Context, article *entity.Article) error
	CreateWithTx(ctx context.Context, tx Transaction, article *entity.Article) error
	GetByID(ctx context.Context, id uint) (*entity.Article, error)
	GetAll(ctx context.Context) ([]*entity.Article, error)
	GetByStatus(ctx context.Context, status string) ([]*entity.Article, error)
	Update(ctx context.Context, article *entity.Article) error
	UpdateWithTx(ctx context.Context, tx Transaction, article *entity.Article) error
	Delete(ctx context.Context, id uint) error
	DeleteWithTx(ctx context.Context, tx Transaction, id uint) error
}