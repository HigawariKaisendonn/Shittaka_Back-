package dto

// SignUpRequest はサインアップリクエストのDTO
type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

// SignInRequest はサインインリクエストのDTO
type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse は認証レスポンスのDTO
type AuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	User         UserDTO `json:"user"`
	ExpiresAt    int64  `json:"expires_at"`
}

// UserDTO はユーザー情報のDTO
type UserDTO struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}