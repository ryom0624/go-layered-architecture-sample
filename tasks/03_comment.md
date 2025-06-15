# コメント機能実装計画

## 概要
既存のClean Architectureに従い、記事に対する階層的コメント機能とモデレーション機能を実装します。

## 実装対象機能
- 記事へのコメント投稿・編集・削除
- 階層的返信機能（最大3階層）
- コメント承認システム
- コメントモデレーション機能

## データ設計

### Comment エンティティ
```go
type Comment struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Content     string    `json:"content" gorm:"not null"`
    AuthorID    uint      `json:"author_id" gorm:"not null"`
    ArticleID   uint      `json:"article_id" gorm:"not null"`
    ParentID    *uint     `json:"parent_id,omitempty" gorm:"index"` // 返信の場合の親コメントID
    Status      string    `json:"status" gorm:"default:'pending'"` // pending, approved, rejected
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    
    // Relations
    Author   User      `json:"author" gorm:"foreignKey:AuthorID"`
    Article  Article   `json:"article" gorm:"foreignKey:ArticleID"`
    Parent   *Comment  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
    Replies  []Comment `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
}
```

## API設計

### エンドポイント

#### コメント管理
- `POST /api/v1/articles/:id/comments` - 記事にコメント投稿
- `POST /api/v1/comments/:id/replies` - コメントに返信
- `GET /api/v1/articles/:id/comments` - 記事のコメント一覧取得（返信含む）
- `GET /api/v1/comments/:id` - 特定コメント詳細取得
- `PUT /api/v1/comments/:id` - コメント更新（投稿者のみ）
- `DELETE /api/v1/comments/:id` - コメント削除（投稿者のみ）

#### コメントモデレーション
- `PUT /api/v1/comments/:id/approve` - コメント承認（管理者）
- `PUT /api/v1/comments/:id/reject` - コメント拒否（管理者）
- `GET /api/v1/comments/pending` - 承認待ちコメント一覧（管理者）

### リクエスト/レスポンス例
```json
// POST /api/v1/articles/1/comments
{
  "content": "素晴らしい記事でした！"
}

// Response
{
  "id": 1,
  "content": "素晴らしい記事でした！",
  "author_id": 2,
  "article_id": 1,
  "parent_id": null,
  "status": "pending",
  "created_at": "2025-06-14T10:00:00Z",
  "updated_at": "2025-06-14T10:00:00Z",
  "author": {
    "id": 2,
    "name": "田中太郎",
    "email": "tanaka@example.com"
  }
}
```

## アーキテクチャ実装計画

### 1. Domain Layer（内側の層）

#### 1.1 エンティティ作成
- `internal/domain/entity/comment.go` - Comment構造体定義

#### 1.2 リポジトリインターフェース
- `internal/domain/repository/comment_repository.go` - CommentRepository定義

```go
type CommentRepository interface {
    Create(ctx context.Context, comment *entity.Comment) error
    CreateWithTx(ctx context.Context, tx Transaction, comment *entity.Comment) error
    GetByID(ctx context.Context, id uint) (*entity.Comment, error)
    GetByArticleID(ctx context.Context, articleID uint) ([]*entity.Comment, error)
    GetReplies(ctx context.Context, parentID uint) ([]*entity.Comment, error)
    Update(ctx context.Context, comment *entity.Comment) error
    UpdateWithTx(ctx context.Context, tx Transaction, comment *entity.Comment) error
    Delete(ctx context.Context, id uint) error
    GetByStatus(ctx context.Context, status string) ([]*entity.Comment, error)
}
```

### 2. Infrastructure Layer（外側の層）

#### 2.1 リポジトリ実装
- `internal/infrastructure/repository/comment_repository_impl.go`
- GORMを使用した階層的コメント取得の実装
- Author・Articleリレーションのプリロード対応

#### 2.2 統合テスト
- `internal/infrastructure/repository/comment_repository_impl_test.go`
- go-sqlmockを使用した統合テスト
- 階層的コメント取得のテスト

### 3. Application Layer（ユースケース層）

#### 3.1 CommentUsecase実装
- `internal/usecase/comment_usecase.go`

```go
type CommentUsecase interface {
    CreateComment(ctx context.Context, content string, authorID, articleID uint) (*entity.Comment, error)
    CreateReply(ctx context.Context, content string, authorID, parentID uint) (*entity.Comment, error)
    GetArticleComments(ctx context.Context, articleID uint) ([]*entity.Comment, error)
    GetComment(ctx context.Context, id uint) (*entity.Comment, error)
    UpdateComment(ctx context.Context, id uint, content string, userID uint) (*entity.Comment, error)
    DeleteComment(ctx context.Context, id uint, userID uint) error
    ApproveComment(ctx context.Context, id uint) (*entity.Comment, error)
    RejectComment(ctx context.Context, id uint) (*entity.Comment, error)
    GetPendingComments(ctx context.Context) ([]*entity.Comment, error)
}
```

#### 3.2 単体テスト
- `internal/usecase/comment_usecase_test.go`
- カスタムモックを使用したテスト
- コメント検証・承認ロジックのテスト

### 4. Presentation Layer（プレゼンテーション層）

#### 4.1 ハンドラー実装
- `internal/presentation/handler/comment_handler.go`
- ネストしたコメント構造のレスポンス処理

#### 4.2 HTTPテスト
- `internal/presentation/handler/comment_handler_test.go`
- httptestを使用したHTTPテスト

#### 4.3 ルーティング設定
- `internal/presentation/router/router.go`
- コメントAPI エンドポイントの追加

### 5. 依存関係の配線

#### 5.1 main.go更新
```go
// コメントリポジトリ初期化
commentRepo := repository.NewCommentRepository(database.DB)

// コメントユースケース初期化
commentUsecase := usecase.NewCommentUsecase(commentRepo, userRepo, articleRepo, transactionManager)

// コメントハンドラー初期化
commentHandler := handler.NewCommentHandler(commentUsecase)

// ルーター設定
router := router.SetupRouter(userHandler, articleHandler, commentHandler)
```

## 実装順序

1. **Domain Layer**
   - Comment entity作成
   - CommentRepository interface作成

2. **Infrastructure Layer**
   - CommentRepository実装（階層的クエリ対応）
   - Database connection更新（AutoMigrate追加）

3. **Application Layer**
   - CommentUsecase実装（ビジネスロジック）

4. **Presentation Layer**
   - CommentHandler実装
   - Router更新

5. **統合**
   - main.goで依存関係配線
   - 動作確認

6. **テスト作成**
   - 各レイヤーのテスト実装
   - 階層的コメント機能の検証

## データ検証ルール

### 入力検証
- Content: 必須、1文字以上1000文字以下
- AuthorID: 必須、usersテーブルに存在
- ArticleID: 必須、articlesテーブルに存在
- ParentID: 任意、存在する場合はcommentsテーブルに存在
- Status: 'pending', 'approved', 'rejected'のいずれか

## ビジネスロジック

### コメント投稿
- 新規コメントのデフォルトステータスは'pending'
- 公開済み記事にのみコメント可能
- 階層の深さ制限（最大3階層）

### 権限管理
- ユーザーは自身のコメントのみ編集・削除可能
- 管理者はすべてのコメントを承認・拒否可能

### データ整合性
- 親コメント削除時の子コメント処理
- コメント削除は論理削除（ステータス変更）を推奨

## パフォーマンス考慮事項

### データベース最適化
- ArticleID、ParentID、Statusにインデックス追加
- 階層的クエリのパフォーマンス最適化
- N+1問題の回避（Preload使用）

### キャッシュ戦略
- 記事ごとのコメント数キャッシュ
- 承認済みコメントの一時キャッシュ

この計画に従って段階的に実装を進めることで、スケーラブルなコメント機能を構築できます。