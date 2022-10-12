package handlers

import (
	"encoding/json"
	"github.com/DmitriyV003/bonus/cmd/gophermart/application_errors"
	"github.com/DmitriyV003/bonus/cmd/gophermart/container"
	"github.com/DmitriyV003/bonus/cmd/gophermart/requests"
	"github.com/DmitriyV003/bonus/cmd/gophermart/services"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func RegisterHandler(container *container.Container) http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		var regRequest requests.RegistrationRequest

		validate := validator.New()

		if err := json.NewDecoder(request.Body).Decode(&regRequest); err != nil {
			application_errors.SwitchError(&res, err)
			return
		}

		if err := validate.Struct(&regRequest); err != nil {
			application_errors.WriteHTTPError(&res, http.StatusBadRequest, err)
			return
		}

		service := services.NewUserService(container)
		err := service.Create(&regRequest)
		if err != nil {
			application_errors.SwitchError(&res, err)
			return
		}

		res.WriteHeader(http.StatusOK)
	}
}
