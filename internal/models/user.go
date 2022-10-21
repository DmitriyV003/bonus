package models

import "time"

type User struct {
	Id        int64      `json:"id,omitempty"`
	Login     string     `json:"login,omitempty"`
	Password  string     `json:"-"`
	Balance   int64      `json:"balance,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
