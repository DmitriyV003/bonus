package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/DmitriyV003/bonus/internal/applicationerrors"
	"github.com/DmitriyV003/bonus/internal/resources"
	"github.com/DmitriyV003/bonus/internal/services"
	"github.com/DmitriyV003/bonus/internal/services/interfaces"
)

type UserWithdawsHandler struct {
	userService interfaces.UserService
}

func NewUserWithdawsHandler(us interfaces.UserService) *UserWithdawsHandler {
	return &UserWithdawsHandler{
		userService: us,
	}
}

// Handle Get all user withdrawals
func (h *UserWithdawsHandler) Handle() http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		payments, err := h.userService.AllWithdrawsByUser(request.Context(), services.GetLoggedInUser())
		if err != nil {
			applicationerrors.SwitchError(&res, err, nil)
			return
		}

		var paymentsToReturn []*resources.PaymentResource
		for _, payment := range payments {
			paymentResource := resources.NewPaymentResource(payment.OrderNumber, payment.Amount, payment.CreatedAt)
			paymentsToReturn = append(paymentsToReturn, paymentResource)
		}

		data, err := json.Marshal(paymentsToReturn)
		if err != nil {
			applicationerrors.SwitchError(&res, err, nil)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.Write(data)
		res.WriteHeader(http.StatusOK)
	}
}
