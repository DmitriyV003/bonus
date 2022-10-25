package handlers

import (
	"encoding/json"
	"github.com/DmitriyV003/bonus/internal/applicationerrors"
	"github.com/DmitriyV003/bonus/internal/requests"
	"github.com/DmitriyV003/bonus/internal/services"
	"net/http"
)

type WithdrawHandler struct {
	userService *services.UserService
}

func NewWithdrawHandler(userService *services.UserService) *WithdrawHandler {
	return &WithdrawHandler{
		userService: userService,
	}
}

func (h *WithdrawHandler) Handle() http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		var withdrawReq requests.WithdrawRequest
		if err := json.NewDecoder(request.Body).Decode(&withdrawReq); err != nil {
			applicationerrors.SwitchError(&res, err)
			return
		}

		err := h.userService.Withdraw(services.GetLoggedInUser(), withdrawReq.Order, withdrawReq.Sum)
		if err != nil {
			applicationerrors.SwitchError(&res, err)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
	}
}
