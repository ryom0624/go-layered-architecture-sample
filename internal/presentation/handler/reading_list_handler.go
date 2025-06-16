package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"layered-architecture-template/internal/usecase"
)

type ReadingListHandler struct {
	readingListUsecase usecase.ReadingListUsecase
}

type CreateReadingListRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
}

type UpdateReadingListRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
}

type AddArticleToListRequest struct {
	Notes string `json:"notes"`
}

type UpdateReadingListItemRequest struct {
	Notes string `json:"notes"`
}

func NewReadingListHandler(readingListUsecase usecase.ReadingListUsecase) *ReadingListHandler {
	return &ReadingListHandler{
		readingListUsecase: readingListUsecase,
	}
}

func (h *ReadingListHandler) CreateReadingList(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID required"})
		return
	}

	var req CreateReadingListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	readingList, err := h.readingListUsecase.CreateReadingList(
		c.Request.Context(),
		req.Name,
		req.Description,
		userID.(uint),
		req.IsPublic,
	)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "reading list name is required" ||
			err.Error() == "reading list name must be 100 characters or less" ||
			err.Error() == "description must be 500 characters or less" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, readingList)
}

func (h *ReadingListHandler) GetUserReadingLists(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		// If no user_id provided, try to get from auth context
		if userID, exists := c.Get("userID"); exists {
			userIDStr = strconv.FormatUint(uint64(userID.(uint)), 10)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID required"})
			return
		}
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	readingLists, err := h.readingListUsecase.GetUserReadingLists(c.Request.Context(), uint(userID))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reading_lists": readingLists})
}

func (h *ReadingListHandler) GetReadingList(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		userID = uint(0) // Anonymous user for public lists
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reading list ID"})
		return
	}

	readingList, err := h.readingListUsecase.GetReadingList(c.Request.Context(), uint(id), userID.(uint))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "reading list not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "unauthorized: this reading list is private" {
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, readingList)
}

func (h *ReadingListHandler) UpdateReadingList(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID required"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reading list ID"})
		return
	}

	var req UpdateReadingListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	readingList, err := h.readingListUsecase.UpdateReadingList(
		c.Request.Context(),
		uint(id),
		req.Name,
		req.Description,
		userID.(uint),
		req.IsPublic,
	)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "reading list not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "unauthorized: you can only update your own reading lists" {
			statusCode = http.StatusForbidden
		} else if err.Error() == "reading list name is required" ||
			err.Error() == "reading list name must be 100 characters or less" ||
			err.Error() == "description must be 500 characters or less" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, readingList)
}

func (h *ReadingListHandler) DeleteReadingList(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID required"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reading list ID"})
		return
	}

	err = h.readingListUsecase.DeleteReadingList(c.Request.Context(), uint(id), userID.(uint))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "reading list not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "unauthorized: you can only delete your own reading lists" {
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "読書リストを削除しました"})
}

func (h *ReadingListHandler) GetPublicReadingLists(c *gin.Context) {
	readingLists, err := h.readingListUsecase.GetPublicReadingLists(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reading_lists": readingLists})
}

func (h *ReadingListHandler) AddArticleToList(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID required"})
		return
	}

	idStr := c.Param("id")
	readingListID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reading list ID"})
		return
	}

	var req AddArticleToListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get article ID from request body or URL
	articleIDStr := c.PostForm("article_id")
	if articleIDStr == "" {
		// Try to get from JSON body
		type ArticleRequest struct {
			ArticleID uint   `json:"article_id" binding:"required"`
			Notes     string `json:"notes"`
		}
		var articleReq ArticleRequest
		if err := c.ShouldBindJSON(&articleReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Article ID is required"})
			return
		}
		req.Notes = articleReq.Notes
		articleIDStr = strconv.FormatUint(uint64(articleReq.ArticleID), 10)
	}

	articleID, err := strconv.ParseUint(articleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	item, err := h.readingListUsecase.AddArticleToList(
		c.Request.Context(),
		uint(readingListID),
		uint(articleID),
		userID.(uint),
		req.Notes,
	)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "reading list not found" || err.Error() == "article not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "unauthorized: you can only add articles to your own reading lists" {
			statusCode = http.StatusForbidden
		} else if err.Error() == "article is already in this reading list" ||
			err.Error() == "notes must be 1000 characters or less" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "記事を読書リストに追加しました",
		"item":    item,
	})
}

func (h *ReadingListHandler) RemoveArticleFromList(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID required"})
		return
	}

	idStr := c.Param("id")
	readingListID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reading list ID"})
		return
	}

	articleIDStr := c.Param("article_id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	err = h.readingListUsecase.RemoveArticleFromList(
		c.Request.Context(),
		uint(readingListID),
		uint(articleID),
		userID.(uint),
	)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "reading list not found" || err.Error() == "article not found in this reading list" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "unauthorized: you can only remove articles from your own reading lists" {
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "記事を読書リストから削除しました"})
}

func (h *ReadingListHandler) GetReadingListArticles(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		userID = uint(0) // Anonymous user for public lists
	}

	idStr := c.Param("id")
	readingListID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reading list ID"})
		return
	}

	items, err := h.readingListUsecase.GetReadingListItems(c.Request.Context(), uint(readingListID), userID.(uint))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "reading list not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "unauthorized: this reading list is private" {
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"articles": items})
}

func (h *ReadingListHandler) UpdateReadingListItem(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID required"})
		return
	}

	idStr := c.Param("id")
	readingListID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reading list ID"})
		return
	}

	articleIDStr := c.Param("article_id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	var req UpdateReadingListItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.readingListUsecase.UpdateReadingListItem(
		c.Request.Context(),
		uint(readingListID),
		uint(articleID),
		userID.(uint),
		req.Notes,
	)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "reading list not found" || err.Error() == "article not found in this reading list" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "unauthorized: you can only update items in your own reading lists" {
			statusCode = http.StatusForbidden
		} else if err.Error() == "notes must be 1000 characters or less" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "読書リストアイテムを更新しました",
		"item":    item,
	})
}