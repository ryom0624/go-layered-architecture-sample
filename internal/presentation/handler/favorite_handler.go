package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"layered-architecture-template/internal/usecase"
)

type FavoriteHandler struct {
	favoriteUsecase usecase.FavoriteUsecase
}

func NewFavoriteHandler(favoriteUsecase usecase.FavoriteUsecase) *FavoriteHandler {
	return &FavoriteHandler{
		favoriteUsecase: favoriteUsecase,
	}
}

func (h *FavoriteHandler) AddFavorite(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID required"})
		return
	}

	articleIDStr := c.Param("id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	favorite, err := h.favoriteUsecase.AddFavorite(c.Request.Context(), userID.(uint), uint(articleID))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" || err.Error() == "article not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "cannot favorite unpublished article" || err.Error() == "article already favorited" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "記事をお気に入りに追加しました",
		"favorite": favorite,
	})
}

func (h *FavoriteHandler) RemoveFavorite(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID required"})
		return
	}

	articleIDStr := c.Param("id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	err = h.favoriteUsecase.RemoveFavorite(c.Request.Context(), userID.(uint), uint(articleID))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "favorite not found" || err.Error() == "article not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "お気に入りから削除しました"})
}

func (h *FavoriteHandler) GetUserFavorites(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	favorites, err := h.favoriteUsecase.GetUserFavorites(c.Request.Context(), uint(userID))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"favorites": favorites})
}

func (h *FavoriteHandler) GetArticleFavorites(c *gin.Context) {
	articleIDStr := c.Param("id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	favorites, err := h.favoriteUsecase.GetArticleFavorites(c.Request.Context(), uint(articleID))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "article not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"favorites": favorites})
}

func (h *FavoriteHandler) CheckFavoriteStatus(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID required"})
		return
	}

	articleIDStr := c.Param("id")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	isFavorited, err := h.favoriteUsecase.IsFavorited(c.Request.Context(), userID.(uint), uint(articleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"is_favorited": isFavorited})
}