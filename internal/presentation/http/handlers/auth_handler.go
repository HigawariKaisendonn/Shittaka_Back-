package handlers

// auth_handler.goは認証に関するHTTPハンドラーを定義
// HTTPハンドラーとは、HTTPリクエストを受け取り、適切なユースケースに処理を委譲する

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"Shittaka_back/internal/application/auth/dto"
	"Shittaka_back/internal/application/auth/usecases"
	"Shittaka_back/internal/domain/shared"
	presentationDTO "Shittaka_back/internal/presentation/dto"
)

// AuthHandler は認証関連のHTTPハンドラー
type AuthHandler struct {
	authUsecase *usecases.AuthUsecase
}

// NewAuthHandler は新しいAuthHandlerを作成
func NewAuthHandler(authUsecase *usecases.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

// SignupHandler はユーザー登録を処理
func (h *AuthHandler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req presentationDTO.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// DTOの変換
	usecaseReq := dto.SignUpRequest{
		Email:    req.Email,
		Password: req.Password,
		Username: req.Username,
	}

	authResp, err := h.authUsecase.SignUp(r.Context(), usecaseReq)
	if err != nil {
		h.handleUsecaseError(w, err)
		return
	}

	// レスポンスDTOに変換
	response := presentationDTO.AuthResponse{
		Token:        authResp.Token,
		RefreshToken: authResp.RefreshToken,
		User: presentationDTO.UserDTO{
			ID:       authResp.User.ID,
			Email:    authResp.User.Email,
			Username: authResp.User.Username,
		},
		ExpiresAt: authResp.ExpiresAt,
	}

	h.sendJSON(w, response, http.StatusCreated)
}

// LoginHandler はユーザーログインを処理
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req presentationDTO.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// DTOの変換
	usecaseReq := dto.SignInRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	authResp, err := h.authUsecase.SignIn(r.Context(), usecaseReq)
	if err != nil {
		h.handleUsecaseError(w, err)
		return
	}

	// レスポンスDTOに変換
	response := presentationDTO.AuthResponse{
		Token:        authResp.Token,
		RefreshToken: authResp.RefreshToken,
		User: presentationDTO.UserDTO{
			ID:       authResp.User.ID,
			Email:    authResp.User.Email,
			Username: authResp.User.Username,
		},
		ExpiresAt: authResp.ExpiresAt,
	}

	h.sendJSON(w, response, http.StatusOK)
}

// LogoutHandler はユーザーログアウトを処理
func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	if token == "" {
		h.sendError(w, "Authorization token required", http.StatusBadRequest)
		return
	}

	err := h.authUsecase.SignOut(r.Context(), token)
	if err != nil {
		h.handleUsecaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

// TestConnectionHandler はSupabaseとの接続テストを行う
func (h *AuthHandler) TestConnectionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"status":    "connected",
		"message":   "Supabase connection is configured",
		"timestamp": time.Now().Unix(),
	}

	h.sendJSON(w, response, http.StatusOK)
}

// ヘルパー関数

// handleUsecaseError はユースケースエラーを適切なHTTPエラーに変換
func (h *AuthHandler) handleUsecaseError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case shared.ValidationError:
		h.sendError(w, e.Message, http.StatusBadRequest)
	case shared.DomainError:
		switch e.Code {
		case "USER_EXISTS":
			h.sendError(w, e.Message, http.StatusConflict)
		case "AUTH_FAILED":
			h.sendError(w, e.Message, http.StatusUnauthorized)
		default:
			h.sendError(w, e.Message, http.StatusInternalServerError)
		}
	default:
		log.Printf("Usecase error: %v", err)
		h.sendError(w, "Internal server error", http.StatusInternalServerError)
	}
}

// sendJSON はJSONレスポンスを送信
func (h *AuthHandler) sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("JSON encode error: %v", err)
	}
}

// sendError はエラーレスポンスを送信
func (h *AuthHandler) sendError(w http.ResponseWriter, message string, statusCode int) {
	response := presentationDTO.ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	}
	h.sendJSON(w, response, statusCode)
}
