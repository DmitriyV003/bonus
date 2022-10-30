package interfaces

import (
	"context"
	"github.com/DmitriyV003/bonus/internal/models"
	"github.com/DmitriyV003/bonus/internal/requests"
)

type UserService interface {
	Create(ctx context.Context, request *requests.RegistrationRequest) (*Token, error)
	Withdraw(ctx context.Context, user *models.User, orderNumber string, sum float64) error
	AllWithdrawsByUser(ctx context.Context, user *models.User) ([]*models.Payment, error)
}
