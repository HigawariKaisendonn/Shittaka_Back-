package handlers

// choice_handler.goは選択肢に関するHTTPハンドラーを定義

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"Shittaka_back/internal/domain/choices/entities"
	"Shittaka_back/internal/domain/choices/services"
	"Shittaka_back/internal/domain/shared"
	presentationDTO "Shittaka_back/internal/presentation/dto"
)

// ChoiceHandler は選択肢関連のHTTPハンドラー
type ChoiceHandler struct {
	choiceService *services.ChoiceService
}

// NewChoiceHandler は新しいChoiceHandlerを作成
func NewChoiceHandler(choiceService *services.ChoiceService) *ChoiceHandler {
	return &ChoiceHandler{
		choiceService: choiceService,
	}
}

// GetChoicesHandler は問題IDに紐づく選択肢を取得
func (h *ChoiceHandler) GetChoicesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URLから問題IDを取得 (/api/choices/{questionID})
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		h.sendError(w, "Question ID is required", http.StatusBadRequest)
		return
	}

	questionID, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		h.sendError(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	choices, err := h.choiceService.GetChoices(r.Context(), questionID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	// レスポンスDTOに変換
	var choiceResponses []presentationDTO.ChoiceResponse
	for _, choice := range choices {
		choiceResponses = append(choiceResponses, presentationDTO.ChoiceResponse{
			ID:         choice.ID,
			QuestionID: choice.QuestionID,
			Text:       choice.Text,
			IsCorrect:  choice.IsCorrect,
		})
	}

	response := presentationDTO.ChoicesResponse{
		Choices: choiceResponses,
	}

	h.sendJSON(w, response, http.StatusOK)
}

// CreateChoiceHandler は新しい選択肢を作成
func (h *ChoiceHandler) CreateChoiceHandler(w http.ResponseWriter, r *http.Request) {
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
	
	// デバッグ用ログ
	if len(userToken) > 20 {
		log.Printf("Choice Handler - UserToken: %s...", userToken[:20])
	} else {
		log.Printf("Choice Handler - UserToken: %s", userToken)
	}

	var req presentationDTO.CreateChoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// エンティティに変換
	choice := entities.Choice{
		QuestionID: req.QuestionID,
		Text:       req.Text,
		IsCorrect:  req.IsCorrect,
	}

	createdChoice, err := h.choiceService.CreateChoiceWithAuth(r.Context(), choice, userToken)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	// レスポンスDTOに変換
	response := presentationDTO.ChoiceResponse{
		ID:         createdChoice.ID,
		QuestionID: createdChoice.QuestionID,
		Text:       createdChoice.Text,
		IsCorrect:  createdChoice.IsCorrect,
	}

	h.sendJSON(w, response, http.StatusCreated)
}

// UpdateChoiceHandler は既存の選択肢を更新
func (h *ChoiceHandler) UpdateChoiceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req presentationDTO.UpdateChoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// エンティティに変換
	choice := entities.Choice{
		ID:         req.ID,
		QuestionID: req.QuestionID,
		Text:       req.Text,
		IsCorrect:  req.IsCorrect,
	}

	updatedChoice, err := h.choiceService.UpdateChoice(r.Context(), choice)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	// レスポンスDTOに変換
	response := presentationDTO.ChoiceResponse{
		ID:         updatedChoice.ID,
		QuestionID: updatedChoice.QuestionID,
		Text:       updatedChoice.Text,
		IsCorrect:  updatedChoice.IsCorrect,
	}

	h.sendJSON(w, response, http.StatusOK)
}

// DeleteChoiceHandler は選択肢を削除
func (h *ChoiceHandler) DeleteChoiceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URLから選択肢IDを取得 (/api/choices/{id}/delete)
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		h.sendError(w, "Choice ID is required", http.StatusBadRequest)
		return
	}

	choiceID, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		h.sendError(w, "Invalid choice ID", http.StatusBadRequest)
		return
	}

	if err := h.choiceService.DeleteChoice(r.Context(), choiceID); err != nil {
		h.handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ヘルパー関数

// extractToken はリクエストからトークンを抽出
func (h *ChoiceHandler) extractToken(r *http.Request) (string, error) {
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

// handleServiceError はサービスエラーを適切なHTTPエラーに変換
func (h *ChoiceHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case shared.ValidationError:
		h.sendError(w, e.Message, http.StatusBadRequest)
	case shared.DomainError:
		switch e.Code {
		case "NOT_FOUND":
			h.sendError(w, e.Message, http.StatusNotFound)
		case "CHOICE_EXISTS":
			h.sendError(w, e.Message, http.StatusConflict)
		default:
			h.sendError(w, e.Message, http.StatusInternalServerError)
		}
	default:
		log.Printf("Choice service error: %v", err)
		h.sendError(w, "Internal server error", http.StatusInternalServerError)
	}
}

// sendJSON はJSONレスポンスを送信
func (h *ChoiceHandler) sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("JSON encode error: %v", err)
	}
}

// sendError はエラーレスポンスを送信
func (h *ChoiceHandler) sendError(w http.ResponseWriter, message string, statusCode int) {
	response := presentationDTO.ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	}
	h.sendJSON(w, response, statusCode)
}