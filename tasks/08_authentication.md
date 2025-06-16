# Task 8: 認証システム (Authentication System)

## 概要
JWT（JSON Web Token）を用いたユーザー認証システムを実装します。ユーザー登録、ログイン、トークン管理、認証ミドルウェアを含む包括的な認証機能を提供します。

## 技術要件

### セキュリティ要件
- **パスワードハッシュ化**: bcryptによる安全なパスワード保存
- **JWT トークン**: アクセストークンによる認証
- **リフレッシュトークン**: セキュアなトークン更新機能
- **トークン有効期限**: 設定可能なアクセス・リフレッシュトークン期限

### 機能要件
1. **ユーザー登録** (Register)
   - 名前、メールアドレス、パスワードによる新規ユーザー作成
   - メールアドレス重複チェック
   - パスワード強度バリデーション（8文字以上）

2. **ログイン** (Login)
   - メールアドレスとパスワードによる認証
   - JWTアクセストークンとリフレッシュトークンの発行

3. **トークン更新** (Refresh Token)
   - リフレッシュトークンによるアクセストークン更新
   - 古いリフレッシュトークンの無効化

4. **ログアウト** (Logout)
   - 単一デバイスからのログアウト（リフレッシュトークン削除）
   - 全デバイスからのログアウト（全リフレッシュトークン削除）

5. **認証ミドルウェア**
   - 必須認証ミドルウェア（AuthMiddleware）
   - オプション認証ミドルウェア（OptionalAuthMiddleware）

## データベース設計

### Userエンティティ拡張
```go
type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name" gorm:"not null"`
    Email     string    `json:"email" gorm:"uniqueIndex;not null"`
    Password  string    `json:"-" gorm:"not null"`  // 新規追加
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### RefreshTokenエンティティ
```go
type RefreshToken struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    UserID    uint      `json:"user_id" gorm:"not null;index"`
    Token     string    `json:"token" gorm:"not null;uniqueIndex"`
    ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
    CreatedAt time.Time `json:"created_at"`
    
    User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
```

## API設計

### エンドポイント
- `POST /api/v1/auth/register` - ユーザー登録
- `POST /api/v1/auth/login` - ログイン
- `POST /api/v1/auth/refresh` - トークン更新
- `POST /api/v1/auth/logout` - ログアウト
- `POST /api/v1/auth/logout-all` - 全デバイスログアウト（認証必須）

### リクエスト/レスポンス例

#### 登録/ログイン成功レスポンス
```json
{
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### 認証ヘッダー形式
```
Authorization: Bearer <access_token>
```

## 実装順序（Clean Architecture準拠）

### 1. Domain Layer
- [x] Userエンティティにパスワードフィールド追加
- [x] 認証関連エンティティ定義（auth.go）
- [x] AuthRepositoryインターフェース定義

### 2. Infrastructure Layer
- [x] AuthRepositoryImpl実装（パスワードハッシュ化含む）
- [x] JWTユーティリティ実装（pkg/jwt/jwt.go）
- [x] 設定ファイル拡張（認証設定追加）

### 3. Application Layer
- [x] AuthUsecase実装（ビジネスロジック）
- [x] バリデーションとエラーハンドリング

### 4. Presentation Layer
- [x] AuthHandler実装（HTTPハンドラー）
- [x] 認証ミドルウェア実装
- [x] ルーティング設定

### 5. Dependency Injection
- [ ] main.goでの依存関係配線
- [ ] データベースマイグレーション更新

## セキュリティ考慮事項

### トークン管理
- **アクセストークン**: 短期間（デフォルト15分）
- **リフレッシュトークン**: 長期間（デフォルト7日）
- **トークンローテーション**: リフレッシュ時に新しいリフレッシュトークン発行

### パスワードセキュリティ
- bcryptによるハッシュ化（コスト値: DefaultCost）
- パスワード平文の応答除外（JSONタグ: `json:"-"`）

### 環境変数設定
```env
JWT_SECRET=your-super-secret-jwt-key-here
ACCESS_TOKEN_DURATION_MINUTES=15
REFRESH_TOKEN_DURATION_DAYS=7
```

## テスト要件

### ユニットテスト
- [ ] AuthUsecase: 全メソッドのテスト
- [ ] AuthRepository: データベース操作テスト
- [ ] JWT Manager: トークン生成・検証テスト

### 統合テスト
- [ ] AuthHandler: HTTPリクエスト・レスポンステスト
- [ ] 認証ミドルウェア: トークン検証テスト

### セキュリティテスト
- [ ] 無効なトークンでのアクセス拒否
- [ ] 期限切れトークンの処理
- [ ] 不正なリフレッシュトークンの処理

## 実装後の確認項目

### 機能確認
- [ ] ユーザー登録が正常に動作する
- [ ] ログインでトークンが発行される
- [ ] リフレッシュトークンでアクセストークンが更新される
- [ ] ログアウトでトークンが無効化される
- [ ] 認証ミドルウェアが正常に動作する

### セキュリティ確認
- [ ] パスワードがハッシュ化されて保存される
- [ ] 無効なトークンでアクセスが拒否される
- [ ] トークンの有効期限が適切に管理される
- [ ] リフレッシュトークンの適切なローテーション

## 関連ファイル

### 新規作成
- `internal/domain/entity/auth.go`
- `internal/domain/repository/auth_repository.go`
- `internal/infrastructure/repository/auth_repository_impl.go`
- `internal/usecase/auth_usecase.go`
- `internal/presentation/handler/auth_handler.go`
- `internal/presentation/middleware/auth.go`
- `pkg/jwt/jwt.go`

### 変更
- `internal/domain/entity/user.go` (パスワードフィールド追加)
- `pkg/config/config.go` (認証設定追加)
- `internal/infrastructure/database/connection.go` (マイグレーション更新)
- `internal/presentation/router/router.go` (認証エンドポイント追加)
- `cmd/main.go` (依存関係配線)