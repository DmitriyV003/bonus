package applicationerrors

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"net/http"
	"runtime"
	"strings"
)

var ErrNotFound = errors.New("not found")
var ErrInternalServer = errors.New("internal server error")
var ErrConflict = errors.New("conflict")
var ErrInvalidOrderNumber = errors.New("invalid order number")
var ErrModelAlreadyCreated = errors.New("model already created")
var ErrLowUserABalance = errors.New("low user balance")
var ErrServiceUnavailable = errors.New("service unavailable")

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

func SwitchError(w *http.ResponseWriter, err error, fields map[string]interface{}, messages ...string) {
	pc, file, line, ok := runtime.Caller(1)
	errorString := ""
	if ok {
		errorString = fmt.Sprintf("Called from %s, line #%d, func: %v\n", file, line, runtime.FuncForPC(pc).Name())
	}
	log.Error().Fields(fields).Err(err).Msg(strings.Join(messages, "; ") + errorString)
	switch {
	case errors.Is(err, ErrNotFound):
		WriteHTTPError(w, http.StatusNotFound, nil)
	case errors.Is(err, ErrConflict):
		WriteHTTPError(w, http.StatusConflict, nil)
	case errors.Is(err, ErrInvalidOrderNumber):
		WriteHTTPError(w, http.StatusUnprocessableEntity, nil)
	case errors.Is(err, ErrModelAlreadyCreated):
		(*w).WriteHeader(http.StatusOK)
	case errors.Is(err, ErrLowUserABalance):
		(*w).WriteHeader(http.StatusPaymentRequired)
	case errors.Is(err, ErrServiceUnavailable):
		(*w).WriteHeader(http.StatusAccepted)
	default:
		WriteHTTPError(w, http.StatusInternalServerError, nil)
	}
}
