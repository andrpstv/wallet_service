package repository

import (
	"database/sql"
	"fmt"

	"wallet_service/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type WalletRepositoryInterface interface {
	CreateWallet() (*models.Wallet, error)
	GetWalletByID(id uuid.UUID) (*models.Wallet, error)
	UpdateWalletBalance(id uuid.UUID, newBalance float64) error
	Deposit(walletID uuid.UUID, amount float64) error
	Withdraw(walletID uuid.UUID, amount float64) error
	PerformOperation(walletID uuid.UUID, operationType models.OperationType, amount float64) error
	Transfer(fromWalletID, toWalletID uuid.UUID, amount float64) error
}

type WalletRepository struct {
	db *sqlx.DB
}

func NewWalletRepository(db *sqlx.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) CreateWallet() (*models.Wallet, error) {
	wallet := &models.Wallet{
		ID:      uuid.New(),
		Balance: 0.0,
	}

	query := `INSERT INTO wallets (id, balance) VALUES ($1, $2)`
	_, err := r.db.Exec(query, wallet.ID, wallet.Balance)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	return wallet, nil
}

func (r *WalletRepository) GetWalletByID(id uuid.UUID) (*models.Wallet, error) {
	var wallet models.Wallet
	query := `SELECT id, balance, created_at, updated_at FROM wallets WHERE id = $1`
	err := r.db.Get(&wallet, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("wallet not found")
		}
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}
	return &wallet, nil
}

func (r *WalletRepository) UpdateWalletBalance(id uuid.UUID, newBalance float64) error {
	query := `UPDATE wallets SET balance = $1, updated_at = NOW() WHERE id = $2`
	result, err := r.db.Exec(query, newBalance, id)
	if err != nil {
		return fmt.Errorf("failed to update wallet balance: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("wallet not found")
	}

	return nil
}

func (r *WalletRepository) Deposit(walletID uuid.UUID, amount float64) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Lock the wallet row for update
	var currentBalance float64
	query := `SELECT balance FROM wallets WHERE id = $1 FOR UPDATE`
	err = tx.Get(&currentBalance, query, walletID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("wallet not found")
		}
		return fmt.Errorf("failed to get wallet balance: %w", err)
	}

	newBalance := currentBalance + amount
	if newBalance < 0 {
		return fmt.Errorf("insufficient funds")
	}

	updateQuery := `UPDATE wallets SET balance = $1, updated_at = NOW() WHERE id = $2`
	_, err = tx.Exec(updateQuery, newBalance, walletID)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *WalletRepository) Withdraw(walletID uuid.UUID, amount float64) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Lock the wallet row for update
	var currentBalance float64
	query := `SELECT balance FROM wallets WHERE id = $1 FOR UPDATE`
	err = tx.Get(&currentBalance, query, walletID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("wallet not found")
		}
		return fmt.Errorf("failed to get wallet balance: %w", err)
	}

	newBalance := currentBalance - amount
	if newBalance < 0 {
		return fmt.Errorf("insufficient funds")
	}

	updateQuery := `UPDATE wallets SET balance = $1, updated_at = NOW() WHERE id = $2`
	_, err = tx.Exec(updateQuery, newBalance, walletID)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *WalletRepository) PerformOperation(walletID uuid.UUID, operationType models.OperationType, amount float64) error {
	switch operationType {
	case models.DEPOSIT:
		return r.Deposit(walletID, amount)
	case models.WITHDRAW:
		return r.Withdraw(walletID, amount)
	default:
		return fmt.Errorf("invalid operation type")
	}
}

func (r *WalletRepository) Transfer(fromWalletID, toWalletID uuid.UUID, amount float64) error {
	if fromWalletID == toWalletID {
		return fmt.Errorf("cannot transfer to the same wallet")
	}

	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Lock the source wallet row for update
	var fromBalance float64
	query := `SELECT balance FROM wallets WHERE id = $1 FOR UPDATE`
	err = tx.Get(&fromBalance, query, fromWalletID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("source wallet not found")
		}
		return fmt.Errorf("failed to get source wallet balance: %w", err)
	}

	if fromBalance < amount {
		return fmt.Errorf("insufficient funds")
	}

	// Lock the destination wallet row for update
	var toBalance float64
	err = tx.Get(&toBalance, query, toWalletID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("destination wallet not found")
		}
		return fmt.Errorf("failed to get destination wallet balance: %w", err)
	}

	// Update balances
	newFromBalance := fromBalance - amount
	newToBalance := toBalance + amount

	updateQuery := `UPDATE wallets SET balance = $1, updated_at = NOW() WHERE id = $2`
	_, err = tx.Exec(updateQuery, newFromBalance, fromWalletID)
	if err != nil {
		return fmt.Errorf("failed to update source wallet balance: %w", err)
	}

	_, err = tx.Exec(updateQuery, newToBalance, toWalletID)
	if err != nil {
		return fmt.Errorf("failed to update destination wallet balance: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
