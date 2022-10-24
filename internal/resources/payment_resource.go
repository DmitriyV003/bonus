package resources

import (
	"time"
)

type PaymentResource struct {
	Order     string    `json:"order"`
	Sum       float64   `json:"sum"`
	CreatedAt time.Time `json:"processed_at"`
}

func NewPaymentResource(order string, sum int64, createdAt time.Time) *PaymentResource {
	//createdAtParsed, err := time.Parse(time.RFC3339, createdAt.String())
	//if err != nil {
	//	fmt.Println(err)
	//	return nil
	//}

	return &PaymentResource{
		Order:     order,
		Sum:       float64(sum) / 10000,
		CreatedAt: createdAt,
	}
}
