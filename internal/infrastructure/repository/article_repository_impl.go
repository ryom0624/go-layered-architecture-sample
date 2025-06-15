package repository

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"
	"layered-architecture-template/internal/domain/entity"
	"layered-architecture-template/internal/domain/repository"

	"gorm.io/gorm"
)

type articleRepositoryImpl struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) repository.ArticleRepository {
	return &articleRepositoryImpl{db: db}
}

func (r *articleRepositoryImpl) Create(ctx context.Context, article *entity.Article) error {
	return r.db.WithContext(ctx).Create(article).Error
}

func (r *articleRepositoryImpl) CreateWithTx(ctx context.Context, tx repository.Transaction, article *entity.Article) error {
	db := tx.GetDB().(*gorm.DB)
	return db.WithContext(ctx).Create(article).Error
}

func (r *articleRepositoryImpl) GetByID(ctx context.Context, id uint) (*entity.Article, error) {
	var article entity.Article
	err := r.db.WithContext(ctx).Preload("Author").First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepositoryImpl) GetAll(ctx context.Context) ([]*entity.Article, error) {
	var articles []*entity.Article
	err := r.db.WithContext(ctx).Preload("Author").Find(&articles).Error
	return articles, err
}

func (r *articleRepositoryImpl) GetByStatus(ctx context.Context, status string) ([]*entity.Article, error) {
	var articles []*entity.Article
	err := r.db.WithContext(ctx).Preload("Author").Where("status = ?", status).Find(&articles).Error
	return articles, err
}

func (r *articleRepositoryImpl) Update(ctx context.Context, article *entity.Article) error {
	return r.db.WithContext(ctx).Save(article).Error
}

func (r *articleRepositoryImpl) UpdateWithTx(ctx context.Context, tx repository.Transaction, article *entity.Article) error {
	db := tx.GetDB().(*gorm.DB)
	return db.WithContext(ctx).Save(article).Error
}

func (r *articleRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Article{}, id).Error
}

func (r *articleRepositoryImpl) DeleteWithTx(ctx context.Context, tx repository.Transaction, id uint) error {
	db := tx.GetDB().(*gorm.DB)
	return db.WithContext(ctx).Delete(&entity.Article{}, id).Error
}

// Search implements comprehensive article search and filtering
func (r *articleRepositoryImpl) Search(ctx context.Context, params *entity.ArticleSearchParams) (*entity.ArticleSearchResult, error) {
	query := r.db.WithContext(ctx).Model(&entity.Article{}).Preload("Author")
	
	// Apply filters
	query = r.applyFilters(query, params)
	
	// Count total records
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count articles: %w", err)
	}
	
	// Apply sorting
	query = r.applySorting(query, params)
	
	// Apply pagination
	offset := params.GetOffset()
	query = query.Offset(offset).Limit(params.Limit)
	
	// Execute query
	var articles []*entity.Article
	if err := query.Find(&articles).Error; err != nil {
		return nil, fmt.Errorf("failed to search articles: %w", err)
	}
	
	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))
	
	return &entity.ArticleSearchResult{
		Articles:   articles,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// SearchWithFilters searches articles with custom filters
func (r *articleRepositoryImpl) SearchWithFilters(ctx context.Context, searchQuery string, filters map[string]interface{}) ([]*entity.Article, error) {
	query := r.db.WithContext(ctx).Model(&entity.Article{}).Preload("Author")
	
	// Apply full-text search
	if searchQuery != "" {
		query = query.Where("title ILIKE ? OR content ILIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}
	
	// Apply custom filters
	for key, value := range filters {
		switch key {
		case "author_id":
			query = query.Where("author_id = ?", value)
		case "status":
			query = query.Where("status = ?", value)
		case "date_from":
			query = query.Where("created_at >= ?", value)
		case "date_to":
			query = query.Where("created_at <= ?", value)
		}
	}
	
	var articles []*entity.Article
	err := query.Find(&articles).Error
	return articles, err
}

// GetByDateRange retrieves articles within a date range
func (r *articleRepositoryImpl) GetByDateRange(ctx context.Context, from, to time.Time) ([]*entity.Article, error) {
	var articles []*entity.Article
	err := r.db.WithContext(ctx).Preload("Author").
		Where("created_at BETWEEN ? AND ?", from, to).
		Order("created_at DESC").
		Find(&articles).Error
	return articles, err
}

// GetPopular retrieves popular articles (for now, ordered by created_at, can be enhanced with view counts)
func (r *articleRepositoryImpl) GetPopular(ctx context.Context, limit int) ([]*entity.Article, error) {
	var articles []*entity.Article
	err := r.db.WithContext(ctx).Preload("Author").
		Where("status = ?", "published").
		Order("created_at DESC").
		Limit(limit).
		Find(&articles).Error
	return articles, err
}

// CountByFilters counts articles matching given filters
func (r *articleRepositoryImpl) CountByFilters(ctx context.Context, filters map[string]interface{}) (int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.Article{})
	
	for key, value := range filters {
		switch key {
		case "author_id":
			query = query.Where("author_id = ?", value)
		case "status":
			query = query.Where("status = ?", value)
		case "date_from":
			query = query.Where("created_at >= ?", value)
		case "date_to":
			query = query.Where("created_at <= ?", value)
		}
	}
	
	var count int64
	err := query.Count(&count).Error
	return count, err
}

// applyFilters applies search and filter conditions to the query
func (r *articleRepositoryImpl) applyFilters(query *gorm.DB, params *entity.ArticleSearchParams) *gorm.DB {
	// Full-text search
	if params.Query != "" {
		searchTerm := "%" + strings.ToLower(params.Query) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(content) LIKE ?", searchTerm, searchTerm)
	}
	
	// Author filter
	if params.AuthorID != 0 {
		query = query.Where("author_id = ?", params.AuthorID)
	}
	
	// Status filter
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	
	// Date range filter
	if !params.DateFrom.IsZero() {
		query = query.Where("created_at >= ?", params.DateFrom)
	}
	if !params.DateTo.IsZero() {
		query = query.Where("created_at <= ?", params.DateTo)
	}
	
	return query
}

// applySorting applies sorting to the query
func (r *articleRepositoryImpl) applySorting(query *gorm.DB, params *entity.ArticleSearchParams) *gorm.DB {
	switch params.SortBy {
	case "title":
		if params.SortOrder == "asc" {
			query = query.Order("title ASC")
		} else {
			query = query.Order("title DESC")
		}
	case "author":
		if params.SortOrder == "asc" {
			query = query.Joins("JOIN users ON users.id = articles.author_id").Order("users.name ASC")
		} else {
			query = query.Joins("JOIN users ON users.id = articles.author_id").Order("users.name DESC")
		}
	case "status":
		if params.SortOrder == "asc" {
			query = query.Order("status ASC")
		} else {
			query = query.Order("status DESC")
		}
	case "updated_at":
		if params.SortOrder == "asc" {
			query = query.Order("updated_at ASC")
		} else {
			query = query.Order("updated_at DESC")
		}
	case "relevance":
		// For relevance, we'd need more sophisticated full-text search
		// For now, fallback to created_at desc
		query = query.Order("created_at DESC")
	default: // created_at
		if params.SortOrder == "asc" {
			query = query.Order("created_at ASC")
		} else {
			query = query.Order("created_at DESC")
		}
	}
	
	return query
}