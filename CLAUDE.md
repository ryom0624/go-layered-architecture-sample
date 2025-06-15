# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Commands

### Development
- `go mod tidy` - Download and clean up dependencies
- `go run cmd/main.go` - Start the application locally
- `go build -o bin/app cmd/main.go` - Build the application binary

### Testing
- `go test ./...` - Run all tests
- `go test ./internal/usecase/...` - Run tests for a specific package
- `go test -v ./...` - Run tests with verbose output
- `go test -cover ./...` - Run tests with coverage
- `docker compose exec app go test ./... -v` - Run tests in Docker container

### Environment Setup
- `cp .env.example .env` - Copy environment configuration template
- `docker compose up -d postgres` - Start PostgreSQL database only
- `docker compose up` - Start full application stack with database

### Database Seeding
- `go run cmd/seed/main.go` - Seed database with sample data locally
- `./scripts/seed.sh` - Run seed script locally
- `./scripts/docker-seed.sh` - Run seed script in Docker container
- `docker compose run --rm app go run cmd/seed/main.go` - Direct Docker seeding

## Architecture Overview

This is a Clean Architecture implementation with strict dependency rules:

### Dependency Flow (Inner → Outer)
1. **Domain Layer** (`internal/domain/`) - Core business entities and repository interfaces
2. **Application Layer** (`internal/usecase/`) - Business logic and use cases
3. **Infrastructure Layer** (`internal/infrastructure/`) - Database connections and repository implementations
4. **Presentation Layer** (`internal/presentation/`) - HTTP handlers and routing

### Key Architectural Rules
- Dependencies flow inward only (outer layers depend on inner layers, never vice versa)
- Domain layer has no external dependencies
- Repository interfaces are defined in domain, implemented in infrastructure
- Dependency injection occurs in `cmd/main.go` where all layers are wired together

### Configuration Management
- Environment variables loaded via `pkg/config/config.go`
- Database connection supports both PostgreSQL and MySQL via `DB_DRIVER` env var
- Auto-migration handled by GORM in `internal/infrastructure/database/connection.go`

### Adding New Features
When adding new entities, follow this sequence:
1. Entity in `internal/domain/entity/`
2. Repository interface in `internal/domain/repository/`  
3. Repository implementation in `internal/infrastructure/repository/`
4. Use case in `internal/usecase/`
5. Handler in `internal/presentation/handler/`
6. Route registration in `internal/presentation/router/`
7. Wire dependencies in `cmd/main.go`

### Database Operations
- GORM handles migrations automatically on startup
- Repository pattern abstracts database operations
- Context is passed through all database operations for timeout/cancellation support

### Testing Architecture
The codebase includes comprehensive tests for all layers:
- **Unit Tests**: Use case layer with custom mocks (no external dependencies)
- **Integration Tests**: Repository layer with go-sqlmock for database operations
- **HTTP Tests**: Handler layer with httptest and Gin test mode
- **Config Tests**: Environment variable loading and validation
- Test files follow the pattern `*_test.go` alongside source files
- Custom mock implementations instead of external mocking libraries
- Standard Go testing package without assertion libraries

### API Endpoints
RESTful API with the following endpoints:

#### User Management
- `POST /api/v1/users` - Create user (requires name, email)
- `GET /api/v1/users` - Get all users
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user (partial updates supported)
- `DELETE /api/v1/users/:id` - Delete user

#### Article Management
- `POST /api/v1/articles` - Create article (requires title, content, author_id)
- `GET /api/v1/articles` - Get all articles (supports filtering via query parameters)
- `GET /api/v1/articles/published` - Get published articles only
- `GET /api/v1/articles/popular` - Get popular articles (ordered by publication date)
- `GET /api/v1/articles/recent` - Get recent published articles
- `GET /api/v1/articles/:id` - Get article by ID
- `PUT /api/v1/articles/:id` - Update article (partial updates supported)
- `DELETE /api/v1/articles/:id` - Delete article
- `PUT /api/v1/articles/:id/publish` - Publish article (uses transaction)
- `PUT /api/v1/articles/:id/unpublish` - Unpublish article (uses transaction)
- `POST /api/v1/articles/:id/comments` - Create comment on article
- `GET /api/v1/articles/:id/comments` - Get all comments for article

#### Search and Filtering
- `GET /api/v1/search` - Search articles with comprehensive filtering
  - Query parameters:
    - `query` - Full-text search query (title and content)
    - `author_id` - Filter by author ID
    - `status` - Filter by status (published, draft)
    - `date_from` - Filter articles from date (YYYY-MM-DD format)
    - `date_to` - Filter articles to date (YYYY-MM-DD format)
    - `sort_by` - Sort field (created_at, updated_at, title, author, status, relevance)
    - `sort_order` - Sort order (asc, desc)
    - `page` - Page number for pagination (default: 1)
    - `limit` - Items per page (default: 20, max: 100)

#### Comment Management
- `GET /api/v1/users/:id/comments` - Get all comments by user
- `GET /api/v1/comments/:id` - Get comment by ID
- `PUT /api/v1/comments/:id` - Update comment (requires same author)
- `DELETE /api/v1/comments/:id` - Delete comment (requires same author)
- `POST /api/v1/comments/:id/replies` - Create reply to comment (hierarchical, max 3 levels)
- `GET /api/v1/comments/pending` - Get all pending comments (moderation)
- `PUT /api/v1/comments/:id/approve` - Approve comment (moderation)
- `PUT /api/v1/comments/:id/reject` - Reject comment (moderation)

### Database Seeding Data
The seed command creates sample data:
- **10 Users**: John Doe, Jane Smith, Alice Johnson, etc.
- **15 Articles**: Technical articles on various programming topics
- **25+ Comments**: Sample comments with hierarchical structure and different statuses
- **Mixed Status**: Some articles published, some in draft status; comments with pending/approved/rejected status
- **Relations**: Articles assigned to users, comments linked to articles and users with parent-child relationships