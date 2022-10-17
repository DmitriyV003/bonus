package container

import "github.com/DmitriyV003/bonus/cmd/gophermart/repository"

type Container struct {
	Users  *repository.UserRepository
	Orders *repository.OrderRepository
}
