package entities

// genre.goはジャンルのドメインエンティティを定義

// Genre はジャンルエンティティ
type Genre struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// NewGenre は新しいGenreエンティティを作成
func NewGenre(name string) *Genre {
	return &Genre{
		Name: name,
	}
}