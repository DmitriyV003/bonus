package services

import (
	"context"
	"fmt"

	"github.com/DmitriyV003/bonus/internal/models"
	"github.com/DmitriyV003/bonus/internal/repository/interfaces"
	"github.com/DmitriyV003/bonus/internal/resources"
)

type BalanceService struct {
	payments interfaces.PaymentRepository
	users    interfaces.UserRepository
}

func NewBalanceService(payments interfaces.PaymentRepository, users interfaces.UserRepository) *BalanceService {
	return &BalanceService{
		payments: payments,
		users:    users,
	}
}

func (bs *BalanceService) Balance(ctx context.Context, user *models.User) (*resources.UserBalanceResource, error) {
	withdrawn, err := bs.payments.WithdrawnAmountByUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("error to withdraw from balance: %w", err)
	}

	resource := resources.NewUserBalanceResource(user.Balance, withdrawn)

	return resource, nil
}

func (bs *BalanceService) Withdraw(ctx context.Context, payment *models.Payment, user *models.User) error {
	balance := user.Balance
	user.Balance = balance - payment.Amount
	err := bs.users.UpdateBalance(ctx, user)
	if err != nil {
		return fmt.Errorf("error to update user balance: %w", err)
	}

	return nil
}

func (bs *BalanceService) Accrual(ctx context.Context, payment *models.Payment, user *models.User) error {
	balance := user.Balance
	user.Balance = balance + payment.Amount
	err := bs.users.UpdateBalance(ctx, user)
	if err != nil {
		return fmt.Errorf("error to update user balance: %w", err)
	}

	return nil
}
