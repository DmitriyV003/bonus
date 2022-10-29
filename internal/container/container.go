package container

import (
	"github.com/DmitriyV003/bonus/internal/repository/interfaces"
	"github.com/DmitriyV003/bonus/internal/services"
)

type Repositories struct {
	Users    interfaces.UserRepository
	Orders   interfaces.OrderRepository
	Payments interfaces.PaymentRepository
}

type Services struct {
	OrderService *services.OrderService
}
