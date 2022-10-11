package main

import (
	"github.com/DmitriyV003/bonus/cmd/gophermart/config"
	"github.com/DmitriyV003/bonus/cmd/gophermart/container"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	application := config.App{
		Conf:      config.Config{},
		Container: &container.Container{},
	}
	application.Config()

	defer application.Close()

	log.Info().Msgf("server is starting at %s", application.Conf.Address)
	srv := &http.Server{
		Addr:    application.Conf.Address,
		Handler: application.Start(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic().Err(err)
	}
}
