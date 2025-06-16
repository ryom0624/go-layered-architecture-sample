package router

import (
	"layered-architecture-template/internal/presentation/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *handler.UserHandler, articleHandler *handler.ArticleHandler, commentHandler *handler.CommentHandler, searchHandler *handler.SearchHandler, favoriteHandler *handler.FavoriteHandler, readingListHandler *handler.ReadingListHandler) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
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
		}

		articles := api.Group("/articles")
		{
			articles.POST("", articleHandler.CreateArticle)
			articles.GET("", articleHandler.GetAllArticles)
			articles.GET("/published", articleHandler.GetPublishedArticles)
			articles.GET("/popular", searchHandler.GetPopularArticles)
			articles.GET("/recent", searchHandler.GetRecentArticles)
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
	}

	return r
}
