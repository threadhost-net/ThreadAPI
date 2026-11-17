package httpapi

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

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
		r.Post("/auth/email/verify/confirm", emailVerifyConfirm)
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

type ctxKey string

const userIDKey ctxKey = "userID"

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Error(w, "Не авторизован", http.StatusUnauthorized)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		client := toolbox.Redis.Get("sessions")
		n, err := client.Exists(ctx, "session:"+cookie.Value).Result()
		if err != nil {
			http.Error(w, "Ошибка сервера. Обратитесь в поддержку", http.StatusInternalServerError)
			return
		}
		if n == 0 {
			http.Error(w, "Ваша сессия истекла. Войдите в учётную запись заново", http.StatusForbidden)
			return
		}

		userIDStr, _, ok := strings.Cut(cookie.Value, ":")
		if !ok || userIDStr == "" {
			http.Error(w, "Не авторизован", http.StatusUnauthorized)
			return
		}

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Не авторизован", http.StatusUnauthorized)
			return
		}
		ctx = context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func isHX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func renderSite(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("sus"))
}
