package entities

type Choice struct {
	ID         int64  `json:"id"`          // 選択肢ID (PK)
	QuestionID int64  `json:"question_id"` // 紐づく問題のID (FK -> questions.id)
	Text       string `json:"text"`        // 選択肢の本文
	IsCorrect  bool   `json:"is_correct"`  // 正解かどうか
}
