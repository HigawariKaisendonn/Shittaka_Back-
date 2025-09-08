package entities

import (
	"Shittaka_back/internal/domain/shared"
	"time"
)

// Question は問題のドメインエンティティ
type Question struct {
	ID             int64     `json:"id"`
	GenreID        int64     `json:"genre_id"`
	UserID         string    `json:"user_id"`
	Title          string    `json:"title"`
	Body           string    `json:"body"`
	Explanation    string    `json:"explanation"`
	CreatedAt      time.Time `json:"created_at"`
	Views          int       `json:"views"`
	CorrectCount   int       `json:"correct_count"`
	IncorrectCount int       `json:"incorrect_count"`
}

// NewQuestion は新しいQuestionエンティティを作成
func NewQuestion(genreID int64, userID, title, body, explanation string) *Question {
	return &Question{
		GenreID:        genreID,
		UserID:         userID,
		Title:          title,
		Body:           body,
		Explanation:    explanation,
		CreatedAt:      time.Now(),
		Views:          0,
		CorrectCount:   0,
		IncorrectCount: 0,
	}
}

// Validate はQuestionエンティティのバリデーションを行う
func (q *Question) Validate() error {
	if q.GenreID == 0 {
		return shared.NewValidationError("genre_id", "genre_id is required")
	}
	if q.UserID == "" {
		return shared.NewValidationError("user_id", "user_id is required")
	}
	if q.Title == "" {
		return shared.NewValidationError("title", "title is required")
	}
	return nil
}

// IncrementViews は閲覧数をインクリメント
func (q *Question) IncrementViews() {
	q.Views++
}

// IncrementCorrectCount は正解数をインクリメント
func (q *Question) IncrementCorrectCount() {
	q.CorrectCount++
}

// IncrementIncorrectCount は不正解数をインクリメント
func (q *Question) IncrementIncorrectCount() {
	q.IncorrectCount++
}