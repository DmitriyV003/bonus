package services

import (
	"context"
	"errors"
	"github.com/DmitriyV003/bonus/cmd/gophermart/container"
	"github.com/DmitriyV003/bonus/cmd/gophermart/models"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

type AuthService struct {
	container *container.Container
	secret    string
}

type Token struct {
	Value string
}

func NewAuthService(container *container.Container, secret string) *AuthService {
	return &AuthService{
		container: container,
		secret:    secret,
	}
}

func (myself *AuthService) LoginByUser(user *models.User) (*Token, error) {
	token, err := myself.generateJwt(user)

	return token, err
}

func (myself *AuthService) Login(login string, password string) (*Token, error) {
	user, err := myself.container.Users.GetByLogin(context.Background(), login)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	token, err := myself.generateJwt(user)

	return token, err
}

func (myself *AuthService) ValidateToken(token string) (bool, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("d")
		}

		return []byte(myself.secret), nil
	})

	return parsedToken.Valid, err
}

func (myself *AuthService) generateJwt(user *models.User) (*Token, error) {
	token := jwt.New(jwt.SigningMethodHS512)

	claims := jwt.MapClaims{}
	claims["exp"] = time.Now().Add(24 * 60 * 5 * time.Minute).Unix()
	claims["iat"] = time.Now().Unix()
	claims["sub"] = strconv.FormatInt(user.Id, 10)
	token.Claims = claims

	genToken, err := token.SignedString([]byte(myself.secret))
	if err != nil {
		return nil, err
	}

	return &Token{Value: genToken}, nil
}
