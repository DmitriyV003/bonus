package handlers

import (
	"github.com/DmitriyV003/bonus/cmd/gophermart/config"
	"net/http"
)

func RegisterHandler(container *config.Container) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}
