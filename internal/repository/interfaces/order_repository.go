package interfaces

import (
	"context"

	"github.com/DmitriyV003/bonus/internal/models"
)

type OrderRepository interface {
	Create(ctx context.Context, order *models.Order) (*models.Order, error)
	UpdateByID(ctx context.Context, order *models.Order) error
	GetByIDWithUser(ctx context.Context, id int64) (*models.Order, error)
	GetByNumber(ctx context.Context, number string) (*models.Order, error)
	OrdersByUser(ctx context.Context, user *models.User) ([]*models.Order, error)
	AllPending(ctx context.Context) ([]*models.Order, error)
}
