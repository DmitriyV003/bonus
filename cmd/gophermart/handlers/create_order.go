package handlers

import (
	"encoding/json"
	"github.com/DmitriyV003/bonus/cmd/gophermart/application_errors"
	"github.com/DmitriyV003/bonus/cmd/gophermart/container"
	"github.com/DmitriyV003/bonus/cmd/gophermart/models"
	"github.com/DmitriyV003/bonus/cmd/gophermart/services"
	"io/ioutil"
	"net/http"
)

func CreateOrderHandler(container *container.Container) http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		response, err := ioutil.ReadAll(request.Body)
		if err != nil {
			application_errors.SwitchError(&res, err)
			return
		}
		defer request.Body.Close()

		if string(response) == "" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		orderService := services.NewOrderService(container, services.NewLuhnOrderNumberValidator())
		order, err := orderService.Create(request.Context().Value("user").(*models.User), string(response))
		if err != nil {
			application_errors.SwitchError(&res, err)
			return
		}

		data, err := json.Marshal(order)
		if err != nil {
			application_errors.SwitchError(&res, err)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.Write(data)
		res.WriteHeader(http.StatusAccepted)
	}
}
