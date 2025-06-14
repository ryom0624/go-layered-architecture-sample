# 記事投稿機能実装計画（トランザクション対応版）

## 概要
既存のClean Architectureに従い、データベーストランザクション機能を含む記事投稿機能を実装します。

## 実装対象機能
- 記事の作成、読み取り、更新、削除（CRUD）
- 記事の公開/非公開切り替え
- データベーストランザクション管理
- 複数テーブル操作時の整合性保証

## データ設計

### Article エンティティ
```go
type Article struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Title     string    `json:"title" gorm:"not null"`
    Content   string    `json:"content" gorm:"type:text"`
    AuthorID  uint      `json:"author_id" gorm:"not null"`
    Author    User      `json:"author" gorm:"foreignKey:AuthorID"`
    Status    string    `json:"status" gorm:"default:'draft'"` // draft, published
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

## API設計

### エンドポイント
- `POST /api/v1/articles` - 記事作成
- `GET /api/v1/articles` - 記事一覧取得（公開済みのみ）
- `GET /api/v1/articles/:id` - 記事詳細取得
- `PUT /api/v1/articles/:id` - 記事更新
- `DELETE /api/v1/articles/:id` - 記事削除
- `PUT /api/v1/articles/:id/publish` - 記事公開
- `PUT /api/v1/articles/:id/unpublish` - 記事非公開

### リクエスト/レスポンス例
```json
// POST /api/v1/articles
{
  "title": "記事タイトル",
  "content": "記事本文",
  "author_id": 1
}

// Response
{
  "id": 1,
  "title": "記事タイトル",
  "content": "記事本文",
  "author_id": 1,
  "status": "draft",
  "created_at": "2025-06-14T10:00:00Z",
  "updated_at": "2025-06-14T10:00:00Z"
}
```

## アーキテクチャ実装計画

### 1. トランザクション基盤整備

#### 1.1 トランザクション管理インターフェース
```go
// internal/domain/repository/transaction.go
type Transaction interface {
    Begin() (Transaction, error)
    Commit() error
    Rollback() error
    GetDB() interface{}
}

type TransactionManager interface {
    WithTransaction(ctx context.Context, fn func(tx Transaction) error) error
}
```

#### 1.2 GORM トランザクション実装
```go
// internal/infrastructure/database/transaction.go
type gormTransaction struct {
    tx *gorm.DB
}

type gormTransactionManager struct {
    db *gorm.DB
}
```

### 2. Domain Layer（内側の層）

#### 2.1 エンティティ作成
- `internal/domain/entity/article.go` - Article構造体定義

#### 2.2 リポジトリインターフェース
- `internal/domain/repository/article_repository.go` - ArticleRepository定義

```go
type ArticleRepository interface {
    Create(ctx context.Context, article *entity.Article) error
    CreateWithTx(ctx context.Context, tx Transaction, article *entity.Article) error
    GetByID(ctx context.Context, id uint) (*entity.Article, error)
    GetAll(ctx context.Context, status string) ([]*entity.Article, error)
    Update(ctx context.Context, article *entity.Article) error
    UpdateWithTx(ctx context.Context, tx Transaction, article *entity.Article) error
    Delete(ctx context.Context, id uint) error
    DeleteWithTx(ctx context.Context, tx Transaction, id uint) error
}
```

### 3. Infrastructure Layer（外側の層）

#### 3.1 リポジトリ実装
- `internal/infrastructure/repository/article_repository_impl.go`
- GORMを使用したデータベース操作実装
- トランザクション対応メソッド実装

#### 3.2 データベース接続更新
- `internal/infrastructure/database/connection.go`にArticleのAutoMigrate追加

### 4. Application Layer（ユースケース層）

#### 4.1 ArticleUsecase実装
- `internal/usecase/article_usecase.go`

```go
type ArticleUsecase interface {
    CreateArticle(ctx context.Context, title, content string, authorID uint) (*entity.Article, error)
    GetArticle(ctx context.Context, id uint) (*entity.Article, error)
    GetPublishedArticles(ctx context.Context) ([]*entity.Article, error)
    UpdateArticle(ctx context.Context, id uint, title, content string) (*entity.Article, error)
    DeleteArticle(ctx context.Context, id uint) error
    PublishArticle(ctx context.Context, id uint) (*entity.Article, error)
    UnpublishArticle(ctx context.Context, id uint) (*entity.Article, error)
}
```

#### 4.2 トランザクション活用例
```go
func (a *articleUsecase) PublishArticle(ctx context.Context, id uint) (*entity.Article, error) {
    return a.transactionManager.WithTransaction(ctx, func(tx Transaction) error {
        // 1. 記事の存在確認
        article, err := a.articleRepo.GetByID(ctx, id)
        if err != nil {
            return err
        }
        
        // 2. ステータス更新
        article.Status = "published"
        if err := a.articleRepo.UpdateWithTx(ctx, tx, article); err != nil {
            return err
        }
        
        // 3. 追加処理（例：通知、ログ記録など）
        // 複数のテーブル操作をトランザクション内で実行
        
        return nil
    })
}
```

### 5. Presentation Layer（プレゼンテーション層）

#### 5.1 ハンドラー実装
- `internal/presentation/handler/article_handler.go`
- リクエスト/レスポンス構造体定義
- HTTPハンドラー実装

#### 5.2 ルーティング設定
- `internal/presentation/router/router.go`にarticlesエンドポイント追加

### 6. 依存関係の配線

#### 6.1 main.go更新
```go
// トランザクションマネージャー初期化
transactionManager := database.NewGormTransactionManager(database.DB)

// リポジトリ初期化
articleRepo := repository.NewArticleRepository(database.DB)

// ユースケース初期化
articleUsecase := usecase.NewArticleUsecase(articleRepo, transactionManager)

// ハンドラー初期化
articleHandler := handler.NewArticleHandler(articleUsecase)

// ルーター設定
router := router.SetupRouter(userHandler, articleHandler)
```

## テスト戦略

### 7.1 単体テスト
- **Entity**: データ検証ロジックのテスト
- **UseCase**: ビジネスロジックのテスト（モックリポジトリ使用）
- **Handler**: HTTPハンドラーのテスト（httptest使用）

### 7.2 統合テスト
- **Repository**: データベース操作のテスト（testcontainerまたはsqlmock）
- **Transaction**: トランザクション動作の検証

### 7.3 トランザクションテスト例
```go
func TestArticleUsecase_PublishArticle_TransactionRollback(t *testing.T) {
    // エラーが発生した場合のロールバック動作を検証
    // 複数テーブル操作の整合性確認
}
```

## 実装順序

1. **トランザクション基盤整備**
   - Transaction interface作成
   - TransactionManager実装

2. **Domain Layer**
   - Article entity作成
   - ArticleRepository interface作成

3. **Infrastructure Layer**
   - ArticleRepository実装（トランザクション対応）
   - Database connection更新

4. **Application Layer**
   - ArticleUsecase実装（トランザクション活用）

5. **Presentation Layer**
   - ArticleHandler実装
   - Router更新

6. **統合**
   - main.goで依存関係配線
   - 動作確認

7. **テスト作成**
   - 各レイヤーのテスト実装
   - トランザクション動作検証

## 注意事項

### データ整合性
- 記事作成時のAuthor存在確認
- 記事更新時の権限チェック（将来実装）
- 複数テーブル操作時のトランザクション必須

### パフォーマンス
- 記事一覧取得時のページネーション（将来実装）
- N+1問題の回避（Preload使用）

### セキュリティ
- SQL インジェクション対策（GORM使用で自動対応）
- 入力値検証の実装

この計画に従って段階的に実装を進めることで、保守性と拡張性の高い記事投稿機能を構築できます。