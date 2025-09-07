package usecases

// genre_usecase.goはジャンル関連のユースケースを定義

import (
	"context"
	"strings"

	"Shittaka_back/internal/application/genre/dto"
	"Shittaka_back/internal/domain/genre/entities"
	"Shittaka_back/internal/domain/genre/repositories"
	"Shittaka_back/internal/domain/shared"
)

// GenreUsecase はジャンルユースケース
type GenreUsecase struct {
	genreRepo repositories.GenreRepository
}

// NewGenreUsecase は新しいGenreUsecaseを作成
func NewGenreUsecase(genreRepo repositories.GenreRepository) *GenreUsecase {
	return &GenreUsecase{
		genreRepo: genreRepo,
	}
}

// CreateGenre は新しいジャンルを作成する
func (u *GenreUsecase) CreateGenre(ctx context.Context, req dto.CreateGenreRequest) (*dto.GenreResponse, error) {
	// バリデーション
	if err := u.validateCreateGenreRequest(req); err != nil {
		return nil, err
	}

	// 同名のジャンルが既に存在するかチェック
	existingGenre, err := u.genreRepo.FindByName(ctx, req.Name)
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}
	if existingGenre != nil {
		return nil, shared.NewDomainError("GENRE_EXISTS", "ジャンルが既に存在します")
	}

	// ジャンルエンティティを作成
	genre := entities.NewGenre(req.Name)

	// リポジトリに保存
	createdGenre, err := u.genreRepo.Create(ctx, genre)
	if err != nil {
		return nil, err
	}

	// レスポンスDTOに変換
	return &dto.GenreResponse{
		ID:   createdGenre.ID,
		Name: createdGenre.Name,
	}, nil
}

// validateCreateGenreRequest はジャンル作成リクエストをバリデーション
func (u *GenreUsecase) validateCreateGenreRequest(req dto.CreateGenreRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return shared.NewValidationError("ジャンル名は必須です")
	}

	if len(req.Name) > 50 {
		return shared.NewValidationError("ジャンル名は50文字以内で入力してください")
	}

	return nil
}

// isNotFoundError はエラーがNot Foundエラーかどうかを判定
func isNotFoundError(err error) bool {
	if domainErr, ok := err.(shared.DomainError); ok {
		return domainErr.Code == "NOT_FOUND"
	}
	return false
}