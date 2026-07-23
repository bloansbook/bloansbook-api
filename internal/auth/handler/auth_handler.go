package handler

import (
	"time"

	"github.com/bloansbook/bloansbook-api/internal/auth"
	"github.com/bloansbook/bloansbook-api/internal/auth/middleware"
	"github.com/bloansbook/bloansbook-api/internal/auth/usecase"
	"github.com/bloansbook/bloansbook-api/pkg/response"
	"github.com/bloansbook/bloansbook-api/pkg/sysmsg"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type AuthHandler struct {
	usecase *usecase.AuthUsecase
}

func NewAuthHandler(u *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		usecase: u,
	}
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var payload auth.LoginDTO

	if err := c.Bind().Body(&payload); err != nil {
		return response.Error(c, sysmsg.BadRequest, fiber.StatusBadRequest)
	}

	tokens, err := h.usecase.Login(c.Context(), payload)
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusUnauthorized)
	}

	accessExp := time.Unix(tokens.ExpiresAt, 0)
	refreshExp := time.Now().Add(7 * 24 * 60 * 60 * 1000 * 1000 * 1000) // 7 days

	middleware.SetAuthCookies(c, tokens.AccessToken, tokens.RefreshToken, accessExp, refreshExp)
	return response.Success(c, sysmsg.LoginSuccess, nil, fiber.StatusOK)
}

func (h *AuthHandler) RefreshToken(c fiber.Ctx) error {
	refreshToken := middleware.GetRefreshToken(c)
	if refreshToken == "" {
		return response.Error(c, sysmsg.Unauthorized, fiber.StatusUnauthorized)
	}

	tokens, err := h.usecase.RefreshToken(c.Context(), refreshToken)
	if err != nil {
		return response.Error(c, sysmsg.TokenInvalid, fiber.StatusUnauthorized)
	}

	accessExp := time.Unix(tokens.ExpiresAt, 0)
	refreshExp := time.Now().Add(7 * 24 * 60 * 60 * 1000 * 1000 * 1000)

	middleware.ClearAuthCookies(c)
	middleware.SetAuthCookies(c, tokens.AccessToken, tokens.RefreshToken, accessExp, refreshExp)
	return response.Success(c, "", nil, fiber.StatusOK)
}

func (h *AuthHandler) GetProfile(c fiber.Ctx) error {
	staffID, ok := middleware.CallerStaffID(c)
	if !ok {
		return response.Error(c, sysmsg.Unauthorized, fiber.StatusUnauthorized)
	}

	id, ok := staffID.(uuid.UUID)
	if !ok {
		return response.Error(c, sysmsg.Unauthorized, fiber.StatusUnauthorized)
	}

	profile, err := h.usecase.GetProfile(c.Context(), id)
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusUnauthorized)
	}

	return response.Success(c, sysmsg.ProfileFetched, profile, fiber.StatusOK)
}
