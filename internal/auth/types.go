package auth

// フロントエンド↔GoAPI間で使用する型定義-----------------------------------------------

// AuthRequest は認証リクエストの構造体
type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username,omitempty"`
}

// AuthResponse は認証レスポンスの構造体
type AuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
	ExpiresAt    int64  `json:"expires_at"`
}

// User はユーザー情報の構造体
type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// Supabase↔GoAPI間で使用する型定義----------------------------------------------------

// SignupRequest はユーザー登録リクエストの構造体
type SignupRequest struct {
	Email    string                 `json:"email"`
	Password string                 `json:"password"`
	Data     map[string]interface{} `json:"data,omitempty"`
	Phone    string                 `json:"phone,omitempty"`
}

// SignInRequest はユーザーログインリクエストの構造体
type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone,omitempty"`
}

// Supabaseから返される構造体
type AuthenticatedDetails struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
	ExpiresIn    int    `json:"expires_in"`
}

// ErrorResponse はエラーレスポンスの構造体
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
