package seed

import (
	"context"
	"log"

	"layered-architecture-template/internal/domain/repository"
)

type Seeder struct {
	userSeeder    *UserSeeder
	articleSeeder *ArticleSeeder
}

func NewSeeder(userRepo repository.UserRepository, articleRepo repository.ArticleRepository) *Seeder {
	return &Seeder{
		userSeeder:    NewUserSeeder(userRepo),
		articleSeeder: NewArticleSeeder(articleRepo, userRepo),
	}
}

func (s *Seeder) SeedAll(ctx context.Context) error {
	log.Println("Starting database seeding...")

	// Seed users first (articles depend on users)
	log.Println("Seeding users...")
	if err := s.userSeeder.SeedUsers(ctx); err != nil {
		log.Printf("Failed to seed users: %v", err)
		return err
	}

	// Seed articles
	log.Println("Seeding articles...")
	if err := s.articleSeeder.SeedArticles(ctx); err != nil {
		log.Printf("Failed to seed articles: %v", err)
		return err
	}

	log.Println("Database seeding completed successfully!")
	return nil
}

func (s *Seeder) SeedUsers(ctx context.Context) error {
	return s.userSeeder.SeedUsers(ctx)
}

func (s *Seeder) SeedArticles(ctx context.Context) error {
	return s.articleSeeder.SeedArticles(ctx)
}