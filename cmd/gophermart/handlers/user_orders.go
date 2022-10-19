package handlers

import (
	"encoding/json"
	"github.com/DmitriyV003/bonus/cmd/gophermart/application_errors"
	"github.com/DmitriyV003/bonus/cmd/gophermart/container"
	"github.com/DmitriyV003/bonus/cmd/gophermart/models"
	"github.com/DmitriyV003/bonus/cmd/gophermart/resources"
	"github.com/DmitriyV003/bonus/cmd/gophermart/services"
	"net/http"
)

func UserOrdersHandler(container *container.Container) http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		orderService := services.NewOrderService(container, nil, nil)
		orders, err := orderService.OrdersByUser(request.Context().Value("user").(*models.User))
		if err != nil {
			application_errors.SwitchError(&res, err)
			return
		}

		var ordersToReturn []*resources.OrderResource
		for _, order := range orders {
			orderResource := resources.NewOrderResource(order.Number, *order.Status, float64(order.Amount)/10000, order.CreatedAt)
			ordersToReturn = append(ordersToReturn, orderResource)
		}

		data, err := json.Marshal(ordersToReturn)
		if err != nil {
			application_errors.SwitchError(&res, err)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.Write(data)
		res.WriteHeader(http.StatusAccepted)
	}
}
