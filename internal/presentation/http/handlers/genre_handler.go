package handlers

// genre_handler.goはジャンルに関するHTTPハンドラーを定義

import (
	"encoding/json"
	"log"
	"net/http"

	genreDto "Shittaka_back/internal/application/genre/dto"
	"Shittaka_back/internal/application/genre/usecases"
	"Shittaka_back/internal/domain/shared"
	presentationDTO "Shittaka_back/internal/presentation/dto"
)

// GenreHandler はジャンル関連のHTTPハンドラー
type GenreHandler struct {
	genreUsecase *usecases.GenreUsecase
}

// NewGenreHandler は新しいGenreHandlerを作成
func NewGenreHandler(genreUsecase *usecases.GenreUsecase) *GenreHandler {
	return &GenreHandler{
		genreUsecase: genreUsecase,
	}
}

// CreateGenreHandler はジャンル作成を処理
func (h *GenreHandler) CreateGenreHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 認証トークンの取得
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		h.sendError(w, "認証が必要です", http.StatusUnauthorized)
		return
	}
	
	// "Bearer " プレフィックスを除去
	userToken := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		userToken = authHeader[7:]
	}

	var req presentationDTO.CreateGenreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// DTOの変換
	usecaseReq := genreDto.CreateGenreRequest{
		Name: req.Name,
	}

	genreResp, err := h.genreUsecase.CreateGenre(r.Context(), usecaseReq, userToken)
	if err != nil {
		h.handleUsecaseError(w, err)
		return
	}

	// レスポンスDTOに変換
	response := presentationDTO.GenreResponse{
		ID:   genreResp.ID,
		Name: genreResp.Name,
	}

	h.sendJSON(w, response, http.StatusCreated)
}

// ヘルパー関数

// handleUsecaseError はユースケースエラーを適切なHTTPエラーに変換
func (h *GenreHandler) handleUsecaseError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case shared.ValidationError:
		h.sendError(w, e.Message, http.StatusBadRequest)
	case shared.DomainError:
		switch e.Code {
		case "GENRE_EXISTS":
			h.sendError(w, e.Message, http.StatusConflict)
		case "NOT_FOUND":
			h.sendError(w, e.Message, http.StatusNotFound)
		default:
			h.sendError(w, e.Message, http.StatusInternalServerError)
		}
	default:
		log.Printf("Genre usecase error: %v", err)
		h.sendError(w, "Internal server error", http.StatusInternalServerError)
	}
}

// sendJSON はJSONレスポンスを送信
func (h *GenreHandler) sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("JSON encode error: %v", err)
	}
}

// sendError はエラーレスポンスを送信
func (h *GenreHandler) sendError(w http.ResponseWriter, message string, statusCode int) {
	response := presentationDTO.ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	}
	h.sendJSON(w, response, statusCode)
}