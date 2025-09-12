package di

import (
	"Shittaka_back/internal/application/answer/usecases"
	"Shittaka_back/internal/infrastructure/answer/supabase"
	"Shittaka_back/internal/presentation/http/handlers"
)

// NewAnswerHandler は新しいAnswerHandlerを作成
func NewAnswerHandler() *handlers.AnswerHandler {
	// 依存関係を構築（外側から内側へ）
	answerRepo := supabase.NewAnswerRepository()
	answerUsecase := usecases.NewAnswerUsecase(answerRepo)
	answerHandler := handlers.NewAnswerHandler(answerUsecase)

	return answerHandler
}