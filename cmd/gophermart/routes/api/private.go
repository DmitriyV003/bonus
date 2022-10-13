package routes

import (
	"context"
	"github.com/DmitriyV003/bonus/cmd/gophermart/config"
	"github.com/DmitriyV003/bonus/cmd/gophermart/container"
	"github.com/DmitriyV003/bonus/cmd/gophermart/handlers"
	"github.com/go-chi/chi/v5"
)

type Private struct {
	Container *container.Container
	Conf      *config.Config
}

func (p *Private) Routes(ctx context.Context) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/user/register", handlers.RegisterHandler(p.Container, p.Conf))
	r.Post("/user/login", handlers.LoginHandler(p.Container, p.Conf))

	return r
}
