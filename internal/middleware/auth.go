package middleware

import (
    "context"
    "log"
    "net/http"
    "strings"

    "booking/internal/utils"
)

// Используем простые строковые константы
const (
    UserIDKey = "userID"
    RoleKey   = "role"
)

// Auth middleware проверяет JWT и сохраняет userID и role в контекст
func Auth(jwtSecret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing authorization header")
                return
            }
            parts := strings.Split(authHeader, " ")
            if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
                writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid authorization header format")
                return
            }
            token := parts[1]
            claims, err := utils.ValidateJWT(token, jwtSecret)
            if err != nil {
                writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid token")
                return
            }
            // Сохраняем в контекст с ключами-строками
            ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
            ctx = context.WithValue(ctx, RoleKey, claims.Role)
            log.Printf("Auth: userID=%s, role=%s", claims.UserID, claims.Role)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func Logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}

func writeError(w http.ResponseWriter, status int, code, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    w.Write([]byte(`{"error":{"code":"` + code + `","message":"` + message + `"}}`))
}