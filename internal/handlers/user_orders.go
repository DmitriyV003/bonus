package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/DmitriyV003/bonus/internal/applicationerrors"
	"github.com/DmitriyV003/bonus/internal/resources"
	"github.com/DmitriyV003/bonus/internal/services"
	"github.com/DmitriyV003/bonus/internal/services/interfaces"
)

type UserOrdersHandler struct {
	orderService interfaces.OrderService
}

func NewUserOrdersHandler(orderService interfaces.OrderService) *UserOrdersHandler {
	return &UserOrdersHandler{
		orderService: orderService,
	}
}

// Handle Get all user orders
// Available for authorized users
func (h *UserOrdersHandler) Handle() http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		orders, err := h.orderService.OrdersByUser(request.Context(), services.GetLoggedInUser())
		if err != nil {
			applicationerrors.SwitchError(&res, err, nil)
			return
		}

		var ordersToReturn []*resources.OrderResource
		for _, order := range orders {
			orderResource := resources.NewOrderResource(order.Number, order.Status, order.Amount, order.CreatedAt)
			ordersToReturn = append(ordersToReturn, orderResource)
		}

		data, err := json.Marshal(ordersToReturn)
		if err != nil {
			applicationerrors.SwitchError(&res, err, nil)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.Write(data)
		res.WriteHeader(http.StatusAccepted)
	}
}
