package services

type LuhnOrderNumberValidator struct {
}

func NewLuhnOrderNumberValidator() *LuhnOrderNumberValidator {
	return &LuhnOrderNumberValidator{}
}

func (validator *LuhnOrderNumberValidator) Validate(orderNumber int64) bool {
	return (orderNumber%10+checksum(orderNumber/10))%10 == 0
}

func checksum(number int64) int64 {
	var luhn int64

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 {
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}
