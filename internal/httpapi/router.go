package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/threadhost-net/ThreadAPI/internal/applications"
	"github.com/threadhost-net/ThreadAPI/platform"
)

var toolbox *applications.Toolbox

func InitRouter(stack *platform.Stack, domain string) http.Handler {
	toolbox = &applications.Toolbox{Stack: stack}

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
		r.Get("/", HealthCheck)
		r.Post("/auth/register", Register)                      // auth
		r.Post("/auth/login", Login)                            // auth
		r.Post("/auth/logout", Logout)                          // auth
		r.Post("/auth/confirm_mail", ConfirmEmail)              // auth
		r.Post("/auth/resend_verification", ResendVerification) // auth

		r.Post("/user/forgot_password", ForgotPassword)                // user
		r.Post("/user/forgot_password/confirm", ConfirmForgotPassword) // user
		r.Post("/user/change_new_email/confirm", ConfirmNewEmail)      // user
	})

	// Приватные пути
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware)                                  // точно ли в toolbox?
		r.Get("/user/info", UserInfo)                          // user
		r.Post("/user/change_new_password", ChangeNewPassword) // user
		r.Post("/user/change_new_email", ChangeNewEmail)       // user
		r.Delete("/user/delete", UserDelete)                   // user
	})

	return r
}
