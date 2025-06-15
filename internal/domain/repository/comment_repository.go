package repository

import (
	"context"
	"layered-architecture-template/internal/domain/entity"
)

type CommentRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, comment *entity.Comment) error
	CreateWithTx(ctx context.Context, tx Transaction, comment *entity.Comment) error
	GetByID(ctx context.Context, id uint) (*entity.Comment, error)
	Update(ctx context.Context, comment *entity.Comment) error
	UpdateWithTx(ctx context.Context, tx Transaction, comment *entity.Comment) error
	Delete(ctx context.Context, id uint) error

	// Query operations
	GetByArticleID(ctx context.Context, articleID uint) ([]*entity.Comment, error)
	GetReplies(ctx context.Context, parentID uint) ([]*entity.Comment, error)
	GetByStatus(ctx context.Context, status string) ([]*entity.Comment, error)
	GetByUserID(ctx context.Context, userID uint) ([]*entity.Comment, error)

	// Hierarchical operations
	GetCommentsWithReplies(ctx context.Context, articleID uint) ([]*entity.Comment, error)
	GetCommentDepth(ctx context.Context, commentID uint) (int, error)

	// Statistics
	CountByArticleID(ctx context.Context, articleID uint) (int64, error)
	CountByStatus(ctx context.Context, status string) (int64, error)
}