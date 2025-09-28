package services

import (
	"github.com/sirupsen/logrus"
)

// AuthService handles authentication operations
type AuthService struct {
	jwtSecret string
	logger    *logrus.Logger
}

// NewAuthService creates a new auth service
func NewAuthService(jwtSecret string, logger *logrus.Logger) *AuthService {
	return &AuthService{
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}