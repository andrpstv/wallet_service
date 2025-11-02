package models

import (
	"testing"

	"github.com/google/uuid"
)

func TestWalletOperation(t *testing.T) {
	walletID := uuid.New()
	operation := WalletOperation{
		WalletID:      walletID,
		OperationType: DEPOSIT,
		Amount:        100.0,
	}

	if operation.WalletID != walletID {
		t.Errorf("Expected wallet ID %v, got %v", walletID, operation.WalletID)
	}

	if operation.OperationType != DEPOSIT {
		t.Errorf("Expected operation type DEPOSIT, got %v", operation.OperationType)
	}

	if operation.Amount != 100.0 {
		t.Errorf("Expected amount 100.0, got %v", operation.Amount)
	}
}

func TestOperationType(t *testing.T) {
	if DEPOSIT != "DEPOSIT" {
		t.Errorf("Expected DEPOSIT to be 'DEPOSIT', got %v", DEPOSIT)
	}

	if WITHDRAW != "WITHDRAW" {
		t.Errorf("Expected WITHDRAW to be 'WITHDRAW', got %v", WITHDRAW)
	}
}
