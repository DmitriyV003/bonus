package middlewares

import (
	"context"
	"github.com/DmitriyV003/bonus/internal/repository"
	"github.com/DmitriyV003/bonus/internal/services"
	"net/http"
	"strconv"
	"strings"
)

type AuthMiddleware struct {
	authService *services.AuthService
	users       *repository.UserRepository
}

func NewAuthMiddleware(authService *services.AuthService, users *repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		users:       users,
	}
}

func (m *AuthMiddleware) Pipe() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenHeader := r.Header.Get("Authorization")
			if tokenHeader == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			token := strings.Split(tokenHeader, " ")
			isValid, err := m.authService.ValidateToken(token[1])
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
			err = m.authService.ParseTokenWithClaims(&newToken)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			parsedUserId, err := strconv.ParseInt(newToken.Claims["user_id"].(string), 10, 64)
			user, err := m.users.GetById(context.Background(), parsedUserId)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			services.SetLoggedInUser(user)

			next.ServeHTTP(w, r.WithContext(context.WithValue(context.Background(), "user", user)))
		})
	}
}
