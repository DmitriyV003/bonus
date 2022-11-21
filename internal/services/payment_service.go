package services

import (
	"context"
	"fmt"
	"time"

	"github.com/DmitriyV003/bonus/internal/models"
	"github.com/DmitriyV003/bonus/internal/repository/interfaces"
	"github.com/rs/zerolog/log"
)

type PaymentService struct {
	payments       interfaces.PaymentRepository
	users          interfaces.UserRepository
	balanceService *BalanceService
}

func NewPaymentService(payments interfaces.PaymentRepository, users interfaces.UserRepository, balanceService *BalanceService) *PaymentService {
	return &PaymentService{
		payments:       payments,
		users:          users,
		balanceService: balanceService,
	}
}

func (ps *PaymentService) CreateWithdrawPayment(ctx context.Context, user *models.User, amount int64, orderNumber string) (*models.Payment, error) {
	payment := models.Payment{
		Type:            models.WithdrawType,
		TransactionType: models.CREDIT,
		OrderNumber:     orderNumber,
		Amount:          amount,
		User:            user,
		CreatedAt:       time.Now(),
	}

	createdPayment, err := ps.create(ctx, &payment)
	if err != nil {
		return nil, fmt.Errorf("error to create payment in db: %w", err)
	}

	log.Info().Fields(map[string]interface{}{
		"user_id":          user.ID,
		"payment_id":       createdPayment.ID,
		"amount":           createdPayment.Amount,
		"type":             createdPayment.Type,
		"transaction_type": createdPayment.TransactionType,
	}).Msg("Withdraw payment created")

	return createdPayment, nil
}

func (ps *PaymentService) CreateAccrualPayment(ctx context.Context, user *models.User, amount int64, orderNumber string) error {
	payment := models.Payment{
		Type:            models.AccrualType,
		TransactionType: models.DEBIT,
		OrderNumber:     orderNumber,
		Amount:          amount,
		User:            user,
		CreatedAt:       time.Now(),
	}

	createdPayment, err := ps.create(ctx, &payment)
	if err != nil {
		return fmt.Errorf("error to save payment in db: %w", err)
	}

	err = ps.balanceService.Accrual(ctx, createdPayment, user)
	if err != nil {
		return fmt.Errorf("error to change user palance: %w", err)
	}
	log.Info().Fields(map[string]interface{}{
		"user_id":          user.ID,
		"payment_id":       createdPayment.ID,
		"amount":           createdPayment.Amount,
		"type":             createdPayment.Type,
		"transaction_type": createdPayment.TransactionType,
	}).Msg("Accrual payment created")

	return nil
}

func (ps *PaymentService) create(ctx context.Context, payment *models.Payment) (*models.Payment, error) {
	createdPayment, err := ps.payments.Create(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("error to create payment in db: %w", err)
	}

	return createdPayment, nil
}
