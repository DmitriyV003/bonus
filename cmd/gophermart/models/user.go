package models

import "time"

type User struct {
	Id        int64      `json:"id"`
	Login     string     `json:"login"`
	Password  string     `json:"-"`
	Balance   int64      `json:"balance"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
