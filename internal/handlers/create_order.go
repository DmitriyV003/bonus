package handlers

import (
	"encoding/json"
	"github.com/DmitriyV003/bonus/internal/applicationerrors"
	"github.com/DmitriyV003/bonus/internal/services"
	"github.com/DmitriyV003/bonus/internal/services/interfaces"
	"io"
	"net/http"
)

type CreateOrderHandler struct {
	orderService interfaces.OrderService
}

func NewCreateOrderHandler(orderService interfaces.OrderService) *CreateOrderHandler {
	return &CreateOrderHandler{
		orderService: orderService,
	}
}

func (h *CreateOrderHandler) Handle() http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		response, err := io.ReadAll(request.Body)
		if err != nil {
			applicationerrors.SwitchError(&res, err, nil, "error to read body")
			return
		}
		defer request.Body.Close()

		if len(response) == 0 {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		order, err := h.orderService.Create(request.Context(), services.GetLoggedInUser(), string(response))
		if err != nil {
			applicationerrors.SwitchError(&res, err, map[string]interface{}{
				"user_id": services.GetLoggedInUser().ID,
			}, "error to create order")
			return
		}

		data, err := json.Marshal(order)

		if err != nil {
			applicationerrors.SwitchError(&res, err, map[string]interface{}{
				"user_id":  services.GetLoggedInUser().ID,
				"order_id": order.ID,
			}, "error to marshal order")
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusAccepted)
		res.Write(data)
	}
}
