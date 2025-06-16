# Task: CORS設定実装

## 概要
SPAに対応したCORS（Cross-Origin Resource Sharing）設定を実装し、フロントエンド開発に必要なHTTPヘッダー設定を行う。

## 要件

### 機能要件
- SPAからのAPIアクセスを許可するCORS設定
- 開発環境と本番環境での柔軟なOrigin設定
- 認証情報（JWT）を含むリクエストの対応
- プリフライト（OPTIONS）リクエストの適切な処理

### 技術要件
- Ginミドルウェアとしての実装
- 環境変数による設定管理
- セキュリティを考慮したOrigin制限

## 実装計画

### 1. CORSミドルウェア実装
- `internal/presentation/middleware/cors.go`
- 環境変数ベースの設定読み込み
- 適切なヘッダー設定

### 2. 設定管理
- `pkg/config/config.go` - CORS設定構造体追加
- 環境変数による設定項目管理

### 3. ルーター統合
- `internal/presentation/router/router.go` - ミドルウェア適用

### 4. 環境設定
- `.env.example` - 設定サンプル追加

## 設定項目

### 環境変数
- `CORS_ALLOWED_ORIGINS` - 許可するオリジン（カンマ区切り、デフォルト: "*"）
- `CORS_ALLOWED_METHODS` - 許可するHTTPメソッド（デフォルト: "GET,POST,PUT,DELETE,OPTIONS"）
- `CORS_ALLOWED_HEADERS` - 許可するヘッダー（デフォルト: "Content-Type,Authorization,X-Requested-With"）
- `CORS_ALLOW_CREDENTIALS` - 認証情報の送信許可（デフォルト: true）

### セキュリティ考慮事項
- 本番環境では適切なOrigin制限
- 認証情報送信時のOrigin厳格化
- 不要なヘッダーの露出防止