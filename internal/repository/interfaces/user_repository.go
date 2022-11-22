package interfaces

import (
	"context"

	"github.com/DmitriyV003/bonus/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByLogin(ctx context.Context, login string) (*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
	UpdateBalance(ctx context.Context, user *models.User) error
}
