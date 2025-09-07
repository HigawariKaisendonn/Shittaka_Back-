# Shittaka Backend

Go言語で構築されたSupabase認証システムのバックエンドAPIです。

## 機能

- ユーザー登録（サインアップ）
- ユーザーログイン
- ユーザーログアウト
- Supabaseとの疎通テスト
- CORS対応

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

# サーバー設定
PORT=8088
```

### 3. Supabaseプロジェクトの設定

1. [Supabase](https://supabase.com)でプロジェクトを作成
2. プロジェクト設定から以下を取得：
   - Project URL (`SUPABASE_URL`)
   - Service Role Key (`SUPABASE_SERVICE_ROLE_KEY`)

### 4.server.exeの作成

``` bash
go build -o server.exe cmd/server/main.go
```

### 5.接続テスト

""の中をenvのものに置き換えてください

```powershell
# powershell
$env:SUPABASE_URL = ".env参照"
$env:SUPABASE_SERVICE_ROLE_KEY = ".env参照"
$env:PORT = "8088"

# サーバーを起動
./server.exe
```

```powershell
# ヘルスチェック
Invoke-WebRequest -Uri "http://localhost:8088/health" -Method GET

# Supabase接続テスト
Invoke-WebRequest -Uri "http://localhost:8088/api/auth/test" -Method GET
```




## API エンドポイント

### 認証関連

- `POST /api/auth/signup` - ユーザー登録
- `POST /api/auth/login` - ユーザーログイン
- `POST /api/auth/logout` - ユーザーログアウト
- `GET /api/auth/test` - Supabase接続テスト

### その他

- `GET /health` - ヘルスチェック

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
├── cmd/
│   └── server/
│       └── main.go          # メインサーバーファイル
├── internal/
│   └── auth/
│       ├── client.go        # Supabaseクライアント
│       ├── handler.go       # HTTPハンドラー
│       └── types.go         # 型定義
├── go.mod
├── go.sum
└── README.md
```

### 依存関係

- `github.com/supabase-community/gotrue-go` - Supabase認証クライアント
- `github.com/joho/godotenv` - 環境変数管理
