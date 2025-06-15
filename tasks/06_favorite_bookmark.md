# お気に入り・ブックマーク機能実装計画

## 概要
既存のClean Architectureに従い、ユーザーが記事を後で読むために保存し、個人的な読書リストを管理できるお気に入り・ブックマークシステムを実装します。

## 実装対象機能
- 記事のお気に入り追加・削除機能
- 読書リスト作成・管理機能
- 読書リストへの記事追加・削除
- 公開・非公開読書リスト機能
- お気に入り記事一覧表示

## データ設計

### Favorite エンティティ
```go
type Favorite struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    UserID    uint      `json:"user_id" gorm:"not null;index"`
    ArticleID uint      `json:"article_id" gorm:"not null;index"`
    CreatedAt time.Time `json:"created_at"`
    
    // Relations
    User    User    `json:"user" gorm:"foreignKey:UserID"`
    Article Article `json:"article" gorm:"foreignKey:ArticleID"`
}

// 重複お気に入り防止のための複合ユニーク制約
// gorm:"uniqueIndex:idx_user_article_favorite"
```

### ReadingList エンティティ
```go
type ReadingList struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name" gorm:"not null"`
    Description string    `json:"description"`
    UserID      uint      `json:"user_id" gorm:"not null;index"`
    IsPublic    bool      `json:"is_public" gorm:"default:false"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    
    // Relations
    User           User                `json:"user" gorm:"foreignKey:UserID"`
    ReadingListItems []ReadingListItem `json:"items,omitempty" gorm:"foreignKey:ReadingListID"`
}

type ReadingListItem struct {
    ID            uint      `json:"id" gorm:"primaryKey"`
    ReadingListID uint      `json:"reading_list_id" gorm:"not null;index"`
    ArticleID     uint      `json:"article_id" gorm:"not null;index"`
    AddedAt       time.Time `json:"added_at"`
    Notes         string    `json:"notes"` // 記事に関する個人メモ
    
    // Relations
    ReadingList ReadingList `json:"reading_list" gorm:"foreignKey:ReadingListID"`
    Article     Article     `json:"article" gorm:"foreignKey:ArticleID"`
}
```

### User エンティティ更新
```go
// 既存のUser構造体に追加
type User struct {
    // ... 既存フィールド
    
    // Relations
    Favorites    []Favorite    `json:"favorites,omitempty" gorm:"foreignKey:UserID"`
    ReadingLists []ReadingList `json:"reading_lists,omitempty" gorm:"foreignKey:UserID"`
}
```

### Article エンティティ更新
```go
// 既存のArticle構造体に追加
type Article struct {
    // ... 既存フィールド
    FavoriteCount int `json:"favorite_count" gorm:"default:0"` // パフォーマンス用の非正規化カウント
    
    // Relations
    Favorites         []Favorite         `json:"favorites,omitempty" gorm:"foreignKey:ArticleID"`
    ReadingListItems  []ReadingListItem  `json:"reading_list_items,omitempty" gorm:"foreignKey:ArticleID"`
}
```

## API設計

### エンドポイント

#### お気に入り管理
- `POST /api/v1/articles/:id/favorite` - 記事をお気に入りに追加
- `DELETE /api/v1/articles/:id/favorite` - 記事をお気に入りから削除
- `GET /api/v1/users/:id/favorites` - ユーザーのお気に入り記事一覧
- `GET /api/v1/articles/:id/favorites` - 記事をお気に入りしたユーザー一覧
- `GET /api/v1/articles/:id/favorite-status` - 現在のユーザーが記事をお気に入りしているかチェック

#### 読書リスト管理
- `POST /api/v1/reading-lists` - 読書リスト作成
- `GET /api/v1/reading-lists` - ユーザーの読書リスト一覧
- `GET /api/v1/reading-lists/:id` - 読書リスト詳細取得
- `PUT /api/v1/reading-lists/:id` - 読書リスト更新
- `DELETE /api/v1/reading-lists/:id` - 読書リスト削除

#### 読書リストアイテム管理
- `POST /api/v1/reading-lists/:id/articles` - 読書リストに記事追加
- `DELETE /api/v1/reading-lists/:id/articles/:article_id` - 読書リストから記事削除
- `GET /api/v1/reading-lists/:id/articles` - 読書リスト内の記事一覧
- `PUT /api/v1/reading-lists/:id/articles/:article_id` - 読書リストアイテム更新（メモ）

#### 公開読書リスト
- `GET /api/v1/reading-lists/public` - 公開読書リスト一覧
- `GET /api/v1/users/:id/reading-lists/public` - ユーザーの公開読書リスト一覧

### リクエスト/レスポンス例
```json
// POST /api/v1/articles/1/favorite
// レスポンス
{
  "message": "記事をお気に入りに追加しました",
  "favorite": {
    "id": 1,
    "user_id": 2,
    "article_id": 1,
    "created_at": "2025-06-14T10:00:00Z"
  }
}

// POST /api/v1/reading-lists
{
  "name": "プログラミング学習",
  "description": "プログラミングの基礎を学ぶための記事集",
  "is_public": false
}

// レスポンス
{
  "id": 1,
  "name": "プログラミング学習",
  "description": "プログラミングの基礎を学ぶための記事集",
  "user_id": 2,
  "is_public": false,
  "created_at": "2025-06-14T10:00:00Z",
  "updated_at": "2025-06-14T10:00:00Z"
}
```

## アーキテクチャ実装計画

### 1. Domain Layer（内側の層）

#### 1.1 エンティティ作成
- `internal/domain/entity/favorite.go` - Favorite構造体定義
- `internal/domain/entity/reading_list.go` - ReadingList・ReadingListItem構造体定義

#### 1.2 リポジトリインターフェース
- `internal/domain/repository/favorite_repository.go` - FavoriteRepository定義

```go
type FavoriteRepository interface {
    Create(ctx context.Context, favorite *entity.Favorite) error
    CreateWithTx(ctx context.Context, tx Transaction, favorite *entity.Favorite) error
    Delete(ctx context.Context, userID, articleID uint) error
    DeleteWithTx(ctx context.Context, tx Transaction, userID, articleID uint) error
    GetByUserID(ctx context.Context, userID uint) ([]*entity.Favorite, error)
    GetByArticleID(ctx context.Context, articleID uint) ([]*entity.Favorite, error)
    IsFavorited(ctx context.Context, userID, articleID uint) (bool, error)
    GetFavoriteCount(ctx context.Context, articleID uint) (int, error)
}
```

- `internal/domain/repository/reading_list_repository.go` - ReadingListRepository定義

```go
type ReadingListRepository interface {
    Create(ctx context.Context, readingList *entity.ReadingList) error
    Update(ctx context.Context, readingList *entity.ReadingList) error
    Delete(ctx context.Context, id uint) error
    GetByUserID(ctx context.Context, userID uint) ([]*entity.ReadingList, error)
    GetByID(ctx context.Context, id uint) (*entity.ReadingList, error)
    GetPublic(ctx context.Context) ([]*entity.ReadingList, error)
    AddItem(ctx context.Context, item *entity.ReadingListItem) error
    RemoveItem(ctx context.Context, readingListID, articleID uint) error
    GetItems(ctx context.Context, readingListID uint) ([]*entity.ReadingListItem, error)
    UpdateItem(ctx context.Context, item *entity.ReadingListItem) error
}
```

### 2. Infrastructure Layer（外側の層）

#### 2.1 リポジトリ実装
- `internal/infrastructure/repository/favorite_repository_impl.go`
- お気に入り数カウント更新の処理
- 重複チェック機能

- `internal/infrastructure/repository/reading_list_repository_impl.go`
- 読書リストアイテム管理
- プライバシー設定対応

#### 2.2 統合テスト
- `internal/infrastructure/repository/favorite_repository_impl_test.go`
- `internal/infrastructure/repository/reading_list_repository_impl_test.go`
- go-sqlmockを使用した統合テスト

### 3. Application Layer（ユースケース層）

#### 3.1 FavoriteUsecase実装
- `internal/usecase/favorite_usecase.go`

```go
type FavoriteUsecase interface {
    AddFavorite(ctx context.Context, userID, articleID uint) (*entity.Favorite, error)
    RemoveFavorite(ctx context.Context, userID, articleID uint) error
    GetUserFavorites(ctx context.Context, userID uint) ([]*entity.Favorite, error)
    GetArticleFavorites(ctx context.Context, articleID uint) ([]*entity.Favorite, error)
    IsFavorited(ctx context.Context, userID, articleID uint) (bool, error)
}
```

#### 3.2 ReadingListUsecase実装
- `internal/usecase/reading_list_usecase.go`

```go
type ReadingListUsecase interface {
    CreateReadingList(ctx context.Context, name, description string, userID uint, isPublic bool) (*entity.ReadingList, error)
    UpdateReadingList(ctx context.Context, id uint, name, description string, userID uint, isPublic bool) (*entity.ReadingList, error)
    DeleteReadingList(ctx context.Context, id, userID uint) error
    GetUserReadingLists(ctx context.Context, userID uint) ([]*entity.ReadingList, error)
    GetReadingList(ctx context.Context, id, userID uint) (*entity.ReadingList, error)
    GetPublicReadingLists(ctx context.Context) ([]*entity.ReadingList, error)
    AddArticleToList(ctx context.Context, readingListID, articleID, userID uint, notes string) (*entity.ReadingListItem, error)
    RemoveArticleFromList(ctx context.Context, readingListID, articleID, userID uint) error
    UpdateReadingListItem(ctx context.Context, readingListID, articleID, userID uint, notes string) (*entity.ReadingListItem, error)
}
```

#### 3.3 単体テスト
- `internal/usecase/favorite_usecase_test.go`
- `internal/usecase/reading_list_usecase_test.go`
- カスタムモックを使用したテスト

### 4. Presentation Layer（プレゼンテーション層）

#### 4.1 ハンドラー実装
- `internal/presentation/handler/favorite_handler.go`
- `internal/presentation/handler/reading_list_handler.go`
- お気に入り・読書リスト操作のHTTPハンドラー

#### 4.2 HTTPテスト
- `internal/presentation/handler/favorite_handler_test.go`
- `internal/presentation/handler/reading_list_handler_test.go`
- httptestを使用したHTTPテスト

#### 4.3 ルーティング設定
- `internal/presentation/router/router.go`
- お気に入り・読書リストAPIエンドポイントの追加

### 5. 依存関係の配線

#### 5.1 main.go更新
```go
// お気に入りリポジトリ初期化
favoriteRepo := repository.NewFavoriteRepository(database.DB)

// 読書リストリポジトリ初期化
readingListRepo := repository.NewReadingListRepository(database.DB)

// ユースケース初期化
favoriteUsecase := usecase.NewFavoriteUsecase(favoriteRepo, articleRepo, transactionManager)
readingListUsecase := usecase.NewReadingListUsecase(readingListRepo, articleRepo, userRepo)

// ハンドラー初期化
favoriteHandler := handler.NewFavoriteHandler(favoriteUsecase)
readingListHandler := handler.NewReadingListHandler(readingListUsecase)

// ルーター設定
router := router.SetupRouter(userHandler, articleHandler, commentHandler, categoryHandler, tagHandler, searchHandler, favoriteHandler, readingListHandler)
```

### 6. データベースシード

#### 6.1 シード作成
- `internal/infrastructure/seed/favorite_seed.go` - サンプルお気に入り
- `internal/infrastructure/seed/reading_list_seed.go` - サンプル読書リスト

## 実装順序

1. **Domain Layer**
   - Favorite・ReadingList entity作成
   - FavoriteRepository・ReadingListRepository interface作成

2. **Infrastructure Layer**
   - Favorite・ReadingListRepository実装
   - Database connection更新（AutoMigrate追加）

3. **Application Layer**
   - Favorite・ReadingListUsecase実装（ビジネスロジック）

4. **Presentation Layer**
   - Favorite・ReadingListHandler実装
   - Router更新

5. **統合**
   - main.goで依存関係配線
   - 動作確認

6. **データベースシード**
   - お気に入り・読書リストのサンプルデータ作成

7. **テスト作成**
   - 各レイヤーのテスト実装
   - お気に入り・読書リスト機能の検証

## データ検証ルール

### Favorite
- UserID: 必須、usersテーブルに存在
- ArticleID: 必須、articlesテーブルに存在
- ユニーク制約: ユーザーあたり記事一つにつき一つのお気に入り

### ReadingList
- Name: 必須、1文字以上100文字以下
- Description: 任意、500文字以下
- UserID: 必須、usersテーブルに存在
- IsPublic: ブール値、デフォルトfalse

### ReadingListItem
- ReadingListID: 必須、存在し、ユーザーに属する必要あり
- ArticleID: 必須、articlesテーブルに存在
- Notes: 任意、1000文字以下
- ユニーク制約: 読書リストあたり記事一つにつき一つのアイテム

## ビジネスロジック

### お気に入り機能
- ユーザーは公開済み記事のみお気に入り可能
- ユーザーは自分の記事をお気に入りできない（オプション）
- お気に入り追加時にarticleのfavorite_countをインクリメント
- お気に入り削除時にfavorite_countをデクリメント
- お気に入り数はトランザクション内で更新

### 読書リスト機能
- ユーザーは複数の読書リストを作成可能
- デフォルトの読書リストは非公開
- ユーザーは自分の読書リストのみ変更可能
- 公開読書リストは全ユーザーに表示
- 同じ記事を複数の読書リストに追加可能
- 記事には個人メモを付けて追加可能

### プライバシー・アクセス制御
- 非公開読書リストは所有者のみ閲覧可能
- 公開読書リストは全ユーザーが閲覧可能
- ユーザーは自分の読書リストのみアイテム追加・削除可能
- 読書リストの所有権は譲渡不可

## パフォーマンス考慮事項

### データベース最適化
- articleのfavorite_count非正規化による高速アクセス
- 複数お気に入りのバッチ操作
- 大きな読書リストのページネーション
- 外部キーと複合制約のインデックス
- 人気読書リストのキャッシュ

### キャッシュ戦略
- ユーザーのお気に入り一覧キャッシュ
- 公開読書リスト一覧キャッシュ
- 記事のお気に入り数キャッシュ

この計画に従って段階的に実装を進めることで、使いやすいお気に入り・ブックマーク機能を構築できます。