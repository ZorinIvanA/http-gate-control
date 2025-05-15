package middleware

import (
	"net/http"
)

// BasicAuthMiddleware добавляет защиту через Basic Auth
func BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем логин и пароль из заголовка
		username, password, ok := r.BasicAuth()

		// Проверяем, что авторизация предоставлена и логин/пароль верны
		if !ok || username != "admin" || password != "password" {
			// Устанавливаем заголовок WWW-Authenticate
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
			// Возвращаем ошибку 401
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Если авторизация успешна, продолжаем обработку
		next.ServeHTTP(w, r)
	})
}
