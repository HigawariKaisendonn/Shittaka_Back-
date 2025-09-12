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

	"Shittaka_back/internal/domain/question/entities"
	"Shittaka_back/internal/domain/question/repositories"
	"Shittaka_back/internal/domain/shared"
)

// QuestionRepositoryImpl はSupabaseを使用したQuestionRepositoryの実装
type QuestionRepositoryImpl struct{}

// NewQuestionRepository は新しいQuestionRepositoryImplを作成
func NewQuestionRepository() repositories.QuestionRepository {
	return &QuestionRepositoryImpl{}
}

// Create は新しい問題を作成（RLS適用のためユーザートークンを使用）
func (r *QuestionRepositoryImpl) Create(ctx context.Context, question *entities.Question, userToken string) (*entities.Question, error) {
	questionData := map[string]interface{}{
		"genre_id":    question.GenreID,
		"user_id":     question.UserID,
		"title":       question.Title,
		"body":        question.Body,
		"explanation": question.Explanation,
	}

	jsonData, err := json.Marshal(questionData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal question data: %w", err)
	}

	url := os.Getenv("SUPABASE_URL") + "/rest/v1/questions"
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
		return nil, fmt.Errorf("create question failed with status %d: %s", resp.StatusCode, string(body))
	}

	var questionList []map[string]interface{}
	if err := json.Unmarshal(body, &questionList); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(questionList) == 0 {
		return nil, fmt.Errorf("no question returned from create operation")
	}

	questionResp := questionList[0]
	return mapToQuestion(questionResp), nil
}

// GetByID はIDで問題を検索
func (r *QuestionRepositoryImpl) GetByID(ctx context.Context, id int64) (*entities.Question, error) {
	url := fmt.Sprintf("%s/rest/v1/questions?id=eq.%d", os.Getenv("SUPABASE_URL"), id)
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
		return nil, fmt.Errorf("find question failed with status %d: %s", resp.StatusCode, string(body))
	}

	var questionList []map[string]interface{}
	if err := json.Unmarshal(body, &questionList); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(questionList) == 0 {
		return nil, shared.NewDomainError("NOT_FOUND", "問題が見つかりません")
	}

	return mapToQuestion(questionList[0]), nil
}

// GetByUserID はユーザーIDで問題一覧を取得
func (r *QuestionRepositoryImpl) GetByUserID(ctx context.Context, userID string, userToken string) ([]*entities.Question, error) {
	url := fmt.Sprintf("%s/rest/v1/questions?user_id=eq.%s", os.Getenv("SUPABASE_URL"), userID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))
	req.Header.Set("Authorization", "Bearer "+userToken)

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
		return nil, fmt.Errorf("find questions by user failed with status %d: %s", resp.StatusCode, string(body))
	}

	var questionList []map[string]interface{}
	if err := json.Unmarshal(body, &questionList); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	questions := make([]*entities.Question, len(questionList))
	for i, questionData := range questionList {
		questions[i] = mapToQuestion(questionData)
	}

	return questions, nil
}

// Update は問題を更新（RLS適用のためユーザートークンを使用）
func (r *QuestionRepositoryImpl) Update(ctx context.Context, question *entities.Question, userToken string) error {
	questionData := map[string]interface{}{
		"title":       question.Title,
		"body":        question.Body,
		"explanation": question.Explanation,
	}

	jsonData, err := json.Marshal(questionData)
	if err != nil {
		return fmt.Errorf("failed to marshal question data: %w", err)
	}

	url := fmt.Sprintf("%s/rest/v1/questions?id=eq.%d", os.Getenv("SUPABASE_URL"), question.ID)
	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))
	req.Header.Set("Authorization", "Bearer "+userToken)

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
		return fmt.Errorf("update question failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// Delete は問題を削除（RLS適用のためユーザートークンを使用）
func (r *QuestionRepositoryImpl) Delete(ctx context.Context, id int64, userToken string) error {
	url := fmt.Sprintf("%s/rest/v1/questions?id=eq.%d", os.Getenv("SUPABASE_URL"), id)
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))
	req.Header.Set("Authorization", "Bearer "+userToken)

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
		return fmt.Errorf("delete question failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetAll は全ての問題を取得
func (r *QuestionRepositoryImpl) GetAll(ctx context.Context) ([]*entities.Question, error) {
	url := os.Getenv("SUPABASE_URL") + "/rest/v1/questions"
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
		return nil, fmt.Errorf("find all questions failed with status %d: %s", resp.StatusCode, string(body))
	}

	var questionList []map[string]interface{}
	if err := json.Unmarshal(body, &questionList); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	questions := make([]*entities.Question, len(questionList))
	for i, questionData := range questionList {
		questions[i] = mapToQuestion(questionData)
	}

	return questions, nil
}

// mapToQuestion は map[string]interface{} を Question エンティティに変換
func mapToQuestion(m map[string]interface{}) *entities.Question {
	return &entities.Question{
		ID:             getInt64(m, "id"),
		GenreID:        getInt64(m, "genre_id"),
		UserID:         getString(m, "user_id"),
		Title:          getString(m, "title"),
		Body:           getString(m, "body"),
		Explanation:    getString(m, "explanation"),
		CreatedAt:      getTime(m, "created_at"),
		Views:          getInt(m, "views"),
		CorrectCount:   getInt(m, "correct_count"),
		IncorrectCount: getInt(m, "incorrect_count"),
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

// getInt は map から int を安全に取得
func getInt(m map[string]interface{}, key string) int {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case float64:
			return int(v)
		case int64:
			return int(v)
		case int:
			return v
		case string:
			if i, err := strconv.Atoi(v); err == nil {
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