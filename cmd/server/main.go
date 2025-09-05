package main

import (
	"Shittaka_back/internal/auth"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// 環境変数を読み込み
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	// 必要な環境変数をチェック
	if os.Getenv("SUPABASE_URL") == "" {
		log.Fatal("SUPABASE_URL is required")
	}
	if os.Getenv("SUPABASE_SERVICE_ROLE_KEY") == "" {
		log.Fatal("SUPABASE_SERVICE_ROLE_KEY is required")
	}

	// Supabaseクライアントとハンドラーを作成
	client := auth.NewClient()
	authHandler := auth.NewAuthHandler(client)

	// ルーティング設定
	mux := http.NewServeMux()

	// CORS対応のミドルウェア
	corsHandler := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		}
	}

	// 認証関連のエンドポイント
	mux.HandleFunc("/api/auth/signup", corsHandler(authHandler.SignupHandler))
	mux.HandleFunc("/api/auth/login", corsHandler(authHandler.LoginHandler))
	mux.HandleFunc("/api/auth/logout", corsHandler(authHandler.LogoutHandler))
	mux.HandleFunc("/api/auth/test", corsHandler(authHandler.TestConnectionHandler))

	// ヘルスチェック用エンドポイント
	mux.HandleFunc("/health", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))

	// 静的ファイル配信
	fs := http.FileServer(http.Dir("./static/"))
	mux.Handle("/", http.StripPrefix("/", fs))

	// ポート設定
	port := os.Getenv("PORT")
	if port == "" {
		port = "8088"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Supabase URL: %s", os.Getenv("SUPABASE_URL"))

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
