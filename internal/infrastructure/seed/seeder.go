package seed

import (
	"context"
	"log"

	"layered-architecture-template/internal/domain/repository"
	"gorm.io/gorm"
)

type Seeder struct {
	userSeeder    *UserSeeder
	articleSeeder *ArticleSeeder
	db            *gorm.DB
}

func NewSeeder(userRepo repository.UserRepository, articleRepo repository.ArticleRepository, db *gorm.DB) *Seeder {
	return &Seeder{
		userSeeder:    NewUserSeeder(userRepo),
		articleSeeder: NewArticleSeeder(articleRepo, userRepo),
		db:            db,
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

	// Seed comments (comments depend on users and articles)
	log.Println("Seeding comments...")
	if err := SeedComments(s.db); err != nil {
		log.Printf("Failed to seed comments: %v", err)
		return err
	}

	// Seed enhanced comments for better discussions
	log.Println("Seeding enhanced comments...")
	if err := SeedEnhancedComments(s.db); err != nil {
		log.Printf("Failed to seed enhanced comments: %v", err)
		return err
	}

	// Seed favorites
	log.Println("Seeding favorites...")
	if err := SeedFavorites(s.db); err != nil {
		log.Printf("Failed to seed favorites: %v", err)
		return err
	}

	// Seed reading lists
	log.Println("Seeding reading lists...")
	if err := SeedReadingLists(s.db); err != nil {
		log.Printf("Failed to seed reading lists: %v", err)
		return err
	}

	// Seed view history and reading progress
	log.Println("Seeding view history...")
	if err := SeedViewHistory(s.db); err != nil {
		log.Printf("Failed to seed view history: %v", err)
		return err
	}

	// Seed statistics
	log.Println("Seeding statistics...")
	if err := SeedStatistics(s.db); err != nil {
		log.Printf("Failed to seed statistics: %v", err)
		return err
	}

	// Seed authentication tokens
	log.Println("Seeding auth tokens...")
	if err := SeedAuthTokens(s.db); err != nil {
		log.Printf("Failed to seed auth tokens: %v", err)
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