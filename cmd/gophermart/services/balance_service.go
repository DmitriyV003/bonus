package services

import (
	"github.com/DmitriyV003/bonus/cmd/gophermart/container"
	"github.com/DmitriyV003/bonus/cmd/gophermart/models"
	"github.com/DmitriyV003/bonus/cmd/gophermart/resources"
)

type BalanceService struct {
	container *container.Container
	user      *models.User
}

func NewBalanceService(container *container.Container, user *models.User) *BalanceService {
	return &BalanceService{
		container: container,
		user:      user,
	}
}

func (bs *BalanceService) Balance() (*resources.UserBalanceResource, error) {
	resource := resources.NewUserBalanceResource(float64(bs.user.Balance)/10000, 0)

	return resource, nil
}
