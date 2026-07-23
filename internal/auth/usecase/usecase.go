package usecase

import (
	"context"
	"fmt"

	"github.com/bloansbook/bloansbook-api/internal/auth"
	"github.com/bloansbook/bloansbook-api/internal/auth/repository"
	sr "github.com/bloansbook/bloansbook-api/internal/staff/repository"
	"github.com/bloansbook/bloansbook-api/internal/staff/usecase"
	"github.com/google/uuid"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
	ExpiresAt    int64
}

type AuthUsecase struct {
	authRepo     *repository.AuthRepository
	staffUsecase *usecase.StaffUsecase
	staffRepo    *sr.StaffRepository
}

func NewAuthUsecase(ar *repository.AuthRepository, sr *sr.StaffRepository, su *usecase.StaffUsecase) *AuthUsecase {
	return &AuthUsecase{
		authRepo:     ar,
		staffUsecase: su,
		staffRepo:    sr,
	}
}

func (u *AuthUsecase) Login(ctx context.Context, payload auth.LoginDTO) (*Tokens, error) {
	staff, err := u.staffRepo.GetStaffByStaffID(ctx, payload.StaffID)
	if err != nil {
		return nil, fmt.Errorf("Staff not found: %w", err)
	}

	if !staff.HasLogin {
		return nil, fmt.Errorf("Staff has no login access: %w", err)
	}

	tokenResp, err := u.authRepo.SignInWithPassword(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("Invalid credentials: %w", err)
	}

	return &Tokens{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresIn:    tokenResp.ExpiresIn,
		ExpiresAt:    tokenResp.ExpiresAt,
	}, nil
}

func (u *AuthUsecase) RefreshToken(ctx context.Context, refreshToken string) (*Tokens, error) {
	tokenResp, err := u.authRepo.RefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &Tokens{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresIn:    tokenResp.ExpiresIn,
		ExpiresAt:    tokenResp.ExpiresAt,
	}, nil
}

func (u *AuthUsecase) GetProfile(ctx context.Context, userID uuid.UUID) (*auth.ProfileDTO, error) {
	me, err := u.authRepo.GetMe(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return me, nil
}
