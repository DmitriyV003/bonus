package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DmitriyV003/bonus/internal/applicationerrors"
	"github.com/DmitriyV003/bonus/internal/requests"
	serviceinterfaces "github.com/DmitriyV003/bonus/internal/services/interfaces"
	"github.com/go-playground/validator/v10"
)

type LoginHandler struct {
	authService serviceinterfaces.AuthService
}

func NewLoginHandler(authService serviceinterfaces.AuthService) *LoginHandler {
	return &LoginHandler{
		authService: authService,
	}
}

func (h LoginHandler) Handle() http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		var loginRequest requests.LoginRequest

		validate := validator.New()

		if err := json.NewDecoder(request.Body).Decode(&loginRequest); err != nil {
			applicationerrors.SwitchError(&res, err, nil)
			return
		}

		if err := validate.Struct(&loginRequest); err != nil {
			applicationerrors.WriteHTTPError(&res, http.StatusBadRequest, err)
			return
		}

		token, err := h.authService.Login(request.Context(), loginRequest.Login, loginRequest.Password)
		if err != nil {
			applicationerrors.SwitchError(&res, err, nil)
			return
		}

		res.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token.Value))
		res.WriteHeader(http.StatusOK)
	}
}
