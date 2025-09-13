# Shittaka Backend

Go言語で構築されたSupabase認証システムのバックエンドAPIです。

## 機能

- ユーザー登録（サインアップ）
- ユーザーログイン
- ユーザーログアウト
- Supabaseとの疎通テスト
- CORS対応
- ジャンル作成
- 問題作成
- 問題編集
- 問題削除

## セットアップ

### 1. 依存関係のインストール

```bash
go mod tidy
```

### 2. 環境変数の設定

`env.example`を参考に`.env`ファイルを作成し、Supabaseの設定を行ってください：

```bash
cp env.example .env
```

`.env`ファイルに以下の情報を設定：

```env
# Supabase設定
SUPABASE_URL=https://your-project-id.supabase.co
SUPABASE_SERVICE_ROLE_KEY=your-service-role-key-here
SUPABASE_ANON_KEY=your-anon-key-here
SUPABASE_JWT_KEY=your-anon-jwt-here

# サーバー設定
PORT=8088
APP_ENV=developmenL

# 開発環境用の設定
GIN_MODE=debug
```

### 3. Supabaseプロジェクトの設定

1. [Supabase](https://supabase.com)でプロジェクトを作成
2. プロジェクト設定から以下を取得：
   - Project URL (`SUPABASE_URL`)
   - Service Role Key (`SUPABASE_SERVICE_ROLE_KEY`)



### 5.接続テスト

""の中をenvのものに置き換えてください

```powershell
# powershell
$env:SUPABASE_URL = ".env参照"
$env:SUPABASE_SERVICE_ROLE_KEY = ".env参照"
$env:SUPABASE_ANON_KEY=".env参照"
$env:PORT = "8088"

# サーバーを起動
./server.exe
```

```bash
# bash
export SUPABASE_URL=".env参照"
export SUPABASE_SERVICE_ROLE_KEY=".env参照"
export SUPABASE_ANON_KEY=".env参照"
export  PORT="8088"
```
```
# 実行
go run cmd/server/main.go
```


```powershell
# powershell
# ヘルスチェック
Invoke-WebRequest -Uri "http://localhost:8088/health" -Method GET

# Supabase接続テスト
Invoke-WebRequest -Uri "http://localhost:8088/api/auth/test" -Method GET
```

```bash
# bash
# ヘルスチェック
curl -X GET "http://localhost:8088/health"

# Supabase接続テスト
curl -X GET "http://localhost:8088/api/auth/test"
```




## API エンドポイント

### 認証関連


● 現在完成しているAPIと名前は以下の通りです：

  認証関連 (Auth Handler)

  1. POST /api/auth/signup - ユーザー登録
  2. POST /api/auth/login - ユーザーログイン
  3. POST /api/auth/logout - ユーザーログアウト
  4. GET /api/auth/test - Supabase接続テスト

  ジャンル関連 (Genre Handler)

  5. POST /api/genres - ジャンル作成

  問題関連 (Question Handler)

  6. POST /api/questions - 問題作成
  7. GET /api/questions - 問題一覧取得
  8. GET /api/questions/{id} - 特定の問題取得
  9. PUT /api/questions/{id} - 問題更新
  10. DELETE /api/questions/{id} - 問題削除
  11. GET /api/my-questions - ユーザーの問題一覧取得

      回答関連（Answer Handler）

  13. POST /api/answers - 問題に対する自分の回答

      選択肢関連（Choices Handler）

  15. GET /api/choices/{questionID} - 選択肢取得
  16. POST /api/choices/create - 選択肢作成
  17. PUT /api/choices/update - 選択肢更新
  18. DELETE /api/choices/delete/{id} - 選択肢削除

      その他

  19. POST /api/answers - 自身の回答
  20. / - 静的ファイル配信


## 使用例

### ユーザー登録

```bash
curl -X POST http://localhost:8088/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "username": "testuser"
  }'
```

### ユーザーログイン

```bash
curl -X POST http://localhost:8088/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```




## 開発

### プロジェクト構造


```
│  .env
│  .gitignore
│  env.example
│  go.mod
│  go.sum
│  README.md
│  server.exe
│
├─.claude
│      settings.local.json
│
├─.vscode
├─cmd
│  └─server
│          main.go
│
├─internal
│  ├─application
│  │  └─auth
│  │      ├─dto
│  │      │      auth_dto.go
│  │      │
│  │      └─usecases
│  │              auth_usecase.go
│  │
│  ├─domain
│  │  ├─auth
│  │  │  ├─entities
│  │  │  │      user.go
│  │  │  │
│  │  │  ├─repositories
│  │  │  │      user_repository.go
│  │  │  │
│  │  │  └─services
│  │  │          auth_service.go
│  │  │
│  │  └─shared
│  │          errors.go
│  │
│  ├─infrastructure
│  │  ├─auth
│  │  │  └─supabase
│  │  │          user_repository_impl.go
│  │  │
│  │  ├─config
│  │  │      config.go
│  │  │
│  │  ├─database
│  │  └─di
│  │          container.go
│  │
│  └─presentation
│      ├─dto
│      │      auth_dto.go
│      │
│      └─http
│          ├─handlers
│          │      auth_handler.go
│          │
│          ├─middleware
│          │      cors.go
│          │
│          └─router
│                  router.go
│
└─static
        index.html

```

### 依存関係

- `github.com/supabase-community/gotrue-go` - Supabase認証クライアント
- `github.com/joho/godotenv` - 環境変数管理
