package entity

import (
	"time"
	"layered-architecture-template/pkg/constants"
)

// ArticleSearchParams defines parameters for searching and filtering articles
type ArticleSearchParams struct {
	Query     string    `json:"query" form:"query"`           // Full-text search query
	AuthorID  uint      `json:"author_id" form:"author_id"`   // Author filter
	Status    string    `json:"status" form:"status"`         // Status filter
	DateFrom  time.Time `json:"date_from" form:"date_from"`   // Date range start
	DateTo    time.Time `json:"date_to" form:"date_to"`       // Date range end
	SortBy    string    `json:"sort_by" form:"sort_by"`       // Sort field
	SortOrder string    `json:"sort_order" form:"sort_order"` // asc/desc
	Page      int       `json:"page" form:"page"`             // Page number
	Limit     int       `json:"limit" form:"limit"`           // Page size
}

// ArticleSearchResult represents the result of article search
type ArticleSearchResult struct {
	Articles   []*Article `json:"articles"`
	Total      int64      `json:"total"`
	Page       int        `json:"page"`
	Limit      int        `json:"limit"`
	TotalPages int        `json:"total_pages"`
}

// SetDefaults sets default values for search parameters
func (p *ArticleSearchParams) SetDefaults() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = constants.DefaultPageSize
	} else if p.Limit > constants.MaxPageSize {
		p.Limit = constants.MaxPageSize
	}
	if p.SortBy == "" {
		if p.Query != "" {
			p.SortBy = "relevance"
		} else {
			p.SortBy = "created_at"
		}
	}
	if p.SortOrder == "" {
		p.SortOrder = "desc"
	}
}

// Validate validates search parameters
func (p *ArticleSearchParams) Validate() error {
	if p.Query != "" && len(p.Query) > constants.MaxSearchQueryLength {
		return &ValidationError{Field: "query", Message: "query must be 100 characters or less"}
	}
	if p.Page < 1 {
		return &ValidationError{Field: "page", Message: "page must be 1 or greater"}
	}
	if p.Limit < 1 || p.Limit > constants.MaxPageSize {
		return &ValidationError{Field: "limit", Message: "limit must be between 1 and 100"}
	}
	if p.SortOrder != "" && p.SortOrder != "asc" && p.SortOrder != "desc" {
		return &ValidationError{Field: "sort_order", Message: "sort_order must be 'asc' or 'desc'"}
	}
	if !p.DateFrom.IsZero() && !p.DateTo.IsZero() && p.DateFrom.After(p.DateTo) {
		return &ValidationError{Field: "date_range", Message: "date_from must be before or equal to date_to"}
	}
	
	// Validate sort_by field
	validSortFields := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"title":      true,
		"author":     true,
		"status":     true,
		"relevance":  true,
	}
	if p.SortBy != "" && !validSortFields[p.SortBy] {
		return &ValidationError{Field: "sort_by", Message: "invalid sort field"}
	}
	
	return nil
}

// GetOffset calculates the offset for pagination
func (p *ArticleSearchParams) GetOffset() int {
	return (p.Page - 1) * p.Limit
}