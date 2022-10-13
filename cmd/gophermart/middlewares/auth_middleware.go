package middlewares

import (
	"context"
	"github.com/DmitriyV003/bonus/cmd/gophermart/config"
	"github.com/DmitriyV003/bonus/cmd/gophermart/container"
	"github.com/DmitriyV003/bonus/cmd/gophermart/services"
	"net/http"
	"strings"
)

func AuthMiddleware(container *container.Container, conf *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenHeader := r.Header.Get("Authorization")
			if tokenHeader == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			token := strings.Split(tokenHeader, " ")
			authService := services.NewAuthService(container, conf.JwtSecret)
			isValid, err := authService.ValidateToken(token[1])
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if isValid {
				next.ServeHTTP(w, r.WithContext(context.Background()))
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		})
	}
}
