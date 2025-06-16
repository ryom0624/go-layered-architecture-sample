# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 📋 Development Flow Guide
このプロジェクトでは標準化された開発フローを使用します。新しい機能開発やタスクを依頼する際は、`DEVELOPMENT_FLOW_GUIDE.md`を参照し、以下の指示形式を使用してください：

```
[タスク名]をDevelopment Workflowに従ってください。
```

これにより、Clean Architectureの原則に従った高品質で一貫性のある開発が保証されます。詳細な手順とベストプラクティスについては `DEVELOPMENT_FLOW_GUIDE.md` をご覧ください。

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

## Development Workflow

### Feature Implementation Process
新機能開発時の標準的な開発フローです。Clean Architectureの依存関係ルールに従って内側から外側へ実装します。

#### 1. 計画・設計フェーズ
- `tasks/` ディレクトリの実装計画書を確認
- データ設計（エンティティ、リレーション）の検討
- API設計（エンドポイント、リクエスト/レスポンス）の定義
- アーキテクチャ設計（各レイヤーの責務）の明確化

#### 2. ブランチ作成
```bash
git checkout -b feature/task-name
```

#### 3. 実装順序（Clean Architecture準拠）
**依存関係ルール**: 内側の層は外側の層に依存しない

1. **Domain Layer（内側）**
   - `internal/domain/entity/` - エンティティ定義
   - `internal/domain/repository/` - リポジトリインターフェース定義

2. **Infrastructure Layer（外側）**
   - `internal/infrastructure/repository/` - リポジトリ実装
   - データベース操作、トランザクション対応

3. **Application Layer（ユースケース）**
   - `internal/usecase/` - ビジネスロジック実装
   - バリデーション、エラーハンドリング

4. **Presentation Layer（外側）**
   - `internal/presentation/handler/` - HTTPハンドラー実装
   - `internal/presentation/router/` - ルーティング設定

5. **Dependency Injection（配線）**
   - `cmd/main.go` - 依存関係の配線

#### 4. テスト実装
各レイヤーでテストを実装（カスタムモック使用）
テストの実行が成功するまでテストコードを修正する
```bash
go test ./...  # 全テスト実行
```

#### 5. ドキュメント更新
- `CLAUDE.md` - API エンドポイント追加
- `TEST_PLAN.md` - テスト計画追加
- `README.md` - 機能説明、使用例追加

#### 6. マージ
```bash
git checkout task-documentation
git merge feature/task-name
```

### 実装済み機能の開発履歴

#### ✅ Task 1-3: 基本機能実装
- **User Management**: CRUD操作、バリデーション
- **Article Management**: CRUD操作、公開/非公開、トランザクション対応
- **Comment System**: 階層構造、承認機能、モデレーション

#### ✅ Task 5: 検索・フィルタリング機能
- **Full-text Search**: タイトル・本文の全文検索
- **Advanced Filtering**: 著者・ステータス・日付範囲フィルタ
- **Sorting & Pagination**: 複数ソートオプション、ページネーション
- **Popular/Recent Articles**: 人気記事・最新記事エンドポイント

**実装パターン**:
- SearchParams エンティティでパラメータ管理
- Repository層で動的クエリ構築
- Usecase層でバリデーションとビジネスロジック
- Handler層でHTTPパラメータ解析

#### ✅ Task 6: お気に入り・ブックマーク機能
- **Favorite System**: 記事お気に入り機能
- **Reading Lists**: カスタム読書リスト作成・管理
- **Public/Private Lists**: 公開・非公開リスト機能

**実装パターン**:
- Favorite エンティティでお気に入り管理
- ReadingList, ReadingListItem エンティティで読書リスト管理
- トランザクション内でお気に入り数更新
- アクセス制御（公開・非公開読書リスト）
- 複合ユニーク制約で重複防止

#### ✅ Task 7: 閲覧履歴・統計機能
- **View Tracking**: 記事閲覧の自動トラッキング
- **Reading Analytics**: ユーザー読書分析  
- **Platform Statistics**: プラットフォーム統計情報

**実装パターン**:
- ArticleView エンティティで個別閲覧トラッキング
- UserReadingHistory エンティティで読書履歴管理
- ArticleStatistics エンティティで記事統計
- DailyStatistics エンティティで日次統計
- ViewTrackingMiddleware で自動閲覧記録
- 統計情報のリアルタイム集計とキャッシュ

#### ✅ Task 8: 認証システム
- **JWT Authentication**: アクセストークンによる認証
- **User Registration/Login**: 安全なユーザー登録・ログイン
- **Token Management**: リフレッシュトークンによる自動更新
- **Multi-device Support**: デバイス別ログアウト機能

**実装パターン**:
- RefreshToken エンティティでトークン管理
- bcrypt による安全なパスワードハッシュ化
- JWT によるステートレス認証
- AuthMiddleware による認証制御
- オプション認証ミドルウェアによる柔軟な認証

### 開発ガイドライン

#### コーディング規約
- **Clean Architecture**: 依存関係ルールの厳守
- **Error Handling**: 適切なエラーハンドリングとバリデーション
- **Transaction Management**: データ整合性を保つトランザクション使用
- **Testing**: 各レイヤーでのユニットテスト・統合テスト実装

#### 命名規約
- **Interface**: `XxxRepository`, `XxxUsecase`
- **Implementation**: `xxxRepositoryImpl`, `xxxUsecaseImpl`
- **Handler**: `XxxHandler`
- **Entity**: パスカルケース（`User`, `Article`）

#### ファイル構成
- **Test Files**: `*_test.go` でソースファイルと同階層
- **Mock Implementation**: カスタムモック、外部ライブラリ不使用
- **Documentation**: 各機能の実装計画とAPIドキュメント維持

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
- Authentication configuration:
  - `JWT_SECRET` - Secret key for JWT token signing
  - `ACCESS_TOKEN_DURATION_MINUTES` - Access token expiry (default: 15 minutes)
  - `REFRESH_TOKEN_DURATION_DAYS` - Refresh token expiry (default: 7 days)

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

#### Authentication
- `POST /api/v1/auth/register` - User registration (requires name, email, password)
- `POST /api/v1/auth/login` - User login (requires email, password)
- `POST /api/v1/auth/refresh` - Refresh access token (requires refresh_token)
- `POST /api/v1/auth/logout` - Logout from single device (requires refresh_token)
- `POST /api/v1/auth/logout-all` - Logout from all devices (requires authentication)

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

#### Favorite Management
- `POST /api/v1/articles/:id/favorite` - Add article to favorites
- `DELETE /api/v1/articles/:id/favorite` - Remove article from favorites
- `GET /api/v1/users/:id/favorites` - Get user's favorite articles
- `GET /api/v1/articles/:id/favorites` - Get users who favorited article
- `GET /api/v1/articles/:id/favorite-status` - Check if current user favorited article

#### Reading List Management
- `POST /api/v1/reading-lists` - Create reading list
- `GET /api/v1/reading-lists` - Get user's reading lists
- `GET /api/v1/reading-lists/:id` - Get reading list details
- `PUT /api/v1/reading-lists/:id` - Update reading list
- `DELETE /api/v1/reading-lists/:id` - Delete reading list

#### Reading List Items
- `POST /api/v1/reading-lists/:id/articles` - Add article to reading list
- `DELETE /api/v1/reading-lists/:id/articles/:article_id` - Remove article from reading list
- `GET /api/v1/reading-lists/:id/articles` - Get articles in reading list
- `PUT /api/v1/reading-lists/:id/articles/:article_id` - Update reading list item notes

#### Public Reading Lists
- `GET /api/v1/reading-lists/public` - Get public reading lists
- `GET /api/v1/users/:id/reading-lists/public` - Get user's public reading lists

#### View Tracking and Reading Progress
- `PUT /api/v1/progress/users/:user_id/articles/:article_id` - Track reading progress
- `GET /api/v1/progress/users/:user_id/articles/:article_id` - Get reading progress
- `GET /api/v1/users/:id/reading-history` - Get user's reading history
- `GET /api/v1/articles/trending` - Get trending articles

#### Analytics and Statistics
- `GET /api/v1/analytics/overview` - Get platform overview (admin)
- `GET /api/v1/analytics/daily` - Get daily statistics
- `GET /api/v1/analytics/daily/range` - Get daily statistics range
- `GET /api/v1/analytics/articles` - Get all article statistics
- `GET /api/v1/analytics/articles/popular` - Get popular articles
- `GET /api/v1/users/:id/analytics` - Get user analytics
- `GET /api/v1/articles/:id/statistics` - Get article statistics
- `POST /api/v1/articles/:id/statistics/recalculate` - Recalculate article statistics

### Database Seeding Data
The seed command creates sample data:
- **10 Users**: John Doe, Jane Smith, Alice Johnson, etc.
- **15 Articles**: Technical articles on various programming topics
- **25+ Comments**: Sample comments with hierarchical structure and different statuses
- **Mixed Status**: Some articles published, some in draft status; comments with pending/approved/rejected status
- **Relations**: Articles assigned to users, comments linked to articles and users with parent-child relationships
