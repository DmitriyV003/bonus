package main

import (
	"context"
	"github.com/DmitriyV003/bonus/internal/application"
	"github.com/DmitriyV003/bonus/internal/clients"
	"github.com/DmitriyV003/bonus/internal/config"
	"github.com/DmitriyV003/bonus/internal/container"
	"github.com/DmitriyV003/bonus/internal/services"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	app := application.App{
		Conf:      config.Config{},
		Container: &container.Container{},
	}
	app.Config()

	defer app.Close()

	log.Info().Msgf("server is starting at %s", app.Conf.Address)
	srv := &http.Server{
		Addr:    app.Conf.Address,
		Handler: app.CreateHandler(),
	}

	orderService := services.NewOrderService(app.Container, nil, clients.NewBonusClient(app.Conf.AccrualAddress))
	go orderService.PollPendingOrders(context.Background())

	if err := srv.ListenAndServe(); err != nil {
		log.Panic().Err(err)
	}
}
