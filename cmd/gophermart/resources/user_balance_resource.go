package resources

type UserBalanceResource struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func NewUserBalanceResource(curr float64, withdrawn float64) *UserBalanceResource {
	return &UserBalanceResource{
		Current:   curr,
		Withdrawn: withdrawn,
	}
}
