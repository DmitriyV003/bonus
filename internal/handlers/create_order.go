package handlers

import (
	"encoding/json"
	"github.com/DmitriyV003/bonus/internal/applicationerrors"
	"github.com/DmitriyV003/bonus/internal/services"
	"io/ioutil"
	"net/http"
)

type CreateOrderHandler struct {
	orderService *services.OrderService
}

func NewCreateOrderHandler(orderService *services.OrderService) *CreateOrderHandler {
	return &CreateOrderHandler{
		orderService: orderService,
	}
}

func (h *CreateOrderHandler) Handle() http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		response, err := ioutil.ReadAll(request.Body)
		if err != nil {
			applicationerrors.SwitchError(&res, err)
			return
		}
		defer request.Body.Close()

		if len(response) == 0 {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		order, err := h.orderService.Create(request.Context(), services.GetLoggedInUser(), string(response))
		if err != nil {
			applicationerrors.SwitchError(&res, err)
			return
		}

		data, err := json.Marshal(order)
		if err != nil {
			applicationerrors.SwitchError(&res, err)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusAccepted)
		res.Write(data)
	}
}
