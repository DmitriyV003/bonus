package services

import (
	"context"
	"github.com/DmitriyV003/bonus/internal/models"
	"github.com/DmitriyV003/bonus/internal/repository"
	"github.com/rs/zerolog/log"
	"time"
)

type PaymentService struct {
	payments       *repository.PaymentRepository
	users          *repository.UserRepository
	balanceService *BalanceService
}

func NewPaymentService(payments *repository.PaymentRepository, users *repository.UserRepository, balanceService *BalanceService) *PaymentService {
	return &PaymentService{
		payments:       payments,
		users:          users,
		balanceService: balanceService,
	}
}

func (ps *PaymentService) CreateWithdrawPayment(user *models.User, amount int64, orderNumber string) error {
	payment := models.Payment{
		Type:            models.WITHDRAW_TYPE,
		TransactionType: models.CREDIT,
		OrderNumber:     orderNumber,
		Amount:          amount,
		User:            user,
		CreatedAt:       time.Now(),
	}

	createdPayment, err := ps.create(&payment)
	if err != nil {
		return err
	}

	err = ps.balanceService.Withdraw(createdPayment, user)
	if err != nil {
		return err
	}
	log.Info().Fields(map[string]interface{}{
		"user_id":          user.Id,
		"payment_id":       createdPayment.Id,
		"amount":           createdPayment.Amount,
		"type":             createdPayment.Type,
		"transaction_type": createdPayment.TransactionType,
	}).Msg("Withdraw payment created")

	return nil
}

func (ps *PaymentService) CreateAccrualPayment(user *models.User, amount int64, orderNumber string) error {
	payment := models.Payment{
		Type:            models.ACCRUAL_TYPE,
		TransactionType: models.DEBIT,
		OrderNumber:     orderNumber,
		Amount:          amount,
		User:            user,
		CreatedAt:       time.Now(),
	}

	createdPayment, err := ps.create(&payment)
	if err != nil {
		return err
	}

	err = ps.balanceService.Accrual(createdPayment, user)
	if err != nil {
		return err
	}
	log.Info().Fields(map[string]interface{}{
		"user_id":          user.Id,
		"payment_id":       createdPayment.Id,
		"amount":           createdPayment.Amount,
		"type":             createdPayment.Type,
		"transaction_type": createdPayment.TransactionType,
	}).Msg("Accrual payment created")

	return nil
}

func (ps *PaymentService) create(payment *models.Payment) (*models.Payment, error) {
	createdPayment, err := ps.payments.Create(context.Background(), payment)
	if err != nil {
		return nil, err
	}

	return createdPayment, nil
}
