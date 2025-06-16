package usecase

import (
	"context"
	"layered-architecture-template/internal/domain/entity"
	"layered-architecture-template/internal/domain/repository"
	"layered-architecture-template/pkg/constants"
)

type SearchUsecase interface {
	SearchArticles(ctx context.Context, params *entity.ArticleSearchParams) (*entity.ArticleSearchResult, error)
	GetPopularArticles(ctx context.Context, limit int) ([]*entity.Article, error)
	GetRecentArticles(ctx context.Context, limit int) ([]*entity.Article, error)
	ValidateSearchParams(params *entity.ArticleSearchParams) error
}

type searchUsecaseImpl struct {
	articleRepo repository.ArticleRepository
}

func NewSearchUsecase(articleRepo repository.ArticleRepository) SearchUsecase {
	return &searchUsecaseImpl{
		articleRepo: articleRepo,
	}
}

func (u *searchUsecaseImpl) SearchArticles(ctx context.Context, params *entity.ArticleSearchParams) (*entity.ArticleSearchResult, error) {
	// Set default values
	params.SetDefaults()
	
	// Validate parameters
	if err := params.Validate(); err != nil {
		return nil, err
	}
	
	// Perform search
	return u.articleRepo.Search(ctx, params)
}

func (u *searchUsecaseImpl) GetPopularArticles(ctx context.Context, limit int) ([]*entity.Article, error) {
	if limit <= 0 {
		limit = constants.DefaultItemsLimit
	}
	if limit > constants.MaxItemsLimit {
		limit = constants.MaxItemsLimit
	}
	
	return u.articleRepo.GetPopular(ctx, limit)
}

func (u *searchUsecaseImpl) GetRecentArticles(ctx context.Context, limit int) ([]*entity.Article, error) {
	if limit <= 0 {
		limit = constants.DefaultItemsLimit
	}
	if limit > constants.MaxItemsLimit {
		limit = constants.MaxItemsLimit
	}
	
	// Get recent articles by filtering with status = published and ordering by created_at desc
	params := &entity.ArticleSearchParams{
		Status:    "published",
		SortBy:    "created_at",
		SortOrder: "desc",
		Page:      1,
		Limit:     limit,
	}
	
	result, err := u.articleRepo.Search(ctx, params)
	if err != nil {
		return nil, err
	}
	
	return result.Articles, nil
}

func (u *searchUsecaseImpl) ValidateSearchParams(params *entity.ArticleSearchParams) error {
	if params == nil {
		return &entity.ValidationError{Field: "params", Message: "search parameters are required"}
	}
	
	params.SetDefaults()
	return params.Validate()
}