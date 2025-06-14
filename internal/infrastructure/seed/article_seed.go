package seed

import (
	"context"
	"log"

	"layered-architecture-template/internal/domain/entity"
	"layered-architecture-template/internal/domain/repository"
)

type ArticleSeeder struct {
	articleRepo repository.ArticleRepository
	userRepo    repository.UserRepository
}

func NewArticleSeeder(articleRepo repository.ArticleRepository, userRepo repository.UserRepository) *ArticleSeeder {
	return &ArticleSeeder{
		articleRepo: articleRepo,
		userRepo:    userRepo,
	}
}

func (s *ArticleSeeder) SeedArticles(ctx context.Context) error {
	// Get existing users to use as authors
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	if len(users) == 0 {
		log.Println("No users found. Please seed users first.")
		return nil
	}

	articles := []*entity.Article{
		{
			Title:    "Introduction to Clean Architecture",
			Content:  "Clean Architecture is a software design philosophy that separates the elements of a design into ring levels. The main rule of clean architecture is that code dependencies can only come from the outer levels inward. Code on the inner layers can have no knowledge of functions on the outer layers.",
			AuthorID: users[0].ID,
			Status:   "published",
		},
		{
			Title:    "Getting Started with Go",
			Content:  "Go is an open source programming language that makes it easy to build simple, reliable, and efficient software. This article will guide you through the basics of Go programming language and help you write your first Go program.",
			AuthorID: users[1].ID,
			Status:   "published",
		},
		{
			Title:    "Understanding Database Transactions",
			Content:  "Database transactions are a fundamental concept in database management systems. They ensure data integrity and consistency by grouping multiple operations into a single unit of work that either succeeds completely or fails completely.",
			AuthorID: users[2].ID,
			Status:   "published",
		},
		{
			Title:    "RESTful API Design Best Practices",
			Content:  "REST (Representational State Transfer) is an architectural style for designing networked applications. This article covers best practices for designing RESTful APIs including proper HTTP methods, status codes, and resource naming conventions.",
			AuthorID: users[0].ID,
			Status:   "published",
		},
		{
			Title:    "Microservices Architecture Patterns",
			Content:  "Microservices architecture is a method of developing software systems that focuses on building single-function modules with well-defined interfaces and operations. This article explores common patterns and practices in microservices design.",
			AuthorID: users[3].ID,
			Status:   "draft",
		},
		{
			Title:    "Testing Strategies for Modern Applications",
			Content:  "Testing is a crucial part of software development. This article discusses various testing strategies including unit testing, integration testing, and end-to-end testing, with practical examples and best practices.",
			AuthorID: users[1].ID,
			Status:   "published",
		},
		{
			Title:    "Docker Containerization Guide",
			Content:  "Docker is a platform that uses OS-level virtualization to deliver software in packages called containers. This comprehensive guide covers Docker basics, best practices, and real-world use cases.",
			AuthorID: users[4].ID,
			Status:   "published",
		},
		{
			Title:    "SOLID Principles in Software Design",
			Content:  "SOLID is an acronym for five design principles intended to make software designs more understandable, flexible, and maintainable. This article explains each principle with practical examples.",
			AuthorID: users[2].ID,
			Status:   "draft",
		},
		{
			Title:    "Advanced Git Workflows",
			Content:  "Git is a distributed version control system that tracks changes in any set of files. This article covers advanced Git workflows including feature branches, rebasing, and collaborative development strategies.",
			AuthorID: users[5].ID,
			Status:   "published",
		},
		{
			Title:    "Database Design and Normalization",
			Content:  "Database normalization is the process of structuring a relational database in accordance with a series of normal forms to reduce data redundancy and improve data integrity. This article explains the normalization process step by step.",
			AuthorID: users[3].ID,
			Status:   "published",
		},
		{
			Title:    "Introduction to Machine Learning",
			Content:  "Machine Learning is a subset of artificial intelligence that provides systems the ability to automatically learn and improve from experience without being explicitly programmed. This article introduces basic ML concepts and algorithms.",
			AuthorID: users[6].ID,
			Status:   "draft",
		},
		{
			Title:    "Cloud Computing Fundamentals",
			Content:  "Cloud computing is the delivery of computing services including servers, storage, databases, networking, software, analytics, and intelligence over the Internet. This article covers the fundamentals of cloud computing.",
			AuthorID: users[4].ID,
			Status:   "published",
		},
		{
			Title:    "Security Best Practices in Web Development",
			Content:  "Web application security is critical for protecting user data and maintaining trust. This article covers essential security practices including authentication, authorization, input validation, and protection against common vulnerabilities.",
			AuthorID: users[7].ID,
			Status:   "published",
		},
		{
			Title:    "Performance Optimization Techniques",
			Content:  "Application performance optimization is crucial for user experience and system efficiency. This article discusses various optimization techniques including caching, database optimization, and code profiling.",
			AuthorID: users[5].ID,
			Status:   "draft",
		},
		{
			Title:    "DevOps Culture and Practices",
			Content:  "DevOps is a set of practices that combines software development and IT operations. It aims to shorten the systems development life cycle and provide continuous delivery with high software quality.",
			AuthorID: users[8].ID,
			Status:   "published",
		},
	}

	for _, article := range articles {
		err := s.articleRepo.Create(ctx, article)
		if err != nil {
			log.Printf("Failed to create article '%s': %v", article.Title, err)
			return err
		}
		log.Printf("Created article: %s (Author ID: %d, Status: %s)", article.Title, article.AuthorID, article.Status)
	}

	log.Printf("Successfully seeded %d articles", len(articles))
	return nil
}

func (s *ArticleSeeder) GetSeedArticleTitles() []string {
	return []string{
		"Introduction to Clean Architecture",
		"Getting Started with Go",
		"Understanding Database Transactions",
		"RESTful API Design Best Practices",
		"Microservices Architecture Patterns",
		"Testing Strategies for Modern Applications",
		"Docker Containerization Guide",
		"SOLID Principles in Software Design",
		"Advanced Git Workflows",
		"Database Design and Normalization",
		"Introduction to Machine Learning",
		"Cloud Computing Fundamentals",
		"Security Best Practices in Web Development",
		"Performance Optimization Techniques",
		"DevOps Culture and Practices",
	}
}