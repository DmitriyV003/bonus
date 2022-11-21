package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/DmitriyV003/bonus/internal/applicationerrors"
	"github.com/DmitriyV003/bonus/internal/requests"
	"github.com/DmitriyV003/bonus/internal/services"
	"github.com/DmitriyV003/bonus/internal/services/interfaces"
)

type WithdrawHandler struct {
	userService interfaces.UserService
}

func NewWithdrawHandler(userService interfaces.UserService) *WithdrawHandler {
	return &WithdrawHandler{
		userService: userService,
	}
}

func (h *WithdrawHandler) Handle() http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		var withdrawReq requests.WithdrawRequest
		if err := json.NewDecoder(request.Body).Decode(&withdrawReq); err != nil {
			applicationerrors.SwitchError(&res, err, nil)
			return
		}

		err := h.userService.Withdraw(request.Context(), services.GetLoggedInUser(), withdrawReq.Order, withdrawReq.Sum)
		if err != nil {
			applicationerrors.SwitchError(&res, err, nil)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
	}
}
