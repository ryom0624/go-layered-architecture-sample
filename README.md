# Go Layered Architecture Template

Clean Architecture（クリーンアーキテクチャ）に基づいたGo言語のレイヤードアーキテクチャテンプレートです。
ユーザー管理機能と記事投稿機能を実装し、データベーストランザクション対応の包括的なWebアプリケーションを提供します。

## プロジェクト構造

```
.
├── cmd/                        # アプリケーションのエントリーポイント
│   └── main.go
├── internal/                   # プライベートなアプリケーションコード
│   ├── domain/                 # ドメイン層
│   │   ├── entity/            # エンティティ（User, Article）
│   │   └── repository/        # リポジトリインターフェース（トランザクション対応）
│   ├── usecase/               # アプリケーション層（ユースケース）
│   ├── infrastructure/        # インフラストラクチャ層
│   │   ├── database/         # データベース接続・トランザクション管理
│   │   └── repository/       # リポジトリ実装
│   └── presentation/          # プレゼンテーション層
│       ├── handler/          # HTTPハンドラー
│       └── router/           # ルーター
├── pkg/                       # 外部から利用可能な公開コード
│   └── config/               # 設定管理
├── tasks/                     # タスクドキュメント
│   ├── 01_user.md            # User機能実装ドキュメント
│   ├── 02_article.md         # Article機能実装ドキュメント
│   └── 03_comment.md         # Comment機能実装ドキュメント
├── .env.example              # 環境変数のサンプル
├── docker-compose.yml        # Docker Compose設定
├── Dockerfile               # Docker設定
└── go.mod                   # Go modules
```

## アーキテクチャの特徴

### 1. ドメイン層 (Domain Layer)
- **Entity**: ビジネスエンティティとルール（User, Article, Comment）
- **Repository Interface**: データアクセスの抽象化（トランザクション対応）

### 2. アプリケーション層 (Application Layer)
- **Usecase**: ビジネスロジックとアプリケーションルール
- **Transaction Management**: データ整合性保証

### 3. インフラストラクチャ層 (Infrastructure Layer)
- **Database**: データベース接続管理・トランザクション制御
- **Repository Implementation**: データアクセスの具体実装

### 4. プレゼンテーション層 (Presentation Layer)
- **Handler**: HTTPリクエストの処理
- **Router**: ルーティング設定

## 主要機能

### ✅ ユーザー管理システム
- ユーザーのCRUD操作
- 重複メール検証
- 入力値バリデーション

### ✅ 記事投稿システム
- 記事のCRUD操作
- 記事公開/非公開機能
- 著者情報との関連付け
- データベーストランザクション対応

### ✅ コメントシステム
- 階層コメント機能（最大3レベルの返信）
- コメント承認・モデレーション機能
- ユーザー・記事との関連付け
- データベーストランザクション対応

### ✅ 検索・フィルタリング機能
- 記事の全文検索（タイトル・本文）
- 高度なフィルタリング（著者・ステータス・日付範囲）
- ソート機能（作成日・更新日・タイトル・著者・関連度）
- ページネーション対応
- 人気記事・最新記事取得

### ✅ お気に入り・ブックマーク機能
- 記事お気に入り機能（追加・削除・一覧表示）
- カスタム読書リスト作成・管理
- 読書リストへの記事追加・削除・メモ機能
- 公開・非公開読書リスト設定
- お気に入り数の自動カウント・トランザクション対応

### ✅ 閲覧履歴・統計機能
- 記事閲覧の自動トラッキング
- ユーザー読書履歴管理
- プラットフォーム統計情報（日次統計・記事統計）
- 人気記事・トレンド記事分析
- リアルタイム統計集計

### ✅ 認証システム
- JWT（JSON Web Token）による認証
- ユーザー登録・ログイン機能
- リフレッシュトークンによる自動更新
- マルチデバイス対応（デバイス別ログアウト）
- 安全なパスワードハッシュ化（bcrypt）
- 認証ミドルウェアによるアクセス制御

### ✅ トランザクション機能
- 複数テーブル操作の整合性保証
- エラー時の自動ロールバック
- データ整合性の確保

## セットアップ

### 前提条件
- Go 1.21以上
- Docker & Docker Compose（オプション）

### 1. 依存関係のインストール
```bash
go mod download
```

### 2. 環境変数の設定
```bash
cp .env.example .env
# .envファイルを編集して適切な値を設定
```

### 3. データベースの起動（Docker使用の場合）
```bash
docker compose up -d postgres
```

### 4. アプリケーションの起動
```bash
go run cmd/main.go
```

または、Docker Composeを使用：
```bash
docker compose up
```

## API エンドポイント

### 認証システム

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | ユーザー登録 |
| POST | `/api/v1/auth/login` | ログイン |
| POST | `/api/v1/auth/refresh` | アクセストークン更新 |
| POST | `/api/v1/auth/logout` | 単一デバイスログアウト |
| POST | `/api/v1/auth/logout-all` | 全デバイスログアウト（認証必須） |

### ユーザー管理

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/users` | ユーザー作成 |
| GET | `/api/v1/users` | 全ユーザー取得 |
| GET | `/api/v1/users/:id` | ユーザー取得 |
| PUT | `/api/v1/users/:id` | ユーザー更新 |
| DELETE | `/api/v1/users/:id` | ユーザー削除 |

### 記事管理

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/articles` | 記事作成 |
| GET | `/api/v1/articles` | 全記事取得（フィルタリング対応） |
| GET | `/api/v1/articles/published` | 公開記事のみ取得 |
| GET | `/api/v1/articles/popular` | 人気記事取得 |
| GET | `/api/v1/articles/recent` | 最新記事取得 |
| GET | `/api/v1/articles/:id` | 記事詳細取得 |
| PUT | `/api/v1/articles/:id` | 記事更新 |
| DELETE | `/api/v1/articles/:id` | 記事削除 |
| PUT | `/api/v1/articles/:id/publish` | 記事公開 |
| PUT | `/api/v1/articles/:id/unpublish` | 記事非公開 |
| POST | `/api/v1/articles/:id/comments` | 記事へのコメント作成 |
| GET | `/api/v1/articles/:id/comments` | 記事のコメント取得 |

### 検索・フィルタリング

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/search` | 記事の包括的検索・フィルタリング |

### コメント管理

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/users/:id/comments` | ユーザーのコメント取得 |
| GET | `/api/v1/comments/:id` | コメント詳細取得 |
| PUT | `/api/v1/comments/:id` | コメント更新 |
| DELETE | `/api/v1/comments/:id` | コメント削除 |
| POST | `/api/v1/comments/:id/replies` | 返信コメント作成 |
| GET | `/api/v1/comments/pending` | 承認待ちコメント取得 |
| PUT | `/api/v1/comments/:id/approve` | コメント承認 |
| PUT | `/api/v1/comments/:id/reject` | コメント拒否 |

### お気に入り管理

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/articles/:id/favorite` | 記事をお気に入りに追加 |
| DELETE | `/api/v1/articles/:id/favorite` | 記事をお気に入りから削除 |
| GET | `/api/v1/users/:id/favorites` | ユーザーのお気に入り記事取得 |
| GET | `/api/v1/articles/:id/favorites` | 記事をお気に入りしたユーザー取得 |
| GET | `/api/v1/articles/:id/favorite-status` | お気に入りステータス確認 |

### 読書リスト管理

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/reading-lists` | 読書リスト作成 |
| GET | `/api/v1/reading-lists` | ユーザーの読書リスト取得 |
| GET | `/api/v1/reading-lists/:id` | 読書リスト詳細取得 |
| PUT | `/api/v1/reading-lists/:id` | 読書リスト更新 |
| DELETE | `/api/v1/reading-lists/:id` | 読書リスト削除 |
| POST | `/api/v1/reading-lists/:id/articles` | 読書リストに記事追加 |
| GET | `/api/v1/reading-lists/:id/articles` | 読書リスト内の記事取得 |
| DELETE | `/api/v1/reading-lists/:id/articles/:article_id` | 読書リストから記事削除 |
| PUT | `/api/v1/reading-lists/:id/articles/:article_id` | 読書リストアイテム更新 |
| GET | `/api/v1/reading-lists/public` | 公開読書リスト取得 |
| GET | `/api/v1/users/:id/reading-lists/public` | ユーザーの公開読書リスト取得 |

### リクエスト例

#### ユーザー登録（認証）
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com", "password": "securepassword123"}'
```

#### ログイン
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "john@example.com", "password": "securepassword123"}'
```

#### アクセストークン更新
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "your_refresh_token_here"}'
```

#### ログアウト
```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "your_refresh_token_here"}'
```

#### 全デバイスログアウト
```bash
curl -X POST http://localhost:8080/api/v1/auth/logout-all \
  -H "Authorization: Bearer your_access_token_here"
```

#### 認証が必要なエンドポイントへのアクセス
```bash
curl -H "Authorization: Bearer your_access_token_here" \
  http://localhost:8080/api/v1/protected-endpoint
```

#### ユーザー作成
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'
```

#### ユーザー一覧取得
```bash
curl http://localhost:8080/api/v1/users
```

#### 記事作成
```bash
curl -X POST http://localhost:8080/api/v1/articles \
  -H "Content-Type: application/json" \
  -d '{"title": "My First Article", "content": "Article content here...", "author_id": 1}'
```

#### 記事公開
```bash
curl -X PUT http://localhost:8080/api/v1/articles/1/publish
```

#### 公開記事一覧取得
```bash
curl http://localhost:8080/api/v1/articles/published
```

#### コメント作成
```bash
curl -X POST http://localhost:8080/api/v1/articles/1/comments \
  -H "Content-Type: application/json" \
  -d '{"content": "Great article!", "author_id": 1}'
```

#### 返信コメント作成
```bash
curl -X POST http://localhost:8080/api/v1/comments/1/replies \
  -H "Content-Type: application/json" \
  -d '{"content": "I agree!", "author_id": 2}'
```

#### 記事のコメント取得
```bash
curl http://localhost:8080/api/v1/articles/1/comments
```

#### コメント承認
```bash
curl -X PUT http://localhost:8080/api/v1/comments/1/approve
```

#### 記事検索
```bash
# 全文検索
curl "http://localhost:8080/api/v1/search?query=Go言語"

# フィルタリング付き検索
curl "http://localhost:8080/api/v1/search?query=programming&author_id=1&status=published"

# ページネーション付き検索
curl "http://localhost:8080/api/v1/search?query=tutorial&page=1&limit=5"

# ソート付き検索
curl "http://localhost:8080/api/v1/search?sort_by=created_at&sort_order=desc"

# 日付範囲フィルタ
curl "http://localhost:8080/api/v1/search?date_from=2023-01-01&date_to=2023-12-31"
```

#### 人気記事取得
```bash
curl "http://localhost:8080/api/v1/articles/popular?limit=10"
```

#### 最新記事取得
```bash
curl "http://localhost:8080/api/v1/articles/recent?limit=5"
```

#### お気に入り追加
```bash
curl -X POST http://localhost:8080/api/v1/articles/1/favorite \
  -H "userID: 1"
```

#### お気に入り削除
```bash
curl -X DELETE http://localhost:8080/api/v1/articles/1/favorite \
  -H "userID: 1"
```

#### ユーザーのお気に入り記事取得
```bash
curl http://localhost:8080/api/v1/users/1/favorites
```

#### 読書リスト作成
```bash
curl -X POST http://localhost:8080/api/v1/reading-lists \
  -H "Content-Type: application/json" \
  -H "userID: 1" \
  -d '{"name": "プログラミング学習", "description": "Go言語とアーキテクチャ", "is_public": false}'
```

#### 読書リストに記事追加
```bash
curl -X POST http://localhost:8080/api/v1/reading-lists/1/articles \
  -H "Content-Type: application/json" \
  -H "userID: 1" \
  -d '{"article_id": 1, "notes": "後で詳しく読む"}'
```

#### 公開読書リスト取得
```bash
curl http://localhost:8080/api/v1/reading-lists/public
```

## 環境変数

| 変数名 | 説明 | デフォルト値 |
|--------|------|-------------|
| SERVER_PORT | サーバーポート | 8080 |
| DB_DRIVER | データベースドライバー (postgres/mysql) | postgres |
| DB_HOST | データベースホスト | localhost |
| DB_PORT | データベースポート | 5432 |
| DB_USER | データベースユーザー | user |
| DB_PASSWORD | データベースパスワード | password |
| DB_NAME | データベース名 | database |
| JWT_SECRET | JWT署名用秘密鍵 | your-secret-key |
| ACCESS_TOKEN_DURATION_MINUTES | アクセストークン有効期限（分） | 15 |
| REFRESH_TOKEN_DURATION_DAYS | リフレッシュトークン有効期限（日） | 7 |

## 使用技術

- **Web Framework**: Gin
- **ORM**: GORM
- **Database**: PostgreSQL, MySQL対応
- **Authentication**: JWT (golang-jwt/jwt), bcrypt
- **Configuration**: godotenv
- **Containerization**: Docker
- **Testing**: 標準testingパッケージ, go-sqlmock, httptest

## コード品質

### マジックナンバー排除
本プロジェクトでは、保守性とビジネスルールの透明性向上のため、すべてのマジックナンバーを排除し、中央集権的な定数管理を実装しています。

#### 定数パッケージ構造
```
pkg/constants/
├── config.go      # アプリケーション設定定数（ポート、トークン期間）
├── validation.go  # バリデーション制限値（文字数制限、パスワード要件）
├── pagination.go  # ページネーション・データ取得制限
├── business.go    # ビジネスロジック定数（完読判定閾値、期間設定）
└── parsing.go     # データ解析定数（進数、ビット数）
```

#### 主要定数
- `DefaultServerPort = "8080"` - アプリケーションサーバーデフォルトポート
- `DefaultAccessTokenDurationMinutes = 15` - JWTアクセストークン有効期限
- `MaxCommentLength = 1000` - コメント最大文字数制限  
- `MaxCommentDepth = 2` - コメント階層最大深度（3レベル対応）
- `ReadCompletionThreshold = 90.0` - 記事完読判定閾値（90%）
- `DefaultPageSize = 20` - デフォルトページサイズ
- `MaxPageSize = 100` - 最大ページサイズ

#### メリット
- **保守性向上**: ビジネス要件変更時に単一箇所での変更が可能
- **可読性向上**: 数値の意味が明確で理解しやすいコード
- **一貫性確保**: プロジェクト全体で統一された制限値使用
- **ドキュメント化**: 定数名によるビジネスルールの自己文書化

## 拡張方法

### 新しいエンティティの追加

1. `internal/domain/entity/` に新しいエンティティを作成
2. `internal/domain/repository/` にリポジトリインターフェースを作成
3. `internal/infrastructure/repository/` にリポジトリ実装を作成
4. `internal/usecase/` にユースケースを作成
5. `internal/presentation/handler/` にハンドラーを作成
6. `internal/presentation/router/` にルートを追加

### ミドルウェアの追加

`internal/presentation/router/router.go` でGinのミドルウェアを追加できます。

### バリデーションの追加

ハンドラーでGinのバインディング機能を使用してリクエストバリデーションを行えます。

## テスト

このプロジェクトは包括的なテストスイートを含んでいます：

### テスト実行方法

```bash
# 全テスト実行
go test ./...

# 詳細出力付きテスト実行
go test -v ./...

# カバレッジ付きテスト実行
go test -cover ./...

# Dockerコンテナ内でのテスト実行
docker compose exec app go test ./... -v
```

### テストアーキテクチャ

- **ユースケース層テスト**: ビジネスロジックの単体テスト（カスタムモック使用）
- **ハンドラー層テスト**: HTTP エンドポイントのテスト（httptest使用）
- **リポジトリ層テスト**: データベース操作のテスト（go-sqlmock使用）
- **設定パッケージテスト**: 環境変数読み込みのテスト

### テストファイル

#### ユーザー機能
- `internal/usecase/user_usecase_test.go` - ユーザービジネスロジックテスト
- `internal/presentation/handler/user_handler_test.go` - ユーザーHTTPハンドラーテスト
- `internal/infrastructure/repository/user_repository_impl_test.go` - ユーザーリポジトリテスト

#### 記事機能
- `internal/usecase/article_usecase_test.go` - 記事ビジネスロジックテスト（トランザクション含む）
- `internal/presentation/handler/article_handler_test.go` - 記事HTTPハンドラーテスト
- `internal/infrastructure/repository/article_repository_impl_test.go` - 記事リポジトリテスト

#### コメント機能
- `internal/usecase/comment_usecase_test.go` - コメントビジネスロジックテスト（階層構造・承認機能含む）
- `internal/presentation/handler/comment_handler_test.go` - コメントHTTPハンドラーテスト
- `internal/infrastructure/repository/comment_repository_impl_test.go` - コメントリポジトリテスト

#### 検索・フィルタリング機能
- `internal/usecase/search_usecase_test.go` - 検索ビジネスロジックテスト（検索・ページネーション・バリデーション含む）
- `internal/presentation/handler/search_handler_test.go` - 検索HTTPハンドラーテスト

#### 認証システム
- `internal/usecase/auth_usecase_test.go` - 認証ビジネスロジックテスト（JWT・パスワードハッシュ・トークン管理含む）
- `internal/presentation/handler/auth_handler_test.go` - 認証HTTPハンドラーテスト
- `internal/infrastructure/repository/auth_repository_impl_test.go` - 認証リポジトリテスト

#### 共通機能
- `pkg/config/config_test.go` - 設定管理テスト
- `pkg/jwt/jwt_test.go` - JWTユーティリティテスト

詳細なテスト計画については `TEST_PLAN.md` を参照してください。

## ライセンス

MIT License
