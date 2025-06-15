package main

import (
	"context"
	"log"

	"layered-architecture-template/internal/infrastructure/database"
	"layered-architecture-template/internal/infrastructure/repository"
	"layered-architecture-template/internal/infrastructure/seed"
	"layered-architecture-template/pkg/config"
)

func main() {
	log.Println("Starting seed command...")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB)
	articleRepo := repository.NewArticleRepository(db.DB)

	// Initialize seeder
	seeder := seed.NewSeeder(userRepo, articleRepo, db.DB)

	// Run seeding
	ctx := context.Background()
	if err := seeder.SeedAll(ctx); err != nil {
		log.Fatal("Failed to seed database:", err)
	}

	log.Println("Seed command completed successfully!")
}