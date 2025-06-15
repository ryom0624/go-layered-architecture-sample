package handler

import (
	"net/http"
	"strconv"
	"time"
	"layered-architecture-template/internal/domain/entity"
	"layered-architecture-template/internal/usecase"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	searchUsecase usecase.SearchUsecase
}

func NewSearchHandler(searchUsecase usecase.SearchUsecase) *SearchHandler {
	return &SearchHandler{
		searchUsecase: searchUsecase,
	}
}

// SearchArticles handles article search requests
// @Summary Search articles
// @Description Search articles with various filters
// @Tags search
// @Accept json
// @Produce json
// @Param query query string false "Search query"
// @Param author_id query int false "Author ID filter"
// @Param status query string false "Status filter"
// @Param date_from query string false "Date from (YYYY-MM-DD)"
// @Param date_to query string false "Date to (YYYY-MM-DD)"
// @Param sort_by query string false "Sort field"
// @Param sort_order query string false "Sort order (asc/desc)"
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Success 200 {object} entity.ArticleSearchResult
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/search [get]
func (h *SearchHandler) SearchArticles(c *gin.Context) {
	params := &entity.ArticleSearchParams{}
	
	// Parse query parameters
	params.Query = c.Query("query")
	params.Status = c.Query("status")
	
	// Parse author_id
	if authorIDStr := c.Query("author_id"); authorIDStr != "" {
		if authorID, err := strconv.ParseUint(authorIDStr, 10, 32); err == nil {
			params.AuthorID = uint(authorID)
		}
	}
	
	// Parse date range
	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		if dateFrom, err := time.Parse("2006-01-02", dateFromStr); err == nil {
			params.DateFrom = dateFrom
		}
	}
	if dateToStr := c.Query("date_to"); dateToStr != "" {
		if dateTo, err := time.Parse("2006-01-02", dateToStr); err == nil {
			params.DateTo = dateTo
		}
	}
	
	// Parse sorting
	params.SortBy = c.Query("sort_by")
	params.SortOrder = c.Query("sort_order")
	
	// Parse pagination
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			params.Limit = limit
		}
	}
	
	// Perform search
	result, err := h.searchUsecase.SearchArticles(c.Request.Context(), params)
	if err != nil {
		if validationErr, ok := err.(*entity.ValidationError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, result)
}

// GetPopularArticles handles requests for popular articles
// @Summary Get popular articles
// @Description Get popular articles
// @Tags search
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Success 200 {array} entity.Article
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/articles/popular [get]
func (h *SearchHandler) GetPopularArticles(c *gin.Context) {
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	
	articles, err := h.searchUsecase.GetPopularArticles(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, articles)
}

// GetRecentArticles handles requests for recent articles
// @Summary Get recent articles
// @Description Get recent published articles
// @Tags search
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Success 200 {array} entity.Article
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/articles/recent [get]
func (h *SearchHandler) GetRecentArticles(c *gin.Context) {
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	
	articles, err := h.searchUsecase.GetRecentArticles(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, articles)
}