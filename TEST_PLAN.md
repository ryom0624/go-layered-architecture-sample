# テスト実装計画

## 概要
Clean Architecture パターンに基づくGoアプリケーションの包括的なテスト戦略とその実装状況。
ユーザー管理機能と記事投稿機能（トランザクション対応）の完全なテストスイートを提供。

## 実装済みテスト

## ユーザー機能テスト

### ✅ ユースケース層テスト - User
**ファイル**: `internal/usecase/user_usecase_test.go`

**テスト対象**: ビジネスロジック層の検証
- `TestUserUsecase_CreateUser`
  - 正常なユーザー作成
  - 空の名前バリデーション
  - 空のメールバリデーション  
  - 重複メールバリデーション
- `TestUserUsecase_GetUser`
  - 存在するユーザーの取得
  - 存在しないユーザーの取得
- `TestUserUsecase_UpdateUser`
  - 正常なユーザー更新
  - ID=0での更新エラー
- `TestUserUsecase_DeleteUser` 
  - 正常なユーザー削除
  - ID=0での削除エラー
  - 存在しないユーザーの削除
- `TestUserUsecase_GetAllUsers`
  - 空のユーザーリスト取得
  - データ有りのユーザーリスト取得
- `TestUserUsecase_GetUserByEmail`
  - 存在するメールでの検索
  - 存在しないメールでの検索

**モック実装**: `MockUserRepository` - カスタム実装

### ✅ ハンドラー層テスト - User
**ファイル**: `internal/presentation/handler/user_handler_test.go`

**テスト対象**: HTTP エンドポイントの検証
- `TestUserHandler_CreateUser`
  - 正常なユーザー作成 (201 Created)
  - 無効なJSONリクエスト (400 Bad Request)
  - 必須フィールド不足 (400 Bad Request)
  - 重複メール (500 Internal Server Error)
- `TestUserHandler_GetUser`
  - 存在するユーザーの取得 (200 OK)
  - 存在しないユーザー (404 Not Found)
  - 無効なユーザーID (400 Bad Request)
- `TestUserHandler_GetAllUsers`
  - 空のユーザーリスト (200 OK)
  - データ有りのユーザーリスト (200 OK)
- `TestUserHandler_UpdateUser`
  - 正常なユーザー更新 (200 OK)
  - 部分更新 (200 OK)
  - 存在しないユーザー (404 Not Found)
  - 無効なユーザーID (400 Bad Request)
- `TestUserHandler_DeleteUser`
  - 正常なユーザー削除 (204 No Content)
  - 存在しないユーザー (500 Internal Server Error)
  - 無効なユーザーID (400 Bad Request)

**モック実装**: `MockUserUsecase` - カスタム実装
**テストルーター**: Gin テストモードでのHTTPテスト

### ✅ リポジトリ層テスト - User
**ファイル**: `internal/infrastructure/repository/user_repository_impl_test.go`

**テスト対象**: データベース操作層の検証
- `TestUserRepositoryImpl_Create`
  - 正常なユーザー作成
  - データベースエラー時の作成失敗
- `TestUserRepositoryImpl_GetByID`
  - 正常なID検索
  - 存在しないユーザーの検索
- `TestUserRepositoryImpl_GetByEmail`
  - 正常なメール検索
  - 存在しないメールの検索
- `TestUserRepositoryImpl_GetAll`
  - 全ユーザー取得（データあり）
  - 空のユーザーリスト取得
- `TestUserRepositoryImpl_Update`
  - 正常なユーザー更新
  - データベースエラー時の更新失敗
- `TestUserRepositoryImpl_Delete`
  - 正常なユーザー削除
  - データベースエラー時の削除失敗
  - 存在しないユーザーの削除

**使用技術**: `github.com/DATA-DOG/go-sqlmock`, GORM

## 記事機能テスト

### ✅ ユースケース層テスト - Article
**ファイル**: `internal/usecase/article_usecase_test.go`

**テスト対象**: 記事ビジネスロジック層の検証（トランザクション対応）
- `TestArticleUsecase_CreateArticle`
  - 正常な記事作成
  - 空のタイトルバリデーション
  - 空のコンテンツバリデーション
  - AuthorID=0バリデーション
  - 存在しない著者バリデーション
- `TestArticleUsecase_GetArticle`
  - 存在する記事の取得
  - 存在しない記事の取得
- `TestArticleUsecase_UpdateArticle`
  - 正常な記事更新
  - 部分更新（空フィールド処理）
  - 存在しない記事の更新
  - ID=0での更新エラー
- `TestArticleUsecase_PublishArticle`
  - 正常な記事公開（トランザクション使用）
  - すでに公開済み記事の処理
  - 存在しない記事の公開
- `TestArticleUsecase_UnpublishArticle`
  - 正常な記事非公開（トランザクション使用）
  - すでに非公開記事の処理
- `TestArticleUsecase_GetPublishedArticles`
  - 公開記事のみ取得
  - ステータスフィルタリングの検証
- `TestArticleUsecase_DeleteArticle`
  - 正常な記事削除
  - ID=0での削除エラー
  - 存在しない記事の削除

**モック実装**: `MockArticleRepository`, `MockUserRepository`, `MockTransactionManager`

### ✅ ハンドラー層テスト - Article
**ファイル**: `internal/presentation/handler/article_handler_test.go`

**テスト対象**: 記事HTTP エンドポイントの検証
- `TestArticleHandler_CreateArticle`
  - 正常な記事作成 (201 Created)
  - 無効なJSONリクエスト (400 Bad Request)
  - 必須フィールド不足 (400 Bad Request)
- `TestArticleHandler_GetArticle`
  - 存在する記事の取得 (200 OK)
  - 存在しない記事 (404 Not Found)
  - 無効な記事ID (400 Bad Request)
- `TestArticleHandler_GetAllArticles`
  - 空の記事リスト (200 OK)
  - データ有りの記事リスト (200 OK)
- `TestArticleHandler_GetPublishedArticles`
  - 公開記事のみ取得 (200 OK)
  - ステータスフィルタリング検証
- `TestArticleHandler_UpdateArticle`
  - 正常な記事更新 (200 OK)
  - 存在しない記事 (500 Internal Server Error)
  - 無効な記事ID (400 Bad Request)
- `TestArticleHandler_PublishArticle`
  - 正常な記事公開 (200 OK)
  - 存在しない記事 (500 Internal Server Error)
- `TestArticleHandler_UnpublishArticle`
  - 正常な記事非公開 (200 OK)
- `TestArticleHandler_DeleteArticle`
  - 正常な記事削除 (204 No Content)
  - 存在しない記事 (500 Internal Server Error)
  - 無効な記事ID (400 Bad Request)

**モック実装**: `MockArticleUsecase` - カスタム実装
**テストルーター**: Gin テストモードでのHTTPテスト

### ✅ リポジトリ層テスト - Article
**ファイル**: `internal/infrastructure/repository/article_repository_impl_test.go`

**テスト対象**: 記事データベース操作層の検証
- `TestArticleRepositoryImpl_Create`
  - 正常な記事作成
  - データベースエラー時の作成失敗
- `TestArticleRepositoryImpl_GetByID`
  - 正常なID検索（Author Preload含む）
  - 存在しない記事の検索
- `TestArticleRepositoryImpl_GetAll`
  - 全記事取得（Author Preload含む）
  - 空の記事リスト取得
- `TestArticleRepositoryImpl_GetByStatus`
  - ステータス別記事取得
  - 公開/非公開フィルタリング
- `TestArticleRepositoryImpl_Update`
  - 正常な記事更新
  - データベースエラー時の更新失敗
- `TestArticleRepositoryImpl_Delete`
  - 正常な記事削除
  - データベースエラー時の削除失敗
  - 存在しない記事の削除
- `TestArticleRepositoryImpl_CreateWithTx`
  - トランザクション内での記事作成
- `TestArticleRepositoryImpl_UpdateWithTx`
  - トランザクション内での記事更新
- `TestArticleRepositoryImpl_DeleteWithTx`
  - トランザクション内での記事削除

**使用技術**: `github.com/DATA-DOG/go-sqlmock`, GORM, カスタムTransaction Mock

## 共通機能テスト

### ✅ 設定パッケージテスト
**ファイル**: `pkg/config/config_test.go`

**テスト対象**: 設定管理の検証
- `TestLoadConfig`
  - デフォルト値での設定読み込み
  - 環境変数での設定読み込み
  - 部分的な環境変数での設定読み込み
- `TestGetEnv`
  - 存在する環境変数の取得
  - 存在しない環境変数のデフォルト値取得
  - 空の環境変数のデフォルト値取得
- `TestGetEnvInt`
  - 有効な整数環境変数の取得
  - 存在しない整数環境変数のデフォルト値取得
  - 無効な整数環境変数のデフォルト値取得
  - 空の整数環境変数のデフォルト値取得
  - ゼロ値の整数環境変数の取得
  - 負の整数環境変数の取得

## トランザクション機能テスト

### ✅ 実装済みトランザクションテスト
- **Transaction Manager**: `MockTransactionManager` による単体テスト
- **Repository Transaction Methods**: WithTx メソッドのテスト
- **UseCase Transaction Integration**: 複数操作のトランザクション検証
- **Error Rollback**: エラー時の自動ロールバック検証

### トランザクション対応機能
- 記事公開/非公開処理
- 複数テーブル操作の整合性保証
- エラー発生時のデータ整合性維持

## テスト統計

### 実装済みテストファイル数: 7ファイル
#### ユーザー機能: 3ファイル
- UseCase テスト
- Handler テスト  
- Repository テスト

#### 記事機能: 3ファイル
- UseCase テスト（トランザクション含む）
- Handler テスト
- Repository テスト（トランザクション含む）

#### 共通機能: 1ファイル
- Config テスト

### テストケース総数: 約80+テストケース
- 正常系テスト
- 異常系テスト
- バリデーションテスト
- トランザクションテスト

## 未実装テスト（将来拡張）

## テスト実行方法

### Dockerコンテナ内での実行
```bash
# 全テスト実行
docker compose exec app go test ./... -v

# 特定パッケージのテスト
docker compose exec app go test ./internal/usecase/... -v
docker compose exec app go test ./internal/presentation/handler/... -v

# カバレッジ付きテスト
docker compose exec app go test ./... -cover

# 短縮版テスト
docker compose exec app go test ./...
```

### ローカル実行
```bash
# 全テスト実行
go test ./... -v

# 特定テスト実行
go test ./internal/usecase -v
go test ./internal/presentation/handler -v
```

## テスト設計原則

### 使用技術
- **標準testingパッケージ**: `assert`ライブラリ不使用
- **カスタムモック**: 各層に応じた専用モック実装
- **httptest**: HTTP エンドポイントテスト
- **gin.TestMode**: Ginフレームワークのテストモード

### Clean Architecture 準拠
- **依存関係の分離**: 各層のテストで外部依存をモック化
- **単一責任**: 各テストは特定の責務のみを検証
- **境界の明確化**: レイヤー間の境界を意識したテスト設計

### テストカバレッジ
- **正常系**: 期待される動作の検証
- **異常系**: エラーハンドリングの検証
- **境界値**: エッジケースの検証
- **バリデーション**: 入力値検証の確認

## 実装されたテストファイル一覧

### ユーザー機能テスト
1. `internal/usecase/user_usecase_test.go` - ユーザービジネスロジックテスト
2. `internal/presentation/handler/user_handler_test.go` - ユーザーHTTPハンドラーテスト  
3. `internal/infrastructure/repository/user_repository_impl_test.go` - ユーザーリポジトリテスト

### 記事機能テスト
4. `internal/usecase/article_usecase_test.go` - 記事ビジネスロジックテスト（トランザクション含む）
5. `internal/presentation/handler/article_handler_test.go` - 記事HTTPハンドラーテスト
6. `internal/infrastructure/repository/article_repository_impl_test.go` - 記事リポジトリテスト（トランザクション含む）

### 共通機能テスト
7. `pkg/config/config_test.go` - 設定管理テスト

## 今後の拡張予定

1. **統合テスト** - レイヤー間の連携テスト
2. **E2Eテスト** - アプリケーション全体のテスト
3. **パフォーマンステスト** - 負荷テスト
4. **セキュリティテスト** - 脆弱性テスト
5. **認証・認可機能テスト** - JWT認証等のセキュリティテスト

## 注意事項

- 現在のテストは外部データベースに依存しない設計
- モックによる完全な依存関係の分離を実現
- 実際のHTTPリクエスト/レスポンスを検証
- Clean Architecture の各層の責務を適切に分離してテスト