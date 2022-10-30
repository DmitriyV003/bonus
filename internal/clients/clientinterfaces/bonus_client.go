package clientinterfaces

type Response struct {
	Code int
}

type OrderDetailsResponse struct {
	Order  string  `json:"order,omitempty"`
	Status string  `json:"status,omitempty"`
	Amount float64 `json:"accrual,omitempty"`
}

type BonusClient interface {
	CreateOrder(orderNumber string) (*Response, error)
	GetOrderDetails(orderNumber string) (*OrderDetailsResponse, error)
}
