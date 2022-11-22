package interfaces

import (
	"context"

	"github.com/DmitriyV003/bonus/internal/models"
)

type PaymentService interface {
	CreateWithdrawPayment(ctx context.Context, user *models.User, amount int64, orderNumber string) (*models.Payment, error)
	CreateAccrualPayment(ctx context.Context, user *models.User, amount int64, orderNumber string) error
}
