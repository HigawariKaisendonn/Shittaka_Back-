package repositories

import (
	"context"
	"Shittaka_back/internal/domain/question/entities"
)

// QuestionRepository は問題リポジトリのインターフェース
type QuestionRepository interface {
	Create(ctx context.Context, question *entities.Question, userToken string) (*entities.Question, error)
	GetByID(ctx context.Context, id int64) (*entities.Question, error)
	GetByUserID(ctx context.Context, userID string, userToken string) ([]*entities.Question, error)
	Update(ctx context.Context, question *entities.Question, userToken string) error
	Delete(ctx context.Context, id int64, userToken string) error
	GetAll(ctx context.Context) ([]*entities.Question, error)
}