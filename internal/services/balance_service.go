package services

import (
	"context"
	"fmt"
	"github.com/DmitriyV003/bonus/internal/models"
	"github.com/DmitriyV003/bonus/internal/repository"
	"github.com/DmitriyV003/bonus/internal/resources"
)

type BalanceService struct {
	payments *repository.PaymentRepository
	users    *repository.UserRepository
}

func NewBalanceService(payments *repository.PaymentRepository, users *repository.UserRepository) *BalanceService {
	return &BalanceService{
		payments: payments,
		users:    users,
	}
}

func (bs *BalanceService) Balance(user *models.User) (*resources.UserBalanceResource, error) {
	withdrawn, err := bs.payments.WithdrawnAmountByUser(context.Background(), user)
	if err != nil {
		return nil, fmt.Errorf("error to withdraw from balance: %w", err)
	}

	resource := resources.NewUserBalanceResource(user.Balance, withdrawn)

	return resource, nil
}

func (bs *BalanceService) Withdraw(payment *models.Payment, user *models.User) error {
	balance := user.Balance
	user.Balance = balance - payment.Amount
	err := bs.users.UpdateBalance(context.Background(), user)
	if err != nil {
		return fmt.Errorf("error to update user balance: %w", err)
	}

	return nil
}

func (bs *BalanceService) Accrual(payment *models.Payment, user *models.User) error {
	balance := user.Balance
	user.Balance = balance + payment.Amount
	err := bs.users.UpdateBalance(context.Background(), user)
	if err != nil {
		return fmt.Errorf("error to update user balance: %w", err)
	}

	return nil
}
