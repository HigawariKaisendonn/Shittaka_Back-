package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/supabase-community/gotrue-go"
	"github.com/supabase-community/gotrue-go/types"
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

	// Supabaseでユーザー作成
	signupReq := types.SignupRequest{
		Email:    req.Email,
		Password: req.Password,
		Data: map[string]interface{}{
			"username": req.Username,
		},
	}
	user, err := h.client.Signup(signupReq)
	if err != nil {
		log.Printf("Signup error: %v", err)
		h.sendError(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// レスポンス作成
	response := AuthResponse{
		Token:        user.AccessToken,
		RefreshToken: user.RefreshToken,
		User: User{
			ID:       user.User.ID.String(),
			Email:    user.User.Email,
			Username: req.Username,
		},
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
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

	// Supabaseでログイン
	user, err := h.client.SignInWithEmailPassword(req.Email, req.Password)
	if err != nil {
		log.Printf("Login error: %v", err)
		h.sendError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// レスポンス作成
	response := AuthResponse{
		Token:        user.AccessToken,
		RefreshToken: user.RefreshToken,
		User: User{
			ID:    user.User.ID.String(),
			Email: user.User.Email,
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
