package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/DmitriyV003/bonus/internal/application_errors"
	"github.com/DmitriyV003/bonus/internal/config"
	"github.com/DmitriyV003/bonus/internal/container"
	"github.com/DmitriyV003/bonus/internal/requests"
	"github.com/DmitriyV003/bonus/internal/services"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func LoginHandler(container *container.Container, conf *config.Config) http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		var loginRequest requests.LoginRequest

		validate := validator.New()

		if err := json.NewDecoder(request.Body).Decode(&loginRequest); err != nil {
			application_errors.SwitchError(&res, err)
			return
		}

		if err := validate.Struct(&loginRequest); err != nil {
			application_errors.WriteHTTPError(&res, http.StatusBadRequest, err)
			return
		}

		service := services.NewAuthService(container, conf.JwtSecret)
		token, err := service.Login(loginRequest.Login, loginRequest.Password)
		if err != nil {
			application_errors.SwitchError(&res, err)
			return
		}

		res.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token.Value))
		res.WriteHeader(http.StatusOK)
	}
}
