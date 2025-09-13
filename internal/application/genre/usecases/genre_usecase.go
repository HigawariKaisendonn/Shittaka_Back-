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

// CreateGenre は新しいジャンルを作成する（認証が必要）
func (u *GenreUsecase) CreateGenre(ctx context.Context, req dto.CreateGenreRequest, userToken string) (*dto.GenreResponse, error) {
	// バリデーション
	if err := u.validateCreateGenreRequest(req); err != nil {
		return nil, err
	}

	// 同名のジャンルが既に存在するかチェック
	existingGenre, err := u.genreRepo.FindByName(ctx, req.Name, userToken)
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}
	if existingGenre != nil {
		return nil, shared.NewDomainError("GENRE_EXISTS", "ジャンルが既に存在します")
	}

	// ジャンルエンティティを作成
	genre := entities.NewGenre(req.Name)

	// リポジトリに保存（ユーザートークンを渡してRLS適用）
	createdGenre, err := u.genreRepo.Create(ctx, genre, userToken)
	if err != nil {
		return nil, err
	}

	// レスポンスDTOに変換
	return &dto.GenreResponse{
		ID:   createdGenre.ID,
		Name: createdGenre.Name,
	}, nil
}

// GetAllGenres は全てのジャンルを取得する
func (u *GenreUsecase) GetAllGenres(ctx context.Context) ([]*dto.GenreResponse, error) {
	genres, err := u.genreRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.GenreResponse, len(genres))
	for i, genre := range genres {
		responses[i] = &dto.GenreResponse{
			ID:   genre.ID,
			Name: genre.Name,
		}
	}

	return responses, nil
}

// validateCreateGenreRequest はジャンル作成リクエストをバリデーション
func (u *GenreUsecase) validateCreateGenreRequest(req dto.CreateGenreRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return shared.NewValidationError("name", "ジャンル名は必須です")
	}

	if len(req.Name) > 50 {
		return shared.NewValidationError("name", "ジャンル名は50文字以内で入力してください")
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