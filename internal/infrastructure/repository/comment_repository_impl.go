package repository

import (
	"context"
	"errors"
	"layered-architecture-template/internal/domain/entity"
	"layered-architecture-template/internal/domain/repository"

	"gorm.io/gorm"
)

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) repository.CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(ctx context.Context, comment *entity.Comment) error {
	if err := comment.ValidateContent(); err != nil {
		return err
	}
	if err := comment.ValidateStatus(); err != nil {
		return err
	}

	return r.db.WithContext(ctx).Create(comment).Error
}

func (r *commentRepository) CreateWithTx(ctx context.Context, tx repository.Transaction, comment *entity.Comment) error {
	if err := comment.ValidateContent(); err != nil {
		return err
	}
	if err := comment.ValidateStatus(); err != nil {
		return err
	}

	gormTx := tx.GetDB().(*gorm.DB)
	return gormTx.WithContext(ctx).Create(comment).Error
}

func (r *commentRepository) GetByID(ctx context.Context, id uint) (*entity.Comment, error) {
	var comment entity.Comment
	err := r.db.WithContext(ctx).
		Preload("Author").
		Preload("Article").
		Preload("Parent").
		Preload("Replies").
		First(&comment, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("comment not found")
		}
		return nil, err
	}

	return &comment, nil
}

func (r *commentRepository) Update(ctx context.Context, comment *entity.Comment) error {
	if err := comment.ValidateContent(); err != nil {
		return err
	}
	if err := comment.ValidateStatus(); err != nil {
		return err
	}

	return r.db.WithContext(ctx).Save(comment).Error
}

func (r *commentRepository) UpdateWithTx(ctx context.Context, tx repository.Transaction, comment *entity.Comment) error {
	if err := comment.ValidateContent(); err != nil {
		return err
	}
	if err := comment.ValidateStatus(); err != nil {
		return err
	}

	gormTx := tx.GetDB().(*gorm.DB)
	return gormTx.WithContext(ctx).Save(comment).Error
}

func (r *commentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Comment{}, id).Error
}

func (r *commentRepository) GetByArticleID(ctx context.Context, articleID uint) ([]*entity.Comment, error) {
	var comments []*entity.Comment
	err := r.db.WithContext(ctx).
		Preload("Author").
		Preload("Replies").
		Preload("Replies.Author").
		Where("article_id = ? AND parent_id IS NULL", articleID).
		Order("created_at ASC").
		Find(&comments).Error

	return comments, err
}

func (r *commentRepository) GetReplies(ctx context.Context, parentID uint) ([]*entity.Comment, error) {
	var replies []*entity.Comment
	err := r.db.WithContext(ctx).
		Preload("Author").
		Where("parent_id = ?", parentID).
		Order("created_at ASC").
		Find(&replies).Error

	return replies, err
}

func (r *commentRepository) GetByStatus(ctx context.Context, status string) ([]*entity.Comment, error) {
	var comments []*entity.Comment
	err := r.db.WithContext(ctx).
		Preload("Author").
		Preload("Article").
		Where("status = ?", status).
		Order("created_at DESC").
		Find(&comments).Error

	return comments, err
}

func (r *commentRepository) GetByUserID(ctx context.Context, userID uint) ([]*entity.Comment, error) {
	var comments []*entity.Comment
	err := r.db.WithContext(ctx).
		Preload("Article").
		Where("author_id = ?", userID).
		Order("created_at DESC").
		Find(&comments).Error

	return comments, err
}

func (r *commentRepository) GetCommentsWithReplies(ctx context.Context, articleID uint) ([]*entity.Comment, error) {
	var comments []*entity.Comment
	err := r.db.WithContext(ctx).
		Preload("Author").
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Author").Order("created_at ASC")
		}).
		Where("article_id = ? AND parent_id IS NULL AND status = ?", articleID, "approved").
		Order("created_at ASC").
		Find(&comments).Error

	return comments, err
}

func (r *commentRepository) GetCommentDepth(ctx context.Context, commentID uint) (int, error) {
	var comment entity.Comment
	err := r.db.WithContext(ctx).Select("parent_id").First(&comment, commentID).Error
	if err != nil {
		return 0, err
	}

	depth := 0
	currentParentID := comment.ParentID

	for currentParentID != nil {
		depth++
		if depth > 10 { // Prevent infinite recursion
			return depth, errors.New("comment hierarchy too deep")
		}

		var parent entity.Comment
		err := r.db.WithContext(ctx).Select("parent_id").First(&parent, *currentParentID).Error
		if err != nil {
			break
		}
		currentParentID = parent.ParentID
	}

	return depth, nil
}

func (r *commentRepository) CountByArticleID(ctx context.Context, articleID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Comment{}).
		Where("article_id = ? AND status = ?", articleID, "approved").
		Count(&count).Error

	return count, err
}

func (r *commentRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Comment{}).
		Where("status = ?", status).
		Count(&count).Error

	return count, err
}