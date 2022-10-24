package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/DmitriyV003/bonus/internal/application_errors"
	"github.com/DmitriyV003/bonus/internal/models"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		db: pool,
	}
}

func (orders *OrderRepository) Create(ctx context.Context, order *models.Order) (*models.Order, error) {
	sql := `INSERT INTO orders (number, amount, status, user_id, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	log.Info().Fields(map[string]interface{}{
		"order": order,
	}).Msg("Creating order in db")
	var id int64
	err := orders.db.QueryRow(ctx, sql, order.Number, order.Amount, order.Status, order.User.Id, order.CreatedAt).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("unable to insert order to db: %w", err)
	}

	order, err = orders.GetByIdWithUser(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error to get order by id with user: %w", err)
	}

	return order, nil
}

func (orders *OrderRepository) GetByIdWithUser(ctx context.Context, id int64) (*models.Order, error) {
	sql := `SELECT 
       orders.id, 
       orders.number,
       orders.status,
       orders.amount,
       orders.created_at, 
       orders.updated_at, 
       u.id, 
       u.login, 
       u.balance,
       u.created_at, 
       u.updated_at FROM orders JOIN users u on u.id = orders.user_id WHERE orders.id = $1`
	var user models.User
	var order models.Order

	row := orders.db.QueryRow(ctx, sql, id)
	err := row.Scan(
		&order.Id,
		&order.Number,
		&order.Status,
		&order.Amount,
		&order.CreatedAt,
		&order.UpdatedAt,
		&user.Id,
		&user.Login,
		&user.Balance,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	order.User = &user
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, application_errors.ErrNotFound
	}

	return &order, nil
}

func (orders *OrderRepository) GetByNumber(ctx context.Context, number string) (*models.Order, error) {
	sql := `SELECT 
       orders.id, 
       orders.number,
       orders.status,
       orders.amount,
       orders.user_id,
       orders.created_at, 
       orders.updated_at
       FROM orders WHERE orders.number = $1`
	var order models.Order
	var user models.User

	row := orders.db.QueryRow(ctx, sql, number)
	err := row.Scan(
		&order.Id,
		&order.Number,
		&order.Status,
		&order.Amount,
		&user.Id,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	order.User = &user

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, application_errors.ErrNotFound
	}

	return &order, nil
}

func (orders *OrderRepository) OrdersByUser(ctx context.Context, user *models.User) ([]*models.Order, error) {
	sql := `SELECT number, status, amount, created_at FROM orders WHERE user_id = $1`
	var selectedOrders []*models.Order

	rows, err := orders.db.Query(ctx, sql, user.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.Number, &order.Status, &order.Amount, &order.CreatedAt)
		if err != nil {
			return nil, err
		}

		selectedOrders = append(selectedOrders, &order)
	}

	return selectedOrders, nil
}
