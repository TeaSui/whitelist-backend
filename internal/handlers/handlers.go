package handlers

import (
	"context"
	"net/http"
	"time"

	"whitelist-token-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Handlers contains all HTTP handlers
type Handlers struct {
	whitelistService  *services.WhitelistService
	authService      *services.AuthService
	analyticsService *services.AnalyticsService
	blockchainService *services.BlockchainService
	logger           *logrus.Logger
}

// NewHandlers creates a new handlers instance
func NewHandlers(
	whitelistService *services.WhitelistService,
	authService *services.AuthService,
	analyticsService *services.AnalyticsService,
	blockchainService *services.BlockchainService,
	logger *logrus.Logger,
) *Handlers {
	return &Handlers{
		whitelistService:  whitelistService,
		authService:      authService,
		analyticsService: analyticsService,
		blockchainService: blockchainService,
		logger:           logger,
	}
}

// HealthCheck returns the health status of the API
func (h *Handlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"service": "whitelist-token-backend",
		"version": "1.0.0",
	})
}

// Metrics returns basic metrics (placeholder)
func (h *Handlers) Metrics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"metrics": gin.H{
			"uptime": "placeholder",
			"requests": "placeholder",
		},
	})
}

// Auth handlers
func (h *Handlers) Login(c *gin.Context) {
	var req struct {
		Address   string `json:"address" binding:"required"`
		Message   string `json:"message" binding:"required"`
		Signature string `json:"signature" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// For demo purposes, check if address is the admin address (deployer)
	adminAddress := "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
	if req.Address != adminAddress {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Not authorized as admin",
		})
		return
	}

	// In production, you would verify the signature here
	// For demo, we'll just issue a simple JWT token
	token := "demo-admin-token-" + req.Address + "-" + string(rune(time.Now().Unix()))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"token": token,
			"address": req.Address,
			"role": "admin",
		},
	})
}

func (h *Handlers) VerifySignature(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "verify signature endpoint"})
}

// Whitelist handlers
func (h *Handlers) GetWhitelistStatus(c *gin.Context) {
	address := c.Param("address")
	
	// Validate address format
	if len(address) != 42 || address[:2] != "0x" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Ethereum address format",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// For now, use empty merkle proof - this can be enhanced later
	isWhitelisted, err := h.blockchainService.IsWhitelisted(ctx, address, []string{})
	if err != nil {
		h.logger.WithError(err).WithField("address", address).Error("Failed to check whitelist status")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check whitelist status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"address": address,
			"isWhitelisted": isWhitelisted,
		},
	})
}

func (h *Handlers) VerifyWhitelist(c *gin.Context) {
	address := c.Param("address")
	c.JSON(http.StatusOK, gin.H{
		"address": address,
		"verified": false,
	})
}

// Sale handlers
func (h *Handlers) GetSaleInfo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get token information instead of sale info since we don't have a sale contract
	// Get token balance for the configured addresses
	_, err := h.blockchainService.GetTokenBalance(ctx, "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266")
	if err != nil {
		h.logger.WithError(err).Error("Failed to get token balance")
	}

	// Return mock sale info with token data
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"tokenAddress": "0x0165878A594ca255338adfa4d48449f69242Eb8F",
			"treasury": "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266",
			"tokenPrice": "0", // No sale price
			"minPurchase": "0", 
			"maxPurchase": "1000000000000000000000000", // 1M tokens
			"maxSupply": "1000000000000000000000000000", // 1B tokens
			"startTime": 0,
			"endTime": 0,
			"isPaused": false,
			"isActive": false, // No sale active
			"totalSold": "0",
			"claimEnabled": false,
			"claimStartTime": 0,
		},
	})
}

func (h *Handlers) GetUserPurchases(c *gin.Context) {
	address := c.Param("address")
	c.JSON(http.StatusOK, gin.H{
		"address": address,
		"purchases": []interface{}{},
	})
}

func (h *Handlers) GetSaleStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "sale stats endpoint"})
}

// Analytics handlers (placeholders)
func (h *Handlers) GetAnalyticsOverview(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "analytics overview endpoint"})
}

func (h *Handlers) GetSalesAnalytics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "sales analytics endpoint"})
}

func (h *Handlers) GetUserAnalytics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "user analytics endpoint"})
}

// Admin handlers
func (h *Handlers) AddToWhitelist(c *gin.Context) {
	var req struct {
		Address string `json:"address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate address format
	if len(req.Address) != 42 || req.Address[:2] != "0x" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Ethereum address format",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := h.blockchainService.AddToWhitelist(ctx, []string{req.Address})
	if err != nil {
		h.logger.WithError(err).WithField("address", req.Address).Error("Failed to add to whitelist")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add address to whitelist",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Address added to whitelist successfully",
		"data": gin.H{
			"address": req.Address,
		},
	})
}

func (h *Handlers) RemoveFromWhitelist(c *gin.Context) {
	var req struct {
		Address string `json:"address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate address format
	if len(req.Address) != 42 || req.Address[:2] != "0x" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Ethereum address format",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := h.blockchainService.RemoveFromWhitelist(ctx, []string{req.Address})
	if err != nil {
		h.logger.WithError(err).WithField("address", req.Address).Error("Failed to remove from whitelist")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to remove address from whitelist",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Address removed from whitelist successfully",
		"data": gin.H{
			"address": req.Address,
		},
	})
}

func (h *Handlers) BatchUpdateWhitelist(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "batch update whitelist endpoint"})
}

func (h *Handlers) GetAllUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get all users endpoint"})
}

func (h *Handlers) UpdateSaleConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "update sale config endpoint"})
}

func (h *Handlers) PauseSale(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pause sale endpoint"})
}

func (h *Handlers) UnpauseSale(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "unpause sale endpoint"})
}