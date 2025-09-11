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
	"time"

	"Shittaka_back/internal/domain/answer/entities"
	"Shittaka_back/internal/domain/answer/repositories"
)

// AnswerRepositoryImpl はSupabaseを使用したAnswerRepositoryの実装
type AnswerRepositoryImpl struct{}

// NewAnswerRepository は新しいAnswerRepositoryImplを作成
func NewAnswerRepository() repositories.AnswerRepository {
	return &AnswerRepositoryImpl{}
}

// Create は新しい回答を作成（RLS適用のためユーザートークンを使用）
func (r *AnswerRepositoryImpl) Create(ctx context.Context, answer *entities.Answer, userToken string) (*entities.Answer, error) {
	answerData := map[string]interface{}{
		"user_id":     answer.UserID,
		"question_id": answer.QuestionID,
		"choice_id":   answer.ChoiceID,
	}

	jsonData, err := json.Marshal(answerData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal answer data: %w", err)
	}

	url := os.Getenv("SUPABASE_URL") + "/rest/v1/answers"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))
	req.Header.Set("Authorization", "Bearer "+userToken)
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
		return nil, fmt.Errorf("create answer failed with status %d: %s", resp.StatusCode, string(body))
	}

	var answerList []map[string]interface{}
	if err := json.Unmarshal(body, &answerList); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(answerList) == 0 {
		return nil, fmt.Errorf("no answer returned from create operation")
	}

	answerResp := answerList[0]
	return mapToAnswer(answerResp), nil
}

// GetByUserID はユーザーIDで回答一覧を取得
func (r *AnswerRepositoryImpl) GetByUserID(ctx context.Context, userID string) ([]*entities.Answer, error) {
	url := fmt.Sprintf("%s/rest/v1/answers?user_id=eq.%s", os.Getenv("SUPABASE_URL"), userID)
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
		return nil, fmt.Errorf("find answers by user failed with status %d: %s", resp.StatusCode, string(body))
	}

	var answerList []map[string]interface{}
	if err := json.Unmarshal(body, &answerList); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	answers := make([]*entities.Answer, len(answerList))
	for i, answerData := range answerList {
		answers[i] = mapToAnswer(answerData)
	}

	return answers, nil
}

// GetByQuestionID は問題IDで回答一覧を取得
func (r *AnswerRepositoryImpl) GetByQuestionID(ctx context.Context, questionID int64) ([]*entities.Answer, error) {
	url := fmt.Sprintf("%s/rest/v1/answers?question_id=eq.%d", os.Getenv("SUPABASE_URL"), questionID)
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
		return nil, fmt.Errorf("find answers by question failed with status %d: %s", resp.StatusCode, string(body))
	}

	var answerList []map[string]interface{}
	if err := json.Unmarshal(body, &answerList); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	answers := make([]*entities.Answer, len(answerList))
	for i, answerData := range answerList {
		answers[i] = mapToAnswer(answerData)
	}

	return answers, nil
}

// mapToAnswer は map[string]interface{} を Answer エンティティに変換
func mapToAnswer(m map[string]interface{}) *entities.Answer {
	return &entities.Answer{
		ID:         getInt64(m, "id"),
		UserID:     getString(m, "user_id"),
		QuestionID: getInt64(m, "question_id"),
		ChoiceID:   getInt64(m, "choice_id"),
		AnsweredAt: getTime(m, "answered_at"),
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

// getTime は map から time.Time を安全に取得
func getTime(m map[string]interface{}, key string) time.Time {
	if val, ok := m[key]; ok {
		if timeStr, ok := val.(string); ok {
			if t, err := time.Parse(time.RFC3339, timeStr); err == nil {
				return t
			}
		}
	}
	return time.Time{}
}