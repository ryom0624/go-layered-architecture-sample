package middleware

import (
	"strconv"
	"strings"

	"layered-architecture-template/internal/usecase"

	"github.com/gin-gonic/gin"
)

func ViewTrackingMiddleware(viewUsecase usecase.ViewUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" && strings.Contains(c.Request.URL.Path, "/articles/") {
			articleIDStr := c.Param("id")
			if articleIDStr != "" {
				if articleID, err := strconv.ParseUint(articleIDStr, 10, 32); err == nil {
					go func() {
						userIDHeader := c.GetHeader("X-User-ID")
						var userID *uint
						if userIDHeader != "" {
							if uid, err := strconv.ParseUint(userIDHeader, 10, 32); err == nil {
								uidVal := uint(uid)
								userID = &uidVal
							}
						}
						
						ipAddress := c.ClientIP()
						userAgent := c.GetHeader("User-Agent")
						
						viewUsecase.TrackView(c.Request.Context(), uint(articleID), userID, ipAddress, userAgent)
					}()
				}
			}
		}
		
		c.Next()
	}
}