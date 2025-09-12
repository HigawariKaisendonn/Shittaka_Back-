package dto

import "time"

// CreateQuestionRequest は問題作成リクエストのHTTP DTO
type CreateQuestionRequest struct {
	GenreID     int64  `json:"genre_id"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	Explanation string `json:"explanation"`
}

// UpdateQuestionRequest は問題更新リクエストのHTTP DTO
type UpdateQuestionRequest struct {
	Title       string `json:"title"`
	Body        string `json:"body"`
	Explanation string `json:"explanation"`
}

// QuestionResponse は問題レスポンスのHTTP DTO
type QuestionResponse struct {
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