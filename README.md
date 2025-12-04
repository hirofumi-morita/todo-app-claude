# TODO管理ウェブアプリケーション

フルスタックのTODO管理アプリケーションです。Next.js、Go、Gin、Hasura、PostgreSQLを使用して構築されています。

## 技術スタック

### フロントエンド
- **Next.js 14** - Reactフレームワーク
- **TypeScript** - 型安全性
- **Apollo Client** - GraphQLクライアント
- **GraphQL** - APIクエリ言語

### バックエンド
- **Go 1.21** - カスタムAPIサーバー（認証用）
- **Gin** - Webアプリケーションフレームワーク
- **Hasura** - GraphQLエンジン
- **PostgreSQL** - データベース
- **JWT** - 認証

### インフラ
- **Docker & Docker Compose** - コンテナ化

## 主な機能

### ユーザー機能
- ✅ ユーザー登録（メールアドレス、パスワード）
- ✅ ログイン/ログアウト
- ✅ JWT認証

### TODO機能
- ✅ TODOの作成
- ✅ TODOの表示
- ✅ TODOの更新（完了/未完了の切り替え）
- ✅ TODOの削除

### 管理者機能
- ✅ ユーザー一覧表示
- ✅ ユーザーの削除
- ✅ 管理者権限の付与/解除
- ✅ ユーザーのTODO表示

## 前提条件

以下のツールがインストールされている必要があります：

- Docker & Docker Compose
- Node.js 18+ (ローカル開発の場合)
- Go 1.21+ (ローカル開発の場合)

## セットアップ手順

### 1. リポジトリのクローン

```bash
git clone <repository-url>
cd todo-app-claude
```

### 2. Dockerを使用した起動（推奨）

```bash
# すべてのサービスを起動
docker-compose up -d

# ログを確認
docker-compose logs -f
```

以下のサービスが起動します：
- フロントエンド: http://localhost:3000
- Hasura GraphQL API: http://localhost:8080/v1/graphql
- Hasura Console: http://localhost:8080/console
- カスタムバックエンドAPI（認証用）: http://localhost:8081
- PostgreSQL: localhost:5432

### 3. ローカル開発環境のセットアップ

#### データベースの起動

```bash
docker-compose up -d postgres
```

#### バックエンドの起動

```bash
cd backend

# 依存関係のインストール
go mod download

# サーバーの起動
go run cmd/api/main.go
```

#### フロントエンドの起動

```bash
cd frontend

# 依存関係のインストール
npm install

# 開発サーバーの起動
npm run dev
```

## デフォルトの管理者アカウント

初回起動時に自動的に作成されます：

- **メールアドレス**: admin@example.com
- **パスワード**: admin123

⚠️ **重要**: 本番環境では必ずこのアカウントを削除または変更してください。

## 使用方法

### 1. 新規ユーザー登録

1. http://localhost:3000 にアクセス
2. 「新規登録」をクリック
3. メールアドレスとパスワードを入力
4. 登録ボタンをクリック

### 2. ログイン

1. メールアドレスとパスワードを入力
2. ログインボタンをクリック

### 3. TODOの管理

- TODOを追加: タイトルと説明を入力して「追加」ボタンをクリック
- TODOを完了: チェックボックスをクリック
- TODOを削除: 「削除」ボタンをクリック

### 4. 管理者機能（管理者のみ）

1. ログイン後、「管理者ページ」をクリック
2. ユーザー一覧を表示
3. ユーザーの削除や管理者権限の付与が可能

## API エンドポイント

### GraphQL API（Hasura）

すべてのTODO操作はHasura GraphQL APIを通じて実行されます：

**エンドポイント**: `http://localhost:8080/v1/graphql`

主なクエリとミューテーション：
- **Queries**: `todos`, `todos_by_pk`, `users`, `users_by_pk`
- **Mutations**: `insert_todos_one`, `update_todos_by_pk`, `delete_todos_by_pk`

GraphQL APIは自動的にJWTトークンを検証し、ユーザーごとのアクセス制御を行います。

### REST API（認証用）

認証機能はカスタムGoバックエンドで提供されます：

```
POST   /api/register          - ユーザー登録
POST   /api/login             - ログイン
GET    /api/me                - 現在のユーザー情報取得（要認証）
```

### 管理者機能

管理者機能もカスタムバックエンドで提供されます：

```
GET    /api/admin/users           - ユーザー一覧取得（要管理者権限）
GET    /api/admin/users/:id       - ユーザー詳細取得（要管理者権限）
DELETE /api/admin/users/:id       - ユーザー削除（要管理者権限）
PUT    /api/admin/users/:id/role  - 管理者権限の変更（要管理者権限）
GET    /api/admin/users/:id/todos - ユーザーのTODO取得（要管理者権限）
```

## プロジェクト構成

```
.
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go              # エントリーポイント
│   ├── internal/
│   │   ├── database/
│   │   │   └── database.go          # DB接続とマイグレーション
│   │   ├── handlers/
│   │   │   ├── auth.go              # 認証ハンドラー（Gin）
│   │   │   ├── todo.go              # TODOハンドラー（Gin）
│   │   │   └── admin.go             # 管理者ハンドラー（Gin）
│   │   ├── middleware/
│   │   │   ├── auth.go              # 認証ミドルウェア（Gin）
│   │   │   └── password.go          # パスワードハッシュ
│   │   └── models/
│   │       └── user.go              # データモデル
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── src/
│   │   ├── app/
│   │   │   ├── admin/
│   │   │   │   └── page.tsx         # 管理者ページ
│   │   │   ├── login/
│   │   │   │   └── page.tsx         # ログインページ
│   │   │   ├── register/
│   │   │   │   └── page.tsx         # 登録ページ
│   │   │   ├── todos/
│   │   │   │   └── page.tsx         # TODOページ
│   │   │   ├── layout.tsx           # レイアウト
│   │   │   ├── providers.tsx        # Apollo Provider
│   │   │   ├── page.tsx             # ホームページ
│   │   │   └── globals.css          # グローバルスタイル
│   │   ├── lib/
│   │   │   ├── api.ts               # REST APIクライアント
│   │   │   ├── apollo-client.ts     # Apollo Clientセットアップ
│   │   │   ├── graphql.ts           # GraphQLクエリ/ミューテーション
│   │   │   └── auth.ts              # 認証ユーティリティ
│   │   └── types/
│   │       └── index.ts             # TypeScript型定義
│   ├── Dockerfile
│   ├── package.json
│   ├── tsconfig.json
│   └── next.config.js
├── hasura/
│   └── metadata/
│       ├── databases.yaml           # Hasuraデータベース設定
│       ├── actions.yaml             # Hasuraアクション定義
│       └── version.yaml             # Hasuraメタデータバージョン
├── docker-compose.yml
├── .env.example
├── .gitignore
└── README.md
```

## データベーススキーマ

### users テーブル

| カラム名   | 型        | 説明              |
|-----------|-----------|-------------------|
| id        | SERIAL    | ユーザーID (主キー) |
| email     | VARCHAR   | メールアドレス      |
| password  | VARCHAR   | ハッシュ化パスワード |
| is_admin  | BOOLEAN   | 管理者フラグ        |
| created_at| TIMESTAMP | 作成日時           |
| updated_at| TIMESTAMP | 更新日時           |

### todos テーブル

| カラム名    | 型        | 説明              |
|------------|-----------|-------------------|
| id         | SERIAL    | TODO ID (主キー)   |
| user_id    | INTEGER   | ユーザーID (外部キー)|
| title      | VARCHAR   | タイトル           |
| description| TEXT      | 説明              |
| completed  | BOOLEAN   | 完了フラグ         |
| created_at | TIMESTAMP | 作成日時           |
| updated_at | TIMESTAMP | 更新日時           |

## 環境変数

`.env.example`をコピーして`.env`を作成し、必要に応じて値を変更してください。

```bash
cp .env.example .env
```

## トラブルシューティング

### ポートが既に使用されている

別のサービスがポート3000、8080、5432を使用している場合、docker-compose.ymlのポート設定を変更してください。

### データベース接続エラー

1. PostgreSQLコンテナが起動しているか確認:
   ```bash
   docker-compose ps
   ```

2. データベースのログを確認:
   ```bash
   docker-compose logs postgres
   ```

### フロントエンドからバックエンドに接続できない

`frontend/src/lib/api.ts`のAPI_URLが正しく設定されているか確認してください。

## 開発のヒント

### バックエンドのテスト

```bash
cd backend
go test ./...
```

### フロントエンドのビルド

```bash
cd frontend
npm run build
```

### データベースのリセット

```bash
docker-compose down -v
docker-compose up -d
```

## アーキテクチャの特徴

このアプリケーションは、HasuraとカスタムGoバックエンドを組み合わせたハイブリッドアーキテクチャを採用しています：

### Hasuraの役割
- **GraphQL API提供**: TODO操作のCRUD機能を自動生成
- **リアルタイム機能**: GraphQL Subscriptionsによるリアルタイム更新
- **権限管理**: JWTベースの行レベルセキュリティ
- **パフォーマンス**: 最適化されたクエリとキャッシング

### カスタムGoバックエンド（Gin）の役割
- **認証**: ユーザー登録・ログイン・JWTトークン発行
- **パスワード管理**: bcryptによる安全なハッシュ化
- **管理者機能**: 複雑なビジネスロジックの実装
- **カスタムエンドポイント**: Hasuraで対応できない特殊な処理

この構成により、開発速度とカスタマイズ性の両立を実現しています。

## セキュリティに関する注意

- 本番環境では、以下のシークレットキーを必ず変更してください：
  - JWTシークレットキー (`backend/internal/middleware/auth.go`)
  - Hasura Admin Secret (`docker-compose.yml`)
  - Hasura JWT Secret (`docker-compose.yml`)
- HTTPSを使用してください
- デフォルトの管理者アカウントを削除または変更してください
- 環境変数を`.env`ファイルで管理し、`.gitignore`に追加してください
- Hasura Consoleへのアクセスを本番環境では無効化してください

## ライセンス

MIT License

## 貢献

プルリクエストを歓迎します。大きな変更の場合は、まずissueを開いて変更内容を議論してください。
