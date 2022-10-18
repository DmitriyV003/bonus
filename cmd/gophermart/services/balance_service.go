package services

import (
	"context"
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
	withdrawn, err := bs.container.Payments.WithdrawnAmountByUser(context.Background(), bs.user)
	if err != nil {
		return nil, err
	}

	resource.Withdrawn = float64(withdrawn) / 10000

	return resource, nil
}

func (bs *BalanceService) Withdraw(payment *models.Payment) error {
	balance := bs.user.Balance * 10000
	bs.user.Balance = balance - payment.Amount*10000
	bs.user.Balance /= 10000
	err := bs.container.Users.UpdateBalance(context.Background(), bs.user)
	if err != nil {
		return err
	}

	return nil
}

func (bs *BalanceService) Accrual(payment *models.Payment) error {
	balance := bs.user.Balance * 10000
	bs.user.Balance = balance + payment.Amount*10000
	bs.user.Balance /= 10000
	err := bs.container.Users.UpdateBalance(context.Background(), bs.user)
	if err != nil {
		return err
	}

	return nil
}
