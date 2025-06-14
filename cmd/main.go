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

	userRepo := repository.NewUserRepository(db.DB)

	userUsecase := usecase.NewUserUsecase(userRepo)

	userHandler := handler.NewUserHandler(userUsecase)

	r := router.SetupRouter(userHandler)

	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}