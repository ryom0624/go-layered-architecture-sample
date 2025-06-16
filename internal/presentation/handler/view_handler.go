package handler

import (
	"net/http"
	"strconv"

	"layered-architecture-template/internal/usecase"

	"github.com/gin-gonic/gin"
)

type ViewHandler struct {
	viewUsecase usecase.ViewUsecase
}

func NewViewHandler(viewUsecase usecase.ViewUsecase) *ViewHandler {
	return &ViewHandler{
		viewUsecase: viewUsecase,
	}
}

type TrackReadingProgressRequest struct {
	Progress    float32 `json:"progress" binding:"required,min=0,max=100"`
	ReadingTime int     `json:"reading_time" binding:"required,min=0"`
}

func (h *ViewHandler) TrackReadingProgress(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	articleIDStr := c.Param("article_id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	var req TrackReadingProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.viewUsecase.TrackReadingProgress(c.Request.Context(), uint(userID), uint(articleID), req.Progress, req.ReadingTime); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reading progress updated successfully"})
}

func (h *ViewHandler) GetReadingHistory(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	articleIDStr := c.Param("article_id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	history, err := h.viewUsecase.GetReadingHistory(c.Request.Context(), uint(userID), uint(articleID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reading history not found"})
		return
	}

	c.JSON(http.StatusOK, history)
}

func (h *ViewHandler) GetUserReadingHistories(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	histories, err := h.viewUsecase.GetUserReadingHistories(c.Request.Context(), uint(userID), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  histories,
		"page":  page,
		"limit": limit,
	})
}

func (h *ViewHandler) GetPopularArticles(c *gin.Context) {
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	days := 7
	if daysStr := c.Query("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	articles, err := h.viewUsecase.GetPopularArticles(c.Request.Context(), limit, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  articles,
		"limit": limit,
		"days":  days,
	})
}

func (h *ViewHandler) GetTrendingArticles(c *gin.Context) {
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	hours := 24
	if hoursStr := c.Query("hours"); hoursStr != "" {
		if h, err := strconv.Atoi(hoursStr); err == nil && h > 0 {
			hours = h
		}
	}

	articles, err := h.viewUsecase.GetTrendingArticles(c.Request.Context(), limit, hours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  articles,
		"limit": limit,
		"hours": hours,
	})
}