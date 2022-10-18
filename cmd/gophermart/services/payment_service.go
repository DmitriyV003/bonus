package services

import (
	"context"
	"github.com/DmitriyV003/bonus/cmd/gophermart/container"
	"github.com/DmitriyV003/bonus/cmd/gophermart/models"
	"time"
)

type PaymentService struct {
	container *container.Container
	payment   *models.Payment
}

func NewPaymentService(container *container.Container, payment *models.Payment) *PaymentService {
	return &PaymentService{
		container: container,
		payment:   payment,
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
	err = balanceService.Withdraw(createdPayment)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PaymentService) create(payment *models.Payment) (*models.Payment, error) {
	createdPayment, err := ps.container.Payments.Create(context.Background(), payment)
	if err != nil {
		return nil, err
	}

	return createdPayment, nil
}
