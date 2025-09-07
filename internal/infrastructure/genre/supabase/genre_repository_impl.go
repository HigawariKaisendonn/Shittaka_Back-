package supabase

// genre_repository_impl.goはSupabaseを使用したGenreRepositoryの実装

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"Shittaka_back/internal/domain/genre/entities"
	"Shittaka_back/internal/domain/genre/repositories"
	"Shittaka_back/internal/domain/shared"
)

// GenreRepositoryImpl はSupabaseを使用したGenreRepositoryの実装
type GenreRepositoryImpl struct{}

// NewGenreRepository は新しいGenreRepositoryImplを作成
func NewGenreRepository() repositories.GenreRepository {
	return &GenreRepositoryImpl{}
}

// Create は新しいジャンルを作成
func (r *GenreRepositoryImpl) Create(ctx context.Context, genre *entities.Genre) (*entities.Genre, error) {
	genreData := map[string]interface{}{
		"name": genre.Name,
	}

	jsonData, err := json.Marshal(genreData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal genre data: %w", err)
	}

	url := os.Getenv("SUPABASE_URL") + "/rest/v1/ganres"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", os.Getenv("SUPABASE_SERVICE_ROLE_KEY"))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_SERVICE_ROLE_KEY"))
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
		return nil, fmt.Errorf("create genre failed with status %d: %s", resp.StatusCode, string(body))
	}

	var genreList []map[string]interface{}
	if err := json.Unmarshal(body, &genreList); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(genreList) == 0 {
		return nil, fmt.Errorf("no genre returned from create operation")
	}

	genreResp := genreList[0]
	return &entities.Genre{
		ID:   getInt64(genreResp, "id"),
		Name: getString(genreResp, "name"),
	}, nil
}

// FindByID はIDでジャンルを検索
func (r *GenreRepositoryImpl) FindByID(ctx context.Context, id int64) (*entities.Genre, error) {
	url := fmt.Sprintf("%s/rest/v1/ganres?id=eq.%d", os.Getenv("SUPABASE_URL"), id)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("apikey", os.Getenv("SUPABASE_SERVICE_ROLE_KEY"))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_SERVICE_ROLE_KEY"))

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
		return nil, fmt.Errorf("find genre failed with status %d: %s", resp.StatusCode, string(body))
	}

	var genreList []map[string]interface{}
	if err := json.Unmarshal(body, &genreList); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(genreList) == 0 {
		return nil, shared.NewDomainError("NOT_FOUND", "ジャンルが見つかりません")
	}

	genreData := genreList[0]
	return &entities.Genre{
		ID:   getInt64(genreData, "id"),
		Name: getString(genreData, "name"),
	}, nil
}

// FindAll は全てのジャンルを取得
func (r *GenreRepositoryImpl) FindAll(ctx context.Context) ([]*entities.Genre, error) {
	url := os.Getenv("SUPABASE_URL") + "/rest/v1/ganres"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("apikey", os.Getenv("SUPABASE_SERVICE_ROLE_KEY"))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_SERVICE_ROLE_KEY"))

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
		return nil, fmt.Errorf("find all genres failed with status %d: %s", resp.StatusCode, string(body))
	}

	var genreList []map[string]interface{}
	if err := json.Unmarshal(body, &genreList); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	genres := make([]*entities.Genre, len(genreList))
	for i, genreData := range genreList {
		genres[i] = &entities.Genre{
			ID:   getInt64(genreData, "id"),
			Name: getString(genreData, "name"),
		}
	}

	return genres, nil
}

// FindByName は名前でジャンルを検索
func (r *GenreRepositoryImpl) FindByName(ctx context.Context, name string) (*entities.Genre, error) {
	url := fmt.Sprintf("%s/rest/v1/ganres?name=eq.%s", os.Getenv("SUPABASE_URL"), name)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("apikey", os.Getenv("SUPABASE_SERVICE_ROLE_KEY"))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_SERVICE_ROLE_KEY"))

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
		return nil, fmt.Errorf("find genre by name failed with status %d: %s", resp.StatusCode, string(body))
	}

	var genreList []map[string]interface{}
	if err := json.Unmarshal(body, &genreList); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(genreList) == 0 {
		return nil, shared.NewDomainError("NOT_FOUND", "ジャンルが見つかりません")
	}

	genreData := genreList[0]
	return &entities.Genre{
		ID:   getInt64(genreData, "id"),
		Name: getString(genreData, "name"),
	}, nil
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