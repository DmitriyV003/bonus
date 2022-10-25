package models

import "time"

const WithdrawType = "withdraw"
const AccrualType = "accrual"

const DEBIT = "debit"
const CREDIT = "credit"

type Payment struct {
	ID              int64      `json:"id,omitempty"`
	Type            string     `json:"type,omitempty"`
	TransactionType string     `json:"transaction_type,omitempty"`
	OrderNumber     string     `json:"order_number,omitempty"`
	Amount          int64      `json:"amount,omitempty"`
	User            *User      `json:"user,omitempty"`
	CreatedAt       time.Time  `json:"created_at,omitempty"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
}
