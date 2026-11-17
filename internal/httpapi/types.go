package httpapi

import (
	"time"

	"github.com/threadhost-net/ThreadAPI/session"
)

type registerRequest struct {
	Username         string `json:"username"`
	Password         string `json:"password"`
	Email            string `json:"email"`
	HCaptchaResponse string `json:"h-captcha-response"`
}

type loginRequest struct {
	Username         string `json:"username"`
	Password         string `json:"password"`
	HCaptchaResponse string `json:"h-captcha-response"`
}

type emailRequest struct {
	Email            string `json:"email"`
	HCaptchaResponse string `json:"h-captcha-response"`
}

type confirmForgotPasswordRequest struct {
	NewPassword string `json:"new_password"`
}

type changeNewPasswordRequest struct {
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
}

type changeNewEmailRequest struct {
	NewEmail string `json:"new_email"`
	Password string `json:"password"`
}

type confirmNewEmailRequest struct {
	UserID   int64  `json:"user_id"`
	NewEmail string `json:"new_email"`
}

type accountInfoResponse struct {
	ID          int64                 `json:"id"`
	Username    string                `json:"username"`
	Email       string                `json:"email"`
	Balance     int64                 `json:"balance"`
	IsActivated bool                  `json:"is_activated"`
	CreatedAt   time.Time             `json:"created_at"`
	Sessions    []session.SessionInfo `json:"sessions"`
}
