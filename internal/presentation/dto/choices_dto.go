package dto

// choices_dto.goは選択肢関連のHTTP DTOを定義

// CreateChoiceRequest は選択肢作成リクエストのHTTP DTO
type CreateChoiceRequest struct {
	QuestionID int64  `json:"question_id"`
	Text       string `json:"text"`
	IsCorrect  bool   `json:"is_correct"`
}

// UpdateChoiceRequest は選択肢更新リクエストのHTTP DTO
type UpdateChoiceRequest struct {
	ID         int64  `json:"id"`
	QuestionID int64  `json:"question_id"`
	Text       string `json:"text"`
	IsCorrect  bool   `json:"is_correct"`
}

// ChoiceResponse は選択肢レスポンスのHTTP DTO
type ChoiceResponse struct {
	ID         int64  `json:"id"`
	QuestionID int64  `json:"question_id"`
	Text       string `json:"text"`
	IsCorrect  bool   `json:"is_correct"`
}

// ChoicesResponse は複数選択肢のレスポンスのHTTP DTO
type ChoicesResponse struct {
	Choices []ChoiceResponse `json:"choices"`
}