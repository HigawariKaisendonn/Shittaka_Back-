package usecases

import (
	"context"
	"strings"

	"Shittaka_back/internal/application/question/dto"
	"Shittaka_back/internal/domain/question/entities"
	"Shittaka_back/internal/domain/question/repositories"
	"Shittaka_back/internal/domain/shared"
)

// QuestionUsecase は問題ユースケース
type QuestionUsecase struct {
	questionRepo repositories.QuestionRepository
}

// NewQuestionUsecase は新しいQuestionUsecaseを作成
func NewQuestionUsecase(questionRepo repositories.QuestionRepository) *QuestionUsecase {
	return &QuestionUsecase{
		questionRepo: questionRepo,
	}
}

// CreateQuestion は新しい問題を作成する（認証が必要）
func (u *QuestionUsecase) CreateQuestion(ctx context.Context, req dto.CreateQuestionRequest, userID string, userToken string) (*dto.QuestionResponse, error) {
	// バリデーション
	if err := u.validateCreateQuestionRequest(req); err != nil {
		return nil, err
	}

	// 問題エンティティを作成
	question := entities.NewQuestion(req.GenreID, userID, req.Title, req.Body, req.Explanation)

	// エンティティレベルでのバリデーション
	if err := question.Validate(); err != nil {
		return nil, err
	}

	// リポジトリに保存（ユーザートークンを渡してRLS適用）
	createdQuestion, err := u.questionRepo.Create(ctx, question, userToken)
	if err != nil {
		return nil, err
	}

	// レスポンスDTOに変換
	return &dto.QuestionResponse{
		ID:             createdQuestion.ID,
		GenreID:        createdQuestion.GenreID,
		UserID:         createdQuestion.UserID,
		Title:          createdQuestion.Title,
		Body:           createdQuestion.Body,
		Explanation:    createdQuestion.Explanation,
		CreatedAt:      createdQuestion.CreatedAt,
		Views:          createdQuestion.Views,
		CorrectCount:   createdQuestion.CorrectCount,
		IncorrectCount: createdQuestion.IncorrectCount,
	}, nil
}

// UpdateQuestion は問題を更新する（作成者のみ）
func (u *QuestionUsecase) UpdateQuestion(ctx context.Context, id int64, req dto.UpdateQuestionRequest, userID string, userToken string) error {
	// バリデーション
	if err := u.validateUpdateQuestionRequest(req); err != nil {
		return err
	}

	// 既存の問題を取得
	existingQuestion, err := u.questionRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 作成者かどうかチェック
	if existingQuestion.UserID != userID {
		return shared.NewDomainError("FORBIDDEN", "この問題を更新する権限がありません")
	}

	// 問題を更新
	existingQuestion.Title = req.Title
	existingQuestion.Body = req.Body
	existingQuestion.Explanation = req.Explanation

	// リポジトリで更新
	return u.questionRepo.Update(ctx, existingQuestion, userToken)
}

// DeleteQuestion は問題を削除する（作成者のみ）
func (u *QuestionUsecase) DeleteQuestion(ctx context.Context, id int64, userID string, userToken string) error {
	// 既存の問題を取得
	existingQuestion, err := u.questionRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 作成者かどうかチェック
	if existingQuestion.UserID != userID {
		return shared.NewDomainError("FORBIDDEN", "この問題を削除する権限がありません")
	}

	// リポジトリで削除
	return u.questionRepo.Delete(ctx, id, userToken)
}

// GetQuestion は問題を取得する
func (u *QuestionUsecase) GetQuestion(ctx context.Context, id int64) (*dto.QuestionResponse, error) {
	question, err := u.questionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// レスポンスDTOに変換
	return &dto.QuestionResponse{
		ID:             question.ID,
		GenreID:        question.GenreID,
		UserID:         question.UserID,
		Title:          question.Title,
		Body:           question.Body,
		Explanation:    question.Explanation,
		CreatedAt:      question.CreatedAt,
		Views:          question.Views,
		CorrectCount:   question.CorrectCount,
		IncorrectCount: question.IncorrectCount,
	}, nil
}

// GetQuestionsByUser はユーザーの問題一覧を取得する
func (u *QuestionUsecase) GetQuestionsByUser(ctx context.Context, userID string, userToken string) ([]*dto.QuestionResponse, error) {
	questions, err := u.questionRepo.GetByUserID(ctx, userID, userToken)
	if err != nil {
		return nil, err
	}

	// レスポンスDTOに変換
	responses := make([]*dto.QuestionResponse, len(questions))
	for i, question := range questions {
		responses[i] = &dto.QuestionResponse{
			ID:             question.ID,
			GenreID:        question.GenreID,
			UserID:         question.UserID,
			Title:          question.Title,
			Body:           question.Body,
			Explanation:    question.Explanation,
			CreatedAt:      question.CreatedAt,
			Views:          question.Views,
			CorrectCount:   question.CorrectCount,
			IncorrectCount: question.IncorrectCount,
		}
	}

	return responses, nil
}

// GetAllQuestions は全ての問題を取得する
func (u *QuestionUsecase) GetAllQuestions(ctx context.Context) ([]*dto.QuestionResponse, error) {
	questions, err := u.questionRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// レスポンスDTOに変換
	responses := make([]*dto.QuestionResponse, len(questions))
	for i, question := range questions {
		responses[i] = &dto.QuestionResponse{
			ID:             question.ID,
			GenreID:        question.GenreID,
			UserID:         question.UserID,
			Title:          question.Title,
			Body:           question.Body,
			Explanation:    question.Explanation,
			CreatedAt:      question.CreatedAt,
			Views:          question.Views,
			CorrectCount:   question.CorrectCount,
			IncorrectCount: question.IncorrectCount,
		}
	}

	return responses, nil
}

// validateCreateQuestionRequest は問題作成リクエストをバリデーション
func (u *QuestionUsecase) validateCreateQuestionRequest(req dto.CreateQuestionRequest) error {
	if req.GenreID == 0 {
		return shared.NewValidationError("genre_id", "ジャンルIDは必須です")
	}

	if strings.TrimSpace(req.Title) == "" {
		return shared.NewValidationError("title", "問題タイトルは必須です")
	}

	if len(req.Title) > 200 {
		return shared.NewValidationError("title", "問題タイトルは200文字以内で入力してください")
	}

	return nil
}

// validateUpdateQuestionRequest は問題更新リクエストをバリデーション
func (u *QuestionUsecase) validateUpdateQuestionRequest(req dto.UpdateQuestionRequest) error {
	if strings.TrimSpace(req.Title) == "" {
		return shared.NewValidationError("title", "問題タイトルは必須です")
	}

	if len(req.Title) > 200 {
		return shared.NewValidationError("title", "問題タイトルは200文字以内で入力してください")
	}

	return nil
}