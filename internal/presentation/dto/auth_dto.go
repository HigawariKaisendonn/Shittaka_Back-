package dto

// HTTPリクエスト/レスポンス用のDTO

// AuthRequest は認証リクエストのHTTP DTO
type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username,omitempty"`
}

// AuthResponse は認証レスポンスのHTTP DTO
type AuthResponse struct {
	Token        string  `json:"token"`
	RefreshToken string  `json:"refresh_token"`
	User         UserDTO `json:"user"`
	ExpiresAt    int64   `json:"expires_at"`
}

// UserDTO はユーザー情報のHTTP DTO
type UserDTO struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// ErrorResponse はエラーレスポンスのHTTP DTO
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}