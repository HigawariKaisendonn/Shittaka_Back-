package dto

// ganres_dto.goはジャンル関連のHTTP DTOを定義

// CreateGenreRequest はジャンル作成リクエストのHTTP DTO
type CreateGenreRequest struct {
	Name string `json:"name"`
}

// GenreResponse はジャンルレスポンスのHTTP DTO
type GenreResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
