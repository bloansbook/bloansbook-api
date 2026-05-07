package handler

import (
	"github.com/bloansbook/bloansbook-api/internal/auth"
	"github.com/bloansbook/bloansbook-api/internal/auth/usecase"
	"github.com/bloansbook/bloansbook-api/internal/models/staff"
	"github.com/bloansbook/bloansbook-api/pkg/response"
	"github.com/bloansbook/bloansbook-api/pkg/sysmsg"
	"github.com/gofiber/fiber/v3"
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

	staffMember, token, err := h.usecase.Login(c.Context(), payload)
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusUnauthorized)
	}

	var rolesList []string = make([]string, 0)
	for _, r := range staffMember.Roles {
		rolesList = append(rolesList, r.Name)
	}

	resp := auth.LoginResponse{
		AccessToken: token,
		Staff: staff.StaffSummary{
			StaffID:    staffMember.StaffID,
			FirstName:  staffMember.FirstName,
			LastName:   staffMember.LastName,
			Email:      staffMember.Email,
			Phone:      staffMember.Phone,
			Department: staffMember.Department,
			JobTitle:   staffMember.JobTitle,
			Status:     staffMember.Status,
		},
		Roles: rolesList,
	}
	return response.Success(c, sysmsg.LoginSuccess, resp, fiber.StatusOK)
}
