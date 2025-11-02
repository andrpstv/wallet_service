package handler

import (
	"net/http"
	"strings"

	"wallet_service/internal/models"
	"wallet_service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WalletHandler struct {
	walletService *service.WalletService
}

func NewWalletHandler(walletService *service.WalletService) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
	}
}

func (h *WalletHandler) GetWalletBalance(c *gin.Context) {
	walletUUIDStr := c.Param("wallet_uuid")
	walletID, err := uuid.Parse(walletUUIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet UUID"})
		return
	}

	wallet, err := h.walletService.GetWalletBalance(walletID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

func (h *WalletHandler) CreateWallet(c *gin.Context) {
	wallet, err := h.walletService.CreateWallet()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet"})
		return
	}

	c.JSON(http.StatusCreated, wallet)
}

func (h *WalletHandler) PerformWalletOperation(c *gin.Context) {
	var req models.WalletOperation
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be positive"})
		return
	}

	if req.OperationType != models.DEPOSIT && req.OperationType != models.WITHDRAW {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation type"})
		return
	}

	err := h.walletService.PerformWalletOperation(req.WalletID, req.OperationType, req.Amount)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
			return
		}
		if strings.Contains(err.Error(), "insufficient funds") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient funds"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Operation successful"})
}
