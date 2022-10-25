package handlers

import (
	"encoding/json"
	"github.com/DmitriyV003/bonus/internal/applicationerrors"
	"github.com/DmitriyV003/bonus/internal/resources"
	"github.com/DmitriyV003/bonus/internal/services"
	"net/http"
)

type UserOrdersHandler struct {
	orderService *services.OrderService
}

func NewUserOrdersHandler(orderService *services.OrderService) *UserOrdersHandler {
	return &UserOrdersHandler{
		orderService: orderService,
	}
}

func (h *UserOrdersHandler) Handle() http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		orders, err := h.orderService.OrdersByUser(request.Context(), services.GetLoggedInUser())
		if err != nil {
			applicationerrors.SwitchError(&res, err)
			return
		}

		var ordersToReturn []*resources.OrderResource
		for _, order := range orders {
			orderResource := resources.NewOrderResource(order.Number, order.Status, order.Amount, order.CreatedAt)
			ordersToReturn = append(ordersToReturn, orderResource)
		}

		data, err := json.Marshal(ordersToReturn)
		if err != nil {
			applicationerrors.SwitchError(&res, err)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.Write(data)
		res.WriteHeader(http.StatusAccepted)
	}
}
