package models

import "github.com/google/uuid"

type OperationType string

const (
	Deposit  OperationType = "DEPOSIT"
	Withdraw OperationType = "WITHDRAW"
)

type Wallet struct {
	ID      uuid.UUID `json:"walletId"`
	Balance float64   `json:"balance"`
}
