package repositories

// genre_repository.goはジャンルリポジトリのインターフェースを定義

import (
	"context"

	"Shittaka_back/internal/domain/genre/entities"
)

// GenreRepository はジャンルリポジトリのインターフェース
type GenreRepository interface {
	// Create は新しいジャンルを作成する
	Create(ctx context.Context, genre *entities.Genre) (*entities.Genre, error)
	
	// FindByID はIDでジャンルを検索する
	FindByID(ctx context.Context, id int64) (*entities.Genre, error)
	
	// FindAll は全てのジャンルを取得する
	FindAll(ctx context.Context) ([]*entities.Genre, error)
	
	// FindByName は名前でジャンルを検索する
	FindByName(ctx context.Context, name string) (*entities.Genre, error)
}