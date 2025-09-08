package di

// container_ganres.goはジャンル機能の依存関係配線を定義

import (
	genreUsecases "Shittaka_back/internal/application/genre/usecases"
	genreSupabase "Shittaka_back/internal/infrastructure/genre/supabase"
	"Shittaka_back/internal/presentation/http/handlers"
)

// NewGenreHandler はジャンル機能の依存関係を構築し、ハンドラーを返す
func NewGenreHandler() *handlers.GenreHandler {
	// リポジトリ（Supabase 実装）
	genreRepo := genreSupabase.NewGenreRepository()

	// ユースケース
	usecase := genreUsecases.NewGenreUsecase(genreRepo)

	// ハンドラー
	return handlers.NewGenreHandler(usecase)
}
