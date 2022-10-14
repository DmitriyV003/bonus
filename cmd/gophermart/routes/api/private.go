package routes

import (
	"context"
	"fmt"
	"github.com/DmitriyV003/bonus/cmd/gophermart/config"
	"github.com/DmitriyV003/bonus/cmd/gophermart/container"
	"github.com/DmitriyV003/bonus/cmd/gophermart/handlers"
	"github.com/DmitriyV003/bonus/cmd/gophermart/middlewares"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Private struct {
	Container *container.Container
	Conf      *config.Config
}

func (p *Private) Routes(ctx context.Context) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/user", func(r chi.Router) {
		r.Post("/register", handlers.RegisterHandler(p.Container, p.Conf))
		r.Post("/login", handlers.LoginHandler(p.Container, p.Conf))

		r.With(middlewares.AuthMiddleware(p.Container, p.Conf)).Group(func(r chi.Router) {
			//r.Post("/orders", handlers.LoginHandler(p.Container, p.Conf))
			r.Get("/test", func(writer http.ResponseWriter, request *http.Request) {
				fmt.Println("AUTH gone")
			})
		})
	})

	return r
}