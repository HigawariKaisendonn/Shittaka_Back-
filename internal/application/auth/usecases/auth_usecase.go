package usecases

// auth_usecase.goは認証に関するユースケースを定義

import (
	"Shittaka_back/internal/application/auth/dto"
	"Shittaka_back/internal/domain/auth/entities"
	"Shittaka_back/internal/domain/auth/repositories"
	"Shittaka_back/internal/domain/auth/services"
	"context"
)

// AuthUsecase は認証に関するユースケース
type AuthUsecase struct {
	authService *services.AuthService
}

// NewAuthUsecase は新しいAuthUsecaseを作成
func NewAuthUsecase(authService *services.AuthService) *AuthUsecase {
	return &AuthUsecase{
		authService: authService,
	}
}

// SignUp はユーザー登録ユースケース
func (u *AuthUsecase) SignUp(ctx context.Context, req dto.SignUpRequest) (*dto.AuthResponse, error) {
	authResult, err := u.authService.SignUp(ctx, req.Email, req.Password, req.Username)
	if err != nil {
		return nil, err
	}

	return u.toAuthResponse(authResult), nil
}

// SignIn はユーザーログインユースケース
func (u *AuthUsecase) SignIn(ctx context.Context, req dto.SignInRequest) (*dto.AuthResponse, error) {
	authResult, err := u.authService.SignIn(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return u.toAuthResponse(authResult), nil
}

// SignOut はユーザーログアウトユースケース
func (u *AuthUsecase) SignOut(ctx context.Context, token string) error {
	return u.authService.SignOut(ctx, token)
}

// toAuthResponse はドメインのAuthResultをDTOに変換
func (u *AuthUsecase) toAuthResponse(authResult *repositories.AuthResult) *dto.AuthResponse {
	return &dto.AuthResponse{
		Token:        authResult.AccessToken,
		RefreshToken: authResult.RefreshToken,
		User:         u.toUserDTO(authResult.User),
		ExpiresAt:    authResult.ExpiresAt,
	}
}

// toUserDTO はUserエンティティをDTOに変換
func (u *AuthUsecase) toUserDTO(user *entities.User) dto.UserDTO {
	return dto.UserDTO{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}
}
