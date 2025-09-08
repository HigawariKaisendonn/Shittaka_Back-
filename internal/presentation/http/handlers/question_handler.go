package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	questionDto "Shittaka_back/internal/application/question/dto"
	"Shittaka_back/internal/application/question/usecases"
	"Shittaka_back/internal/domain/shared"
	presentationDTO "Shittaka_back/internal/presentation/dto"
)

// QuestionHandler は問題関連のHTTPハンドラー
type QuestionHandler struct {
	questionUsecase *usecases.QuestionUsecase
}

// NewQuestionHandler は新しいQuestionHandlerを作成
func NewQuestionHandler(questionUsecase *usecases.QuestionUsecase) *QuestionHandler {
	return &QuestionHandler{
		questionUsecase: questionUsecase,
	}
}

// CreateQuestionHandler は問題作成を処理
func (h *QuestionHandler) CreateQuestionHandler(w http.ResponseWriter, r *http.Request) {
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

	var req presentationDTO.CreateQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// DTOの変換
	usecaseReq := questionDto.CreateQuestionRequest{
		GenreID:     req.GenreID,
		Title:       req.Title,
		Body:        req.Body,
		Explanation: req.Explanation,
	}

	questionResp, err := h.questionUsecase.CreateQuestion(r.Context(), usecaseReq, userID, userToken)
	if err != nil {
		h.handleUsecaseError(w, err)
		return
	}

	// レスポンスDTOに変換
	response := presentationDTO.QuestionResponse{
		ID:             questionResp.ID,
		GenreID:        questionResp.GenreID,
		UserID:         questionResp.UserID,
		Title:          questionResp.Title,
		Body:           questionResp.Body,
		Explanation:    questionResp.Explanation,
		CreatedAt:      questionResp.CreatedAt,
		Views:          questionResp.Views,
		CorrectCount:   questionResp.CorrectCount,
		IncorrectCount: questionResp.IncorrectCount,
	}

	h.sendJSON(w, response, http.StatusCreated)
}

// UpdateQuestionHandler は問題更新を処理
func (h *QuestionHandler) UpdateQuestionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
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

	// URLから問題IDを取得
	questionID, err := h.getQuestionIDFromPath(r.URL.Path)
	if err != nil {
		h.sendError(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	var req presentationDTO.UpdateQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// DTOの変換
	usecaseReq := questionDto.UpdateQuestionRequest{
		Title:       req.Title,
		Body:        req.Body,
		Explanation: req.Explanation,
	}

	err = h.questionUsecase.UpdateQuestion(r.Context(), questionID, usecaseReq, userID, userToken)
	if err != nil {
		h.handleUsecaseError(w, err)
		return
	}

	h.sendJSON(w, map[string]string{"message": "問題が正常に更新されました"}, http.StatusOK)
}

// DeleteQuestionHandler は問題削除を処理
func (h *QuestionHandler) DeleteQuestionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
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

	// URLから問題IDを取得
	questionID, err := h.getQuestionIDFromPath(r.URL.Path)
	if err != nil {
		h.sendError(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	err = h.questionUsecase.DeleteQuestion(r.Context(), questionID, userID, userToken)
	if err != nil {
		h.handleUsecaseError(w, err)
		return
	}

	h.sendJSON(w, map[string]string{"message": "問題が正常に削除されました"}, http.StatusOK)
}

// GetQuestionHandler は問題取得を処理
func (h *QuestionHandler) GetQuestionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URLから問題IDを取得
	questionID, err := h.getQuestionIDFromPath(r.URL.Path)
	if err != nil {
		h.sendError(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	questionResp, err := h.questionUsecase.GetQuestion(r.Context(), questionID)
	if err != nil {
		h.handleUsecaseError(w, err)
		return
	}

	// レスポンスDTOに変換
	response := presentationDTO.QuestionResponse{
		ID:             questionResp.ID,
		GenreID:        questionResp.GenreID,
		UserID:         questionResp.UserID,
		Title:          questionResp.Title,
		Body:           questionResp.Body,
		Explanation:    questionResp.Explanation,
		CreatedAt:      questionResp.CreatedAt,
		Views:          questionResp.Views,
		CorrectCount:   questionResp.CorrectCount,
		IncorrectCount: questionResp.IncorrectCount,
	}

	h.sendJSON(w, response, http.StatusOK)
}

// GetQuestionsHandler は問題一覧取得を処理
func (h *QuestionHandler) GetQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	questionResp, err := h.questionUsecase.GetAllQuestions(r.Context())
	if err != nil {
		h.handleUsecaseError(w, err)
		return
	}

	// レスポンスDTOに変換
	responses := make([]presentationDTO.QuestionResponse, len(questionResp))
	for i, q := range questionResp {
		responses[i] = presentationDTO.QuestionResponse{
			ID:             q.ID,
			GenreID:        q.GenreID,
			UserID:         q.UserID,
			Title:          q.Title,
			Body:           q.Body,
			Explanation:    q.Explanation,
			CreatedAt:      q.CreatedAt,
			Views:          q.Views,
			CorrectCount:   q.CorrectCount,
			IncorrectCount: q.IncorrectCount,
		}
	}

	h.sendJSON(w, responses, http.StatusOK)
}

// GetMyQuestionsHandler はユーザーの問題一覧取得を処理
func (h *QuestionHandler) GetMyQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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

	questionResp, err := h.questionUsecase.GetQuestionsByUser(r.Context(), userID, userToken)
	if err != nil {
		h.handleUsecaseError(w, err)
		return
	}

	// レスポンスDTOに変換
	responses := make([]presentationDTO.QuestionResponse, len(questionResp))
	for i, q := range questionResp {
		responses[i] = presentationDTO.QuestionResponse{
			ID:             q.ID,
			GenreID:        q.GenreID,
			UserID:         q.UserID,
			Title:          q.Title,
			Body:           q.Body,
			Explanation:    q.Explanation,
			CreatedAt:      q.CreatedAt,
			Views:          q.Views,
			CorrectCount:   q.CorrectCount,
			IncorrectCount: q.IncorrectCount,
		}
	}

	h.sendJSON(w, responses, http.StatusOK)
}

// ヘルパー関数

// extractToken はリクエストからトークンを抽出
func (h *QuestionHandler) extractToken(r *http.Request) (string, error) {
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
func (h *QuestionHandler) getUserIDFromToken(token string) (string, error) {
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

// getQuestionIDFromPath はURLパスから問題IDを取得
func (h *QuestionHandler) getQuestionIDFromPath(path string) (int64, error) {
	// "/api/questions/{id}" の形式から ID を取得
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		return 0, shared.NewDomainError("INVALID_PATH", "Invalid question ID path")
	}

	idStr := parts[len(parts)-1]
	questionID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, shared.NewDomainError("INVALID_ID", "Invalid question ID format")
	}

	return questionID, nil
}

// handleUsecaseError はユースケースエラーを適切なHTTPエラーに変換
func (h *QuestionHandler) handleUsecaseError(w http.ResponseWriter, err error) {
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
		log.Printf("Question usecase error: %v", err)
		h.sendError(w, "Internal server error", http.StatusInternalServerError)
	}
}

// sendJSON はJSONレスポンスを送信
func (h *QuestionHandler) sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("JSON encode error: %v", err)
	}
}

// sendError はエラーレスポンスを送信
func (h *QuestionHandler) sendError(w http.ResponseWriter, message string, statusCode int) {
	response := presentationDTO.ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	}
	h.sendJSON(w, response, statusCode)
}