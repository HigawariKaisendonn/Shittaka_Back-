package dto

// genre_dto.goはジャンル関連のデータ転送オブジェクトを定義

// CreateGenreRequest はジャンル作成リクエスト
type CreateGenreRequest struct {
	Name string `json:"name"`
}

// GenreResponse はジャンルレスポンス
type GenreResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}