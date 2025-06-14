package seed

import (
	"context"
	"log"

	"layered-architecture-template/internal/domain/entity"
	"layered-architecture-template/internal/domain/repository"
)

type UserSeeder struct {
	userRepo repository.UserRepository
}

func NewUserSeeder(userRepo repository.UserRepository) *UserSeeder {
	return &UserSeeder{
		userRepo: userRepo,
	}
}

func (s *UserSeeder) SeedUsers(ctx context.Context) error {
	users := []*entity.User{
		{
			Name:  "John Doe",
			Email: "john.doe@example.com",
		},
		{
			Name:  "Jane Smith",
			Email: "jane.smith@example.com",
		},
		{
			Name:  "Alice Johnson",
			Email: "alice.johnson@example.com",
		},
		{
			Name:  "Bob Wilson",
			Email: "bob.wilson@example.com",
		},
		{
			Name:  "Carol Brown",
			Email: "carol.brown@example.com",
		},
		{
			Name:  "David Miller",
			Email: "david.miller@example.com",
		},
		{
			Name:  "Emma Davis",
			Email: "emma.davis@example.com",
		},
		{
			Name:  "Frank Garcia",
			Email: "frank.garcia@example.com",
		},
		{
			Name:  "Grace Martinez",
			Email: "grace.martinez@example.com",
		},
		{
			Name:  "Henry Lee",
			Email: "henry.lee@example.com",
		},
	}

	for _, user := range users {
		// Check if user already exists
		existingUser, _ := s.userRepo.GetByEmail(ctx, user.Email)
		if existingUser != nil {
			log.Printf("User with email %s already exists, skipping...", user.Email)
			continue
		}

		err := s.userRepo.Create(ctx, user)
		if err != nil {
			log.Printf("Failed to create user %s: %v", user.Email, err)
			return err
		}
		log.Printf("Created user: %s (%s)", user.Name, user.Email)
	}

	log.Printf("Successfully seeded %d users", len(users))
	return nil
}

func (s *UserSeeder) GetSeedUserEmails() []string {
	return []string{
		"john.doe@example.com",
		"jane.smith@example.com",
		"alice.johnson@example.com",
		"bob.wilson@example.com",
		"carol.brown@example.com",
		"david.miller@example.com",
		"emma.davis@example.com",
		"frank.garcia@example.com",
		"grace.martinez@example.com",
		"henry.lee@example.com",
	}
}