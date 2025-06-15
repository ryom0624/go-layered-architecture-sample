# カテゴリ・タグ機能実装計画

## 概要
既存のClean Architectureに従い、記事分類のためのカテゴリ・タグシステムを実装します。記事とタグの多対多関係を含む包括的な分類機能です。

## 実装対象機能
- カテゴリによる記事分類（1対多関係）
- タグによる記事分類（多対多関係）
- カテゴリ・タグによる記事フィルタリング
- スラッグベースのURL対応

## データ設計

### Category エンティティ
```go
type Category struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name" gorm:"unique;not null"`
    Description string    `json:"description"`
    Slug        string    `json:"slug" gorm:"unique;not null"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    
    // Relations
    Articles []Article `json:"articles,omitempty" gorm:"foreignKey:CategoryID"`
}
```

### Tag エンティティ
```go
type Tag struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name" gorm:"unique;not null"`
    Slug      string    `json:"slug" gorm:"unique;not null"`
    Color     string    `json:"color" gorm:"default:'#3B82F6'"` // UIでの表示色（16進数）
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    
    // Relations
    Articles []Article `json:"articles,omitempty" gorm:"many2many:article_tags;"`
}
```

### Article エンティティ更新
```go
// 既存のArticle構造体に追加
type Article struct {
    // ... 既存フィールド
    CategoryID *uint `json:"category_id,omitempty" gorm:"index"`
    
    // Relations
    Category *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
    Tags     []Tag     `json:"tags,omitempty" gorm:"many2many:article_tags;"`
}
```

## API設計

### エンドポイント

#### カテゴリ管理
- `POST /api/v1/categories` - カテゴリ作成
- `GET /api/v1/categories` - カテゴリ一覧取得
- `GET /api/v1/categories/:slug` - スラッグによるカテゴリ取得
- `PUT /api/v1/categories/:id` - カテゴリ更新
- `DELETE /api/v1/categories/:id` - カテゴリ削除
- `GET /api/v1/categories/:slug/articles` - カテゴリ別記事一覧

#### タグ管理
- `POST /api/v1/tags` - タグ作成
- `GET /api/v1/tags` - タグ一覧取得
- `GET /api/v1/tags/:slug` - スラッグによるタグ取得
- `PUT /api/v1/tags/:id` - タグ更新
- `DELETE /api/v1/tags/:id` - タグ削除
- `GET /api/v1/tags/:slug/articles` - タグ別記事一覧
- `GET /api/v1/tags/popular` - 人気タグ一覧

#### 記事フィルタリング
- `GET /api/v1/articles?category=:slug` - カテゴリによるフィルタ
- `GET /api/v1/articles?tags=:tag1,:tag2` - タグによるフィルタ
- `GET /api/v1/articles?category=:slug&tags=:tag1,:tag2` - 組み合わせフィルタ

### リクエスト/レスポンス例
```json
// POST /api/v1/categories
{
  "name": "プログラミング",
  "description": "プログラミング関連の記事"
}

// Response
{
  "id": 1,
  "name": "プログラミング",
  "description": "プログラミング関連の記事",
  "slug": "programming",
  "created_at": "2025-06-14T10:00:00Z",
  "updated_at": "2025-06-14T10:00:00Z"
}

// POST /api/v1/tags
{
  "name": "Go言語",
  "color": "#00ADD8"
}

// Response
{
  "id": 1,
  "name": "Go言語",
  "slug": "golang",
  "color": "#00ADD8",
  "created_at": "2025-06-14T10:00:00Z",
  "updated_at": "2025-06-14T10:00:00Z"
}
```

## アーキテクチャ実装計画

### 1. Domain Layer（内側の層）

#### 1.1 エンティティ作成
- `internal/domain/entity/category.go` - Category構造体定義
- `internal/domain/entity/tag.go` - Tag構造体定義
- `internal/domain/entity/article.go` - Article構造体更新（CategoryID、Tags追加）

#### 1.2 リポジトリインターフェース
- `internal/domain/repository/category_repository.go` - CategoryRepository定義

```go
type CategoryRepository interface {
    Create(ctx context.Context, category *entity.Category) error
    GetByID(ctx context.Context, id uint) (*entity.Category, error)
    GetBySlug(ctx context.Context, slug string) (*entity.Category, error)
    GetAll(ctx context.Context) ([]*entity.Category, error)
    Update(ctx context.Context, category *entity.Category) error
    Delete(ctx context.Context, id uint) error
    GetWithArticleCount(ctx context.Context) ([]*entity.Category, error)
}
```

- `internal/domain/repository/tag_repository.go` - TagRepository定義

```go
type TagRepository interface {
    Create(ctx context.Context, tag *entity.Tag) error
    GetByID(ctx context.Context, id uint) (*entity.Tag, error)
    GetBySlug(ctx context.Context, slug string) (*entity.Tag, error)
    GetAll(ctx context.Context) ([]*entity.Tag, error)
    Update(ctx context.Context, tag *entity.Tag) error
    Delete(ctx context.Context, id uint) error
    GetByNames(ctx context.Context, names []string) ([]*entity.Tag, error)
    GetPopular(ctx context.Context, limit int) ([]*entity.Tag, error)
}
```

### 2. Infrastructure Layer（外側の層）

#### 2.1 リポジトリ実装
- `internal/infrastructure/repository/category_repository_impl.go`
- GORMを使用した実装、記事数カウント機能付き

- `internal/infrastructure/repository/tag_repository_impl.go`
- 多対多関係の処理、人気タグ取得機能

#### 2.2 統合テスト
- `internal/infrastructure/repository/category_repository_impl_test.go`
- `internal/infrastructure/repository/tag_repository_impl_test.go`
- go-sqlmockを使用した統合テスト

### 3. Application Layer（ユースケース層）

#### 3.1 CategoryUsecase実装
- `internal/usecase/category_usecase.go`

```go
type CategoryUsecase interface {
    CreateCategory(ctx context.Context, name, description string) (*entity.Category, error)
    GetCategory(ctx context.Context, slug string) (*entity.Category, error)
    GetAllCategories(ctx context.Context) ([]*entity.Category, error)
    UpdateCategory(ctx context.Context, id uint, name, description string) (*entity.Category, error)
    DeleteCategory(ctx context.Context, id uint) error
    GetCategoriesWithCount(ctx context.Context) ([]*entity.Category, error)
}
```

#### 3.2 TagUsecase実装
- `internal/usecase/tag_usecase.go`

```go
type TagUsecase interface {
    CreateTag(ctx context.Context, name, color string) (*entity.Tag, error)
    GetTag(ctx context.Context, slug string) (*entity.Tag, error)
    GetAllTags(ctx context.Context) ([]*entity.Tag, error)
    UpdateTag(ctx context.Context, id uint, name, color string) (*entity.Tag, error)
    DeleteTag(ctx context.Context, id uint) error
    GetOrCreateTags(ctx context.Context, tagNames []string) ([]*entity.Tag, error)
    GetPopularTags(ctx context.Context, limit int) ([]*entity.Tag, error)
}
```

#### 3.3 ArticleUsecase更新
- `internal/usecase/article_usecase.go`に以下メソッド追加：
  - `GetArticlesByCategory`
  - `GetArticlesByTags`
  - CreateArticle・UpdateArticleにカテゴリ・タグパラメータ追加

#### 3.4 単体テスト
- `internal/usecase/category_usecase_test.go`
- `internal/usecase/tag_usecase_test.go`
- カスタムモックを使用したテスト

### 4. Presentation Layer（プレゼンテーション層）

#### 4.1 ハンドラー実装
- `internal/presentation/handler/category_handler.go`
- `internal/presentation/handler/tag_handler.go`
- スラッグベースのルーティング対応

#### 4.2 ArticleHandler更新
- `internal/presentation/handler/article_handler.go`
- カテゴリ・タグパラメータ受け取り機能追加
- フィルタリングエンドポイント追加

#### 4.3 HTTPテスト
- `internal/presentation/handler/category_handler_test.go`
- `internal/presentation/handler/tag_handler_test.go`
- httptestを使用したHTTPテスト

#### 4.4 ルーティング設定
- `internal/presentation/router/router.go`
- カテゴリ・タグエンドポイント追加
- 記事フィルタリングルート更新

### 5. 依存関係の配線

#### 5.1 main.go更新
```go
// カテゴリリポジトリ初期化
categoryRepo := repository.NewCategoryRepository(database.DB)

// タグリポジトリ初期化
tagRepo := repository.NewTagRepository(database.DB)

// ユースケース初期化
categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
tagUsecase := usecase.NewTagUsecase(tagRepo)
articleUsecase := usecase.NewArticleUsecase(articleRepo, userRepo, categoryRepo, tagRepo, transactionManager)

// ハンドラー初期化
categoryHandler := handler.NewCategoryHandler(categoryUsecase)
tagHandler := handler.NewTagHandler(tagUsecase)

// ルーター設定
router := router.SetupRouter(userHandler, articleHandler, commentHandler, categoryHandler, tagHandler)
```

### 6. データベースシード

#### 6.1 シード作成
- `internal/infrastructure/seed/category_seed.go` - サンプルカテゴリ
- `internal/infrastructure/seed/tag_seed.go` - サンプルタグ
- `internal/infrastructure/seeder.go` - 記事とカテゴリ・タグの関連付け

## 実装順序

1. **Domain Layer**
   - Category・Tag entity作成
   - CategoryRepository・TagRepository interface作成

2. **Infrastructure Layer**
   - Category・TagRepository実装（多対多関係対応）
   - Database connection更新（AutoMigrate追加）

3. **Application Layer**
   - Category・TagUsecase実装（スラッグ生成ロジック）
   - ArticleUsecase更新（フィルタリング機能追加）

4. **Presentation Layer**
   - Category・TagHandler実装
   - ArticleHandler更新（フィルタリング対応）
   - Router更新

5. **統合**
   - main.goで依存関係配線
   - 動作確認

6. **データベースシード**
   - カテゴリ・タグのサンプルデータ作成
   - 記事との関連付け

7. **テスト作成**
   - 各レイヤーのテスト実装
   - 多対多関係の検証

## データ検証ルール

### Category
- Name: 必須、2文字以上100文字以下、ユニーク
- Slug: 名前から自動生成、ユニーク
- Description: 任意、500文字以下

### Tag
- Name: 必須、2文字以上50文字以下、ユニーク
- Slug: 名前から自動生成、ユニーク
- Color: 任意、有効な16進数カラーコード

## ビジネスロジック

### スラッグ生成
- カテゴリ・タグのスラッグは名前から自動生成
- 日本語名は適切なURL形式に変換
- 重複回避のためのナンバリング対応

### 記事との関連
- 記事は1つのカテゴリ（任意）と複数のタグを持つ
- カテゴリ削除時は記事のcategory_idをnullに設定
- タグ削除時は関連テーブルから関連を削除、記事は維持

### 人気タグ
- 使用回数（記事数）に基づく人気度計算
- キャッシュ機能による高速化

## パフォーマンス考慮事項

### データベース最適化
- slug、category_idにインデックス追加
- article_tags中間テーブルの最適化
- N+1問題の回避（Preload使用）

### キャッシュ戦略
- 人気タグのキャッシュ
- カテゴリ一覧のキャッシュ
- 記事数カウントのキャッシュ

この計画に従って段階的に実装を進めることで、柔軟で拡張性の高い分類システムを構築できます。