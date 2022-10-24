package container

import (
	"github.com/DmitriyV003/bonus/internal/repository"
	"github.com/DmitriyV003/bonus/internal/services"
)

type Repositories struct {
	Users    *repository.UserRepository
	Orders   *repository.OrderRepository
	Payments *repository.PaymentRepository
}

type Services struct {
	OrderService *services.OrderService
}
