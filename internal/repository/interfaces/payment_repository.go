package interfaces

import (
	"context"

	"github.com/DmitriyV003/bonus/internal/models"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *models.Payment) (*models.Payment, error)
	GetByIDWithUser(ctx context.Context, id int64) (*models.Payment, error)
	WithdrawnAmountByUser(ctx context.Context, user *models.User) (int64, error)
	GetWithdrawsByUser(ctx context.Context, user *models.User) ([]*models.Payment, error)
}
