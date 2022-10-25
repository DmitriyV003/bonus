package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/DmitriyV003/bonus/internal/applicationerrors"
	"github.com/DmitriyV003/bonus/internal/models"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
	"time"
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
	err := orders.db.QueryRow(ctx, sql, order.Number, order.Amount, order.Status, order.User.ID, order.CreatedAt).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("unable to insert order to db: %w", err)
	}

	order, err = orders.GetByIDWithUser(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error to get order by id with user: %w", err)
	}

	return order, nil
}

func (orders *OrderRepository) UpdateByID(ctx context.Context, order *models.Order) error {
	sql := `UPDATE orders SET amount = $1, status = $2, updated_at = $3 WHERE id = $4`

	_, err := orders.db.Exec(ctx, sql, order.Amount, order.Status, time.Now(), order.ID)
	if err != nil {
		return fmt.Errorf("unable to update order in db: %w", err)
	}

	return nil
}

func (orders *OrderRepository) GetByIDWithUser(ctx context.Context, id int64) (*models.Order, error) {
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
		&order.ID,
		&order.Number,
		&order.Status,
		&order.Amount,
		&order.CreatedAt,
		&order.UpdatedAt,
		&user.ID,
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
		return nil, applicationerrors.ErrNotFound
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
		&order.ID,
		&order.Number,
		&order.Status,
		&order.Amount,
		&user.ID,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	order.User = &user

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, applicationerrors.ErrNotFound
	}

	return &order, nil
}

func (orders *OrderRepository) OrdersByUser(ctx context.Context, user *models.User) ([]*models.Order, error) {
	sql := `SELECT number, status, amount, created_at FROM orders WHERE user_id = $1`
	var selectedOrders []*models.Order

	rows, err := orders.db.Query(ctx, sql, user.ID)
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

func (orders *OrderRepository) AllPending(ctx context.Context) ([]*models.Order, error) {
	sql := `SELECT id, number, user_id, status, created_at FROM orders WHERE status = $1 OR status = $2`
	countSQL := `SELECT count(*) as count FROM orders WHERE status = $1 OR status = $2`
	var count int64

	err := orders.db.QueryRow(ctx, countSQL, models.NewStatus, models.ProcessingStatus).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("error to get pending orders: %w", err)
	}

	if count == 0 {
		return make([]*models.Order, 0), nil
	}

	rows, err := orders.db.Query(ctx, sql, models.NewStatus, models.ProcessingStatus)
	if err != nil {
		return nil, fmt.Errorf("error to get pending orders: %w", err)
	}
	defer rows.Close()

	selectedOrders := make([]*models.Order, 0, count)
	for rows.Next() {
		var order models.Order
		user := models.User{}
		err = rows.Scan(&order.ID, &order.Number, &user.ID, &order.Status, &order.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error to scan pending orders: %w", err)
		}
		order.User = &user

		selectedOrders = append(selectedOrders, &order)
	}

	return selectedOrders, nil
}
