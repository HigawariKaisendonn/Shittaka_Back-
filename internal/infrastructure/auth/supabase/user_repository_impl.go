package supabase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"Shittaka_back/internal/domain/auth/entities"
	"Shittaka_back/internal/domain/auth/repositories"

	"github.com/supabase-community/gotrue-go"
)

// UserRepositoryImpl はSupabaseを使用したUserRepositoryの実装
type UserRepositoryImpl struct {
	client gotrue.Client
}

// NewUserRepository は新しいUserRepositoryImplを作成
func NewUserRepository() *UserRepositoryImpl {
	baseURL := os.Getenv("SUPABASE_URL")
	baseURL = strings.TrimSuffix(baseURL, "/")
	authURL := baseURL + "/auth/v1"

	client := gotrue.New(
		authURL,
		os.Getenv("SUPABASE_SERVICE_ROLE_KEY"),
	)

	return &UserRepositoryImpl{
		client: client,
	}
}

// Create は新しいユーザーを作成
func (r *UserRepositoryImpl) Create(ctx context.Context, email, password string, metadata map[string]interface{}) (*entities.User, error) {
	signupData := map[string]interface{}{
		"email":    email,
		"password": password,
		"data":     metadata,
	}

	jsonData, err := json.Marshal(signupData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal signup data: %w", err)
	}

	authURL := os.Getenv("SUPABASE_URL") + "/auth/v1/signup"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", authURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("apikey", os.Getenv("SUPABASE_SERVICE_ROLE_KEY"))
	httpReq.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_SERVICE_ROLE_KEY"))

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("signup failed with status %d: %s", resp.StatusCode, string(body))
	}

	var supabaseResp map[string]interface{}
	if err := json.Unmarshal(body, &supabaseResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	username := ""
	if metadata != nil {
		if u, ok := metadata["username"].(string); ok {
			username = u
		}
	}

	user := entities.NewUser(
		getString(supabaseResp, "id"),
		getString(supabaseResp, "email"),
		username,
	)

	return user, nil
}

// Authenticate はユーザーの認証を行い、トークンを返す
func (r *UserRepositoryImpl) Authenticate(ctx context.Context, email, password string) (*repositories.AuthResult, error) {
	loginData := map[string]interface{}{
		"email":    email,
		"password": password,
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal login data: %w", err)
	}

	authURL := os.Getenv("SUPABASE_URL") + "/auth/v1/token?grant_type=password"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", authURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("apikey", os.Getenv("SUPABASE_SERVICE_ROLE_KEY"))

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		// メール確認エラーの特別な処理
		if resp.StatusCode == 400 && strings.Contains(string(body), "email_not_confirmed") {
			return nil, fmt.Errorf("email confirmation required: please check your email and click the confirmation link")
		}
		return nil, fmt.Errorf("authentication failed with status %d: %s", resp.StatusCode, string(body))
	}

	var supabaseResp map[string]interface{}
	if err := json.Unmarshal(body, &supabaseResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	userMap := getMap(supabaseResp, "user")
	user := entities.NewUser(
		getString(userMap, "id"),
		getString(userMap, "email"),
		"", // username は別途取得が必要
	)

	return &repositories.AuthResult{
		User:         user,
		AccessToken:  getString(supabaseResp, "access_token"),
		RefreshToken: getString(supabaseResp, "refresh_token"),
		ExpiresAt:    time.Now().Add(time.Hour * 24).Unix(),
	}, nil
}

// FindByID はIDでユーザーを検索
func (r *UserRepositoryImpl) FindByID(ctx context.Context, id string) (*entities.User, error) {
	// Supabase Admin APIを使用してユーザーを取得
	// 実装は省略（必要に応じて追加）
	return nil, fmt.Errorf("not implemented")
}

// FindByEmail はEmailでユーザーを検索
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	// Supabase Admin APIを使用してユーザーを検索
	// 実装は省略（必要に応じて追加）
	return nil, fmt.Errorf("not implemented")
}

// Update はユーザー情報を更新
func (r *UserRepositoryImpl) Update(ctx context.Context, user *entities.User) error {
	return fmt.Errorf("not implemented")
}

// Delete はユーザーを削除
func (r *UserRepositoryImpl) Delete(ctx context.Context, id string) error {
	return fmt.Errorf("not implemented")
}

// Logout はユーザーをログアウトさせる
func (r *UserRepositoryImpl) Logout(ctx context.Context, token string) error {
	clientWithToken := r.client.WithToken(token)
	return clientWithToken.Logout()
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

// getMap は map から map を安全に取得
func getMap(m map[string]interface{}, key string) map[string]interface{} {
	if val, ok := m[key]; ok {
		if mapVal, ok := val.(map[string]interface{}); ok {
			return mapVal
		}
	}
	return make(map[string]interface{})
}
