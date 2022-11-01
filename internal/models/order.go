package models

import "time"

const NewStatus = "NEW"
const ProcessingStatus = "PROCESSING"
const InvalidStatus = "INVALID"
const ProcessedStatus = "PROCESSED"

type Order struct {
	ID        int64      `json:"id,omitempty"`
	Number    string     `json:"number,omitempty"`
	Status    string     `json:"status,omitempty"`
	Amount    int64      `json:"amount,omitempty"`
	User      *User      `json:"user,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

var NewOrder = func(number string, status string, amount int64, user *User) *Order {
	return &Order{
		Number:    number,
		Amount:    amount,
		Status:    status,
		User:      user,
		CreatedAt: time.Now(),
	}
}
