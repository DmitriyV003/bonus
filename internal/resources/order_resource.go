package resources

import "time"

type OrderResource struct {
	Number    string    `json:"number"`
	Status    string    `json:"status"`
	Accrual   float64   `json:"accrual"`
	CreatedAt time.Time `json:"uploaded_at"`
}

func NewOrderResource(number string, status string, accrual int64, createdAt time.Time) *OrderResource {
	return &OrderResource{
		Number:    number,
		Status:    status,
		Accrual:   float64(accrual) / 10000,
		CreatedAt: createdAt,
	}
}
