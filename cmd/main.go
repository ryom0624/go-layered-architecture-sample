package main

import (
	"log"

	"layered-architecture-template/internal/infrastructure/database"
	"layered-architecture-template/internal/infrastructure/repository"
	"layered-architecture-template/internal/presentation/handler"
	"layered-architecture-template/internal/presentation/middleware"
	"layered-architecture-template/internal/presentation/router"
	"layered-architecture-template/internal/usecase"
	"layered-architecture-template/pkg/config"
	"layered-architecture-template/pkg/jwt"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	transactionManager := database.NewGormTransactionManager(db.DB)
	
	userRepo := repository.NewUserRepository(db.DB)
	articleRepo := repository.NewArticleRepository(db.DB)
	commentRepo := repository.NewCommentRepository(db.DB)
	favoriteRepo := repository.NewFavoriteRepository(db.DB)
	readingListRepo := repository.NewReadingListRepository(db.DB)
	viewRepo := repository.NewViewRepository(db.DB)
	statisticsRepo := repository.NewStatisticsRepository(db.DB)
	authRepo := repository.NewAuthRepository(db.DB)
	categoryRepo := repository.NewCategoryRepository(db.DB)
	tagRepo := repository.NewTagRepository(db.DB)

	jwtManager := jwt.NewJWTManager(cfg.Auth.JWTSecret, cfg.Auth.AccessTokenDuration)
	
	userUsecase := usecase.NewUserUsecase(userRepo)
	articleUsecase := usecase.NewArticleUsecase(articleRepo, userRepo, categoryRepo, tagRepo, transactionManager)
	commentUsecase := usecase.NewCommentUsecase(commentRepo, userRepo, articleRepo, transactionManager)
	searchUsecase := usecase.NewSearchUsecase(articleRepo)
	favoriteUsecase := usecase.NewFavoriteUsecase(favoriteRepo, articleRepo, userRepo, transactionManager)
	readingListUsecase := usecase.NewReadingListUsecase(readingListRepo, articleRepo, userRepo)
	viewUsecase := usecase.NewViewUsecase(viewRepo, articleRepo)
	statisticsUsecase := usecase.NewStatisticsUsecase(statisticsRepo, viewRepo, userRepo, articleRepo)
	authUsecase := usecase.NewAuthUsecase(authRepo, jwtManager, cfg.Auth.RefreshTokenDuration)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
	tagUsecase := usecase.NewTagUsecase(tagRepo)

	userHandler := handler.NewUserHandler(userUsecase)
	articleHandler := handler.NewArticleHandler(articleUsecase)
	commentHandler := handler.NewCommentHandler(commentUsecase)
	searchHandler := handler.NewSearchHandler(searchUsecase)
	favoriteHandler := handler.NewFavoriteHandler(favoriteUsecase)
	readingListHandler := handler.NewReadingListHandler(readingListUsecase)
	viewHandler := handler.NewViewHandler(viewUsecase)
	statisticsHandler := handler.NewStatisticsHandler(statisticsUsecase)
	authHandler := handler.NewAuthHandler(authUsecase)
	categoryHandler := handler.NewCategoryHandler(categoryUsecase)
	tagHandler := handler.NewTagHandler(tagUsecase)
	
	authMiddleware := middleware.AuthMiddleware(jwtManager)
	optionalAuthMiddleware := middleware.OptionalAuthMiddleware(jwtManager)

	r := router.SetupRouter(userHandler, articleHandler, commentHandler, searchHandler, favoriteHandler, readingListHandler, viewHandler, statisticsHandler, authHandler, categoryHandler, tagHandler, viewUsecase, authMiddleware, optionalAuthMiddleware)

	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}