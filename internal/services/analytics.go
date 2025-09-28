package services

import (
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// AnalyticsService handles analytics operations
type AnalyticsService struct {
	db     *gorm.DB
	redis  *redis.Client
	logger *logrus.Logger
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(
	db *gorm.DB,
	redis *redis.Client,
	logger *logrus.Logger,
) *AnalyticsService {
	return &AnalyticsService{
		db:     db,
		redis:  redis,
		logger: logger,
	}
}