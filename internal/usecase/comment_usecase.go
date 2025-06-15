package usecase

import (
	"context"
	"errors"
	"layered-architecture-template/internal/domain/entity"
	"layered-architecture-template/internal/domain/repository"
)

type CommentUsecase interface {
	CreateComment(ctx context.Context, content string, authorID, articleID uint) (*entity.Comment, error)
	CreateReply(ctx context.Context, content string, authorID, parentID uint) (*entity.Comment, error)
	GetArticleComments(ctx context.Context, articleID uint) ([]*entity.Comment, error)
	GetComment(ctx context.Context, id uint) (*entity.Comment, error)
	UpdateComment(ctx context.Context, id uint, content string, userID uint) (*entity.Comment, error)
	DeleteComment(ctx context.Context, id uint, userID uint) error
	ApproveComment(ctx context.Context, id uint) (*entity.Comment, error)
	RejectComment(ctx context.Context, id uint) (*entity.Comment, error)
	GetPendingComments(ctx context.Context) ([]*entity.Comment, error)
	GetUserComments(ctx context.Context, userID uint) ([]*entity.Comment, error)
}

type commentUsecase struct {
	commentRepo        repository.CommentRepository
	userRepo           repository.UserRepository
	articleRepo        repository.ArticleRepository
	transactionManager repository.TransactionManager
}

func NewCommentUsecase(
	commentRepo repository.CommentRepository,
	userRepo repository.UserRepository,
	articleRepo repository.ArticleRepository,
	transactionManager repository.TransactionManager,
) CommentUsecase {
	return &commentUsecase{
		commentRepo:        commentRepo,
		userRepo:           userRepo,
		articleRepo:        articleRepo,
		transactionManager: transactionManager,
	}
}

func (u *commentUsecase) CreateComment(ctx context.Context, content string, authorID, articleID uint) (*entity.Comment, error) {
	if content == "" {
		return nil, errors.New("content is required")
	}
	if authorID == 0 {
		return nil, errors.New("author ID is required")
	}
	if articleID == 0 {
		return nil, errors.New("article ID is required")
	}

	// Verify author exists
	_, err := u.userRepo.GetByID(ctx, authorID)
	if err != nil {
		return nil, errors.New("author not found")
	}

	// Verify article exists and is published
	article, err := u.articleRepo.GetByID(ctx, articleID)
	if err != nil {
		return nil, errors.New("article not found")
	}
	if article.Status != "published" {
		return nil, errors.New("comments can only be made on published articles")
	}

	comment := &entity.Comment{
		Content:   content,
		AuthorID:  authorID,
		ArticleID: articleID,
		Status:    "pending",
	}

	err = u.commentRepo.Create(ctx, comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (u *commentUsecase) CreateReply(ctx context.Context, content string, authorID, parentID uint) (*entity.Comment, error) {
	if content == "" {
		return nil, errors.New("content is required")
	}
	if authorID == 0 {
		return nil, errors.New("author ID is required")
	}
	if parentID == 0 {
		return nil, errors.New("parent comment ID is required")
	}

	// Verify author exists
	_, err := u.userRepo.GetByID(ctx, authorID)
	if err != nil {
		return nil, errors.New("author not found")
	}

	// Verify parent comment exists and is approved
	parentComment, err := u.commentRepo.GetByID(ctx, parentID)
	if err != nil {
		return nil, errors.New("parent comment not found")
	}
	if !parentComment.IsApproved() {
		return nil, errors.New("cannot reply to unapproved comment")
	}

	// Check comment depth limit (max 3 levels)
	depth, err := u.commentRepo.GetCommentDepth(ctx, parentID)
	if err != nil {
		return nil, err
	}
	if depth >= 2 { // 0-indexed, so 2 means 3 levels deep
		return nil, errors.New("maximum reply depth exceeded")
	}

	// Verify parent article is published
	article, err := u.articleRepo.GetByID(ctx, parentComment.ArticleID)
	if err != nil {
		return nil, errors.New("article not found")
	}
	if article.Status != "published" {
		return nil, errors.New("comments can only be made on published articles")
	}

	comment := &entity.Comment{
		Content:   content,
		AuthorID:  authorID,
		ArticleID: parentComment.ArticleID,
		ParentID:  &parentID,
		Status:    "pending",
	}

	err = u.commentRepo.Create(ctx, comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (u *commentUsecase) GetArticleComments(ctx context.Context, articleID uint) ([]*entity.Comment, error) {
	if articleID == 0 {
		return nil, errors.New("article ID is required")
	}

	return u.commentRepo.GetCommentsWithReplies(ctx, articleID)
}

func (u *commentUsecase) GetComment(ctx context.Context, id uint) (*entity.Comment, error) {
	if id == 0 {
		return nil, errors.New("comment ID is required")
	}

	return u.commentRepo.GetByID(ctx, id)
}

func (u *commentUsecase) UpdateComment(ctx context.Context, id uint, content string, userID uint) (*entity.Comment, error) {
	if id == 0 {
		return nil, errors.New("comment ID is required")
	}
	if content == "" {
		return nil, errors.New("content is required")
	}
	if userID == 0 {
		return nil, errors.New("user ID is required")
	}

	comment, err := u.commentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("comment not found")
	}

	// Only comment author can update their comment
	if comment.AuthorID != userID {
		return nil, errors.New("only comment author can update the comment")
	}

	comment.Content = content
	comment.Status = "pending" // Reset to pending after edit

	err = u.commentRepo.Update(ctx, comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (u *commentUsecase) DeleteComment(ctx context.Context, id uint, userID uint) error {
	if id == 0 {
		return errors.New("comment ID is required")
	}
	if userID == 0 {
		return errors.New("user ID is required")
	}

	comment, err := u.commentRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("comment not found")
	}

	// Only comment author can delete their comment
	if comment.AuthorID != userID {
		return errors.New("only comment author can delete the comment")
	}

	return u.commentRepo.Delete(ctx, id)
}

func (u *commentUsecase) ApproveComment(ctx context.Context, id uint) (*entity.Comment, error) {
	if id == 0 {
		return nil, errors.New("comment ID is required")
	}

	var result *entity.Comment
	err := u.transactionManager.WithTransaction(ctx, func(tx repository.Transaction) error {
		comment, err := u.commentRepo.GetByID(ctx, id)
		if err != nil {
			return errors.New("comment not found")
		}

		if comment.IsApproved() {
			return errors.New("comment is already approved")
		}

		comment.Status = "approved"
		if err := u.commentRepo.UpdateWithTx(ctx, tx, comment); err != nil {
			return err
		}

		result = comment
		return nil
	})

	return result, err
}

func (u *commentUsecase) RejectComment(ctx context.Context, id uint) (*entity.Comment, error) {
	if id == 0 {
		return nil, errors.New("comment ID is required")
	}

	var result *entity.Comment
	err := u.transactionManager.WithTransaction(ctx, func(tx repository.Transaction) error {
		comment, err := u.commentRepo.GetByID(ctx, id)
		if err != nil {
			return errors.New("comment not found")
		}

		if comment.Status == "rejected" {
			return errors.New("comment is already rejected")
		}

		comment.Status = "rejected"
		if err := u.commentRepo.UpdateWithTx(ctx, tx, comment); err != nil {
			return err
		}

		result = comment
		return nil
	})

	return result, err
}

func (u *commentUsecase) GetPendingComments(ctx context.Context) ([]*entity.Comment, error) {
	return u.commentRepo.GetByStatus(ctx, "pending")
}

func (u *commentUsecase) GetUserComments(ctx context.Context, userID uint) ([]*entity.Comment, error) {
	if userID == 0 {
		return nil, errors.New("user ID is required")
	}

	return u.commentRepo.GetByUserID(ctx, userID)
}