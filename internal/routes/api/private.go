package routes

import (
	"fmt"
	"github.com/DmitriyV003/bonus/internal/clients"
	"github.com/DmitriyV003/bonus/internal/config"
	"github.com/DmitriyV003/bonus/internal/container"
	"github.com/DmitriyV003/bonus/internal/handlers"
	"github.com/DmitriyV003/bonus/internal/middlewares"
	"github.com/DmitriyV003/bonus/internal/services"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Private struct {
	Container *container.Container
	Conf      *config.Config
}

func (p *Private) Routes() *chi.Mux {
	r := chi.NewRouter()
	register := handlers.NewRegisterHandler(p.Conf.JwtSecret)
	withdraw := handlers.NewWithdrawHandler()
	allWithdraw := handlers.NewUserWithdawsHandler(services.NewUserService(p.Container, nil, services.NewLuhnOrderNumberValidator(), services.NewPaymentService(p.Container)))

	r.Route("/user", func(r chi.Router) {
		r.Post("/register", register.Handle(p.Container, p.Conf))
		r.Post("/login", handlers.LoginHandler(p.Container, p.Conf))

		r.With(middlewares.AuthMiddleware(p.Container, p.Conf)).Group(func(r chi.Router) {
			r.Post("/orders", handlers.CreateOrderHandler(p.Container, clients.NewBonusClient(p.Conf.AccrualAddress)))
			r.Get("/orders", handlers.UserOrdersHandler(p.Container))
			r.Get("/balance", handlers.UserBalanceHandler(p.Container))
			r.Post("/balance/withdraw", withdraw.Handle(p.Container))
			r.Get("/balance/withdrawals", allWithdraw.Handle())

			r.Get("/test", func(writer http.ResponseWriter, request *http.Request) {
				fmt.Println("AUTH gone")
				fmt.Println(request.Context().Value("user"))
			})
		})

	})

	return r
}
