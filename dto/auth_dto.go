package dto

type AuthRegisterRequest struct {
	FullName             string   `json:"full_name" binding:"required"`
	Email                string   `json:"email" binding:"required"`
	Phone                string   `json:"phone" binding:"required"`
	Password             string   `json:"password" binding:"required"`
	PasswordConfirmation string   `json:"password_confirmation" binding:"required"`
	Role                 string   `json:"role" binding:"required"`
	Latitude             *float64 `json:"latitude,omitempty"`
	Longitude            *float64 `json:"longitude,omitempty"`
	Address              *string  `json:"address,omitempty"`
}

type AuthLoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthRefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AuthForgotPasswordRequest struct {
	Email string `json:"email" binding:"required"`
}

type AuthResetPasswordRequest struct {
	Token                string `json:"token" binding:"required"`
	Password             string `json:"password" binding:"required"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required"`
}
