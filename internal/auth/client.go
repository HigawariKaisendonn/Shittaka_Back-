package auth

import (
	"os"
	"strings"

	"github.com/supabase-community/gotrue-go"
)

// NewClient は gotrue クライアントを返す
func NewClient() gotrue.Client {
	baseURL := os.Getenv("SUPABASE_URL")

	// URLの末尾にスラッシュがある場合は除去
	baseURL = strings.TrimSuffix(baseURL, "/")

	// /auth/v1を追加
	authURL := baseURL + "/auth/v1"

	return gotrue.New(
		authURL,
		os.Getenv("SUPABASE_SERVICE_ROLE_KEY"),
	)
}
