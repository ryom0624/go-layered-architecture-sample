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
		// AI/ML Series
		{
			Title:    "Deep Learning with TensorFlow 2.0",
			Content:  "TensorFlow 2.0 brings eager execution and simplified APIs for building neural networks. This comprehensive guide covers building your first deep learning model with practical examples and best practices for production deployment.",
			AuthorID: users[13].ID, // Lisa Wang
			Status:   "published",
		},
		{
			Title:    "Natural Language Processing Fundamentals",
			Content:  "NLP is transforming how we interact with technology. Learn about tokenization, word embeddings, transformers, and how to build your own language models using modern frameworks like Hugging Face.",
			AuthorID: users[47].ID, // Sage AI
			Status:   "published",
		},
		{
			Title:    "Computer Vision with OpenCV and PyTorch",
			Content:  "Computer vision enables machines to interpret visual information. This article covers image processing, object detection, facial recognition, and real-time video analysis using OpenCV and PyTorch.",
			AuthorID: users[13].ID,
			Status:   "draft",
		},
		{
			Title:    "MLOps: Production Machine Learning Pipelines",
			Content:  "Taking machine learning models from research to production requires robust MLOps practices. Learn about model versioning, automated training, deployment strategies, and monitoring in production.",
			AuthorID: users[39].ID, // Ben Data
			Status:   "published",
		},
		
		// Web3/Blockchain Series
		{
			Title:    "Smart Contract Development with Solidity",
			Content:  "Smart contracts are self-executing contracts with terms directly written into code. This guide covers Solidity fundamentals, security best practices, and deploying contracts on Ethereum.",
			AuthorID: users[17].ID, // Nina Petrov
			Status:   "published",
		},
		{
			Title:    "DeFi Protocols and Yield Farming",
			Content:  "Decentralized Finance (DeFi) is revolutionizing traditional finance. Explore automated market makers, liquidity pools, yield farming strategies, and risk management in DeFi protocols.",
			AuthorID: users[48].ID, // River Blockchain
			Status:   "published",
		},
		{
			Title:    "NFT Marketplace Development",
			Content:  "Non-Fungible Tokens have created new digital economies. Learn how to build an NFT marketplace from smart contract development to frontend integration with Web3 libraries.",
			AuthorID: users[17].ID,
			Status:   "draft",
		},
		{
			Title:    "Layer 2 Scaling Solutions",
			Content:  "Ethereum's scaling challenges have led to innovative Layer 2 solutions. Compare rollups, sidechains, and state channels, and learn when to use each scaling approach.",
			AuthorID: users[48].ID,
			Status:   "published",
		},
		
		// DevOps/Cloud Series
		{
			Title:    "Kubernetes Cluster Management",
			Content:  "Kubernetes has become the standard for container orchestration. This comprehensive guide covers cluster setup, pod scheduling, service discovery, and advanced networking configurations.",
			AuthorID: users[46].ID, // Storm DevOps
			Status:   "published",
		},
		{
			Title:    "Infrastructure as Code with Terraform",
			Content:  "Managing cloud infrastructure manually is error-prone and inefficient. Learn how to use Terraform to provision, modify, and version your infrastructure across multiple cloud providers.",
			AuthorID: users[25].ID, // Mark Deploy
			Status:   "published",
		},
		{
			Title:    "CI/CD Pipeline Optimization",
			Content:  "Efficient CI/CD pipelines are crucial for rapid software delivery. Explore advanced techniques for pipeline optimization, parallel execution, and automated testing strategies.",
			AuthorID: users[25].ID,
			Status:   "published",
		},
		{
			Title:    "Monitoring and Observability with Prometheus",
			Content:  "Effective monitoring is essential for maintaining reliable services. Learn how to implement comprehensive monitoring using Prometheus, Grafana, and distributed tracing with Jaeger.",
			AuthorID: users[38].ID, // Amy Cloud
			Status:   "published",
		},
		{
			Title:    "Serverless Architecture Patterns",
			Content:  "Serverless computing enables building applications without managing servers. Explore AWS Lambda, Azure Functions, event-driven architectures, and serverless best practices.",
			AuthorID: users[38].ID,
			Status:   "draft",
		},
		
		// Security Series
		{
			Title:    "Web Application Security Testing",
			Content:  "Security testing is crucial for protecting applications from cyber threats. Learn about penetration testing, vulnerability scanning, and automated security testing in CI/CD pipelines.",
			AuthorID: users[16].ID, // Kevin O'Connor
			Status:   "published",
		},
		{
			Title:    "OAuth 2.0 and OpenID Connect",
			Content:  "Modern authentication requires robust and secure protocols. Understand OAuth 2.0 flows, OpenID Connect, JWT tokens, and implementing secure authentication in distributed systems.",
			AuthorID: users[37].ID, // Chris Secure
			Status:   "published",
		},
		{
			Title:    "Container Security Best Practices",
			Content:  "Containers introduce new security considerations. Learn about image scanning, runtime security, network policies, and securing container orchestration platforms like Kubernetes.",
			AuthorID: users[37].ID,
			Status:   "published",
		},
		{
			Title:    "Zero Trust Security Architecture",
			Content:  "Traditional perimeter-based security is insufficient for modern threats. Explore zero trust principles, identity verification, micro-segmentation, and continuous security monitoring.",
			AuthorID: users[16].ID,
			Status:   "draft",
		},
		
		// Frontend/Mobile Series
		{
			Title:    "React Performance Optimization",
			Content:  "Building performant React applications requires understanding rendering behavior and optimization techniques. Learn about memo, useMemo, useCallback, and code splitting for better performance.",
			AuthorID: users[44].ID, // Ruby Web
			Status:   "published",
		},
		{
			Title:    "Modern CSS Grid and Flexbox",
			Content:  "CSS Grid and Flexbox provide powerful layout capabilities for modern web design. Master responsive design patterns, grid template areas, and flexible box model for complex layouts.",
			AuthorID: users[19].ID, // Rachel Green
			Status:   "published",
		},
		{
			Title:    "Progressive Web Apps Development",
			Content:  "PWAs combine the best of web and mobile apps. Learn about service workers, offline functionality, push notifications, and creating app-like experiences in the browser.",
			AuthorID: users[44].ID,
			Status:   "published",
		},
		{
			Title:    "React Native vs Flutter",
			Content:  "Choosing the right cross-platform framework is crucial for mobile development. Compare React Native and Flutter in terms of performance, development experience, and ecosystem.",
			AuthorID: users[18].ID, // Daniel Kim
			Status:   "published",
		},
		{
			Title:    "State Management in React Applications",
			Content:  "Managing state in complex React applications requires careful architecture decisions. Compare Redux, Zustand, Jotai, and Context API for different use cases.",
			AuthorID: users[42].ID, // Luna Mobile
			Status:   "draft",
		},
		
		// Backend/API Series
		{
			Title:    "GraphQL vs REST API Design",
			Content:  "Choosing between GraphQL and REST depends on your specific requirements. Compare query flexibility, performance characteristics, caching strategies, and development complexity.",
			AuthorID: users[45].ID, // Phoenix Backend
			Status:   "published",
		},
		{
			Title:    "Building Microservices with gRPC",
			Content:  "gRPC provides efficient communication between microservices. Learn about Protocol Buffers, service definition, streaming, and implementing resilient service-to-service communication.",
			AuthorID: users[49].ID, // Sky Platform
			Status:   "published",
		},
		{
			Title:    "Database Sharding Strategies",
			Content:  "Scaling databases horizontally requires careful sharding strategies. Explore horizontal partitioning, shard key selection, cross-shard queries, and managing distributed transactions.",
			AuthorID: users[9].ID, // David Miller (from original)
			Status:   "published",
		},
		{
			Title:    "Event-Driven Architecture with Apache Kafka",
			Content:  "Event-driven systems enable loosely coupled, scalable architectures. Learn about Kafka producers, consumers, stream processing, and building resilient event-driven systems.",
			AuthorID: users[45].ID,
			Status:   "published",
		},
		{
			Title:    "API Rate Limiting and Throttling",
			Content:  "Protecting APIs from abuse requires effective rate limiting. Implement token bucket, sliding window, and fixed window algorithms for different rate limiting scenarios.",
			AuthorID: users[49].ID,
			Status:   "draft",
		},
		
		// Data Engineering Series
		{
			Title:    "Data Pipeline Architecture with Apache Airflow",
			Content:  "Building reliable data pipelines requires robust orchestration. Learn about DAGs, task dependencies, error handling, and monitoring in Apache Airflow workflows.",
			AuthorID: users[15].ID, // Maria Gonzalez
			Status:   "published",
		},
		{
			Title:    "Real-time Analytics with Apache Spark",
			Content:  "Processing large-scale data in real-time enables immediate insights. Master Spark Streaming, structured streaming, and building scalable analytics pipelines.",
			AuthorID: users[39].ID, // Ben Data
			Status:   "published",
		},
		{
			Title:    "Data Lake Architecture on AWS",
			Content:  "Data lakes enable storing vast amounts of structured and unstructured data. Learn about S3, Glue, Athena, and building cost-effective data lake solutions on AWS.",
			AuthorID: users[15].ID,
			Status:   "published",
		},
		{
			Title:    "Time Series Database Optimization",
			Content:  "Time series data requires specialized storage and query optimization. Compare InfluxDB, TimescaleDB, and ClickHouse for different time series use cases.",
			AuthorID: users[39].ID,
			Status:   "draft",
		},
		
		// Game Development Series
		{
			Title:    "Unity 3D Game Development Fundamentals",
			Content:  "Creating engaging 3D games requires understanding Unity's component system. Learn about GameObjects, physics, lighting, and creating immersive gaming experiences.",
			AuthorID: users[43].ID, // Felix Game
			Status:   "published",
		},
		{
			Title:    "Multiplayer Game Networking",
			Content:  "Networked multiplayer games present unique challenges. Explore client-server architecture, lag compensation, prediction, and synchronization in real-time games.",
			AuthorID: users[43].ID,
			Status:   "published",
		},
		{
			Title:    "Game Performance Optimization",
			Content:  "Optimizing game performance is crucial for smooth gameplay. Learn about profiling, memory management, draw call optimization, and platform-specific optimizations.",
			AuthorID: users[43].ID,
			Status:   "draft",
		},
		
		// UX/Design Series
		{
			Title:    "Design Systems and Component Libraries",
			Content:  "Consistent design across applications requires systematic approaches. Build scalable design systems, component libraries, and design tokens for cohesive user experiences.",
			AuthorID: users[40].ID, // Zoe Creative
			Status:   "published",
		},
		{
			Title:    "Accessibility in Web Development",
			Content:  "Building accessible applications ensures inclusivity for all users. Learn about WCAG guidelines, screen readers, keyboard navigation, and testing for accessibility compliance.",
			AuthorID: users[40].ID,
			Status:   "published",
		},
		{
			Title:    "User Research and Testing Methods",
			Content:  "Understanding user needs drives successful product development. Explore user interviews, usability testing, A/B testing, and translating research into design decisions.",
			AuthorID: users[41].ID, // Max Product
			Status:   "published",
		},
		
		// Advanced Topics
		{
			Title:    "Quantum Computing Applications",
			Content:  "Quantum computing promises to revolutionize certain computational problems. Explore quantum algorithms, quantum supremacy, and practical applications in cryptography and optimization.",
			AuthorID: users[30].ID, // Dr. Rebecca Foster
			Status:   "published",
		},
		{
			Title:    "WebAssembly for High-Performance Web Apps",
			Content:  "WebAssembly enables near-native performance in web browsers. Learn about compiling languages to WASM, integration with JavaScript, and use cases for computational tasks.",
			AuthorID: users[14].ID, // James Thompson
			Status:   "published",
		},
		{
			Title:    "Edge Computing and IoT Integration",
			Content:  "Edge computing brings computation closer to data sources. Explore edge devices, IoT protocols, real-time processing, and building distributed edge applications.",
			AuthorID: users[38].ID, // Amy Cloud
			Status:   "draft",
		},
		
		// Career and Industry
		{
			Title:    "Technical Leadership and Team Management",
			Content:  "Growing from individual contributor to technical leader requires new skills. Learn about team building, technical decision-making, and balancing hands-on work with leadership.",
			AuthorID: users[31].ID, // Prof. Alan Wright
			Status:   "published",
		},
		{
			Title:    "Open Source Contribution Best Practices",
			Content:  "Contributing to open source projects builds skills and community. Learn about finding projects, making meaningful contributions, and building your open source reputation.",
			AuthorID: users[14].ID, // James Thompson
			Status:   "published",
		},
		{
			Title:    "Remote Work for Software Engineers",
			Content:  "Remote software development requires adapting communication and collaboration practices. Explore tools, techniques, and strategies for effective remote engineering teams.",
			AuthorID: users[32].ID, // Samantha Code
			Status:   "published",
		},
		
		// Emerging Technologies
		{
			Title:    "Augmented Reality Development with ARKit",
			Content:  "AR is transforming how we interact with digital content. Learn about ARKit fundamentals, 3D object placement, gesture recognition, and building immersive AR experiences.",
			AuthorID: users[42].ID, // Luna Mobile
			Status:   "draft",
		},
		{
			Title:    "Voice User Interface Development",
			Content:  "Voice interfaces are becoming ubiquitous in applications. Explore speech recognition, natural language understanding, and building conversational user interfaces.",
			AuthorID: users[47].ID, // Sage AI
			Status:   "published",
		},
		{
			Title:    "5G Network Programming",
			Content:  "5G networks enable new categories of applications with ultra-low latency and high bandwidth. Learn about network slicing, edge computing, and 5G application development.",
			AuthorID: users[18].ID, // Daniel Kim
			Status:   "draft",
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