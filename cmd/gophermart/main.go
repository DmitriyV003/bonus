package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	_ "net/http/pprof"

	"github.com/DmitriyV003/bonus/internal/application"
	"github.com/DmitriyV003/bonus/internal/config"
	"github.com/DmitriyV003/bonus/internal/container"
	"github.com/rs/zerolog/log"
)

func main() {
	app := application.App{
		Conf:         config.Config{},
		Repositories: &container.Repositories{},
		Services:     &container.Services{},
	}
	const pprofPort = ":8082"
	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app.ApplyConfig()

	defer app.Close()

	log.Info().Msgf("server is starting at %s", app.Conf.Address)
	srv := &http.Server{
		Addr:    app.Conf.Address,
		Handler: app.CreateHandler(),
		BaseContext: func(listener net.Listener) context.Context {
			return mainCtx
		},
	}

	g, gCtx := errgroup.WithContext(mainCtx)
	g.Go(func() error {
		return srv.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		log.Warn().Msg("Server down")
		return srv.Shutdown(gCtx)
	})
	g.Go(func() error {
		return app.Services.OrderService.PollPendingOrders(gCtx)
	})

	http.ListenAndServe(pprofPort, nil)

	if err := g.Wait(); err != nil {
		log.Error().Err(err).Msg("Server down")
	}
}
