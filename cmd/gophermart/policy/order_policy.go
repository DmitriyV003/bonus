package policy

import "github.com/DmitriyV003/bonus/cmd/gophermart/models"

type OrderPolicy struct {
	order *models.Order
	user  *models.User
}

func NewOrderPolicy(order *models.Order, user *models.User) *OrderPolicy {
	return &OrderPolicy{
		order: order,
		user:  user,
	}
}

func (op *OrderPolicy) Create() bool {
	return op.order.User.Id == op.user.Id
}
