# Claude Code開発フロー指示書

## 概要
この指示書は、Claude Codeに開発タスクを依頼する際の標準的なフローを定義しています。Clean Architectureの原則に従い、体系的で品質の高い開発を実現するためのガイドラインです。

## 基本指示フォーマット

```
[タスク名]をDevelopment Workflowに従ってください。
```

**例:**
- `Task8をDevelopment Workflowに従ってください。`
- `認証システムの実装をDevelopment Workflowに従ってください。`
- `APIレート制限機能の追加をDevelopment Workflowに従ってください。`

## Development Workflow詳細

### 1. 計画・設計フェーズ
Claude Codeは以下を自動的に実行します：

- **既存コードベースの調査**
  - `tasks/`ディレクトリの実装計画書確認
  - 既存アーキテクチャパターンの分析
  - 依存関係の理解

- **タスク管理**
  - TodoWriteツールで詳細なタスク分解
  - 優先度とステータスの設定
  - 進捗の可視化

### 2. ブランチ戦略
```bash
git checkout -b feature/[task-name]
```

### 3. 実装順序（Clean Architecture準拠）
**依存関係ルール**: 内側の層は外側の層に依存しない

#### 段階的実装:
1. **Domain Layer（最内側）**
   - `internal/domain/entity/` - エンティティ定義
   - `internal/domain/repository/` - リポジトリインターフェース

2. **Infrastructure Layer（外側）**
   - `internal/infrastructure/repository/` - リポジトリ実装
   - `internal/infrastructure/database/` - DB設定更新

3. **Application Layer（ユースケース）**
   - `internal/usecase/` - ビジネスロジック実装

4. **Presentation Layer（最外側）**
   - `internal/presentation/handler/` - HTTPハンドラー
   - `internal/presentation/middleware/` - ミドルウェア
   - `internal/presentation/router/` - ルーティング設定

5. **Dependency Injection**
   - `cmd/main.go` - 依存関係の配線

### 4. テスト実装
- **カスタムモック使用** - 外部ライブラリに依存しない
- **各レイヤーでのテスト** - 単体テスト・統合テスト
- **テスト実行確認** - `go test ./...`で全テスト成功

### 5. ドキュメント更新
- **API仕様** - `CLAUDE.md`のAPIエンドポイント追加
- **実装パターン** - アーキテクチャ説明の更新
- **機能説明** - READMEの機能リスト更新

### 6. マージ戦略
```bash
git checkout task-documentation
git merge feature/[task-name]
```

## Claude Codeの自動実行プロセス

### Phase 1: Analysis & Planning
- [ ] 既存コードベース分析
- [ ] タスク計画書確認（`tasks/`ディレクトリ）
- [ ] TodoWriteでタスク分解・管理開始
- [ ] フィーチャーブランチ作成

### Phase 2: Clean Architecture Implementation
- [ ] Domain Layer実装（エンティティ・インターフェース）
- [ ] Infrastructure Layer実装（リポジトリ・DB）
- [ ] Application Layer実装（ユースケース）
- [ ] Presentation Layer実装（ハンドラー・ルーター）
- [ ] Dependency Injection（main.go配線）

### Phase 3: Quality Assurance
- [ ] 包括的テスト実装
- [ ] 全テスト実行・成功確認
- [ ] ビルド確認
- [ ] 既存テスト修正（必要に応じて）

### Phase 4: Documentation & Integration
- [ ] CLAUDE.md更新（API・機能説明）
- [ ] コミット作成（詳細なコミットメッセージ）
- [ ] task-documentationブランチへのマージ

## 品質基準

### コーディング規約
- **Clean Architecture準拠** - 依存関係ルールの厳守
- **エラーハンドリング** - 適切なバリデーションとエラー処理
- **トランザクション管理** - データ整合性確保
- **命名規約** - 一貫したインターフェース・実装命名

### テスト要件
- **カバレッジ** - 各レイヤーでの包括的テスト
- **モック実装** - 外部依存なしのカスタムモック
- **成功基準** - `go test ./...`で全テスト通過

### ドキュメント要件
- **API仕様** - 新しいエンドポイントの完全な文書化
- **実装パターン** - アーキテクチャ説明の更新
- **使用例** - 実用的な使用方法の説明

## 使用例

### 基本的な機能追加
```
認証機能をDevelopment Workflowに従ってください。
```

### 複雑なシステム実装
```
Task9: リアルタイム通知システムをDevelopment Workflowに従ってください。
```

### 既存機能の拡張
```
検索機能にフィルタリングオプションを追加する実装をDevelopment Workflowに従ってください。
```

## 期待される成果物

### 技術成果物
1. **完全に動作するフィーチャー** - テスト済み・本番準備完了
2. **Clean Architectureコンプライアンス** - 適切な層分離
3. **包括的テストスイート** - 高品質なテストカバレッジ
4. **完全なドキュメント** - API仕様・使用方法

### プロジェクト成果物
1. **詳細なコミット履歴** - 実装プロセスの透明性
2. **更新されたドキュメント** - 最新の機能・API仕様
3. **統合済みフィーチャー** - メインブランチへの安全なマージ

## トラブルシューティング

### よくある問題と解決方法
- **テスト失敗** → Claude Codeが自動的に修正・再実行
- **依存関係エラー** → `go mod tidy`実行
- **ビルドエラー** → 段階的デバッグ・修正
- **マージコンフリクト** → 適切なブランチ戦略で回避

## まとめ

この開発フローにより以下が保証されます：

- ✅ **品質の一貫性** - Clean Architectureの徹底
- ✅ **開発効率** - 体系的な実装プロセス
- ✅ **保守性** - 適切なテスト・ドキュメント
- ✅ **拡張性** - 将来の機能追加への対応

このフローを使用することで、Claude Codeは常に高品質で保守性の高いコードを提供し、プロジェクトの長期的な成功を支援します。