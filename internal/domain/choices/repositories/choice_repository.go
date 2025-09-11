package repositories

import (
	"context"
	"strconv"

	entities "Shittaka_back/internal/domain/choices/entities"

	"github.com/nedpals/supabase-go"
)

// ChoiceRepository は選択肢のリポジトリを表すインターフェース
// Service層から利用され、DB操作の抽象化を担当する
type ChoiceRepository interface {
	GetByQuestionID(ctx context.Context, questionID int64) ([]entities.Choice, error) // 問題IDに紐づく選択肢を取得
	Create(ctx context.Context, choice entities.Choice) (*entities.Choice, error)     // 新しい選択肢を作成
	Update(ctx context.Context, choice entities.Choice) (*entities.Choice, error)     // 既存の選択肢を更新
	Delete(ctx context.Context, id int64) error                                       // 選択肢を削除
}

// choiceRepository は ChoiceRepository インターフェースの実装
type choiceRepository struct {
	client *supabase.Client // Supabase のクライアント
}

// NewChoiceRepository は ChoiceRepository のコンストラクタ
func NewChoiceRepository(client *supabase.Client) ChoiceRepository {
	return &choiceRepository{client: client}
}

// GetByQuestionID は指定された questionID に紐づく選択肢を DB から取得
func (r *choiceRepository) GetByQuestionID(ctx context.Context, questionID int64) ([]entities.Choice, error) {
	var choices []entities.Choice       //選択肢を格納するスライスを準備
	err := r.client.DB.From("choices"). //choicesテーブルを対象にクエリを作成
						Select("*").                                          //全てのカラム（id, question_id, test, is_correct）を取得
						Eq("question_id", strconv.FormatInt(questionID, 10)). // int64 → string に変換して検索
						Execute(&choices)
	if err != nil {
		return nil, err
	}
	return choices, nil
}

// Create は新しい選択肢を DB に追加
func (r *choiceRepository) Create(ctx context.Context, choice entities.Choice) (*entities.Choice, error) {
	var inserted []entities.Choice
	err := r.client.DB.From("choices").
		Insert(choice).
		Execute(&inserted) // Execute の結果をスライスに格納
	if err != nil {
		return nil, err
	}
	return &inserted[0], nil // 追加されたレコードを返す
}

// Update は既存の選択肢を更新
func (r *choiceRepository) Update(ctx context.Context, choice entities.Choice) (*entities.Choice, error) {
	var updated []entities.Choice
	err := r.client.DB.From("choices").
		Update(choice).
		Eq("id", strconv.FormatInt(choice.ID, 10)). // int64 → string
		Execute(&updated)                           // ← Execute にポインタを渡す
	if err != nil {
		return nil, err
	}
	return &updated[0], nil
}

// Delete は選択肢を ID 指定で削除
func (r *choiceRepository) Delete(ctx context.Context, id int64) error {
	return r.client.DB.From("choices").
		Delete().
		Eq("id", strconv.FormatInt(id, 10)). // ID を条件に削除
		Execute(nil)
}
