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

type LoginHandler struct {
	authService *services.AuthService
}

func NewLoginHandler(authService *services.AuthService) *LoginHandler {
	return &LoginHandler{
		authService: authService,
	}
}

func (h LoginHandler) Handle() http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		var loginRequest requests.LoginRequest

		validate := validator.New()

		if err := json.NewDecoder(request.Body).Decode(&loginRequest); err != nil {
			applicationerrors.SwitchError(&res, err)
			return
		}

		if err := validate.Struct(&loginRequest); err != nil {
			applicationerrors.WriteHTTPError(&res, http.StatusBadRequest, err)
			return
		}

		token, err := h.authService.Login(request.Context(), loginRequest.Login, loginRequest.Password)
		if err != nil {
			applicationerrors.SwitchError(&res, err)
			return
		}

		res.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token.Value))
		res.WriteHeader(http.StatusOK)
	}
}
