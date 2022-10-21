package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/DmitriyV003/bonus/internal/application_errors"
	models2 "github.com/DmitriyV003/bonus/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepository struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(pool *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{
		db: pool,
	}
}

func (payments *PaymentRepository) Create(ctx context.Context, payment *models2.Payment) (*models2.Payment, error) {
	sql := `INSERT INTO payments (
    	user_id, 
        type, 
        transaction_type,
        order_number,
        amount,
        created_at
    ) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	var id int64
	err := payments.db.QueryRow(
		ctx,
		sql,
		payment.User.Id,
		payment.Type,
		payment.TransactionType,
		payment.OrderNumber,
		payment.Amount,
		payment.CreatedAt,
	).Scan(&id)
	if err != nil {
		return nil, application_errors.ErrInternalServer
	}

	payment, err = payments.GetByIdWithUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (payments *PaymentRepository) GetByIdWithUser(ctx context.Context, id int64) (*models2.Payment, error) {
	sql := `SELECT 
       p.id, 
       p.type, 
       p.transaction_type,
       p.order_number,
       p.amount,
       p.created_at,
       p.updated_at,
       u.id,
       u.login,
       u.balance,
       u.created_at,
       u.updated_at
	FROM payments as p LEFT JOIN users u on u.id = p.user_id WHERE p.id = $1`

	var payment models2.Payment
	var user models2.User

	row := payments.db.QueryRow(ctx, sql, id)

	err := row.Scan(
		&payment.Id,
		&payment.Type,
		&payment.TransactionType,
		&payment.OrderNumber,
		&payment.Amount,
		&payment.CreatedAt,
		&payment.UpdatedAt,
		&user.Id,
		&user.Login,
		&user.Balance,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	payment.User = &user

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, application_errors.ErrNotFound
	}

	return &payment, nil
}

func (payments *PaymentRepository) WithdrawnAmountByUser(ctx context.Context, user *models2.User) (int64, error) {
	sql := `SELECT COALESCE(SUM(amount), 0) FROM payments WHERE user_id = $1 AND type = 'withdraw'`

	var amount int64
	err := payments.db.QueryRow(ctx, sql, user.Id).Scan(&amount)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return amount, nil
}
