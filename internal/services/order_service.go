package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/DmitriyV003/bonus/internal/application_errors"
	"github.com/DmitriyV003/bonus/internal/clients"
	"github.com/DmitriyV003/bonus/internal/container"
	"github.com/DmitriyV003/bonus/internal/models"
	"github.com/DmitriyV003/bonus/internal/policy"
	"strconv"
)

type OrderService struct {
	container   *container.Container
	validator   OrderValidator
	bonusClient *clients.BonusClient
}

func NewOrderService(
	container *container.Container,
	validator OrderValidator,
	bonusClient *clients.BonusClient,
) *OrderService {
	return &OrderService{
		container:   container,
		validator:   validator,
		bonusClient: bonusClient,
	}
}

func (myself *OrderService) Create(user *models.User, orderNumber string) (*models.Order, error) {
	parsedOderNumber, err := strconv.ParseInt(orderNumber, 10, 64)
	if err != nil {
		return nil, err
	}

	isValid := myself.validator.Validate(parsedOderNumber)
	if !isValid {
		return nil, application_errors.ErrInvalidOrderNumber
	}

	order, err := myself.container.Orders.GetByNumber(context.Background(), orderNumber)
	if err != nil && !errors.Is(err, application_errors.ErrNotFound) {
		return nil, err
	}

	orderPolicy := policy.NewOrderPolicy(order, user)
	if order != nil {
		if orderPolicy.Create() {
			return nil, fmt.Errorf("order already accepted: %w", application_errors.ErrModelAlreadyCreated)
		} else {
			return nil, fmt.Errorf("order already uploaded bt another user: %w", application_errors.ErrConflict)
		}
	}

	order = models.NewOrder(orderNumber, models.NewStatus, 0, user)
	order, err = myself.container.Orders.Create(context.Background(), order)
	if err != nil {
		return nil, fmt.Errorf("unable to create order in db: %w", err)
	}

	_, err = myself.bonusClient.CreateOrder(orderNumber)
	if err != nil {
		return nil, fmt.Errorf("unable to create order in black box: %w", err)
	}

	orderDetails, err := myself.bonusClient.GetOrderDetails(orderNumber)
	if err != nil {
		return nil, fmt.Errorf("unable to get order details: %w", err)
	}

	order.Status = orderDetails.Status
	order.Amount = int64(orderDetails.Amount * 10000)
	err = myself.container.Orders.UpdateById(context.Background(), order)
	if err != nil {
		return nil, fmt.Errorf("unable to create order in db: %w", err)
	}

	if orderDetails.Amount > 0 {
		paymentService := NewPaymentService(myself.container)
		err = paymentService.CreateAccrualPayment(user, int64(orderDetails.Amount*10000), orderNumber)
		if err != nil {
			return nil, fmt.Errorf("unable to create payment: %w", err)
		}
	}

	return order, nil
}

func (myself *OrderService) OrdersByUser(user *models.User) ([]*models.Order, error) {
	orders, err := myself.container.Orders.OrdersByUser(context.Background(), user)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
