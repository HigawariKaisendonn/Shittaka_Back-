package main

// main.goはサーバー起動のメインファイル

import (
	"log"
	"net/http"

	"Shittaka_back/internal/infrastructure/di"
	"Shittaka_back/internal/presentation/http/router"
)

func main() {
	// DIコンテナを初期化
	container := di.NewContainer()

	log.Printf("Server starting on port %s", container.Config.Port)
	log.Printf("Supabase URL: %s", container.Config.SupabaseURL)

	// ルーターを設定
	mux := router.SetupRoutes(container.AuthHandler)

	// サーバーを起動
	if err := http.ListenAndServe(":"+container.Config.Port, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
