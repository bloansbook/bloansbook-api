package auth

import (
	"github.com/bloansbook/bloansbook-api/internal/models/staff"
)

type LoginDTO struct {
	StaffID  string `json:"staffId" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken string             `json:"accessToken"`
	Staff       staff.StaffSummary `json:"staff"`
	Roles       []string           `json:"roles"`
}

// TODO: Setup Login with staffID and Password with supabase auth
// TODO: Return Response which includes the token - should save the role as part of the payload
// TODO: Check the payload to see know role and role_permission
// TODO: Authorise based on the role and permissions

// TODO: Setup resend to send email to newly added users, their password and staffID so that they can login
