package models

import "time"

type Order struct {
	Id        int64      `json:"id,omitempty"`
	Number    string     `json:"number,omitempty"`
	Status    *string    `json:"status,omitempty"`
	Amount    int64      `json:"amount,omitempty"`
	User      *User      `json:"user,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func NewOrder(number string, amount int64, user *User) *Order {
	return &Order{
		Number:    number,
		Amount:    amount,
		User:      user,
		CreatedAt: time.Now(),
	}
}
