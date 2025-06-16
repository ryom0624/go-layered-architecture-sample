package usecase

import (
	"context"
	"errors"
	"time"

	"layered-architecture-template/internal/domain/entity"
	"layered-architecture-template/internal/domain/repository"
)

type FavoriteUsecase interface {
	AddFavorite(ctx context.Context, userID, articleID uint) (*entity.Favorite, error)
	RemoveFavorite(ctx context.Context, userID, articleID uint) error
	GetUserFavorites(ctx context.Context, userID uint) ([]*entity.Favorite, error)
	GetArticleFavorites(ctx context.Context, articleID uint) ([]*entity.Favorite, error)
	IsFavorited(ctx context.Context, userID, articleID uint) (bool, error)
}

type favoriteUsecaseImpl struct {
	favoriteRepo      repository.FavoriteRepository
	articleRepo       repository.ArticleRepository
	userRepo          repository.UserRepository
	transactionManager repository.TransactionManager
}

func NewFavoriteUsecase(
	favoriteRepo repository.FavoriteRepository,
	articleRepo repository.ArticleRepository,
	userRepo repository.UserRepository,
	transactionManager repository.TransactionManager,
) FavoriteUsecase {
	return &favoriteUsecaseImpl{
		favoriteRepo:       favoriteRepo,
		articleRepo:        articleRepo,
		userRepo:           userRepo,
		transactionManager: transactionManager,
	}
}

func (u *favoriteUsecaseImpl) AddFavorite(ctx context.Context, userID, articleID uint) (*entity.Favorite, error) {
	// Check if user exists
	if _, err := u.userRepo.GetByID(ctx, userID); err != nil {
		return nil, errors.New("user not found")
	}

	// Check if article exists and is published
	article, err := u.articleRepo.GetByID(ctx, articleID)
	if err != nil {
		return nil, errors.New("article not found")
	}
	if article.Status != "published" {
		return nil, errors.New("cannot favorite unpublished article")
	}

	// Check if already favorited
	isFavorited, err := u.favoriteRepo.IsFavorited(ctx, userID, articleID)
	if err != nil {
		return nil, err
	}
	if isFavorited {
		return nil, errors.New("article already favorited")
	}

	// Create favorite and update article favorite count in transaction
	var favorite *entity.Favorite
	err = u.transactionManager.WithTransaction(ctx, func(tx repository.Transaction) error {
		favorite = &entity.Favorite{
			UserID:    userID,
			ArticleID: articleID,
			CreatedAt: time.Now(),
		}

		if err := u.favoriteRepo.CreateWithTx(ctx, tx, favorite); err != nil {
			return err
		}

		// Update article favorite count
		article.FavoriteCount++
		return u.articleRepo.UpdateWithTx(ctx, tx, article)
	})

	if err != nil {
		return nil, err
	}

	return favorite, nil
}

func (u *favoriteUsecaseImpl) RemoveFavorite(ctx context.Context, userID, articleID uint) error {
	// Check if favorite exists
	isFavorited, err := u.favoriteRepo.IsFavorited(ctx, userID, articleID)
	if err != nil {
		return err
	}
	if !isFavorited {
		return errors.New("favorite not found")
	}

	// Get article for updating favorite count
	article, err := u.articleRepo.GetByID(ctx, articleID)
	if err != nil {
		return errors.New("article not found")
	}

	// Remove favorite and update article favorite count in transaction
	return u.transactionManager.WithTransaction(ctx, func(tx repository.Transaction) error {
		if err := u.favoriteRepo.DeleteWithTx(ctx, tx, userID, articleID); err != nil {
			return err
		}

		// Update article favorite count
		if article.FavoriteCount > 0 {
			article.FavoriteCount--
		}
		return u.articleRepo.UpdateWithTx(ctx, tx, article)
	})
}

func (u *favoriteUsecaseImpl) GetUserFavorites(ctx context.Context, userID uint) ([]*entity.Favorite, error) {
	// Check if user exists
	if _, err := u.userRepo.GetByID(ctx, userID); err != nil {
		return nil, errors.New("user not found")
	}

	return u.favoriteRepo.GetByUserID(ctx, userID)
}

func (u *favoriteUsecaseImpl) GetArticleFavorites(ctx context.Context, articleID uint) ([]*entity.Favorite, error) {
	// Check if article exists
	if _, err := u.articleRepo.GetByID(ctx, articleID); err != nil {
		return nil, errors.New("article not found")
	}

	return u.favoriteRepo.GetByArticleID(ctx, articleID)
}

func (u *favoriteUsecaseImpl) IsFavorited(ctx context.Context, userID, articleID uint) (bool, error) {
	return u.favoriteRepo.IsFavorited(ctx, userID, articleID)
}