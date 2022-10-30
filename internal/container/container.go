package container

import (
	"github.com/DmitriyV003/bonus/internal/repository/interfaces"
	serviceinterfaces "github.com/DmitriyV003/bonus/internal/services/interfaces"
)

type Repositories struct {
	Users    interfaces.UserRepository
	Orders   interfaces.OrderRepository
	Payments interfaces.PaymentRepository
}

type Services struct {
	OrderService serviceinterfaces.OrderService
}
