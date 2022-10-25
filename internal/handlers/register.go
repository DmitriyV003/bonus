package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/DmitriyV003/bonus/internal/applicationerrors"
	"github.com/DmitriyV003/bonus/internal/requests"
	"github.com/DmitriyV003/bonus/internal/services"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type RegisterHandler struct {
	userService *services.UserService
}

func NewRegisterHandler(userService *services.UserService) *RegisterHandler {
	return &RegisterHandler{
		userService: userService,
	}
}

func (h *RegisterHandler) Handle() http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		var regRequest requests.RegistrationRequest

		validate := validator.New()

		if err := json.NewDecoder(request.Body).Decode(&regRequest); err != nil {
			applicationerrors.SwitchError(&res, err)
			return
		}

		if err := validate.Struct(&regRequest); err != nil {
			applicationerrors.WriteHTTPError(&res, http.StatusBadRequest, err)
			return
		}

		token, err := h.userService.Create(request.Context(), &regRequest)
		if err != nil {
			applicationerrors.SwitchError(&res, err)
			return
		}

		res.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token.Value))
		res.WriteHeader(http.StatusOK)
	}
}
