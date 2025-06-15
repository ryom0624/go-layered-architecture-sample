# 閲覧履歴・統計機能実装計画

## 概要
既存のClean Architectureに従い、記事の閲覧履歴、ユーザーの読書行動を追跡し、コンテンツパフォーマンスの分析機能を提供する包括的な統計システムを実装します。

## 実装対象機能
- 記事閲覧の自動トラッキング
- ユーザー読書履歴管理
- 記事統計情報（閲覧数、平均読了時間等）
- 日次・期間別統計レポート
- 人気記事・トレンド記事分析
- 読書進捗管理機能

## データ設計

### ArticleView エンティティ
```go
type ArticleView struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    ArticleID uint      `json:"article_id" gorm:"not null;index"`
    UserID    *uint     `json:"user_id,omitempty" gorm:"index"` // 匿名ユーザーの場合はnull
    IPAddress string    `json:"ip_address" gorm:"index"`
    UserAgent string    `json:"user_agent"`
    ViewedAt  time.Time `json:"viewed_at" gorm:"index"`
    Duration  int       `json:"duration"` // 読了時間（秒）
    
    // Relations
    Article *Article `json:"article,omitempty" gorm:"foreignKey:ArticleID"`
    User    *User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
```

### UserReadingHistory エンティティ
```go
type UserReadingHistory struct {
    ID         uint      `json:"id" gorm:"primaryKey"`
    UserID     uint      `json:"user_id" gorm:"not null;index"`
    ArticleID  uint      `json:"article_id" gorm:"not null;index"`
    FirstViewAt time.Time `json:"first_view_at"`
    LastViewAt  time.Time `json:"last_view_at"`
    ViewCount   int       `json:"view_count" gorm:"default:1"`
    TotalDuration int     `json:"total_duration"` // 累計読了時間（秒）
    IsCompleted   bool    `json:"is_completed" gorm:"default:false"` // 最後まで読了したか
    Progress      float64 `json:"progress" gorm:"default:0"` // 読了進捗率（%）
    
    // Relations
    User    User    `json:"user" gorm:"foreignKey:UserID"`
    Article Article `json:"article" gorm:"foreignKey:ArticleID"`
}
```

### ArticleStatistics エンティティ
```go
type ArticleStatistics struct {
    ID                uint      `json:"id" gorm:"primaryKey"`
    ArticleID         uint      `json:"article_id" gorm:"unique;not null"`
    TotalViews        int       `json:"total_views" gorm:"default:0"`
    UniqueViews       int       `json:"unique_views" gorm:"default:0"`
    AuthenticatedViews int      `json:"authenticated_views" gorm:"default:0"`
    AnonymousViews    int       `json:"anonymous_views" gorm:"default:0"`
    AverageReadTime   float64   `json:"average_read_time" gorm:"default:0"` // 平均読了時間（秒）
    CompletionRate    float64   `json:"completion_rate" gorm:"default:0"` // 読了率（%）
    LastViewedAt      *time.Time `json:"last_viewed_at,omitempty"`
    UpdatedAt         time.Time `json:"updated_at"`
    
    // Relations
    Article Article `json:"article" gorm:"foreignKey:ArticleID"`
}
```

### DailyStatistics エンティティ
```go
type DailyStatistics struct {
    ID             uint      `json:"id" gorm:"primaryKey"`
    Date           time.Time `json:"date" gorm:"unique;not null;index"`
    TotalViews     int       `json:"total_views" gorm:"default:0"`
    UniqueVisitors int       `json:"unique_visitors" gorm:"default:0"`
    NewUsers       int       `json:"new_users" gorm:"default:0"`
    ArticlesRead   int       `json:"articles_read" gorm:"default:0"`
    AverageSession float64   `json:"average_session" gorm:"default:0"` // 平均セッション時間（分）
    TopArticleID   *uint     `json:"top_article_id,omitempty"`
    
    // Relations
    TopArticle *Article `json:"top_article,omitempty" gorm:"foreignKey:TopArticleID"`
}
```

### Article エンティティ更新
```go
// 既存のArticle構造体に追加
type Article struct {
    // ... 既存フィールド
    ViewCount      int      `json:"view_count" gorm:"default:0"`
    UniqueViews    int      `json:"unique_views" gorm:"default:0"`
    LastViewedAt   *time.Time `json:"last_viewed_at,omitempty"`
    
    // Relations
    Views         []ArticleView      `json:"views,omitempty" gorm:"foreignKey:ArticleID"`
    Statistics    *ArticleStatistics `json:"statistics,omitempty" gorm:"foreignKey:ArticleID"`
    ReadingHistory []UserReadingHistory `json:"reading_history,omitempty" gorm:"foreignKey:ArticleID"`
}
```

### User エンティティ更新
```go
// 既存のUser構造体に追加
type User struct {
    // ... 既存フィールド
    TotalReadingTime int `json:"total_reading_time" gorm:"default:0"` // 累計読了時間（秒）
    ArticlesRead     int `json:"articles_read" gorm:"default:0"`
    
    // Relations
    Views          []ArticleView         `json:"views,omitempty" gorm:"foreignKey:UserID"`
    ReadingHistory []UserReadingHistory  `json:"reading_history,omitempty" gorm:"foreignKey:UserID"`
}
```

## API設計

### エンドポイント

#### 閲覧履歴（自動トラッキング）
- 記事アクセス時にミドルウェアで自動記録
- ユーザーID、IPアドレス、ユーザーエージェント、タイムスタンプを取得

#### 読書履歴管理
- `GET /api/v1/users/:id/reading-history` - ユーザーの読書履歴取得
- `GET /api/v1/users/:id/reading-progress` - 記事の読了進捗取得
- `PUT /api/v1/articles/:id/reading-progress` - 読了進捗更新
- `POST /api/v1/articles/:id/mark-complete` - 記事を読了済みにマーク

#### 記事統計
- `GET /api/v1/articles/:id/statistics` - 記事の詳細統計取得
- `GET /api/v1/articles/:id/views` - 記事の閲覧履歴取得
- `GET /api/v1/articles/popular` - 閲覧数に基づく人気記事取得
- `GET /api/v1/articles/trending` - トレンド記事取得（最近の人気度）

#### ユーザー分析
- `GET /api/v1/users/:id/analytics` - ユーザーの読書分析
- `GET /api/v1/users/:id/most-read-categories` - ユーザーの読書傾向
- `GET /api/v1/users/:id/reading-streaks` - 読書継続記録

#### プラットフォーム分析（管理者）
- `GET /api/v1/analytics/overview` - プラットフォーム概要統計
- `GET /api/v1/analytics/daily` - 日次統計
- `GET /api/v1/analytics/articles/top` - 高パフォーマンス記事
- `GET /api/v1/analytics/users/active` - アクティブユーザー統計

### リクエスト/レスポンス例
```json
// PUT /api/v1/articles/1/reading-progress
{
  "progress": 75.5,
  "duration": 180
}

// GET /api/v1/articles/1/statistics
{
  "article_id": 1,
  "total_views": 1250,
  "unique_views": 980,
  "authenticated_views": 750,
  "anonymous_views": 500,
  "average_read_time": 245.5,
  "completion_rate": 68.5,
  "last_viewed_at": "2025-06-14T15:30:00Z"
}

// GET /api/v1/users/2/analytics
{
  "user_id": 2,
  "total_reading_time": 12500,
  "articles_read": 45,
  "average_read_time": 278,
  "favorite_categories": [
    {"category": "プログラミング", "count": 20},
    {"category": "技術", "count": 15}
  ],
  "reading_streak": {
    "current_streak": 7,
    "longest_streak": 23
  }
}
```

## アーキテクチャ実装計画

### 1. Domain Layer（内側の層）

#### 1.1 エンティティ作成
- `internal/domain/entity/article_view.go` - ArticleView構造体定義
- `internal/domain/entity/user_reading_history.go` - UserReadingHistory構造体定義
- `internal/domain/entity/article_statistics.go` - ArticleStatistics・DailyStatistics構造体定義

#### 1.2 リポジトリインターフェース
- `internal/domain/repository/view_repository.go` - ViewRepository定義

```go
type ViewRepository interface {
    RecordView(ctx context.Context, view *entity.ArticleView) error
    GetViewHistory(ctx context.Context, userID uint) ([]*entity.ArticleView, error)
    GetArticleViews(ctx context.Context, articleID uint) ([]*entity.ArticleView, error)
    GetUserViews(ctx context.Context, userID uint, limit int) ([]*entity.ArticleView, error)
    UpdateReadingProgress(ctx context.Context, userID, articleID uint, progress float64, duration int) error
}
```

- `internal/domain/repository/statistics_repository.go` - StatisticsRepository定義

```go
type StatisticsRepository interface {
    UpdateArticleStats(ctx context.Context, articleID uint) error
    GetArticleStats(ctx context.Context, articleID uint) (*entity.ArticleStatistics, error)
    GetDailyStats(ctx context.Context, date time.Time) (*entity.DailyStatistics, error)
    GetPopularArticles(ctx context.Context, limit int) ([]*entity.Article, error)
    GetTrendingArticles(ctx context.Context, days, limit int) ([]*entity.Article, error)
    CreateDailyStats(ctx context.Context, stats *entity.DailyStatistics) error
}
```

### 2. Infrastructure Layer（外側の層）

#### 2.1 リポジトリ実装
- `internal/infrastructure/repository/view_repository_impl.go`
- 閲覧記録と集約処理
- 重複閲覧検知機能

- `internal/infrastructure/repository/statistics_repository_impl.go`
- 統計計算とキャッシュ処理
- 日次統計バッチ処理

#### 2.2 統合テスト
- `internal/infrastructure/repository/view_repository_impl_test.go`
- `internal/infrastructure/repository/statistics_repository_impl_test.go`
- go-sqlmockを使用した統合テスト

### 3. Application Layer（ユースケース層）

#### 3.1 ViewUsecase実装
- `internal/usecase/view_usecase.go`

```go
type ViewUsecase interface {
    RecordArticleView(ctx context.Context, articleID uint, userID *uint, ipAddress, userAgent string) error
    GetUserReadingHistory(ctx context.Context, userID uint) ([]*entity.UserReadingHistory, error)
    UpdateReadingProgress(ctx context.Context, userID, articleID uint, progress float64, duration int) error
    MarkAsCompleted(ctx context.Context, userID, articleID uint) error
    GetReadingAnalytics(ctx context.Context, userID uint) (*UserAnalytics, error)
}
```

#### 3.2 StatisticsUsecase実装
- `internal/usecase/statistics_usecase.go`

```go
type StatisticsUsecase interface {
    GetArticleStatistics(ctx context.Context, articleID uint) (*entity.ArticleStatistics, error)
    GetPlatformOverview(ctx context.Context) (*PlatformOverview, error)
    GetDailyStatistics(ctx context.Context, from, to time.Time) ([]*entity.DailyStatistics, error)
    GetPopularArticles(ctx context.Context, limit int) ([]*entity.Article, error)
    GetTrendingArticles(ctx context.Context, days, limit int) ([]*entity.Article, error)
    GenerateDailyReport(ctx context.Context, date time.Time) error
}
```

#### 3.3 単体テスト
- `internal/usecase/view_usecase_test.go`
- `internal/usecase/statistics_usecase_test.go`
- カスタムモックを使用したテスト

### 4. Middleware（ミドルウェア）

#### 4.1 閲覧トラッキング
- `internal/presentation/middleware/view_tracking.go`
- 記事アクセス時の自動閲覧記録
- ユーザー情報・IPアドレス抽出

### 5. Presentation Layer（プレゼンテーション層）

#### 5.1 ハンドラー実装
- `internal/presentation/handler/view_handler.go`
- `internal/presentation/handler/statistics_handler.go`
- 閲覧・統計操作のHTTPハンドラー

#### 5.2 HTTPテスト
- `internal/presentation/handler/view_handler_test.go`
- `internal/presentation/handler/statistics_handler_test.go`
- httptestを使用したHTTPテスト

#### 5.3 ルーティング設定
- `internal/presentation/router/router.go`
- 記事ルートに閲覧トラッキングミドルウェア追加
- 統計・分析エンドポイント追加

### 6. 依存関係の配線

#### 6.1 main.go更新
```go
// 閲覧・統計リポジトリ初期化
viewRepo := repository.NewViewRepository(database.DB)
statisticsRepo := repository.NewStatisticsRepository(database.DB)

// ユースケース初期化
viewUsecase := usecase.NewViewUsecase(viewRepo, articleRepo, userRepo)
statisticsUsecase := usecase.NewStatisticsUsecase(statisticsRepo, viewRepo)

// ハンドラー初期化
viewHandler := handler.NewViewHandler(viewUsecase)
statisticsHandler := handler.NewStatisticsHandler(statisticsUsecase)

// ミドルウェア初期化
viewTrackingMiddleware := middleware.NewViewTracking(viewUsecase)

// ルーター設定
router := router.SetupRouter(handlers..., viewTrackingMiddleware)
```

### 7. バックグラウンドジョブ

#### 7.1 統計集約ジョブ
- `internal/infrastructure/jobs/statistics_aggregator.go`
- 日次統計の集約バックグラウンドジョブ
- 記事統計の定期更新

## 実装順序

1. **Domain Layer**
   - ArticleView・UserReadingHistory・ArticleStatistics entity作成
   - ViewRepository・StatisticsRepository interface作成

2. **Infrastructure Layer**
   - View・StatisticsRepository実装（閲覧記録・統計計算）
   - Database connection更新（AutoMigrate追加）

3. **Application Layer**
   - View・StatisticsUsecase実装（ビジネスロジック）

4. **Middleware**
   - 閲覧トラッキングミドルウェア実装

5. **Presentation Layer**
   - View・StatisticsHandler実装
   - Router更新（ミドルウェア適用）

6. **統合**
   - main.goで依存関係配線
   - 動作確認

7. **バックグラウンドジョブ**
   - 統計集約ジョブ実装
   - 定期実行設定

8. **テスト作成**
   - 各レイヤーのテスト実装
   - 閲覧・統計機能の検証

## バックグラウンド処理

### リアルタイム閲覧記録
- 個別閲覧の即座記録
- ユーザー読書履歴のリアルタイム更新
- 読了進捗と完了の追跡

### バッチ統計更新
- 記事統計の時間別集約
- 日次プラットフォーム統計の計算
- トレンド・人気記事ランキングの更新
- 古い閲覧記録のクリーンアップ（データ保持ポリシー）

## データ検証ルール

### 入力検証
- ArticleID: 必須、articlesテーブルに存在
- UserID: 任意（匿名ユーザー対応）、存在する場合はusersテーブルに存在
- Progress: 0〜100の範囲
- Duration: 正の整数値（秒）
- IPAddress: 有効なIPアドレス形式

## ビジネスロジック

### 閲覧トラッキング
- 記事アクセス時の自動記録（ミドルウェア）
- 同一ユーザー・同一記事の重複排除ロジック（時間窓内）
- 匿名ユーザーとログインユーザーの区別
- 読了進捗の段階的更新

### 統計計算
- 記事統計の定期更新（時間別バッチ）
- 日次統計の夜間集約
- 人気記事・トレンド記事の動的ランキング
- ユーザー分析データの生成

### プライバシー・データ保護
- ユーザーの詳細追跡オプトアウト機能
- 読書履歴のユーザー削除機能
- IPアドレスのハッシュ化
- GDPR準拠のデータエクスポート・削除

## パフォーマンス考慮事項

### データベース最適化
```sql
-- 閲覧トラッキング用インデックス
CREATE INDEX idx_article_views_article_id_viewed_at ON article_views(article_id, viewed_at DESC);
CREATE INDEX idx_article_views_user_id ON article_views(user_id);
CREATE INDEX idx_article_views_ip_address ON article_views(ip_address);

-- 読書履歴用インデックス
CREATE INDEX idx_user_reading_history_user_id ON user_reading_history(user_id);
CREATE INDEX idx_user_reading_history_last_view_at ON user_reading_history(last_view_at DESC);

-- 統計用インデックス
CREATE INDEX idx_article_statistics_total_views ON article_statistics(total_views DESC);
CREATE INDEX idx_daily_statistics_date ON daily_statistics(date DESC);
```

### パフォーマンス最適化
- リアルタイムカウンター更新のためのデータベーストリガー使用
- 閲覧重複排除ロジック（同一ユーザー・同一記事・時間窓内）
- パフォーマンス維持のための古い閲覧記録アーカイブ
- 複雑な分析クエリのためのマテリアライズドビュー使用

### キャッシュ戦略
- 人気記事ランキングキャッシュ
- ユーザー分析データキャッシュ
- 日次統計サマリーキャッシュ

## 分析機能

### 記事パフォーマンス
- 時系列閲覧トレンド
- 読了完了率
- 平均読了時間
- トラフィックソース・リファラー分析
- 人気コンテンツトピック

### ユーザーエンゲージメント
- 読書継続記録・習慣
- お気に入りカテゴリ・トピック
- 読書速度・好み
- 時間帯別読書パターン

### プラットフォーム洞察
- 全体エンゲージメント指標
- コンテンツパフォーマンスランキング
- ユーザー継続・成長率
- 人気読書時間・パターン

### データ保持ポリシー
- 匿名閲覧記録: 90日間
- ユーザー読書履歴: 無期限（ユーザーが削除可能）
- 日次統計: 無期限（集約データ）

この計画に従って段階的に実装を進めることで、包括的で高性能な閲覧履歴・統計システムを構築できます。