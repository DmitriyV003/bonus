package services

import (
	"context"
	"fmt"
	"github.com/DmitriyV003/bonus/cmd/gophermart/application_errors"
	"github.com/DmitriyV003/bonus/cmd/gophermart/container"
	"github.com/DmitriyV003/bonus/cmd/gophermart/models"
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

func (myself *OrderService) Store(orderNumber string) error {
	orderNum, err := strconv.ParseInt(orderNumber, 10, 64)
	if err != nil {
		return err
	}

	isValid := myself.validator.Validate(orderNum)
	if !isValid {
		return application_errors.ErrInvalidOrderNumber
	}

	fmt.Println(context.Background().Value("user"))

	return nil
}
