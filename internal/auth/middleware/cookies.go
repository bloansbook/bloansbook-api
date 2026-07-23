package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
)

const (
	AccessTokenCookie  = "access_token"
	RefreshTokenCookie = "refresh_token"
)

func SetAuthCookies(c fiber.Ctx, accessToken, refreshToken string, accessExp, refreshExp time.Time) {
	c.Cookie(&fiber.Cookie{
		Name:     AccessTokenCookie,
		Value:    accessToken,
		Expires:  accessExp,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Path:     "/",
		Domain:   "",
	})

	c.Cookie(&fiber.Cookie{
		Name:     RefreshTokenCookie,
		Value:    refreshToken,
		Expires:  refreshExp,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Path:     "/api/v1/auth/refresh",
		Domain:   "",
	})
}

func ClearAuthCookies(c fiber.Ctx) {
	expired := time.Now().Add(-1 * time.Hour)

	c.Cookie(&fiber.Cookie{
		Name:     AccessTokenCookie,
		Value:    "",
		Expires:  expired,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Path:     "/",
		Domain:   "",
	})

	c.Cookie(&fiber.Cookie{
		Name:     RefreshTokenCookie,
		Value:    "",
		Expires:  expired,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Path:     "/api/v1/auth/refresh",
		Domain:   "",
	})
}

func GetRefreshToken(c fiber.Ctx) string {
	return c.Cookies(RefreshTokenCookie)
}

func GetAccessToken(c fiber.Ctx) string {
	return c.Cookies(AccessTokenCookie)
}
