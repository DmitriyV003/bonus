package main

import (
	"github.com/DmitriyV003/bonus/cmd/gophermart/application"
	"github.com/DmitriyV003/bonus/cmd/gophermart/config"
	"github.com/DmitriyV003/bonus/cmd/gophermart/container"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	application := application.App{
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

	if err := srv.ListenAndServe(); err != nil {
		log.Panic().Err(err)
	}
}
