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
	projectRef  string
}

func NewAuthRepository(cfg *config.Config) *AuthRepository {
	client := s.New(
		cfg.Supabase.ProjectRef,
		cfg.Supabase.AnonKey,
	)
	adminClient := client.WithToken(cfg.Supabase.ServiceRoleKey)

	return &AuthRepository{
		client:      client,
		adminClient: adminClient,
		projectRef:  cfg.Supabase.ProjectRef,
	}
}

func (a *AuthRepository) SignInWithPassword(ctx context.Context, payload auth.LoginDTO) (*types.TokenResponse, error) {
	email := fmt.Sprintf("%s@bloansbook.local", payload.StaffID)

	resp, err := a.client.Token(ctx, types.TokenRequest{
		GrantType: "password",
		Email:     email,
		Password:  payload.Password,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to sign in staff: %w", err)
	}

	return resp, nil
}

func (a *AuthRepository) AdminCreateUser(ctx context.Context, payload auth.LoginDTO) (*types.AdminCreateUserResponse, error) {
	email := fmt.Sprintf("%s@bloansbook.local", payload.StaffID)

	resp, err := a.adminClient.AdminCreateUser(ctx, types.AdminCreateUserRequest{
		Email:        email,
		Password:     &payload.Password,
		EmailConfirm: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create supabase auth for staff: %w", err)
	}

	return resp, nil
}

func (r *AuthRepository) GetUser(ctx context.Context, accessToken string) (*types.UserResponse, error) {
	user, err := r.client.WithToken(accessToken).GetUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
