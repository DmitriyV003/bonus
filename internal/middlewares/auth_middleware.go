package middlewares

import (
	"context"
	"github.com/DmitriyV003/bonus/internal/config"
	"github.com/DmitriyV003/bonus/internal/container"
	"github.com/DmitriyV003/bonus/internal/services"
	"net/http"
	"strconv"
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
			if !isValid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			newToken := services.Token{
				Value:  token[1],
				Claims: map[string]interface{}{},
			}
			err = authService.ParseTokenWithClaims(&newToken)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			parsedUserId, err := strconv.ParseInt(newToken.Claims["user_id"].(string), 10, 64)
			user, err := container.Users.GetById(context.Background(), parsedUserId)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			services.SetLoggedInUser(user)

			next.ServeHTTP(w, r.WithContext(context.WithValue(context.Background(), "user", user)))
		})
	}
}
