package router

import (
	"layered-architecture-template/internal/presentation/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *handler.UserHandler, articleHandler *handler.ArticleHandler) *gin.Engine {
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
		}

		articles := api.Group("/articles")
		{
			articles.POST("", articleHandler.CreateArticle)
			articles.GET("", articleHandler.GetAllArticles)
			articles.GET("/published", articleHandler.GetPublishedArticles)
			articles.GET("/:id", articleHandler.GetArticle)
			articles.PUT("/:id", articleHandler.UpdateArticle)
			articles.DELETE("/:id", articleHandler.DeleteArticle)
			articles.PUT("/:id/publish", articleHandler.PublishArticle)
			articles.PUT("/:id/unpublish", articleHandler.UnpublishArticle)
		}
	}

	return r
}
