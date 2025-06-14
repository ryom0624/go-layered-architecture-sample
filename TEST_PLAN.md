# テスト実装計画

## 概要
Clean Architecture パターンに基づくGoアプリケーションの包括的なテスト戦略とその実装状況。

## 実装済みテスト

### ✅ ユースケース層テスト
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

### ✅ ハンドラー層テスト
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

### ✅ リポジトリ層テスト
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

1. `internal/usecase/user_usecase_test.go` - ビジネスロジックテスト
2. `internal/presentation/handler/user_handler_test.go` - HTTPハンドラーテスト  
3. `internal/infrastructure/repository/user_repository_impl_test.go` - データベースリポジトリテスト
4. `pkg/config/config_test.go` - 設定管理テスト

## 今後の拡張予定

1. **リポジトリ層テスト** - データベース操作の検証
2. **統合テスト** - レイヤー間の連携テスト
3. **E2Eテスト** - アプリケーション全体のテスト
4. **パフォーマンステスト** - 負荷テスト
5. **セキュリティテスト** - 脆弱性テスト

## 注意事項

- 現在のテストは外部データベースに依存しない設計
- モックによる完全な依存関係の分離を実現
- 実際のHTTPリクエスト/レスポンスを検証
- Clean Architecture の各層の責務を適切に分離してテスト