package handlers

import (
	"encoding/json"
	"github.com/DmitriyV003/bonus/internal/application_errors"
	"github.com/DmitriyV003/bonus/internal/clients"
	"github.com/DmitriyV003/bonus/internal/container"
	"github.com/DmitriyV003/bonus/internal/models"
	services2 "github.com/DmitriyV003/bonus/internal/services"
	"io/ioutil"
	"net/http"
)

func CreateOrderHandler(container *container.Container, bonusClient *clients.BonusClient) http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		response, err := ioutil.ReadAll(request.Body)
		if err != nil {
			application_errors.SwitchError(&res, err)
			return
		}
		defer request.Body.Close()

		if len(response) == 0 {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		orderService := services2.NewOrderService(container, services2.NewLuhnOrderNumberValidator(), bonusClient)
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
		res.WriteHeader(http.StatusAccepted)
		res.Write(data)
	}
}
