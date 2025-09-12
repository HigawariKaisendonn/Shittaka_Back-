package di

// container_choices.goは選択肢機能の依存関係配線を定義

import (
	"Shittaka_back/internal/domain/choices/services"
	choiceSupabase "Shittaka_back/internal/infrastructure/choice/supabase"
	"Shittaka_back/internal/presentation/http/handlers"
)

// NewChoiceHandler は選択肢機能の依存関係を構築し、ハンドラーを返す
func NewChoiceHandler() *handlers.ChoiceHandler {
	// リポジトリ（Supabase HTTP実装）
	choiceRepo := choiceSupabase.NewChoiceRepository()

	// サービス
	choiceService := services.NewChoiceService(choiceRepo)

	// ハンドラー
	return handlers.NewChoiceHandler(choiceService)
}