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
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
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

	//_, err = myself.bonusClient.CreateOrder(orderNumber)
	//if err != nil {
	//	return nil, fmt.Errorf("unable to create order in black box: %w", err)
	//}
	//
	//orderDetails, err := myself.bonusClient.GetOrderDetails(orderNumber)
	//if err != nil {
	//	return nil, fmt.Errorf("unable to get order details: %w", err)
	//}
	//
	//order.Status = orderDetails.Status
	//order.Amount = int64(orderDetails.Amount * 10000)
	//err = myself.container.Orders.UpdateById(context.Background(), order)
	//if err != nil {
	//	return nil, fmt.Errorf("unable to create order in db: %w", err)
	//}
	//
	//if orderDetails.Amount > 0 {
	//	paymentService := NewPaymentService(myself.container)
	//	err = paymentService.CreateAccrualPayment(user, int64(orderDetails.Amount*10000), orderNumber)
	//	if err != nil {
	//		return nil, fmt.Errorf("unable to create payment: %w", err)
	//	}
	//}

	err = myself.sendAndUpdateOrder(context.Background(), user, order)
	if err != nil {
		return nil, fmt.Errorf("unable to update order: %w", err)
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

func (myself *OrderService) PollPendingOrders(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	ctx, cancel := context.WithCancel(ctx)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			cancel()
			return
		case <-ticker.C:
			orders, err := myself.container.Orders.AllPending(ctx)
			if err != nil {
				cancel()
				return
			}

			for _, order := range orders {
				log.Info().Fields(map[string]interface{}{
					"order_id": order.Id,
					"status":   order.Status,
					"number":   order.Number,
				}).Msg("Polling order")

				user, err := myself.container.Users.GetById(ctx, order.User.Id)
				if err != nil && !errors.Is(err, application_errors.ErrNotFound) {
					cancel()
					return
				}

				err = myself.sendAndUpdateOrder(ctx, user, order)
				if err != nil && !errors.Is(err, application_errors.ErrServiceUnavailable) {
					log.Error().Err(err)
					cancel()
				}
			}
		}
	}
}

func (myself *OrderService) sendAndUpdateOrder(ctx context.Context, user *models.User, order *models.Order) error {
	//_, err := myself.bonusClient.CreateOrder(order.Number)
	//if err != nil {
	//	log.Error().Err(err)
	//	return fmt.Errorf("unable to create order in black box: %w", err)
	//}

	orderDetails, err := myself.bonusClient.GetOrderDetails(order.Number)
	if err != nil {
		log.Error().Err(err)
		return fmt.Errorf("unable to get order details: %w", err)
	}

	if orderDetails != nil && order.Status != orderDetails.Status {
		order.Status = orderDetails.Status
		order.Amount = int64(orderDetails.Amount * 10000)
		err = myself.container.Orders.UpdateById(context.Background(), order)
		if err != nil {
			log.Error().Err(err)
			return fmt.Errorf("unable to create order in db: %w", err)
		}

		if orderDetails.Amount > 0 {
			paymentService := NewPaymentService(myself.container)
			err = paymentService.CreateAccrualPayment(user, int64(orderDetails.Amount*10000), order.Number)
			if err != nil {
				return fmt.Errorf("unable to create payment: %w", err)
			}
		}
	}

	return nil
}
