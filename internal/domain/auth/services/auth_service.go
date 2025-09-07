package services

import (
	"context"
	"fmt"
	"strings"
	"Shittaka_back/internal/domain/auth/repositories"
	"Shittaka_back/internal/domain/shared"
)

// AuthService は認証に関するドメインサービス
type AuthService struct {
	userRepo repositories.UserRepository
}

// NewAuthService は新しいAuthServiceを作成
func NewAuthService(userRepo repositories.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// SignUp はユーザー登録を行う
func (s *AuthService) SignUp(ctx context.Context, email, password, username string) (*repositories.AuthResult, error) {
	// バリデーション
	if err := s.validateSignUpInput(email, password, username); err != nil {
		return nil, err
	}
	
	// 既存ユーザーチェック
	existingUser, _ := s.userRepo.FindByEmail(ctx, email)
	if existingUser != nil {
		return nil, shared.NewDomainError("USER_EXISTS", "user with this email already exists")
	}
	
	// ユーザー作成
	metadata := map[string]interface{}{
		"username": username,
	}
	
	_, err := s.userRepo.Create(ctx, email, password, metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	
	// 認証（トークン取得）
	authResult, err := s.userRepo.Authenticate(ctx, email, password)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate after signup: %w", err)
	}
	
	return authResult, nil
}

// SignIn はユーザーログインを行う
func (s *AuthService) SignIn(ctx context.Context, email, password string) (*repositories.AuthResult, error) {
	// バリデーション
	if err := s.validateSignInInput(email, password); err != nil {
		return nil, err
	}
	
	// 認証
	authResult, err := s.userRepo.Authenticate(ctx, email, password)
	if err != nil {
		return nil, shared.NewDomainError("AUTH_FAILED", "invalid credentials or email not confirmed")
	}
	
	return authResult, nil
}

// SignOut はユーザーログアウトを行う
func (s *AuthService) SignOut(ctx context.Context, token string) error {
	if token == "" {
		return shared.NewValidationError("token", "token is required")
	}
	
	// "Bearer " プレフィックスを除去
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}
	
	return s.userRepo.Logout(ctx, token)
}

// validateSignUpInput はサインアップ入力をバリデート
func (s *AuthService) validateSignUpInput(email, password, username string) error {
	if email == "" {
		return shared.NewValidationError("email", "email is required")
	}
	if password == "" {
		return shared.NewValidationError("password", "password is required")
	}
	if len(password) < 6 {
		return shared.NewValidationError("password", "password must be at least 6 characters")
	}
	if username == "" {
		return shared.NewValidationError("username", "username is required")
	}
	if !strings.Contains(email, "@") {
		return shared.NewValidationError("email", "invalid email format")
	}
	
	return nil
}

// validateSignInInput はサインイン入力をバリデート
func (s *AuthService) validateSignInInput(email, password string) error {
	if email == "" {
		return shared.NewValidationError("email", "email is required")
	}
	if password == "" {
		return shared.NewValidationError("password", "password is required")
	}
	
	return nil
}