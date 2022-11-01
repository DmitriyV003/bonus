package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/DmitriyV003/bonus/internal/applicationerrors"
	"github.com/DmitriyV003/bonus/internal/clients/clientinterfaces"
	"github.com/DmitriyV003/bonus/internal/models"
	"github.com/DmitriyV003/bonus/internal/policy"
	"github.com/DmitriyV003/bonus/internal/repository/interfaces"
	interfaces2 "github.com/DmitriyV003/bonus/internal/services/interfaces"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
)

type OrderService struct {
	validator      OrderValidator
	bonusClient    clientinterfaces.BonusClient
	orders         interfaces.OrderRepository
	users          interfaces.UserRepository
	paymentService interfaces2.PaymentService
}

func NewOrderService(
	validator OrderValidator,
	bonusClient clientinterfaces.BonusClient,
	orders interfaces.OrderRepository,
	users interfaces.UserRepository,
	paymentService interfaces2.PaymentService,
) *OrderService {
	return &OrderService{
		validator:      validator,
		bonusClient:    bonusClient,
		orders:         orders,
		users:          users,
		paymentService: paymentService,
	}
}

func (myself *OrderService) Create(ctx context.Context, user *models.User, orderNumber string) (*models.Order, error) {
	parsedOderNumber, err := strconv.ParseInt(orderNumber, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error to parse order number: %w", err)
	}

	isValid := myself.validator.Validate(parsedOderNumber)
	if !isValid {
		return nil, fmt.Errorf("invalid order number: %w", applicationerrors.ErrInvalidOrderNumber)
	}

	order, err := myself.orders.GetByNumber(ctx, orderNumber)
	if err != nil && !errors.Is(err, applicationerrors.ErrNotFound) {
		return nil, fmt.Errorf("db error to get order by number: %w", err)
	}

	orderPolicy := policy.NewOrderPolicy(order, user)
	if order != nil {
		if orderPolicy.Create() {
			return nil, fmt.Errorf("order already accepted: %w", applicationerrors.ErrModelAlreadyCreated)
		} else {
			return nil, fmt.Errorf("order already uploaded bt another user: %w", applicationerrors.ErrConflict)
		}
	}

	order = models.NewOrder(orderNumber, models.NewStatus, 0, user)
	fmt.Println(order)
	order, err = myself.orders.Create(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("unable to create order in db: %w", err)
	}

	err = myself.sendAndUpdateOrder(ctx, user, order)
	if err != nil {
		return nil, fmt.Errorf("unable to update order: %w", err)
	}

	return order, nil
}

func (myself *OrderService) OrdersByUser(ctx context.Context, user *models.User) ([]*models.Order, error) {
	orders, err := myself.orders.OrdersByUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("error to get orders by user: %w", err)
	}

	return orders, nil
}

func (myself *OrderService) PollPendingOrders(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return errors.New("context cancelled in polling orders")
		case <-ticker.C:
			orders, err := myself.orders.AllPending(ctx)
			if err != nil {
				return err
			}

			for _, order := range orders {
				log.Info().Fields(map[string]interface{}{
					"order_id": order.ID,
					"status":   order.Status,
					"number":   order.Number,
				}).Msg("Polling order")

				user, err := myself.users.GetByID(ctx, order.User.ID)
				if err != nil && !errors.Is(err, applicationerrors.ErrNotFound) {
					log.Error().Err(err).Msg("error occurred")
					return err
				}

				err = myself.sendAndUpdateOrder(ctx, user, order)
				if err != nil && !errors.Is(err, applicationerrors.ErrServiceUnavailable) {
					log.Error().Err(err).Msg("error occurred")
					return err
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
		err = myself.orders.UpdateByID(context.Background(), order)
		if err != nil {
			log.Error().Err(err)
			return fmt.Errorf("unable to create order in db: %w", err)
		}

		if orderDetails.Amount > 0 {
			err = myself.paymentService.CreateAccrualPayment(ctx, user, int64(orderDetails.Amount*10000), order.Number)
			if err != nil {
				return fmt.Errorf("unable to create payment: %w", err)
			}
		}
	}

	return nil
}
