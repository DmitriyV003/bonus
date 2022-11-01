package interfaces

import (
	"context"
	"github.com/DmitriyV003/bonus/internal/models"
)

type OrderService interface {
	Create(ctx context.Context, user *models.User, orderNumber string) (*models.Order, error)
	OrdersByUser(ctx context.Context, user *models.User) ([]*models.Order, error)
	PollPendingOrders(ctx context.Context) error
}
