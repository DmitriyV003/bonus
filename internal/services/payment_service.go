package services

import (
	"context"
	"github.com/DmitriyV003/bonus/internal/container"
	"github.com/DmitriyV003/bonus/internal/models"
	"github.com/rs/zerolog/log"
	"time"
)

type PaymentService struct {
	container *container.Container
}

func NewPaymentService(container *container.Container) *PaymentService {
	return &PaymentService{
		container: container,
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

	balanceService := NewBalanceService(ps.container, user)
	err = balanceService.Withdraw(createdPayment)
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

	balanceService := NewBalanceService(ps.container, user)
	err = balanceService.Accrual(createdPayment)
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
	createdPayment, err := ps.container.Payments.Create(context.Background(), payment)
	if err != nil {
		return nil, err
	}

	return createdPayment, nil
}
