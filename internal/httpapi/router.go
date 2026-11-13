package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/threadhost-net/ThreadAPI/internal/applications"
)

func InitRouter(toolbox *applications.ToolBox, domain string) http.Handler {

	// Настройки роутера
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://" + domain, "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Публичные пути
	r.Group(func(r chi.Router) {
		r.Get("/", toolbox.HealthCheck)
		r.Post("/auth/register", toolbox.Register)                      // auth
		r.Post("/auth/login", toolbox.Login)                            // auth
		r.Post("/auth/logout", toolbox.Logout)                          // auth
		r.Post("/auth/confirm_mail", toolbox.ConfirmEmail)              // auth
		r.Post("/auth/resend_verification", toolbox.ResendVerification) // auth

		r.Post("/user/forgot_password", toolbox.ForgotPassword)                // user
		r.Post("/user/forgot_password/confirm", toolbox.ConfirmForgotPassword) // user
		r.Post("/user/change_new_email/confirm", toolbox.ConfirmNewEmail)      // user
	})

	// Приватные пути
	r.Group(func(r chi.Router) {
		r.Use(toolbox.AuthMiddleware)                                  // точно ли в toolbox?
		r.Get("/user/info", toolbox.UserInfo)                          // user
		r.Post("/user/change_new_password", toolbox.ChangeNewPassword) // user
		r.Post("/user/change_new_email", toolbox.ChangeNewEmail)       // user
		r.Delete("/user/delete", toolbox.UserDelete)                   // user
	})

	return r
}
