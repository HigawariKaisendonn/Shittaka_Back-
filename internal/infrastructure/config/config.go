package config

import (
	"log"
	"os"
	
	"github.com/joho/godotenv"
)

// Config はアプリケーションの設定を保持
type Config struct {
	SupabaseURL        string
	SupabaseServiceKey string
	Port               string
}

// LoadConfig は設定を読み込む
func LoadConfig() *Config {
	// 環境変数を読み込み
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}
	
	// 必要な環境変数をチェック
	supabaseURL := os.Getenv("SUPABASE_URL")
	if supabaseURL == "" {
		log.Fatal("SUPABASE_URL is required")
	}
	
	supabaseServiceKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")
	if supabaseServiceKey == "" {
		log.Fatal("SUPABASE_SERVICE_ROLE_KEY is required")
	}
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8088"
	}
	
	return &Config{
		SupabaseURL:        supabaseURL,
		SupabaseServiceKey: supabaseServiceKey,
		Port:               port,
	}
}