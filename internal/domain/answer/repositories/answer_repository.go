package repositories

import (
	"context"
	"Shittaka_back/internal/domain/answer/entities"
)

// AnswerRepository は回答履歴リポジトリのインターフェース
type AnswerRepository interface {
	Create(ctx context.Context, answer *entities.Answer, userToken string) (*entities.Answer, error)
	GetByUserID(ctx context.Context, userID string) ([]*entities.Answer, error)
	GetByQuestionID(ctx context.Context, questionID int64) ([]*entities.Answer, error)
}