package router

import (
	"layered-architecture-template/internal/presentation/handler"
	"layered-architecture-template/internal/presentation/middleware"
	"layered-architecture-template/internal/usecase"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *handler.UserHandler, articleHandler *handler.ArticleHandler, commentHandler *handler.CommentHandler, searchHandler *handler.SearchHandler, favoriteHandler *handler.FavoriteHandler, readingListHandler *handler.ReadingListHandler, viewHandler *handler.ViewHandler, statisticsHandler *handler.StatisticsHandler, authHandler *handler.AuthHandler, categoryHandler *handler.CategoryHandler, tagHandler *handler.TagHandler, viewUsecase usecase.ViewUsecase, authMiddleware gin.HandlerFunc, optionalAuthMiddleware gin.HandlerFunc) *gin.Engine {
	r := gin.Default()
	
	// Add view tracking middleware
	r.Use(middleware.ViewTrackingMiddleware(viewUsecase))

	api := r.Group("/api/v1")
	{
		// Authentication routes (no auth required)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", authHandler.Logout)
			auth.POST("/logout-all", authMiddleware, authHandler.LogoutAll)
		}
		users := api.Group("/users")
		{
			users.POST("", userHandler.CreateUser)
			users.GET("", userHandler.GetAllUsers)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
			users.GET("/:id/comments", commentHandler.GetUserComments)
			users.GET("/:id/favorites", favoriteHandler.GetUserFavorites)
			users.GET("/:id/reading-lists/public", readingListHandler.GetUserReadingLists)
			
			// Reading history and analytics routes
			users.GET("/:id/reading-history", viewHandler.GetUserReadingHistories)
			users.GET("/:id/analytics", statisticsHandler.GetUserAnalytics)
		}

		articles := api.Group("/articles")
		{
			articles.POST("", articleHandler.CreateArticle)
			articles.GET("", articleHandler.GetAllArticles)
			articles.GET("/published", articleHandler.GetPublishedArticles)
			articles.GET("/popular", searchHandler.GetPopularArticles)
			articles.GET("/recent", searchHandler.GetRecentArticles)
			articles.GET("/trending", viewHandler.GetTrendingArticles)
			articles.GET("/:id", articleHandler.GetArticle)
			articles.PUT("/:id", articleHandler.UpdateArticle)
			articles.DELETE("/:id", articleHandler.DeleteArticle)
			articles.PUT("/:id/publish", articleHandler.PublishArticle)
			articles.PUT("/:id/unpublish", articleHandler.UnpublishArticle)
			
			// Comment routes for articles
			articles.POST("/:id/comments", commentHandler.CreateCommentOnArticle)
			articles.GET("/:id/comments", commentHandler.GetArticleComments)
			
			// Favorite routes for articles
			articles.POST("/:id/favorite", favoriteHandler.AddFavorite)
			articles.DELETE("/:id/favorite", favoriteHandler.RemoveFavorite)
			articles.GET("/:id/favorites", favoriteHandler.GetArticleFavorites)
			articles.GET("/:id/favorite-status", favoriteHandler.CheckFavoriteStatus)
			
			// Statistics routes for articles
			articles.GET("/:id/statistics", statisticsHandler.GetArticleStatistics)
			articles.POST("/:id/statistics/recalculate", statisticsHandler.RecalculateArticleStatistics)
		}

		comments := api.Group("/comments")
		{
			comments.GET("/pending", commentHandler.GetPendingComments)
			comments.GET("/:id", commentHandler.GetComment)
			comments.PUT("/:id", commentHandler.UpdateComment)
			comments.DELETE("/:id", commentHandler.DeleteComment)
			comments.POST("/:id/replies", commentHandler.CreateReply)
			comments.PUT("/:id/approve", commentHandler.ApproveComment)
			comments.PUT("/:id/reject", commentHandler.RejectComment)
		}

		// Search routes
		api.GET("/search", searchHandler.SearchArticles)
		
		// Reading progress routes
		progress := api.Group("/progress")
		{
			progress.PUT("/users/:user_id/articles/:article_id", viewHandler.TrackReadingProgress)
			progress.GET("/users/:user_id/articles/:article_id", viewHandler.GetReadingHistory)
		}
		
		// Analytics routes
		analytics := api.Group("/analytics")
		{
			analytics.GET("/overview", statisticsHandler.GetPlatformOverview)
			analytics.GET("/daily", statisticsHandler.GetDailyStatistics)
			analytics.GET("/daily/range", statisticsHandler.GetDailyStatisticsRange)
			analytics.GET("/articles", statisticsHandler.GetAllArticleStatistics)
			analytics.GET("/articles/popular", viewHandler.GetPopularArticles)
		}
		
		// Reading list routes
		readingLists := api.Group("/reading-lists")
		{
			readingLists.POST("", readingListHandler.CreateReadingList)
			readingLists.GET("", readingListHandler.GetUserReadingLists)
			readingLists.GET("/public", readingListHandler.GetPublicReadingLists)
			readingLists.GET("/:id", readingListHandler.GetReadingList)
			readingLists.PUT("/:id", readingListHandler.UpdateReadingList)
			readingLists.DELETE("/:id", readingListHandler.DeleteReadingList)
			readingLists.POST("/:id/articles", readingListHandler.AddArticleToList)
			readingLists.GET("/:id/articles", readingListHandler.GetReadingListArticles)
			readingLists.DELETE("/:id/articles/:article_id", readingListHandler.RemoveArticleFromList)
			readingLists.PUT("/:id/articles/:article_id", readingListHandler.UpdateReadingListItem)
		}
		
		// Category routes
		categories := api.Group("/categories")
		{
			categories.POST("", categoryHandler.CreateCategory)
			categories.GET("", categoryHandler.GetAllCategories)
			categories.GET("/with-count", categoryHandler.GetCategoriesWithCount)
			categories.GET("/:slug", categoryHandler.GetCategory)
			categories.PUT("/:id", categoryHandler.UpdateCategory)
			categories.DELETE("/:id", categoryHandler.DeleteCategory)
			categories.GET("/:slug/articles", articleHandler.GetArticlesByCategory)
		}
		
		// Tag routes
		tags := api.Group("/tags")
		{
			tags.POST("", tagHandler.CreateTag)
			tags.GET("", tagHandler.GetAllTags)
			tags.GET("/popular", tagHandler.GetPopularTags)
			tags.GET("/:slug", tagHandler.GetTag)
			tags.PUT("/:id", tagHandler.UpdateTag)
			tags.DELETE("/:id", tagHandler.DeleteTag)
			tags.GET("/articles", articleHandler.GetArticlesByTags)
		}
	}

	return r
}
