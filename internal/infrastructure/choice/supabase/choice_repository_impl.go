package supabase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"Shittaka_back/internal/domain/choices/entities"
	"Shittaka_back/internal/domain/choices/repositories"
)

// ChoiceRepositoryImpl はSupabaseを使用したChoiceRepositoryの実装
type ChoiceRepositoryImpl struct{}

// NewChoiceRepository は新しいChoiceRepositoryImplを作成
func NewChoiceRepository() repositories.ChoiceRepository {
	return &ChoiceRepositoryImpl{}
}

// GetByQuestionID は問題IDで選択肢一覧を取得
func (r *ChoiceRepositoryImpl) GetByQuestionID(ctx context.Context, questionID int64) ([]entities.Choice, error) {
	url := fmt.Sprintf("%s/rest/v1/choices?question_id=eq.%d", os.Getenv("SUPABASE_URL"), questionID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_ANON_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("find choices failed with status %d: %s", resp.StatusCode, string(body))
	}

	var choiceList []map[string]interface{}
	if err := json.Unmarshal(body, &choiceList); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	choices := make([]entities.Choice, len(choiceList))
	for i, choiceData := range choiceList {
		choices[i] = mapToChoice(choiceData)
	}

	return choices, nil
}

// Create は新しい選択肢を作成
func (r *ChoiceRepositoryImpl) Create(ctx context.Context, choice entities.Choice) (*entities.Choice, error) {
	choiceData := map[string]interface{}{
		"question_id": choice.QuestionID,
		"text":        choice.Text,
		"is_correct":  choice.IsCorrect,
	}

	jsonData, err := json.Marshal(choiceData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal choice data: %w", err)
	}

	url := os.Getenv("SUPABASE_URL") + "/rest/v1/choices"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_ANON_KEY"))
	req.Header.Set("Prefer", "return=representation")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("create choice failed with status %d: %s", resp.StatusCode, string(body))
	}

	var choiceList []map[string]interface{}
	if err := json.Unmarshal(body, &choiceList); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(choiceList) == 0 {
		return nil, fmt.Errorf("no choice returned from create operation")
	}

	choiceResp := choiceList[0]
	result := mapToChoice(choiceResp)
	return &result, nil
}

// CreateWithAuth は認証トークンを使って新しい選択肢を作成
func (r *ChoiceRepositoryImpl) CreateWithAuth(ctx context.Context, choice entities.Choice, userToken string) (*entities.Choice, error) {
	choiceData := map[string]interface{}{
		"question_id": choice.QuestionID,
		"text":        choice.Text,
		"is_correct":  choice.IsCorrect,
	}

	jsonData, err := json.Marshal(choiceData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal choice data: %w", err)
	}

	url := os.Getenv("SUPABASE_URL") + "/rest/v1/choices"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))
	req.Header.Set("Authorization", "Bearer "+userToken)
	req.Header.Set("Prefer", "return=representation")
	
	// デバッグ用ログ
	if len(userToken) > 20 {
		fmt.Printf("CreateWithAuth - UserToken: %s...\n", userToken[:20])
	} else {
		fmt.Printf("CreateWithAuth - UserToken: %s\n", userToken)
	}
	fmt.Printf("CreateWithAuth - URL: %s\n", url)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("create choice failed with status %d: %s", resp.StatusCode, string(body))
	}

	var choiceList []map[string]interface{}
	if err := json.Unmarshal(body, &choiceList); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(choiceList) == 0 {
		return nil, fmt.Errorf("no choice returned from create operation")
	}

	choiceResp := choiceList[0]
	result := mapToChoice(choiceResp)
	return &result, nil
}

// Update は選択肢を更新
func (r *ChoiceRepositoryImpl) Update(ctx context.Context, choice entities.Choice) (*entities.Choice, error) {
	choiceData := map[string]interface{}{
		"text":       choice.Text,
		"is_correct": choice.IsCorrect,
	}

	jsonData, err := json.Marshal(choiceData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal choice data: %w", err)
	}

	url := fmt.Sprintf("%s/rest/v1/choices?id=eq.%d", os.Getenv("SUPABASE_URL"), choice.ID)
	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_ANON_KEY"))
	req.Header.Set("Prefer", "return=representation")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("update choice failed with status %d: %s", resp.StatusCode, string(body))
	}

	var choiceList []map[string]interface{}
	if err := json.Unmarshal(body, &choiceList); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(choiceList) == 0 {
		return nil, fmt.Errorf("no choice returned from update operation")
	}

	choiceResp := choiceList[0]
	result := mapToChoice(choiceResp)
	return &result, nil
}

// Delete は選択肢を削除
func (r *ChoiceRepositoryImpl) Delete(ctx context.Context, id int64) error {
	url := fmt.Sprintf("%s/rest/v1/choices?id=eq.%d", os.Getenv("SUPABASE_URL"), id)
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_ANON_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("delete choice failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// mapToChoice は map[string]interface{} を Choice エンティティに変換
func mapToChoice(m map[string]interface{}) entities.Choice {
	return entities.Choice{
		ID:         getInt64(m, "id"),
		QuestionID: getInt64(m, "question_id"),
		Text:       getString(m, "text"),
		IsCorrect:  getBool(m, "is_correct"),
	}
}

// ヘルパー関数

// getString は map から文字列を安全に取得
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// getInt64 は map から int64 を安全に取得
func getInt64(m map[string]interface{}, key string) int64 {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case float64:
			return int64(v)
		case int64:
			return v
		case int:
			return int64(v)
		case string:
			if i, err := strconv.ParseInt(v, 10, 64); err == nil {
				return i
			}
		}
	}
	return 0
}

// getBool は map から bool を安全に取得
func getBool(m map[string]interface{}, key string) bool {
	if val, ok := m[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}