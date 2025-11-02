package service

import (
	"fmt"
	"sync"

	"wallet_service/internal/models"
	"wallet_service/internal/repository"

	"github.com/google/uuid"
)

type WalletService struct {
	repo      repository.WalletRepositoryInterface
	mu        sync.RWMutex
	walletMUs map[uuid.UUID]*sync.Mutex
}

func NewWalletService(repo repository.WalletRepositoryInterface) *WalletService {
	return &WalletService{
		repo:      repo,
		walletMUs: make(map[uuid.UUID]*sync.Mutex),
	}
}

func (s *WalletService) getWalletMutex(walletID uuid.UUID) *sync.Mutex {
	s.mu.Lock()
	defer s.mu.Unlock()

	if mu, exists := s.walletMUs[walletID]; exists {
		return mu
	}

	mu := &sync.Mutex{}
	s.walletMUs[walletID] = mu
	return mu
}

func (s *WalletService) GetWalletBalance(walletID uuid.UUID) (*models.Wallet, error) {
	return s.repo.GetWalletByID(walletID)
}

func (s *WalletService) PerformWalletOperation(walletID uuid.UUID, operationType models.OperationType, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	mu := s.getWalletMutex(walletID)
	mu.Lock()
	defer mu.Unlock()

	return s.repo.PerformOperation(walletID, operationType, amount)
}

func (s *WalletService) CreateWallet() (*models.Wallet, error) {
	return s.repo.CreateWallet()
}

func (s *WalletService) Transfer(fromWalletID, toWalletID uuid.UUID, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	if fromWalletID == toWalletID {
		return fmt.Errorf("cannot transfer to the same wallet")
	}

	var firstMu, secondMu *sync.Mutex
	if fromWalletID.String() < toWalletID.String() {
		firstMu = s.getWalletMutex(fromWalletID)
		secondMu = s.getWalletMutex(toWalletID)
	} else {
		firstMu = s.getWalletMutex(toWalletID)
		secondMu = s.getWalletMutex(fromWalletID)
	}

	firstMu.Lock()
	secondMu.Lock()
	defer firstMu.Unlock()
	defer secondMu.Unlock()

	return s.repo.Transfer(fromWalletID, toWalletID, amount)
}
