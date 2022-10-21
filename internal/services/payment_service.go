package services

import (
	"context"
	"github.com/DmitriyV003/bonus/internal/container"
	models2 "github.com/DmitriyV003/bonus/internal/models"
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

func (ps *PaymentService) CreateWithdrawPayment(user *models2.User, amount int64, orderNumber string) error {
	payment := models2.Payment{
		Type:            models2.WITHDRAW_TYPE,
		TransactionType: models2.CREDIT,
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

func (ps *PaymentService) CreateAccrualPayment(user *models2.User, amount int64, orderNumber string) error {
	payment := models2.Payment{
		Type:            models2.ACCRUAL_TYPE,
		TransactionType: models2.DEBIT,
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

func (ps *PaymentService) create(payment *models2.Payment) (*models2.Payment, error) {
	createdPayment, err := ps.container.Payments.Create(context.Background(), payment)
	if err != nil {
		return nil, err
	}

	return createdPayment, nil
}
