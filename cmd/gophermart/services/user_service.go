package services

import (
	"context"
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

func NewUserService(container *container.Container) *UserService {
	return &UserService{container: container}
}

func (u *UserService) Create(request *requests.RegistrationRequest) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(request.Password), 14)
	if err != nil {
		return err
	}

	user := models.User{
		Login:     request.Login,
		Password:  string(bytes),
		CreatedAt: time.Now(),
	}
	err = u.container.Users.Create(context.Background(), &user)

	return err
}
