package services

type OrderValidator interface {
	Validate(orderNumber int64) bool
}
