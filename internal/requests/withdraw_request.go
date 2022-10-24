package requests

type WithdrawRequest struct {
	Order string  `json:"order" validate:"required"`
	Sum   float64 `json:"sum" validate:"required"`
}
