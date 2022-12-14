package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/DmitriyV003/bonus/internal/models"
	"github.com/DmitriyV003/bonus/internal/repository/interfaces"
	serviceinterfaces "github.com/DmitriyV003/bonus/internal/services/interfaces"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

var loggedInUser *models.User

func SetLoggedInUser(user *models.User) {
	loggedInUser = user
}

func GetLoggedInUser() *models.User {
	return loggedInUser
}

type AuthService struct {
	users  interfaces.UserRepository
	secret string
}

func NewAuthService(secret string, users interfaces.UserRepository) *AuthService {
	return &AuthService{
		secret: secret,
		users:  users,
	}
}

func (myself *AuthService) LoginByUser(user *models.User) (*serviceinterfaces.Token, error) {
	token, err := myself.generateJwt(user)
	if err != nil {
		return nil, fmt.Errorf("error to generate jwt to login user: %w", err)
	}

	return token, nil
}

func (myself *AuthService) Login(ctx context.Context, login string, password string) (*serviceinterfaces.Token, error) {
	user, err := myself.users.GetByLogin(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("error to get user by login: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("error to compare passwords: %w", err)
	}

	token, err := myself.generateJwt(user)
	if err != nil {
		return nil, fmt.Errorf("error to generate jwt token: %w", err)
	}

	return token, nil
}

func (myself *AuthService) ValidateToken(token string) (bool, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("error to validate token")
		}

		return []byte(myself.secret), nil
	})

	return parsedToken.Valid, err
}

func (myself *AuthService) ParseTokenWithClaims(token *serviceinterfaces.Token) error {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(myself.secret), nil
	})
	if err != nil {
		return fmt.Errorf("error to parse jwt token: %w", err)
	}

	token.Claims = claims
	return nil
}

func (myself *AuthService) generateJwt(user *models.User) (*serviceinterfaces.Token, error) {
	token := jwt.New(jwt.SigningMethodHS512)

	claims := jwt.MapClaims{}
	claims["exp"] = time.Now().Add(24 * 60 * 5 * time.Minute).Unix()
	claims["iat"] = time.Now().Unix()
	claims["sub"] = strconv.FormatInt(user.ID, 10)
	claims["user_id"] = strconv.FormatInt(user.ID, 10)
	token.Claims = claims

	genToken, err := token.SignedString([]byte(myself.secret))
	if err != nil {
		return nil, fmt.Errorf("error sign token: %w", err)
	}

	return &serviceinterfaces.Token{Value: genToken}, nil
}
