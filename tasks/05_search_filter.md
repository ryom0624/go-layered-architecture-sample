# 検索・フィルタリング機能実装計画

## 概要
既存のClean Architectureに従い、記事の包括的な検索・フィルタリング機能を実装します。全文検索、日付フィルタ、高度なクエリ機能を含みます。

## 実装対象機能
- 記事タイトル・本文の全文検索
- カテゴリ・タグによるフィルタリング
- 作成者・ステータスによるフィルタリング
- 日付範囲でのフィルタリング
- ソート機能（関連度、日付、タイトル等）
- ページネーション機能

## データ設計

### 検索機能拡張
```go
// Article エンティティに検索対応を追加
type Article struct {
    // ... 既存フィールド
    SearchVector string `json:"-" gorm:"type:text"` // 全文検索最適化用
}
```

### 検索パラメータ構造体
```go
type ArticleSearchParams struct {
    Query        string    `json:"query" form:"query"`                // 全文検索クエリ
    CategorySlug string    `json:"category" form:"category"`          // カテゴリフィルタ
    Tags         []string  `json:"tags" form:"tags"`                  // タグフィルタ（カンマ区切り）
    AuthorID     uint      `json:"author_id" form:"author_id"`        // 作成者フィルタ
    Status       string    `json:"status" form:"status"`              // ステータスフィルタ
    DateFrom     time.Time `json:"date_from" form:"date_from"`        // 日付範囲開始
    DateTo       time.Time `json:"date_to" form:"date_to"`            // 日付範囲終了
    SortBy       string    `json:"sort_by" form:"sort_by"`            // ソートフィールド
    SortOrder    string    `json:"sort_order" form:"sort_order"`      // asc/desc
    Page         int       `json:"page" form:"page"`                  // ページ番号
    Limit        int       `json:"limit" form:"limit"`                // ページサイズ
}

type ArticleSearchResult struct {
    Articles   []*Article `json:"articles"`
    Total      int64      `json:"total"`
    Page       int        `json:"page"`
    Limit      int        `json:"limit"`
    TotalPages int        `json:"total_pages"`
}
```

## API設計

### エンドポイント

#### 検索エンドポイント
- `GET /api/v1/search` - 全体検索
- `GET /api/v1/articles/search` - 記事専用検索
- `GET /api/v1/articles/popular` - 人気記事取得
- `GET /api/v1/articles/recent` - 最新記事取得

#### 拡張記事エンドポイント
- `GET /api/v1/articles?q=:query` - 全文検索
- `GET /api/v1/articles?category=:slug` - カテゴリフィルタ
- `GET /api/v1/articles?tags=:tag1,:tag2` - タグフィルタ
- `GET /api/v1/articles?author_id=:id` - 作成者フィルタ
- `GET /api/v1/articles?status=:status` - ステータスフィルタ
- `GET /api/v1/articles?date_from=:date&date_to=:date` - 日付範囲フィルタ
- `GET /api/v1/articles?sort_by=:field&sort_order=:order` - ソート
- `GET /api/v1/articles?page=:page&limit=:limit` - ページネーション

### 複合フィルタリング例
```
GET /api/v1/articles?q=golang&category=programming&tags=backend,api&sort_by=created_at&sort_order=desc&page=1&limit=10
```

### リクエスト/レスポンス例
```json
// GET /api/v1/articles/search?q=Go言語&category=programming&page=1&limit=5
{
  "articles": [
    {
      "id": 1,
      "title": "Go言語入門",
      "content": "Go言語の基礎について...",
      "author": {
        "id": 1,
        "name": "田中太郎"
      },
      "category": {
        "id": 1,
        "name": "プログラミング",
        "slug": "programming"
      },
      "tags": [
        {
          "id": 1,
          "name": "Go言語",
          "slug": "golang"
        }
      ],
      "created_at": "2025-06-14T10:00:00Z"
    }
  ],
  "total": 25,
  "page": 1,
  "limit": 5,
  "total_pages": 5
}
```

## アーキテクチャ実装計画

### 1. Domain Layer（内側の層）

#### 1.1 リポジトリインターフェース拡張
- `internal/domain/repository/article_repository.go`にメソッド追加

```go
type ArticleRepository interface {
    // ... 既存メソッド
    Search(ctx context.Context, params *ArticleSearchParams) (*ArticleSearchResult, error)
    SearchWithFilters(ctx context.Context, query string, filters map[string]interface{}) ([]*entity.Article, error)
    GetByDateRange(ctx context.Context, from, to time.Time) ([]*entity.Article, error)
    GetPopular(ctx context.Context, limit int) ([]*entity.Article, error)
    CountByFilters(ctx context.Context, filters map[string]interface{}) (int64, error)
}
```

### 2. Infrastructure Layer（外側の層）

#### 2.1 検索機能実装
- `internal/infrastructure/repository/article_repository_impl.go`
- GORMを使用した全文検索機能
- 動的クエリ構築によるフィルタリング
- ページネーション対応

#### 2.2 統合テスト
- `internal/infrastructure/repository/article_repository_impl_test.go`
- 検索機能の統合テスト
- 様々なフィルタ組み合わせのテスト

### 3. Application Layer（ユースケース層）

#### 3.1 SearchUsecase実装
- `internal/usecase/search_usecase.go`

```go
type SearchUsecase interface {
    SearchArticles(ctx context.Context, params *ArticleSearchParams) (*ArticleSearchResult, error)
    GetPopularArticles(ctx context.Context, limit int) ([]*entity.Article, error)
    GetRecentArticles(ctx context.Context, limit int) ([]*entity.Article, error)
    ValidateSearchParams(params *ArticleSearchParams) error
}
```

#### 3.2 ArticleUsecase拡張
- `internal/usecase/article_usecase.go`に検索関連メソッド追加
- GetAllArticlesメソッドをフィルタリング対応に拡張

#### 3.3 単体テスト
- `internal/usecase/search_usecase_test.go`
- カスタムモックを使用したテスト
- 検索パラメータ検証のテスト

### 4. Presentation Layer（プレゼンテーション層）

#### 4.1 SearchHandler実装
- `internal/presentation/handler/search_handler.go`
- 検索操作のHTTPハンドラー
- クエリパラメータの解析と検証

#### 4.2 ArticleHandler拡張
- `internal/presentation/handler/article_handler.go`
- GetAllArticlesをフィルタリング対応に更新
- 検索エンドポイント処理の追加

#### 4.3 HTTPテスト
- `internal/presentation/handler/search_handler_test.go`
- httptestを使用したHTTPテスト

#### 4.4 ルーティング設定
- `internal/presentation/router/router.go`
- 検索ルートの追加
- 記事ルートのクエリパラメータ対応

### 5. 依存関係の配線

#### 5.1 main.go更新
```go
// 検索ユースケース初期化
searchUsecase := usecase.NewSearchUsecase(articleRepo, categoryRepo, tagRepo)

// 検索ハンドラー初期化
searchHandler := handler.NewSearchHandler(searchUsecase)

// ルーター設定
router := router.SetupRouter(userHandler, articleHandler, commentHandler, categoryHandler, tagHandler, searchHandler)
```

## 実装順序

1. **Domain Layer**
   - ArticleRepository interfaceに検索メソッド追加
   - 検索パラメータ構造体定義

2. **Infrastructure Layer**
   - ArticleRepository実装に検索機能追加（全文検索、フィルタリング）
   - ページネーション対応
   - Database最適化（インデックス追加）

3. **Application Layer**
   - SearchUsecase実装（検索ビジネスロジック）
   - ArticleUsecase拡張（フィルタリング機能）

4. **Presentation Layer**
   - SearchHandler実装
   - ArticleHandler拡張（クエリパラメータ対応）
   - Router更新

5. **統合**
   - main.goで依存関係配線
   - 動作確認

6. **テスト作成**
   - 各レイヤーのテスト実装
   - 検索・フィルタリング機能の検証

## 検索機能詳細

### 全文検索
- 記事タイトル・本文の検索
- 引用符によるフレーズ検索対応
- 大文字小文字を区別しない検索
- 部分一致検索
- 関連度による検索結果ランキング

### フィルタリングオプション
- **カテゴリ**: 単一カテゴリスラッグによるフィルタ
- **タグ**: 複数タグによるフィルタ（AND/OR論理）
- **作成者**: 作成者IDまたは名前によるフィルタ
- **ステータス**: 公開状態によるフィルタ
- **日付範囲**: 作成日または更新日によるフィルタ
- **人気度**: 閲覧数やエンゲージメントに基づく

### ソートオプション
- **created_at**: 作成日（新しい順/古い順）
- **updated_at**: 更新日
- **title**: タイトルのアルファベット順
- **author**: 作成者名順
- **status**: 公開状態順
- **relevance**: 検索関連度スコア（検索クエリ時のデフォルト）

### ページネーション
- ページベースのページネーション（設定可能なページサイズ）
- デフォルト制限: 1ページあたり20記事
- 最大制限: 1ページあたり100記事
- レスポンスに総数とページ情報を含む

## データ検証ルール

### 入力検証
- Query: 任意、1文字以上100文字以下
- Page: 最小1、デフォルト1
- Limit: 最小1、最大100、デフォルト20
- SortBy: 有効なフィールド名である必要
- SortOrder: 'asc'または'desc'である必要
- DateFrom/DateTo: 有効な日付形式、DateFrom <= DateTo

## ビジネスロジック

### 検索動作
- 空クエリは全記事を返す（他のフィルタは適用）
- 未公開記事は作成者と管理者のみ閲覧可能
- 検索クエリ提供時は関連度でランキング
- 人気記事は閲覧数やエンゲージメント指標に基づく
- 最新記事はデフォルトで過去30日以内
- 分析・改善のための検索ログ記録

## パフォーマンス考慮事項

### データベース最適化
```sql
-- 全文検索インデックス
CREATE INDEX idx_articles_title_content ON articles USING gin(to_tsvector('japanese', title || ' ' || content));

-- フィルタ用インデックス
CREATE INDEX idx_articles_category_id ON articles(category_id);
CREATE INDEX idx_articles_author_id ON articles(author_id);
CREATE INDEX idx_articles_status ON articles(status);
CREATE INDEX idx_articles_created_at ON articles(created_at);
CREATE INDEX idx_articles_updated_at ON articles(updated_at);

-- よく使用されるフィルタ組み合わせ用複合インデックス
CREATE INDEX idx_articles_status_created_at ON articles(status, created_at DESC);
CREATE INDEX idx_articles_category_status ON articles(category_id, status);
```

### 検索パフォーマンス
- データベース固有の全文検索機能の使用
- 人気クエリの検索結果キャッシュ
- クエリ最適化のための検索分析機能

### キャッシュ戦略
- 頻繁に検索される結果のキャッシュ
- 人気記事ランキングのキャッシュ
- 検索統計情報のキャッシュ

この計画に従って段階的に実装を進めることで、高性能で使いやすい検索システムを構築できます。