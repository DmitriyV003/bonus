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
	Repositories *container.Repositories
	Services     *container.Services
	Conf         *config.Config
}

func (p *Private) Routes() *chi.Mux {
	r := chi.NewRouter()
	luhnValidator := services.NewLuhnOrderNumberValidator()
	bonusClient := clients.NewBonusClient(p.Conf.AccrualAddress)

	authService := services.NewAuthService(p.Conf.JwtSecret, p.Repositories.Users)
	balanceService := services.NewBalanceService(p.Repositories.Payments, p.Repositories.Users)
	paymentService := services.NewPaymentService(p.Repositories.Payments, p.Repositories.Users, balanceService)
	userService := services.NewUserService(luhnValidator, paymentService, p.Repositories.Users, p.Repositories.Payments, authService)
	orderService := services.NewOrderService(luhnValidator, bonusClient, p.Repositories.Orders, p.Repositories.Users, paymentService)

	authMiddleware := middlewares.NewAuthMiddleware(authService, p.Repositories.Users)

	p.Services.OrderService = orderService

	r.Route("/user", func(r chi.Router) {
		r.Post("/register", handlers.NewRegisterHandler(userService).Handle())
		r.Post("/login", handlers.NewLoginHandler(authService).Handle())

		r.With(authMiddleware.Pipe()).Group(func(r chi.Router) {
			r.Post("/orders", handlers.NewCreateOrderHandler(orderService).Handle())
			r.Get("/orders", handlers.NewUserOrdersHandler(orderService).Handle())
			r.Get("/balance", handlers.NewUserBalanceHandler(balanceService).Handle())
			r.Post("/balance/withdraw", handlers.NewWithdrawHandler(userService).Handle())
			r.Get("/withdrawals", handlers.NewUserWithdawsHandler(userService).Handle())

			r.Get("/test", func(writer http.ResponseWriter, request *http.Request) {
				fmt.Println("AUTH gone")
				fmt.Println(services.GetLoggedInUser())
			})
		})

	})

	return r
}
