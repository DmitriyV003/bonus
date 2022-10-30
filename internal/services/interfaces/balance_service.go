package interfaces

import (
	"context"
	"github.com/DmitriyV003/bonus/internal/models"
	"github.com/DmitriyV003/bonus/internal/resources"
)

type BalanceService interface {
	Balance(ctx context.Context, user *models.User) (*resources.UserBalanceResource, error)
	Withdraw(ctx context.Context, payment *models.Payment, user *models.User) error
	Accrual(ctx context.Context, payment *models.Payment, user *models.User) error
}
