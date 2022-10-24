package services

import (
	"context"
	"fmt"
	"github.com/DmitriyV003/bonus/internal/container"
	"github.com/DmitriyV003/bonus/internal/models"
	"github.com/DmitriyV003/bonus/internal/requests"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserService struct {
	container *container.Container
	user      *models.User
}

func NewUserService(container *container.Container, user *models.User) *UserService {
	return &UserService{
		container: container,
		user:      user,
	}
}

func (u *UserService) Create(request *requests.RegistrationRequest, jwtSecret string) (*Token, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(request.Password), 14)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Login:     request.Login,
		Password:  string(bytes),
		CreatedAt: time.Now(),
	}
	err = u.container.Users.Create(context.Background(), &user)
	if err != nil {
		return nil, fmt.Errorf("unable to create user in db: %w", err)
	}

	dbUser, err := u.container.Users.GetByLogin(context.Background(), request.Login)
	if err != nil {
		return nil, fmt.Errorf("error to get user by login: %w", err)
	}

	authService := NewAuthService(u.container, jwtSecret)
	token, err := authService.LoginByUser(dbUser)
	if err != nil {
		return nil, fmt.Errorf("login user programmly: %w", err)
	}

	return token, nil
}
