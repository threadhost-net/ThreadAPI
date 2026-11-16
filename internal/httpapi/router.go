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
		r.Get("/healthcheck", healthCheck)
		r.Get("/test", test)
		r.Post("/auth/register", register) // получить данные и дёрнуть /auth/email/verify
		r.Post("/auth/login", login)       // проверить имя и пароль, после чего выдать токен
		// r.Post("/auth/logout", logout)                         // выйти с текущего устройства
		// r.Post("/auth/email/verify", requestEmailVerification) // запросить письмо для подтверждения почты повторно
		// r.Post("/auth/email/verify/confirm", confirmEmail)     // подтвердить почту по коду
		// r.Post("/auth/recover", requestPasswordReset)          // запросить сброс пароля
		// r.Post("/auth/recover/confirm", confirmPasswordReset)  // подтвердить сброс пароля
		// r.Post("/user/email/new/confirm", confirmNewEmail)     // подтвердить новую почту по коду
	})

	// Приватные пути
	// r.Group(func(r chi.Router) {
	// r.Use(authMiddleware)                           // узнать userID и передать дальше. иначе 403
	// r.Get("/user", getUser)                         // получить инфу аккаунта
	// r.Post("/user/password/new", changePassword)    // установить новый пароль, сначала проверив старый
	// r.Post("/user/email/new", requestEmailChange)   // установить новую почту (вызвать /auth/email/verify)
	// r.Delete("/user/delete", deleteUser)            // удалить пользователя
	// r.Get("/user/sessions", listSessions)           // получить список сессий
	// r.Delete("/user/sessions", deleteOtherSessions) // завершить все сессии, кроме текущей
	// r.Get("/user/payments", listPayments)           // получить историю платежей
	// r.Post("/user/payments", createPayment)         // создать платёж
	// })
	// OAuth пути
	// todo

	return r
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func test(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/test.html")
}

func authMiddleware(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("sus"))
}
