package usecase

import (
	"context"
	"errors"
	"layered-architecture-template/internal/domain/entity"
	"layered-architecture-template/internal/domain/repository"
)

type ArticleUsecase interface {
	CreateArticle(ctx context.Context, title, content string, authorID uint) (*entity.Article, error)
	GetArticle(ctx context.Context, id uint) (*entity.Article, error)
	GetAllArticles(ctx context.Context) ([]*entity.Article, error)
	GetAllArticlesWithFilters(ctx context.Context, params *entity.ArticleSearchParams) (*entity.ArticleSearchResult, error)
	GetPublishedArticles(ctx context.Context) ([]*entity.Article, error)
	UpdateArticle(ctx context.Context, id uint, title, content string) (*entity.Article, error)
	DeleteArticle(ctx context.Context, id uint) error
	PublishArticle(ctx context.Context, id uint) (*entity.Article, error)
	UnpublishArticle(ctx context.Context, id uint) (*entity.Article, error)
}

type articleUsecase struct {
	articleRepo        repository.ArticleRepository
	userRepo           repository.UserRepository
	transactionManager repository.TransactionManager
}

func NewArticleUsecase(
	articleRepo repository.ArticleRepository,
	userRepo repository.UserRepository,
	transactionManager repository.TransactionManager,
) ArticleUsecase {
	return &articleUsecase{
		articleRepo:        articleRepo,
		userRepo:           userRepo,
		transactionManager: transactionManager,
	}
}

func (u *articleUsecase) CreateArticle(ctx context.Context, title, content string, authorID uint) (*entity.Article, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}
	if content == "" {
		return nil, errors.New("content is required")
	}
	if authorID == 0 {
		return nil, errors.New("author ID is required")
	}

	_, err := u.userRepo.GetByID(ctx, authorID)
	if err != nil {
		return nil, errors.New("author not found")
	}

	article := &entity.Article{
		Title:    title,
		Content:  content,
		AuthorID: authorID,
		Status:   "draft",
	}

	err = u.articleRepo.Create(ctx, article)
	if err != nil {
		return nil, err
	}

	return article, nil
}

func (u *articleUsecase) GetArticle(ctx context.Context, id uint) (*entity.Article, error) {
	return u.articleRepo.GetByID(ctx, id)
}

func (u *articleUsecase) GetAllArticles(ctx context.Context) ([]*entity.Article, error) {
	return u.articleRepo.GetAll(ctx)
}

func (u *articleUsecase) GetPublishedArticles(ctx context.Context) ([]*entity.Article, error) {
	return u.articleRepo.GetByStatus(ctx, "published")
}

func (u *articleUsecase) UpdateArticle(ctx context.Context, id uint, title, content string) (*entity.Article, error) {
	if id == 0 {
		return nil, errors.New("article ID is required")
	}

	article, err := u.articleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("article not found")
	}

	if title != "" {
		article.Title = title
	}
	if content != "" {
		article.Content = content
	}

	err = u.articleRepo.Update(ctx, article)
	if err != nil {
		return nil, err
	}

	return article, nil
}

func (u *articleUsecase) DeleteArticle(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("article ID is required")
	}

	_, err := u.articleRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("article not found")
	}

	return u.articleRepo.Delete(ctx, id)
}

func (u *articleUsecase) PublishArticle(ctx context.Context, id uint) (*entity.Article, error) {
	var result *entity.Article
	err := u.transactionManager.WithTransaction(ctx, func(tx repository.Transaction) error {
		article, err := u.articleRepo.GetByID(ctx, id)
		if err != nil {
			return errors.New("article not found")
		}

		if article.Status == "published" {
			return errors.New("article is already published")
		}

		article.Status = "published"
		if err := u.articleRepo.UpdateWithTx(ctx, tx, article); err != nil {
			return err
		}

		result = article
		return nil
	})

	return result, err
}

func (u *articleUsecase) UnpublishArticle(ctx context.Context, id uint) (*entity.Article, error) {
	var result *entity.Article
	err := u.transactionManager.WithTransaction(ctx, func(tx repository.Transaction) error {
		article, err := u.articleRepo.GetByID(ctx, id)
		if err != nil {
			return errors.New("article not found")
		}

		if article.Status == "draft" {
			return errors.New("article is already unpublished")
		}

		article.Status = "draft"
		if err := u.articleRepo.UpdateWithTx(ctx, tx, article); err != nil {
			return err
		}

		result = article
		return nil
	})

	return result, err
}

func (u *articleUsecase) GetAllArticlesWithFilters(ctx context.Context, params *entity.ArticleSearchParams) (*entity.ArticleSearchResult, error) {
	if params == nil {
		// If no filters provided, return all articles with default pagination
		params = &entity.ArticleSearchParams{
			Page:  1,
			Limit: 20,
		}
	}
	
	// Set defaults and validate
	params.SetDefaults()
	if err := params.Validate(); err != nil {
		return nil, err
	}
	
	// Use the repository search method
	return u.articleRepo.Search(ctx, params)
}