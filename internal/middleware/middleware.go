package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Logger middleware
func Logger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log request
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		logger.WithFields(logrus.Fields{
			"status":     statusCode,
			"method":     method,
			"path":       path,
			"ip":         clientIP,
			"latency":    latency,
			"user_agent": c.Request.UserAgent(),
		}).Info("HTTP Request")
	}
}

// Recovery middleware
func Recovery(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.WithFields(logrus.Fields{
					"error": err,
					"path":  c.Request.URL.Path,
					"ip":    c.ClientIP(),
				}).Error("Panic recovered")

				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// CORS middleware
func CORS() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:3000",
		"http://localhost:3001",
		"https://localhost:3000",
		"https://localhost:3001",
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{
		"Origin",
		"Content-Length",
		"Content-Type",
		"Authorization",
		"Accept",
		"X-Requested-With",
	}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	return cors.New(config)
}

// RateLimit middleware (placeholder)
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement rate limiting
		c.Next()
	}
}

// AuthRequired middleware
func AuthRequired(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenString := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization format",
			})
			c.Abort()
			return
		}

		// For demo purposes, accept our demo tokens
		if len(tokenString) >= 16 && tokenString[:16] == "demo-admin-token" {
			// Extract address from token (simplified for demo)
			c.Set("user_address", "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266")
			c.Set("user_role", "admin")
			c.Next()
			return
		}

		// In production, you would verify JWT token here
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
		})
		c.Abort()
	}
}

// AdminRequired middleware
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists || userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Admin access required",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}