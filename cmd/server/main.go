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
	authContainer := di.NewContainer()
	genreHandler := di.NewGenreHandler()
	questionHandler := di.NewQuestionHandler()
	answerHandler := di.NewAnswerHandler()
	choiceHandler := di.NewChoiceHandler()

	log.Printf("Server starting on port %s", authContainer.Config.Port)
	log.Printf("Supabase URL: %s", authContainer.Config.SupabaseURL)

	// ルーターを設定
	mux := router.SetupRoutes(authContainer.AuthHandler, genreHandler, questionHandler, answerHandler, choiceHandler)

	// サーバーを起動
	if err := http.ListenAndServe(":"+authContainer.Config.Port, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
