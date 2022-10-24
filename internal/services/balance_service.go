package services

import (
	"context"
	"github.com/DmitriyV003/bonus/internal/container"
	"github.com/DmitriyV003/bonus/internal/models"
	"github.com/DmitriyV003/bonus/internal/resources"
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
	withdrawn, err := bs.container.Payments.WithdrawnAmountByUser(context.Background(), bs.user)
	if err != nil {
		return nil, err
	}

	resource := resources.NewUserBalanceResource(bs.user.Balance, withdrawn)

	return resource, nil
}

func (bs *BalanceService) Withdraw(payment *models.Payment) error {
	balance := bs.user.Balance
	bs.user.Balance = balance - payment.Amount
	err := bs.container.Users.UpdateBalance(context.Background(), bs.user)
	if err != nil {
		return err
	}

	return nil
}

func (bs *BalanceService) Accrual(payment *models.Payment) error {
	balance := bs.user.Balance
	bs.user.Balance = balance + payment.Amount
	err := bs.container.Users.UpdateBalance(context.Background(), bs.user)
	if err != nil {
		return err
	}

	return nil
}
