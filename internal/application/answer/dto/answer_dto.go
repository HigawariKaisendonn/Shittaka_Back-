package dto

import "time"

// CreateAnswerRequest は回答作成リクエストDTO
type CreateAnswerRequest struct {
	QuestionID int64 `json:"question_id"`
	ChoiceID   int64 `json:"choice_id"`
}

// AnswerResponse は回答レスポンスDTO
type AnswerResponse struct {
	ID         int64     `json:"id"`
	UserID     string    `json:"user_id"`
	QuestionID int64     `json:"question_id"`
	ChoiceID   int64     `json:"choice_id"`
	AnsweredAt time.Time `json:"answered_at"`
}