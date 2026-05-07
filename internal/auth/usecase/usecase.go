package usecase

import (
	"context"
	"fmt"

	"github.com/bloansbook/bloansbook-api/internal/auth"
	"github.com/bloansbook/bloansbook-api/internal/auth/repository"
	"github.com/bloansbook/bloansbook-api/internal/models/staff"
	sr "github.com/bloansbook/bloansbook-api/internal/staff/repository"
	"github.com/bloansbook/bloansbook-api/internal/staff/usecase"
)

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

func (u *AuthUsecase) Login(ctx context.Context, payload auth.LoginDTO) (*staff.StaffDTO, string, error) {
	staff, err := u.staffRepo.GetStaffByStaffID(ctx, payload.StaffID)
	if err != nil {
		return nil, "", fmt.Errorf("Staff not found: %w", err)
	}

	if !staff.HasLogin {
		return nil, "", fmt.Errorf("Staff has no login access: %w", err)
	}

	tokenResp, err := u.authRepo.SignInWithPassword(ctx, payload)
	if err != nil {
		return nil, "", fmt.Errorf("Invalid credentials: %w", err)
	}

	staffDTO, err := u.staffUsecase.GetStaffByID(ctx, staff.ID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch staff: %w", err)
	}

	return staffDTO, tokenResp.AccessToken, nil
}
