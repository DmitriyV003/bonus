package services

import (
	"context"
	"fmt"
	"github.com/DmitriyV003/bonus/cmd/gophermart/container"
	"github.com/DmitriyV003/bonus/cmd/gophermart/models"
	"github.com/DmitriyV003/bonus/cmd/gophermart/requests"
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
		return nil, err
	}

	dbUser, err := u.container.Users.GetByLogin(context.Background(), request.Login)
	if err != nil {
		return nil, err
	}

	authService := NewAuthService(u.container, jwtSecret)
	token, err := authService.LoginByUser(dbUser)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return token, nil
}
