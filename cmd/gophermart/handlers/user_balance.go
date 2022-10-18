package handlers

import (
	"encoding/json"
	"github.com/DmitriyV003/bonus/cmd/gophermart/application_errors"
	"github.com/DmitriyV003/bonus/cmd/gophermart/container"
	"github.com/DmitriyV003/bonus/cmd/gophermart/models"
	"github.com/DmitriyV003/bonus/cmd/gophermart/services"
	"net/http"
)

func UserBalanceHandler(container *container.Container) http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		userService := services.NewBalanceService(container, request.Context().Value("user").(*models.User))
		resource, err := userService.Balance()
		if err != nil {
			application_errors.SwitchError(&res, err)
			return
		}

		data, err := json.Marshal(resource)
		if err != nil {
			application_errors.SwitchError(&res, err)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.Write(data)
		res.WriteHeader(http.StatusAccepted)
	}
}
