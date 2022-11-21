package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/DmitriyV003/bonus/internal/applicationerrors"
	"github.com/DmitriyV003/bonus/internal/services"
	"github.com/DmitriyV003/bonus/internal/services/interfaces"
)

type UserBalanceHandler struct {
	balanceService interfaces.BalanceService
}

func NewUserBalanceHandler(balanceService interfaces.BalanceService) *UserBalanceHandler {
	return &UserBalanceHandler{balanceService: balanceService}
}

func (h *UserBalanceHandler) Handle() http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		resource, err := h.balanceService.Balance(request.Context(), services.GetLoggedInUser())
		if err != nil {
			applicationerrors.SwitchError(&res, err, nil)
			return
		}

		data, err := json.Marshal(resource)
		if err != nil {
			applicationerrors.SwitchError(&res, err, nil)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.Write(data)
		res.WriteHeader(http.StatusAccepted)
	}
}
