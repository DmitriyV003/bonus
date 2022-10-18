package models

import "time"

type Order struct {
	Id        int64      `json:"id"`
	Number    string     `json:"number"`
	Status    *string    `json:"status"`
	Amount    int64      `json:"amount"`
	User      *User      `json:"user"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func NewOrder(number string, amount int64, user *User) *Order {
	return &Order{
		Number:    number,
		Amount:    amount,
		User:      user,
		CreatedAt: time.Now(),
	}
}
