package di

import (
	"Shittaka_back/internal/application/auth/usecases"
	"Shittaka_back/internal/domain/auth/services"
	"Shittaka_back/internal/infrastructure/auth/supabase"
	"Shittaka_back/internal/presentation/http/handlers"
	"Shittaka_back/internal/infrastructure/config"
)

// Container は依存関係のコンテナ
type Container struct {
	Config      *config.Config
	AuthHandler *handlers.AuthHandler
}

// NewContainer は新しいコンテナを作成
func NewContainer() *Container {
	// 設定を読み込み
	cfg := config.LoadConfig()
	
	// 依存関係を構築（外側から内側へ）
	userRepo := supabase.NewUserRepository()
	authService := services.NewAuthService(userRepo)
	authUsecase := usecases.NewAuthUsecase(authService)
	authHandler := handlers.NewAuthHandler(authUsecase)
	
	return &Container{
		Config:      cfg,
		AuthHandler: authHandler,
	}
}