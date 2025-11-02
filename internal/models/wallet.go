package models

import (
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Balance   float64   `json:"balance" db:"balance"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type OperationType string

const (
	DEPOSIT  OperationType = "DEPOSIT"
	WITHDRAW OperationType = "WITHDRAW"
)

type WalletOperation struct {
	WalletID      uuid.UUID     `json:"walletId" db:"wallet_id"`
	OperationType OperationType `json:"operationType" db:"operation_type"`
	Amount        float64       `json:"amount" db:"amount"`
}
