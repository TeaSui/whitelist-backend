package services

import (
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// WhitelistService handles whitelist operations
type WhitelistService struct {
	db               *gorm.DB
	redis            *redis.Client
	blockchainService *BlockchainService
	logger           *logrus.Logger
}

// NewWhitelistService creates a new whitelist service
func NewWhitelistService(
	db *gorm.DB,
	redis *redis.Client,
	blockchainService *BlockchainService,
	logger *logrus.Logger,
) *WhitelistService {
	return &WhitelistService{
		db:               db,
		redis:            redis,
		blockchainService: blockchainService,
		logger:           logger,
	}
}