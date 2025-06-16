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
		// Original users
		{Name: "John Doe", Email: "john.doe@example.com"},
		{Name: "Jane Smith", Email: "jane.smith@example.com"},
		{Name: "Alice Johnson", Email: "alice.johnson@example.com"},
		{Name: "Bob Wilson", Email: "bob.wilson@example.com"},
		{Name: "Carol Brown", Email: "carol.brown@example.com"},
		{Name: "David Miller", Email: "david.miller@example.com"},
		{Name: "Emma Davis", Email: "emma.davis@example.com"},
		{Name: "Frank Garcia", Email: "frank.garcia@example.com"},
		{Name: "Grace Martinez", Email: "grace.martinez@example.com"},
		{Name: "Henry Lee", Email: "henry.lee@example.com"},
		
		// Tech Professionals
		{Name: "Alex Chen", Email: "alex.chen@techcorp.com"},
		{Name: "Sarah Kumar", Email: "sarah.kumar@devstudio.io"},
		{Name: "Mike Rodriguez", Email: "mike.rodriguez@cloudnative.org"},
		{Name: "Lisa Wang", Email: "lisa.wang@airesearch.edu"},
		{Name: "James Thompson", Email: "james.thompson@opensource.dev"},
		{Name: "Maria Gonzalez", Email: "maria.gonzalez@dataanalytics.com"},
		{Name: "Kevin O'Connor", Email: "kevin.oconnor@cybersec.net"},
		{Name: "Nina Petrov", Email: "nina.petrov@blockchain.tech"},
		{Name: "Daniel Kim", Email: "daniel.kim@mobilefirst.app"},
		{Name: "Rachel Green", Email: "rachel.green@webdesign.studio"},
		
		// International Contributors
		{Name: "Yuki Tanaka", Email: "yuki.tanaka@example.jp"},
		{Name: "Pierre Dubois", Email: "pierre.dubois@example.fr"},
		{Name: "Anna Mueller", Email: "anna.mueller@example.de"},
		{Name: "Carlos Silva", Email: "carlos.silva@example.br"},
		{Name: "Priya Sharma", Email: "priya.sharma@example.in"},
		{Name: "Lars Andersson", Email: "lars.andersson@example.se"},
		{Name: "Sofia Rossi", Email: "sofia.rossi@example.it"},
		{Name: "Ahmed Hassan", Email: "ahmed.hassan@example.eg"},
		{Name: "Elena Volkov", Email: "elena.volkov@example.ru"},
		{Name: "Tom Anderson", Email: "tom.anderson@example.au"},
		
		// Specialists
		{Name: "Dr. Rebecca Foster", Email: "rebecca.foster@research.edu"},
		{Name: "Prof. Alan Wright", Email: "alan.wright@university.edu"},
		{Name: "Samantha Code", Email: "samantha.code@freelance.dev"},
		{Name: "Ryan Builder", Email: "ryan.builder@contractor.biz"},
		{Name: "Jessica Debug", Email: "jessica.debug@testing.qa"},
		{Name: "Mark Deploy", Email: "mark.deploy@devops.cloud"},
		{Name: "Lauren Script", Email: "lauren.script@automation.bot"},
		{Name: "Chris Secure", Email: "chris.secure@pentesting.sec"},
		{Name: "Amy Cloud", Email: "amy.cloud@infrastructure.aws"},
		{Name: "Ben Data", Email: "ben.data@analytics.ml"},
		
		// Creative Tech
		{Name: "Zoe Creative", Email: "zoe.creative@design.ux"},
		{Name: "Max Product", Email: "max.product@startup.io"},
		{Name: "Luna Mobile", Email: "luna.mobile@appdev.swift"},
		{Name: "Felix Game", Email: "felix.game@gamedev.unity"},
		{Name: "Ruby Web", Email: "ruby.web@frontend.react"},
		{Name: "Phoenix Backend", Email: "phoenix.backend@api.node"},
		{Name: "Storm DevOps", Email: "storm.devops@kubernetes.helm"},
		{Name: "Sage AI", Email: "sage.ai@machinelearning.py"},
		{Name: "River Blockchain", Email: "river.blockchain@crypto.eth"},
		{Name: "Sky Platform", Email: "sky.platform@microservices.go"},
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