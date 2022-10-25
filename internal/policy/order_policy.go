package policy

import (
	models2 "github.com/DmitriyV003/bonus/internal/models"
)

type OrderPolicy struct {
	order *models2.Order
	user  *models2.User
}

func NewOrderPolicy(order *models2.Order, user *models2.User) *OrderPolicy {
	return &OrderPolicy{
		order: order,
		user:  user,
	}
}

func (op *OrderPolicy) Create() bool {
	return op.order.User.ID == op.user.ID
}
