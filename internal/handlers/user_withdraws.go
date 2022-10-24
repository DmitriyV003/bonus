package handlers

import (
	"encoding/json"
	"github.com/DmitriyV003/bonus/internal/application_errors"
	"github.com/DmitriyV003/bonus/internal/resources"
	"github.com/DmitriyV003/bonus/internal/services"
	"net/http"
)

type UserWithdawsHandler struct {
	userService *services.UserService
}

func NewUserWithdawsHandler(us *services.UserService) *UserWithdawsHandler {
	return &UserWithdawsHandler{
		userService: us,
	}
}

func (h *UserWithdawsHandler) Handle() http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		payments, err := h.userService.AllWithdrawsByUser(services.GetLoggedInUser())
		if err != nil {
			application_errors.SwitchError(&res, err)
			return
		}

		var paymentsToReturn []*resources.PaymentResource
		for _, payment := range payments {
			paymentResource := resources.NewPaymentResource(payment.OrderNumber, payment.Amount, payment.CreatedAt)
			paymentsToReturn = append(paymentsToReturn, paymentResource)
		}

		data, err := json.Marshal(paymentsToReturn)
		if err != nil {
			application_errors.SwitchError(&res, err)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.Write(data)
		res.WriteHeader(http.StatusOK)
	}
}
