package handlers

import (
	"encoding/json"
	"github.com/DmitriyV003/bonus/internal/application_errors"
	"github.com/DmitriyV003/bonus/internal/services"
	"net/http"
)

type UserBalanceHandler struct {
	balanceService *services.BalanceService
}

func NewUserBalanceHandler(balanceService *services.BalanceService) *UserBalanceHandler {
	return &UserBalanceHandler{balanceService: balanceService}
}

func (h *UserBalanceHandler) Handle() http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		resource, err := h.balanceService.Balance(services.GetLoggedInUser())
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
