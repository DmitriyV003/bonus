package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/DmitriyV003/bonus/internal/applicationerrors"
	"github.com/DmitriyV003/bonus/internal/models"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PaymentRepository struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(pool *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{
		db: pool,
	}
}

func (payments *PaymentRepository) Create(ctx context.Context, payment *models.Payment) (*models.Payment, error) {
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
		payment.User.ID,
		payment.Type,
		payment.TransactionType,
		payment.OrderNumber,
		payment.Amount,
		payment.CreatedAt,
	).Scan(&id)
	if err != nil {
		return nil, applicationerrors.ErrInternalServer
	}

	payment, err = payments.GetByIDWithUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (payments *PaymentRepository) GetByIDWithUser(ctx context.Context, id int64) (*models.Payment, error) {
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

	var payment models.Payment
	var user models.User

	row := payments.db.QueryRow(ctx, sql, id)

	err := row.Scan(
		&payment.ID,
		&payment.Type,
		&payment.TransactionType,
		&payment.OrderNumber,
		&payment.Amount,
		&payment.CreatedAt,
		&payment.UpdatedAt,
		&user.ID,
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
		return nil, applicationerrors.ErrNotFound
	}

	return &payment, nil
}

func (payments *PaymentRepository) WithdrawnAmountByUser(ctx context.Context, user *models.User) (int64, error) {
	sql := `SELECT COALESCE(SUM(amount), 0) FROM payments WHERE user_id = $1 AND type = 'withdraw'`

	var amount int64
	err := payments.db.QueryRow(ctx, sql, user.ID).Scan(&amount)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return amount, nil
}

func (payments *PaymentRepository) GetWithdrawsByUser(ctx context.Context, user *models.User) ([]*models.Payment, error) {
	sql := `SELECT id, order_number, amount, created_at 
		FROM payments 
		WHERE user_id = $1 AND type = $2 AND transaction_type = $3
		ORDER BY created_at`
	var selectedPayments []*models.Payment

	rows, err := payments.db.Query(ctx, sql, user.ID, models.WithdrawType, models.CREDIT)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var payment models.Payment
		err = rows.Scan(&payment.ID, &payment.OrderNumber, &payment.Amount, &payment.CreatedAt)
		if err != nil {
			return nil, err
		}

		selectedPayments = append(selectedPayments, &payment)
	}

	return selectedPayments, nil
}
