package service

import (
	"errors"
	"testing"

	"wallet_service/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) CreateWallet() (*models.Wallet, error) {
	args := m.Called()
	return args.Get(0).(*models.Wallet), args.Error(1)
}

func (m *MockWalletRepository) GetWalletByID(id uuid.UUID) (*models.Wallet, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Wallet), args.Error(1)
}

func (m *MockWalletRepository) UpdateWalletBalance(id uuid.UUID, newBalance float64) error {
	args := m.Called(id, newBalance)
	return args.Error(0)
}

func (m *MockWalletRepository) Deposit(walletID uuid.UUID, amount float64) error {
	args := m.Called(walletID, amount)
	return args.Error(0)
}

func (m *MockWalletRepository) Withdraw(walletID uuid.UUID, amount float64) error {
	args := m.Called(walletID, amount)
	return args.Error(0)
}

func (m *MockWalletRepository) PerformOperation(walletID uuid.UUID, operationType models.OperationType, amount float64) error {
	args := m.Called(walletID, operationType, amount)
	return args.Error(0)
}

func (m *MockWalletRepository) Transfer(fromWalletID, toWalletID uuid.UUID, amount float64) error {
	args := m.Called(fromWalletID, toWalletID, amount)
	return args.Error(0)
}

func TestWalletService_GetWalletBalance(t *testing.T) {
	mockRepo := &MockWalletRepository{}
	service := NewWalletService(mockRepo)

	walletID := uuid.New()
	expectedWallet := &models.Wallet{ID: walletID, Balance: 100.0}

	mockRepo.On("GetWalletByID", walletID).Return(expectedWallet, nil)

	wallet, err := service.GetWalletBalance(walletID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if wallet.ID != expectedWallet.ID {
		t.Errorf("Expected wallet ID %v, got %v", expectedWallet.ID, wallet.ID)
	}

	if wallet.Balance != expectedWallet.Balance {
		t.Errorf("Expected balance %v, got %v", expectedWallet.Balance, wallet.Balance)
	}

	mockRepo.AssertExpectations(t)
}

func TestWalletService_PerformWalletOperation_Deposit(t *testing.T) {
	mockRepo := &MockWalletRepository{}
	service := NewWalletService(mockRepo)

	walletID := uuid.New()
	amount := 50.0

	mockRepo.On("PerformOperation", walletID, models.DEPOSIT, amount).Return(nil)

	err := service.PerformWalletOperation(walletID, models.DEPOSIT, amount)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	mockRepo.AssertExpectations(t)
}

func TestWalletService_PerformWalletOperation_InvalidAmount(t *testing.T) {
	mockRepo := &MockWalletRepository{}
	service := NewWalletService(mockRepo)

	walletID := uuid.New()
	amount := -50.0

	err := service.PerformWalletOperation(walletID, models.DEPOSIT, amount)
	if err == nil {
		t.Fatal("Expected error for negative amount, got nil")
	}

	if err.Error() != "amount must be positive" {
		t.Errorf("Expected error message 'amount must be positive', got '%v'", err.Error())
	}

	mockRepo.AssertNotCalled(t, "PerformOperation", mock.Anything, mock.Anything, mock.Anything)
}

func TestWalletService_PerformWalletOperation_Withdraw(t *testing.T) {
	mockRepo := &MockWalletRepository{}
	service := NewWalletService(mockRepo)

	walletID := uuid.New()
	amount := 30.0

	mockRepo.On("PerformOperation", walletID, models.WITHDRAW, amount).Return(nil)

	err := service.PerformWalletOperation(walletID, models.WITHDRAW, amount)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	mockRepo.AssertExpectations(t)
}

func TestWalletService_PerformWalletOperation_InsufficientFunds(t *testing.T) {
	mockRepo := &MockWalletRepository{}
	service := NewWalletService(mockRepo)

	walletID := uuid.New()
	amount := 200.0

	mockRepo.On("PerformOperation", walletID, models.WITHDRAW, amount).Return(errors.New("insufficient funds"))

	err := service.PerformWalletOperation(walletID, models.WITHDRAW, amount)
	if err == nil {
		t.Fatal("Expected error for insufficient funds, got nil")
	}

	if err.Error() != "insufficient funds" {
		t.Errorf("Expected error message 'insufficient funds', got '%v'", err.Error())
	}

	mockRepo.AssertExpectations(t)
}

func TestWalletService_Transfer_Success(t *testing.T) {
	mockRepo := &MockWalletRepository{}
	service := NewWalletService(mockRepo)

	fromWalletID := uuid.New()
	toWalletID := uuid.New()
	amount := 50.0

	mockRepo.On("Transfer", fromWalletID, toWalletID, amount).Return(nil)

	err := service.Transfer(fromWalletID, toWalletID, amount)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	mockRepo.AssertExpectations(t)
}

func TestWalletService_Transfer_InvalidAmount(t *testing.T) {
	mockRepo := &MockWalletRepository{}
	service := NewWalletService(mockRepo)

	fromWalletID := uuid.New()
	toWalletID := uuid.New()
	amount := -50.0

	err := service.Transfer(fromWalletID, toWalletID, amount)
	if err == nil {
		t.Fatal("Expected error for negative amount, got nil")
	}

	if err.Error() != "amount must be positive" {
		t.Errorf("Expected error message 'amount must be positive', got '%v'", err.Error())
	}

	mockRepo.AssertNotCalled(t, "Transfer", mock.Anything, mock.Anything, mock.Anything)
}

func TestWalletService_Transfer_SameWallet(t *testing.T) {
	mockRepo := &MockWalletRepository{}
	service := NewWalletService(mockRepo)

	walletID := uuid.New()
	amount := 50.0

	err := service.Transfer(walletID, walletID, amount)
	if err == nil {
		t.Fatal("Expected error for transferring to same wallet, got nil")
	}

	if err.Error() != "cannot transfer to the same wallet" {
		t.Errorf("Expected error message 'cannot transfer to the same wallet', got '%v'", err.Error())
	}

	mockRepo.AssertNotCalled(t, "Transfer", mock.Anything, mock.Anything, mock.Anything)
}

func TestWalletService_Transfer_InsufficientFunds(t *testing.T) {
	mockRepo := &MockWalletRepository{}
	service := NewWalletService(mockRepo)

	fromWalletID := uuid.New()
	toWalletID := uuid.New()
	amount := 200.0

	mockRepo.On("Transfer", fromWalletID, toWalletID, amount).Return(errors.New("insufficient funds"))

	err := service.Transfer(fromWalletID, toWalletID, amount)
	if err == nil {
		t.Fatal("Expected error for insufficient funds, got nil")
	}

	if err.Error() != "insufficient funds" {
		t.Errorf("Expected error message 'insufficient funds', got '%v'", err.Error())
	}

	mockRepo.AssertExpectations(t)
}

func TestWalletService_Transfer_SourceWalletNotFound(t *testing.T) {
	mockRepo := &MockWalletRepository{}
	service := NewWalletService(mockRepo)

	fromWalletID := uuid.New()
	toWalletID := uuid.New()
	amount := 50.0

	mockRepo.On("Transfer", fromWalletID, toWalletID, amount).Return(errors.New("source wallet not found"))

	err := service.Transfer(fromWalletID, toWalletID, amount)
	if err == nil {
		t.Fatal("Expected error for source wallet not found, got nil")
	}

	if err.Error() != "source wallet not found" {
		t.Errorf("Expected error message 'source wallet not found', got '%v'", err.Error())
	}

	mockRepo.AssertExpectations(t)
}

func TestWalletService_Transfer_DestinationWalletNotFound(t *testing.T) {
	mockRepo := &MockWalletRepository{}
	service := NewWalletService(mockRepo)

	fromWalletID := uuid.New()
	toWalletID := uuid.New()
	amount := 50.0

	mockRepo.On("Transfer", fromWalletID, toWalletID, amount).Return(errors.New("destination wallet not found"))

	err := service.Transfer(fromWalletID, toWalletID, amount)
	if err == nil {
		t.Fatal("Expected error for destination wallet not found, got nil")
	}

	if err.Error() != "destination wallet not found" {
		t.Errorf("Expected error message 'destination wallet not found', got '%v'", err.Error())
	}

	mockRepo.AssertExpectations(t)
}
