package handler

import (
	"net/http"
	"strconv"
	"time"

	"layered-architecture-template/internal/usecase"

	"github.com/gin-gonic/gin"
)

type StatisticsHandler struct {
	statisticsUsecase usecase.StatisticsUsecase
}

func NewStatisticsHandler(statisticsUsecase usecase.StatisticsUsecase) *StatisticsHandler {
	return &StatisticsHandler{
		statisticsUsecase: statisticsUsecase,
	}
}

func (h *StatisticsHandler) GetArticleStatistics(c *gin.Context) {
	articleIDStr := c.Param("article_id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	stats, err := h.statisticsUsecase.GetArticleStatistics(c.Request.Context(), uint(articleID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *StatisticsHandler) RecalculateArticleStatistics(c *gin.Context) {
	articleIDStr := c.Param("article_id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	if err := h.statisticsUsecase.RecalculateArticleStatistics(c.Request.Context(), uint(articleID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article statistics recalculated successfully"})
}

func (h *StatisticsHandler) GetAllArticleStatistics(c *gin.Context) {
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

	stats, err := h.statisticsUsecase.GetAllArticleStatistics(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  stats,
		"page":  page,
		"limit": limit,
	})
}

func (h *StatisticsHandler) GetDailyStatistics(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	stats, err := h.statisticsUsecase.GetDailyStatistics(c.Request.Context(), date)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Daily statistics not found"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *StatisticsHandler) GetDailyStatisticsRange(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required"})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD"})
		return
	}

	stats, err := h.statisticsUsecase.GetDailyStatisticsRange(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       stats,
		"start_date": startDateStr,
		"end_date":   endDateStr,
	})
}

func (h *StatisticsHandler) GetPlatformOverview(c *gin.Context) {
	overview, err := h.statisticsUsecase.GetPlatformOverview(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, overview)
}

func (h *StatisticsHandler) GetUserAnalytics(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	analytics, err := h.statisticsUsecase.GetUserAnalytics(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}