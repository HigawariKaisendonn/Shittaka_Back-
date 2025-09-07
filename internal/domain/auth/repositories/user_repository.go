package repositories

// user_repository.goはユーザーに関するリポジトリのインターフェースを定義
// リポジトリとは、ドメインから見たデータストアのインターフェース

import (
	"Shittaka_back/internal/domain/auth/entities"
	"context"
)

// UserRepository はユーザーに関するリポジトリのインターフェース
type UserRepository interface {
	// Create は新しいユーザーを作成
	Create(ctx context.Context, email, password string, metadata map[string]interface{}) (*entities.User, error)

	// Authenticate はユーザーの認証を行い、トークンを返す
	Authenticate(ctx context.Context, email, password string) (*AuthResult, error)

	// FindByID はIDでユーザーを検索
	FindByID(ctx context.Context, id string) (*entities.User, error)

	// FindByEmail はEmailでユーザーを検索
	FindByEmail(ctx context.Context, email string) (*entities.User, error)

	// Update はユーザー情報を更新
	Update(ctx context.Context, user *entities.User) error

	// Delete はユーザーを削除
	Delete(ctx context.Context, id string) error

	// Logout はユーザーをログアウトさせる
	Logout(ctx context.Context, token string) error
}

// AuthResult は認証結果を表す
type AuthResult struct {
	User         *entities.User
	AccessToken  string
	RefreshToken string
	ExpiresAt    int64
}
