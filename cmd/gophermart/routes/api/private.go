package routes

import (
	"context"
	"github.com/DmitriyV003/bonus/cmd/gophermart/container"
	"github.com/DmitriyV003/bonus/cmd/gophermart/handlers"
	"github.com/go-chi/chi/v5"
)

type Private struct {
	Container *container.Container
}

func (p *Private) Routes(ctx context.Context) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/register", handlers.RegisterHandler(p.Container))

	return r
}
