package handlers

import (
	"encoding/json"
	"github.com/DmitriyV003/bonus/internal/application_errors"
	"github.com/DmitriyV003/bonus/internal/container"
	"github.com/DmitriyV003/bonus/internal/requests"
	"github.com/DmitriyV003/bonus/internal/services"
	"net/http"
)

type WithdrawHandler struct {
}

func NewWithdrawHandler() *WithdrawHandler {
	return &WithdrawHandler{}
}

func (h *WithdrawHandler) Handle(container *container.Container) http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		var withdrawReq requests.WithdrawRequest
		if err := json.NewDecoder(request.Body).Decode(&withdrawReq); err != nil {
			application_errors.SwitchError(&res, err)
			return
		}

		userService := services.NewUserService(container, services.GetLoggedInUser(), services.NewLuhnOrderNumberValidator(), services.NewPaymentService(container))
		err := userService.Withdraw(withdrawReq.Order, withdrawReq.Sum)
		if err != nil {
			application_errors.SwitchError(&res, err)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
	}
}
