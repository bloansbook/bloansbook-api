package auth

type LoginDTO struct {
	StaffID  string `json:"staffId" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken string   `json:"accessToken"`
	Staff       string   `json:"staff"`
	Permissions []string `json:"permissions"`
}
