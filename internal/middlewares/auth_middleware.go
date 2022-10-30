package middlewares

import (
	"context"
	"github.com/DmitriyV003/bonus/internal/repository/interfaces"
	"github.com/DmitriyV003/bonus/internal/services"
	serviceinterfaces "github.com/DmitriyV003/bonus/internal/services/interfaces"
	"net/http"
	"strconv"
	"strings"
)

type AuthMiddleware struct {
	authService serviceinterfaces.AuthService
	users       interfaces.UserRepository
}

func NewAuthMiddleware(authService serviceinterfaces.AuthService, users interfaces.UserRepository) *AuthMiddleware {
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
			if len(token) < 2 {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			isValid, err := m.authService.ValidateToken(token[1])
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if !isValid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			newToken := serviceinterfaces.Token{
				Value:  token[1],
				Claims: map[string]interface{}{},
			}
			err = m.authService.ParseTokenWithClaims(&newToken)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			parsedUserID, err := strconv.ParseInt(newToken.Claims["user_id"].(string), 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			user, err := m.users.GetByID(context.Background(), parsedUserID)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			services.SetLoggedInUser(user)

			next.ServeHTTP(w, r.WithContext(context.Background()))
		})
	}
}
