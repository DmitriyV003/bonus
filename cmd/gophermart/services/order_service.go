package services

import (
	"context"
	"errors"
	"github.com/DmitriyV003/bonus/cmd/gophermart/application_errors"
	"github.com/DmitriyV003/bonus/cmd/gophermart/container"
	"github.com/DmitriyV003/bonus/cmd/gophermart/models"
	"github.com/DmitriyV003/bonus/cmd/gophermart/policy"
	"strconv"
)

type OrderService struct {
	container *container.Container
	order     *models.Order
	validator OrderValidator
}

func NewOrderService(container *container.Container, order *models.Order, validator OrderValidator) *OrderService {
	return &OrderService{
		container: container,
		order:     order,
		validator: validator,
	}
}

func (myself *OrderService) Store(user *models.User, orderNumber string) (*models.Order, error) {
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
			return nil, application_errors.ErrConflict
		} else {
			return nil, application_errors.ErrModelAlreadyCreated
		}
	}

	order = models.NewOrder(orderNumber, 0, user)
	order, err = myself.container.Orders.Create(context.Background(), order)
	if err != nil {
		return nil, err
	}

	return order, nil
}
