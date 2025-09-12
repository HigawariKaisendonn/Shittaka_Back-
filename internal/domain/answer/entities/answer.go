package entities

import (
	"Shittaka_back/internal/domain/shared"
	"time"
)

// Answer は回答履歴のドメインエンティティ
type Answer struct {
	ID         int64     `json:"id"`
	UserID     string    `json:"user_id"`
	QuestionID int64     `json:"question_id"`
	ChoiceID   int64     `json:"choice_id"`
	AnsweredAt time.Time `json:"answered_at"`
}

// NewAnswer は新しいAnswerエンティティを作成
func NewAnswer(userID string, questionID, choiceID int64) *Answer {
	return &Answer{
		UserID:     userID,
		QuestionID: questionID,
		ChoiceID:   choiceID,
		AnsweredAt: time.Now(),
	}
}

// Validate はAnswerエンティティのバリデーションを行う
func (a *Answer) Validate() error {
	if a.UserID == "" {
		return shared.NewValidationError("user_id", "user_id is required")
	}
	if a.QuestionID == 0 {
		return shared.NewValidationError("question_id", "question_id is required")
	}
	if a.ChoiceID == 0 {
		return shared.NewValidationError("choice_id", "choice_id is required")
	}
	return nil
}