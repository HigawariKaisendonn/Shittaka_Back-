package entities

import (
	"time"
	"Shittaka_back/internal/domain/shared"
)

// User はユーザーのドメインエンティティ
type User struct {
	ID        string
	Email     string
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser は新しいUserエンティティを作成
func NewUser(id, email, username string) *User {
	now := time.Now()
	return &User{
		ID:        id,
		Email:     email,
		Username:  username,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Validate はUserエンティティのバリデーションを行う
func (u *User) Validate() error {
	if u.Email == "" {
		return shared.NewValidationError("email", "email is required")
	}
	if u.ID == "" {
		return shared.NewValidationError("id", "id is required")
	}
	return nil
}