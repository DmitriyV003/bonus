package main

import (
	"context"
	"github.com/DmitriyV003/bonus/internal/application"
	"github.com/DmitriyV003/bonus/internal/config"
	"github.com/DmitriyV003/bonus/internal/container"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	app := application.App{
		Conf:         config.Config{},
		Repositories: &container.Repositories{},
		Services:     &container.Services{},
	}
	app.ApplyConfig()

	defer app.Close()

	log.Info().Msgf("server is starting at %s", app.Conf.Address)
	srv := &http.Server{
		Addr:    app.Conf.Address,
		Handler: app.CreateHandler(),
	}

	go app.Services.OrderService.PollPendingOrders(context.Background())

	log.Panic().Err(srv.ListenAndServe())
}
