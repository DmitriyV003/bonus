package interfaces

import (
	"context"
	"github.com/DmitriyV003/bonus/internal/models"
)

type Token struct {
	Value  string
	Claims map[string]interface{}
}

type AuthService interface {
	LoginByUser(user *models.User) (*Token, error)
	Login(ctx context.Context, login string, password string) (*Token, error)
	ValidateToken(token string) (bool, error)
	ParseTokenWithClaims(token *Token) error
}
