package container

import (
	repository2 "github.com/DmitriyV003/bonus/internal/repository"
)

type Container struct {
	Users    *repository2.UserRepository
	Orders   *repository2.OrderRepository
	Payments *repository2.PaymentRepository
}
