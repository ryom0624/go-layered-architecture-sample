# Go Layered Architecture Template

Clean Architecture（クリーンアーキテクチャ）に基づいたGo言語のレイヤードアーキテクチャテンプレートです。

## プロジェクト構造

```
.
├── cmd/                        # アプリケーションのエントリーポイント
│   └── main.go
├── internal/                   # プライベートなアプリケーションコード
│   ├── domain/                 # ドメイン層
│   │   ├── entity/            # エンティティ
│   │   └── repository/        # リポジトリインターフェース
│   ├── usecase/               # アプリケーション層（ユースケース）
│   ├── infrastructure/        # インフラストラクチャ層
│   │   ├── database/         # データベース接続
│   │   └── repository/       # リポジトリ実装
│   └── presentation/          # プレゼンテーション層
│       ├── handler/          # HTTPハンドラー
│       └── router/           # ルーター
├── pkg/                       # 外部から利用可能な公開コード
│   └── config/               # 設定管理
├── .env.example              # 環境変数のサンプル
├── docker-compose.yml        # Docker Compose設定
├── Dockerfile               # Docker設定
└── go.mod                   # Go modules
```

## アーキテクチャの特徴

### 1. ドメイン層 (Domain Layer)
- **Entity**: ビジネスエンティティとルール
- **Repository Interface**: データアクセスの抽象化

### 2. アプリケーション層 (Application Layer)
- **Usecase**: ビジネスロジックとアプリケーションルール

### 3. インフラストラクチャ層 (Infrastructure Layer)
- **Database**: データベース接続管理
- **Repository Implementation**: データアクセスの具体実装

### 4. プレゼンテーション層 (Presentation Layer)
- **Handler**: HTTPリクエストの処理
- **Router**: ルーティング設定

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

### ユーザー管理

- `POST /api/v1/users` - ユーザー作成
- `GET /api/v1/users` - 全ユーザー取得
- `GET /api/v1/users/:id` - ユーザー取得
- `PUT /api/v1/users/:id` - ユーザー更新
- `DELETE /api/v1/users/:id` - ユーザー削除

### リクエスト例

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

## 使用技術

- **Web Framework**: Gin
- **ORM**: GORM
- **Database**: PostgreSQL, MySQL対応
- **Configuration**: godotenv
- **Containerization**: Docker
- **Testing**: 標準testingパッケージ, go-sqlmock, httptest

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

- `internal/usecase/user_usecase_test.go` - ビジネスロジックテスト
- `internal/presentation/handler/user_handler_test.go` - HTTPハンドラーテスト
- `internal/infrastructure/repository/user_repository_impl_test.go` - リポジトリテスト
- `pkg/config/config_test.go` - 設定管理テスト

詳細なテスト計画については `TEST_PLAN.md` を参照してください。

## ライセンス

MIT License
