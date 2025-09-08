package router

// router.goはルーティングを設定

import (
	"net/http"

	"Shittaka_back/internal/presentation/http/handlers"
	"Shittaka_back/internal/presentation/http/middleware"
)

// SetupRoutes はルーティングを設定
func SetupRoutes(authHandler *handlers.AuthHandler, genreHandler *handlers.GenreHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// 認証関連のエンドポイント
	mux.HandleFunc("/api/auth/signup", middleware.CORS(authHandler.SignupHandler))
	mux.HandleFunc("/api/auth/login", middleware.CORS(authHandler.LoginHandler))
	mux.HandleFunc("/api/auth/logout", middleware.CORS(authHandler.LogoutHandler))
	mux.HandleFunc("/api/auth/test", middleware.CORS(authHandler.TestConnectionHandler))

	// ジャンル関連のエンドポイント
	mux.HandleFunc("/api/genres", middleware.CORS(genreHandler.CreateGenreHandler))

	// ヘルスチェック用エンドポイント
	mux.HandleFunc("/health", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))

	// 静的ファイル配信
	fs := http.FileServer(http.Dir("./static/"))
	mux.Handle("/", http.StripPrefix("/", fs))

	return mux
}
