package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	answerDto "Shittaka_back/internal/application/answer/dto"
	"Shittaka_back/internal/application/answer/usecases"
	"Shittaka_back/internal/domain/shared"
	presentationDTO "Shittaka_back/internal/presentation/dto"
)

// AnswerHandler は回答関連のHTTPハンドラー
type AnswerHandler struct {
	answerUsecase *usecases.AnswerUsecase
}

// NewAnswerHandler は新しいAnswerHandlerを作成
func NewAnswerHandler(answerUsecase *usecases.AnswerUsecase) *AnswerHandler {
	return &AnswerHandler{
		answerUsecase: answerUsecase,
	}
}

// CreateAnswerHandler は回答作成を処理
func (h *AnswerHandler) CreateAnswerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 認証トークンの取得
	userToken, err := h.extractToken(r)
	if err != nil {
		h.sendError(w, "認証が必要です", http.StatusUnauthorized)
		return
	}

	// ユーザーIDを取得
	userID, err := h.getUserIDFromToken(userToken)
	if err != nil {
		h.sendError(w, "無効なトークンです", http.StatusUnauthorized)
		return
	}

	var req presentationDTO.CreateAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// DTOの変換
	usecaseReq := answerDto.CreateAnswerRequest{
		QuestionID: req.QuestionID,
		ChoiceID:   req.ChoiceID,
	}

	answerResp, err := h.answerUsecase.CreateAnswer(r.Context(), usecaseReq, userID, userToken)
	if err != nil {
		h.handleUsecaseError(w, err)
		return
	}

	// レスポンスDTOに変換
	response := presentationDTO.AnswerResponse{
		ID:         answerResp.ID,
		UserID:     answerResp.UserID,
		QuestionID: answerResp.QuestionID,
		ChoiceID:   answerResp.ChoiceID,
		AnsweredAt: answerResp.AnsweredAt,
	}

	h.sendJSON(w, response, http.StatusCreated)
}

// ヘルパー関数

// extractToken はリクエストからトークンを抽出
func (h *AnswerHandler) extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", shared.NewDomainError("UNAUTHORIZED", "認証が必要です")
	}

	// "Bearer " プレフィックスを除去
	userToken := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		userToken = authHeader[7:]
	}

	return userToken, nil
}

// getUserIDFromToken はJWTトークンからユーザーIDを取得
func (h *AnswerHandler) getUserIDFromToken(token string) (string, error) {
	// JWTトークンを分割 (header.payload.signature)
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", shared.NewDomainError("INVALID_TOKEN", "Invalid JWT format")
	}

	// payloadをデコード
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("failed to decode JWT payload: %w", err)
	}

	// JSONとしてパース
	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", fmt.Errorf("failed to parse JWT claims: %w", err)
	}

	// subクレームからユーザーIDを取得
	if sub, ok := claims["sub"].(string); ok && sub != "" {
		return sub, nil
	}

	return "", shared.NewDomainError("INVALID_TOKEN", "User ID not found in token")
}

// handleUsecaseError はユースケースエラーを適切なHTTPエラーに変換
func (h *AnswerHandler) handleUsecaseError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case shared.ValidationError:
		h.sendError(w, e.Message, http.StatusBadRequest)
	case shared.DomainError:
		switch e.Code {
		case "NOT_FOUND":
			h.sendError(w, e.Message, http.StatusNotFound)
		case "FORBIDDEN":
			h.sendError(w, e.Message, http.StatusForbidden)
		case "UNAUTHORIZED":
			h.sendError(w, e.Message, http.StatusUnauthorized)
		default:
			h.sendError(w, e.Message, http.StatusInternalServerError)
		}
	default:
		log.Printf("Answer usecase error: %v", err)
		h.sendError(w, "Internal server error", http.StatusInternalServerError)
	}
}

// sendJSON はJSONレスポンスを送信
func (h *AnswerHandler) sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("JSON encode error: %v", err)
	}
}

// sendError はエラーレスポンスを送信
func (h *AnswerHandler) sendError(w http.ResponseWriter, message string, statusCode int) {
	response := presentationDTO.ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	}
	h.sendJSON(w, response, statusCode)
}