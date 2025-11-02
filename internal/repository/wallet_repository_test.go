package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func TestWalletRepository_GetWalletByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewWalletRepository(sqlxDB)

	walletID := uuid.New()
	createdAt := time.Now()
	updatedAt := time.Now()
	rows := sqlmock.NewRows([]string{"id", "balance", "created_at", "updated_at"}).
		AddRow(walletID, 100.0, createdAt, updatedAt)

	mock.ExpectQuery("SELECT id, balance, created_at, updated_at FROM wallets WHERE id = \\$1").
		WithArgs(walletID).
		WillReturnRows(rows)

	wallet, err := repo.GetWalletByID(walletID)
	if err != nil {
		t.Fatalf("Failed to get wallet: %v", err)
	}

	if wallet.ID != walletID {
		t.Errorf("Expected wallet ID %v, got %v", walletID, wallet.ID)
	}

	if wallet.Balance != 100.0 {
		t.Errorf("Expected balance 100.0, got %v", wallet.Balance)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestWalletRepository_UpdateWalletBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewWalletRepository(sqlxDB)

	walletID := uuid.New()
	newBalance := 200.0

	mock.ExpectExec("UPDATE wallets SET balance = \\$1, updated_at = NOW\\(\\) WHERE id = \\$2").
		WithArgs(newBalance, walletID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateWalletBalance(walletID, newBalance)
	if err != nil {
		t.Fatalf("Failed to update wallet balance: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
