package usecases

import (
	"context"

	"Shittaka_back/internal/application/answer/dto"
	"Shittaka_back/internal/domain/answer/entities"
	"Shittaka_back/internal/domain/answer/repositories"
	"Shittaka_back/internal/domain/shared"
)

// AnswerUsecase は回答ユースケース
type AnswerUsecase struct {
	answerRepo repositories.AnswerRepository
}

// NewAnswerUsecase は新しいAnswerUsecaseを作成
func NewAnswerUsecase(answerRepo repositories.AnswerRepository) *AnswerUsecase {
	return &AnswerUsecase{
		answerRepo: answerRepo,
	}
}

// CreateAnswer は新しい回答を作成する（認証が必要）
func (u *AnswerUsecase) CreateAnswer(ctx context.Context, req dto.CreateAnswerRequest, userID string, userToken string) (*dto.AnswerResponse, error) {
	// バリデーション
	if err := u.validateCreateAnswerRequest(req); err != nil {
		return nil, err
	}

	// 回答エンティティを作成
	answer := entities.NewAnswer(userID, req.QuestionID, req.ChoiceID)

	// エンティティレベルでのバリデーション
	if err := answer.Validate(); err != nil {
		return nil, err
	}

	// リポジトリに保存（ユーザートークンを渡してRLS適用）
	createdAnswer, err := u.answerRepo.Create(ctx, answer, userToken)
	if err != nil {
		return nil, err
	}

	// レスポンスDTOに変換
	return &dto.AnswerResponse{
		ID:         createdAnswer.ID,
		UserID:     createdAnswer.UserID,
		QuestionID: createdAnswer.QuestionID,
		ChoiceID:   createdAnswer.ChoiceID,
		AnsweredAt: createdAnswer.AnsweredAt,
	}, nil
}

// GetAnswersByUser はユーザーの回答一覧を取得する
func (u *AnswerUsecase) GetAnswersByUser(ctx context.Context, userID string) ([]*dto.AnswerResponse, error) {
	answers, err := u.answerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// レスポンスDTOに変換
	responses := make([]*dto.AnswerResponse, len(answers))
	for i, answer := range answers {
		responses[i] = &dto.AnswerResponse{
			ID:         answer.ID,
			UserID:     answer.UserID,
			QuestionID: answer.QuestionID,
			ChoiceID:   answer.ChoiceID,
			AnsweredAt: answer.AnsweredAt,
		}
	}

	return responses, nil
}

// GetAnswersByQuestion は問題の回答一覧を取得する
func (u *AnswerUsecase) GetAnswersByQuestion(ctx context.Context, questionID int64) ([]*dto.AnswerResponse, error) {
	answers, err := u.answerRepo.GetByQuestionID(ctx, questionID)
	if err != nil {
		return nil, err
	}

	// レスポンスDTOに変換
	responses := make([]*dto.AnswerResponse, len(answers))
	for i, answer := range answers {
		responses[i] = &dto.AnswerResponse{
			ID:         answer.ID,
			UserID:     answer.UserID,
			QuestionID: answer.QuestionID,
			ChoiceID:   answer.ChoiceID,
			AnsweredAt: answer.AnsweredAt,
		}
	}

	return responses, nil
}

// validateCreateAnswerRequest は回答作成リクエストをバリデーション
func (u *AnswerUsecase) validateCreateAnswerRequest(req dto.CreateAnswerRequest) error {
	if req.QuestionID == 0 {
		return shared.NewValidationError("question_id", "問題IDは必須です")
	}

	if req.ChoiceID == 0 {
		return shared.NewValidationError("choice_id", "選択肢IDは必須です")
	}

	return nil
}