package main

import (
	"log"

	"layered-architecture-template/internal/infrastructure/database"
	"layered-architecture-template/internal/infrastructure/repository"
	"layered-architecture-template/internal/presentation/handler"
	"layered-architecture-template/internal/presentation/router"
	"layered-architecture-template/internal/usecase"
	"layered-architecture-template/pkg/config"
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

	userUsecase := usecase.NewUserUsecase(userRepo)
	articleUsecase := usecase.NewArticleUsecase(articleRepo, userRepo, transactionManager)
	commentUsecase := usecase.NewCommentUsecase(commentRepo, userRepo, articleRepo, transactionManager)
	searchUsecase := usecase.NewSearchUsecase(articleRepo)

	userHandler := handler.NewUserHandler(userUsecase)
	articleHandler := handler.NewArticleHandler(articleUsecase)
	commentHandler := handler.NewCommentHandler(commentUsecase)
	searchHandler := handler.NewSearchHandler(searchUsecase)

	r := router.SetupRouter(userHandler, articleHandler, commentHandler, searchHandler)

	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}