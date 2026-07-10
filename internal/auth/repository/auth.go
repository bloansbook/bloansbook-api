package repository

import (
	"context"
	"fmt"

	"github.com/bloansbook/bloansbook-api/internal/auth"
	"github.com/bloansbook/bloansbook-api/pkg/config"
	s "github.com/mrehanabbasi/supabase-auth-go"
	"github.com/mrehanabbasi/supabase-auth-go/types"
)

type AuthRepository struct {
	client      s.Client
	adminClient s.Client
}

func NewAuthRepository(cfg *config.Config) *AuthRepository {
	client := s.New(cfg.Supabase.ProjectRef, cfg.Supabase.AnonKey)
	return &AuthRepository{
		client:      client,
		adminClient: client.WithToken(cfg.Supabase.ServiceRoleKey),
	}
}

func (a *AuthRepository) SignInWithPassword(ctx context.Context, payload auth.LoginDTO) (*types.TokenResponse, error) {
	resp, err := a.client.Token(ctx, types.TokenRequest{
		GrantType: "password",
		Email:     fmt.Sprintf("%s@bloansbook.local", payload.StaffID),
		Password:  payload.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("invalid credentials: %w", err)
	}
	return resp, nil
}

func (a *AuthRepository) AdminCreateUser(ctx context.Context, payload auth.LoginDTO) (*types.AdminCreateUserResponse, error) {
	resp, err := a.adminClient.AdminCreateUser(ctx, types.AdminCreateUserRequest{
		Email:        fmt.Sprintf("%s@bloansbook.local", payload.StaffID),
		Password:     &payload.Password,
		EmailConfirm: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create auth user: %w", err)
	}
	return resp, nil
}

func (a *AuthRepository) GetUser(ctx context.Context, accessToken string) (*types.UserResponse, error) {
	user, err := a.client.WithToken(accessToken).GetUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}
