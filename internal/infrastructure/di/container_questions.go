package di

import (
	questionUsecases "Shittaka_back/internal/application/question/usecases"
	questionSupabase "Shittaka_back/internal/infrastructure/question/supabase"
	"Shittaka_back/internal/presentation/http/handlers"
)

// NewQuestionHandler は問題機能の依存関係を構築し、ハンドラーを返す
func NewQuestionHandler() *handlers.QuestionHandler {
	// リポジトリ（Supabase 実装）
	questionRepo := questionSupabase.NewQuestionRepository()

	// ユースケース
	usecase := questionUsecases.NewQuestionUsecase(questionRepo)

	// ハンドラー
	return handlers.NewQuestionHandler(usecase)
}