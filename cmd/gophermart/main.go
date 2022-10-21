package main

import (
	"github.com/DmitriyV003/bonus/internal/application"
	"github.com/DmitriyV003/bonus/internal/config"
	"github.com/DmitriyV003/bonus/internal/container"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	applicationApp := application.App{
		Conf:      config.Config{},
		Container: &container.Container{},
	}
	applicationApp.Config()

	defer applicationApp.Close()

	log.Info().Msgf("server is starting at %s", applicationApp.Conf.Address)
	srv := &http.Server{
		Addr:    applicationApp.Conf.Address,
		Handler: applicationApp.CreateHandler(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Panic().Err(err)
	}
}
