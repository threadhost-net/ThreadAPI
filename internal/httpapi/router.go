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
		r.Get("/", renderSite)
		r.Get("/healthcheck", healthCheck)
		r.Post("/auth/register", register)
		r.Post("/auth/login", login)
		r.Post("/auth/logout", logout)
		r.Post("/auth/email/verify", emailVerify)
		r.Post("/auth/email/verify/confirm", emailConfirm)
		r.Post("/auth/recover", recover)
		r.Post("/auth/recover/confirm", recoverConfirm)
		r.Post("/user/email/new/confirm", emailNewConfirm)
	})

	// Приватные пути
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Get("/user", getUser)
		r.Post("/user/password/new", passwordNew)
		r.Post("/user/email/new", emailNew)
		r.Delete("/user/delete", deleteUser)
		r.Get("/user/sessions", listSessions)
		r.Delete("/user/sessions", deleteOtherSessions)
		r.Get("/user/payments", listPayments)
		r.Post("/user/payments", createPayment)
	})
	// OAuth пути
	// todo

	return r
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func authMiddleware(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("sus"))
}

func isHX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func renderSite(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("sus"))
}
