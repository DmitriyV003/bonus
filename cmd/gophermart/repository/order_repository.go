package repository

import (
	"context"
	"github.com/DmitriyV003/bonus/cmd/gophermart/application_errors"
	"github.com/DmitriyV003/bonus/cmd/gophermart/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: pool}
}

func (users *OrderRepository) Create(ctx context.Context, order *models.Order) error {
	sql := `INSERT INTO orders (number, status, amount, user_id, created_at) VALUES ($1, $2, $3, $4, $5)`

	//dbUser, err := users.GetByLogin(ctx, user.Login)
	//if err != nil && !errors.Is(err, application_errors.ErrNotFound) {
	//	return application_errors.ErrInternalServer
	//}

	//if dbUser != nil {
	//	return application_errors.ErrConflict
	//}

	//if errors.Is(err, pgx.ErrNoRows) {
	//	return application_errors.ErrNotFound
	//}

	_, err := users.db.Exec(ctx, sql, order.Number, order.Status, order.Amount, order.User.Id, order.CreatedAt)

	if err != nil {
		return application_errors.ErrInternalServer
	}

	return nil
}
