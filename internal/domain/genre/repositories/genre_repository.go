package repositories

// genre_repository.goはジャンルリポジトリのインターフェースを定義

import (
	"context"

	"Shittaka_back/internal/domain/genre/entities"
)

// GenreRepository はジャンルリポジトリのインターフェース
type GenreRepository interface {
	// Create は新しいジャンルを作成する（認証が必要）
	Create(ctx context.Context, genre *entities.Genre, userToken string) (*entities.Genre, error)
	
	// FindByID はIDでジャンルを検索する
	FindByID(ctx context.Context, id int64) (*entities.Genre, error)
	
	// FindAll は全てのジャンルを取得する
	FindAll(ctx context.Context) ([]*entities.Genre, error)
	
	// FindByName は名前でジャンルを検索する（認証が必要）
	FindByName(ctx context.Context, name string, userToken string) (*entities.Genre, error)
}