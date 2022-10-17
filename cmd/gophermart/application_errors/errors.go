package application_errors

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"net/http"
)

var ErrNotFound = errors.New("not found")
var ErrInternalServer = errors.New("internal server error")
var ErrConflict = errors.New("login already busy")
var ErrInvalidOrderNumber = errors.New("invalid order number")

func WriteHTTPError(w *http.ResponseWriter, status int, errs error) {
	(*w).Header().Set("Content-Type", "application/json")

	if errs == nil {
		http.Error(*w, http.StatusText(status), status)
		return
	}

	errorsMap := map[string]string{}
	for _, er := range errs.(validator.ValidationErrors) {
		errorsMap[er.Field()] = er.Error()
	}

	res, err := json.Marshal(errorsMap)
	if err != nil {
		http.Error(*w, "Internal Server Error", http.StatusInternalServerError)
	}
	(*w).WriteHeader(status)
	(*w).Write(res)
}

func SwitchError(w *http.ResponseWriter, err error) {
	log.Error().Err(err).Msg("Error occurred")
	switch {
	case errors.Is(err, ErrNotFound):
		WriteHTTPError(w, http.StatusNotFound, nil)
	case errors.Is(err, ErrConflict):
		WriteHTTPError(w, http.StatusConflict, nil)
	case errors.Is(err, ErrInvalidOrderNumber):
		WriteHTTPError(w, http.StatusUnprocessableEntity, nil)
	default:
		WriteHTTPError(w, http.StatusInternalServerError, nil)
	}
}
