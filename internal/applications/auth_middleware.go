package applications

import "net/http"

// Заглушка, пропускает всех. Перенести в httpapi
func (t *Toolbox) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// пока пропускаем — все "приватные" роуты сейчас доступны
		next.ServeHTTP(w, r)
	})
}
