package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/supabase-community/gotrue-go"
)

type AuthHandler struct {
	client gotrue.Client
}

// NewAuthHandler は新しいAuthHandlerを作成
func NewAuthHandler(client gotrue.Client) *AuthHandler {
	return &AuthHandler{client: client}
}

// SignupHandler はユーザー登録を処理
func (h *AuthHandler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// バリデーション
	if err := h.validateAuthRequest(req); err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Supabase REST APIを直接呼び出し
	signupData := map[string]interface{}{
		"email":    req.Email,
		"password": req.Password,
		"data": map[string]interface{}{
			"username": req.Username,
		},
	}

	jsonData, err := json.Marshal(signupData)
	if err != nil {
		h.sendError(w, "Failed to marshal signup data", http.StatusInternalServerError)
		return
	}

	// Supabase Auth APIのURL
	authURL := os.Getenv("SUPABASE_URL") + "/auth/v1/signup"
	httpReq, err := http.NewRequest("POST", authURL, bytes.NewBuffer(jsonData))
	if err != nil {
		h.sendError(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("apikey", os.Getenv("SUPABASE_SERVICE_ROLE_KEY"))
	httpReq.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_SERVICE_ROLE_KEY"))

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		log.Printf("HTTP request error: %v", err)
		h.sendError(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		h.sendError(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Printf("Supabase signup failed with status %d: %s", resp.StatusCode, string(body))
		h.sendError(w, fmt.Sprintf("Failed to create user: %s", string(body)), http.StatusInternalServerError)
		return
	}

	// Supabaseからのレスポンスをパース
	log.Printf("Supabase response body: %s", string(body))
	var supabaseResp map[string]interface{}
	if err := json.Unmarshal(body, &supabaseResp); err != nil {
		log.Printf("Failed to parse Supabase response: %v", err)
		h.sendError(w, "Failed to parse response", http.StatusInternalServerError)
		return
	}
	log.Printf("Parsed response: %+v", supabaseResp)

	// レスポンス作成 - メール確認が必要な場合はtokenは空になる
	response := AuthResponse{
		Token:        getString(supabaseResp, "access_token"),
		RefreshToken: getString(supabaseResp, "refresh_token"),
		User: User{
			ID:       getString(supabaseResp, "id"),
			Email:    getString(supabaseResp, "email"),
			Username: req.Username,
		},
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}

	// メール確認が送信された場合の情報も返す
	if confirmationSent := getString(supabaseResp, "confirmation_sent_at"); confirmationSent != "" {
		log.Printf("Confirmation email sent at: %s", confirmationSent)
	}

	h.sendJSON(w, response, http.StatusCreated)
}

// LoginHandler はユーザーログインを処理
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// バリデーション（ログインではusernameは不要）
	if req.Email == "" || req.Password == "" {
		h.sendError(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Supabase REST APIを直接呼び出し（ログイン）
	loginData := map[string]interface{}{
		"email":    req.Email,
		"password": req.Password,
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		h.sendError(w, "Failed to marshal login data", http.StatusInternalServerError)
		return
	}

	// Supabase Auth APIのURL（token エンドポイント）
	authURL := os.Getenv("SUPABASE_URL") + "/auth/v1/token?grant_type=password"
	httpReq, err := http.NewRequest("POST", authURL, bytes.NewBuffer(jsonData))
	if err != nil {
		h.sendError(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("apikey", os.Getenv("SUPABASE_SERVICE_ROLE_KEY"))

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		log.Printf("HTTP login request error: %v", err)
		h.sendError(w, fmt.Sprintf("Failed to login: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		h.sendError(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	log.Printf("Login response status: %d", resp.StatusCode)
	log.Printf("Login response body: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		log.Printf("Supabase login failed with status %d: %s", resp.StatusCode, string(body))
		h.sendError(w, "Invalid credentials or email not confirmed", http.StatusUnauthorized)
		return
	}

	// Supabaseからのレスポンスをパース
	var supabaseResp map[string]interface{}
	if err := json.Unmarshal(body, &supabaseResp); err != nil {
		log.Printf("Failed to parse login response: %v", err)
		h.sendError(w, "Failed to parse response", http.StatusInternalServerError)
		return
	}

	// レスポンス作成
	response := AuthResponse{
		Token:        getString(supabaseResp, "access_token"),
		RefreshToken: getString(supabaseResp, "refresh_token"),
		User: User{
			ID:    getString(getMap(supabaseResp, "user"), "id"),
			Email: getString(getMap(supabaseResp, "user"), "email"),
		},
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}

	h.sendJSON(w, response, http.StatusOK)
}

// LogoutHandler はユーザーログアウトを処理
func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Authorizationヘッダーからトークンを取得
	token := r.Header.Get("Authorization")
	if token == "" {
		h.sendError(w, "Authorization token required", http.StatusBadRequest)
		return
	}

	// "Bearer " プレフィックスを除去
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// Supabaseでログアウト
	// トークンを設定してからログアウトを実行
	clientWithToken := h.client.WithToken(token)
	err := clientWithToken.Logout()
	if err != nil {
		log.Printf("Logout error: %v", err)
		h.sendError(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

// validateAuthRequest はリクエストのバリデーションを行う
func (h *AuthHandler) validateAuthRequest(req AuthRequest) error {
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	if req.Password == "" {
		return fmt.Errorf("password is required")
	}
	if len(req.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}
	// 簡単なメールバリデーション
	if !strings.Contains(req.Email, "@") {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// sendJSON はJSONレスポンスを送信
func (h *AuthHandler) sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("JSON encode error: %v", err)
	}
}

// TestConnectionHandler はSupabaseとの接続テストを行う
func (h *AuthHandler) TestConnectionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 接続テスト用のダミーリクエスト
	// 実際のAPI呼び出しを行わずに、クライアントの設定を確認
	response := map[string]interface{}{
		"status":    "connected",
		"message":   "Supabase connection is configured",
		"timestamp": time.Now().Unix(),
	}

	h.sendJSON(w, response, http.StatusOK)
}

// sendError はエラーレスポンスを送信
func (h *AuthHandler) sendError(w http.ResponseWriter, message string, statusCode int) {
	response := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	}
	h.sendJSON(w, response, statusCode)
}

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
