package models

import "time"

type Order struct {
	Id        int64     `json:"id"`
	Number    int64     `json:"number"`
	Status    string    `json:"status"`
	Amount    int64     `json:"amount"`
	User      *User     `json:"user"`
	UserId    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
