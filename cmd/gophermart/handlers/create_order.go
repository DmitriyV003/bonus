package handlers

import (
	"github.com/DmitriyV003/bonus/cmd/gophermart/application_errors"
	"github.com/DmitriyV003/bonus/cmd/gophermart/config"
	"github.com/DmitriyV003/bonus/cmd/gophermart/container"
	"github.com/DmitriyV003/bonus/cmd/gophermart/services"
	"io/ioutil"
	"net/http"
)

func CreateOrderHandler(container *container.Container, conf *config.Config) http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		response, err := ioutil.ReadAll(request.Body)
		if err != nil {
			application_errors.SwitchError(&res, err)
			return
		}
		defer request.Body.Close()

		orderService := services.NewOrderService(container, nil, services.NewLuhnOrderNumberValidator())
		err = orderService.Store(string(response))
		if err != nil {
			application_errors.SwitchError(&res, err)
			return
		}

		res.WriteHeader(http.StatusOK)
	}
}
