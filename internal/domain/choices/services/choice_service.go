package services

import (
	"context"

	entities "Shittaka_back/internal/domain/choices/entities"
	"Shittaka_back/internal/domain/choices/repositories"
)

// ChoiceService はユースケース層のサービス
// Repository を利用してアプリケーションの処理をまとめる
type ChoiceService struct {
	repo repositories.ChoiceRepository
}

// NewChoiceService は ChoiceService のコンストラクタ
func NewChoiceService(repo repositories.ChoiceRepository) *ChoiceService {
	return &ChoiceService{repo: repo}
}

// GetChoices は問題IDに紐づく選択肢を取得
func (s *ChoiceService) GetChoices(ctx context.Context, questionID int64) ([]entities.Choice, error) {
	return s.repo.GetByQuestionID(ctx, questionID)
}

// CreateChoice は新しい選択肢を作成
func (s *ChoiceService) CreateChoice(ctx context.Context, choice entities.Choice) (*entities.Choice, error) {
	return s.repo.Create(ctx, choice)
}

// CreateChoiceWithAuth は認証付きで新しい選択肢を作成
func (s *ChoiceService) CreateChoiceWithAuth(ctx context.Context, choice entities.Choice, userToken string) (*entities.Choice, error) {
	return s.repo.CreateWithAuth(ctx, choice, userToken)
}

// UpdateChoice は既存の選択肢を更新
func (s *ChoiceService) UpdateChoice(ctx context.Context, choice entities.Choice) (*entities.Choice, error) {
	return s.repo.Update(ctx, choice)
}

// DeleteChoice は選択肢を削除
func (s *ChoiceService) DeleteChoice(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
