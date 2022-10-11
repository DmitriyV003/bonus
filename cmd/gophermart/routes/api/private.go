package routes

import (
	"github.com/DmitriyV003/bonus/cmd/gophermart/container"
	"github.com/go-chi/chi/v5"
)

type Private struct {
	Container *container.Container
}

func (p *Private) Routes() *chi.Mux {
	r := chi.NewRouter()

	return r
}
