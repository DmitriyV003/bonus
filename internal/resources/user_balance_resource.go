package resources

type UserBalanceResource struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func NewUserBalanceResource(curr int64, withdrawn int64) *UserBalanceResource {
	resource := UserBalanceResource{
		Current:   float64(curr) / 10000,
		Withdrawn: float64(withdrawn) / 10000,
	}

	return &resource
}
