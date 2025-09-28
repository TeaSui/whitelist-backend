package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"whitelist-token-backend/internal/config"
	"whitelist-token-backend/internal/database"
	"whitelist-token-backend/internal/handlers"
	"whitelist-token-backend/internal/middleware"
	"whitelist-token-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup logger
	logger := logrus.New()
	if cfg.Environment == "production" {
		logger.SetLevel(logrus.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
	} else {
		logger.SetLevel(logrus.DebugLevel)
		gin.SetMode(gin.DebugMode)
	}

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}

	// Run auto-migrations for now
	if err := database.AutoMigrate(db); err != nil {
		logger.Fatalf("Failed to run auto-migrations: %v", err)
	}

	// Initialize Redis
	redisClient, err := database.InitializeRedis(cfg.RedisURL)
	if err != nil {
		logger.Fatalf("Failed to initialize Redis: %v", err)
	}

	// Initialize blockchain service
	blockchainService, err := services.NewBlockchainService(cfg.BlockchainRPCURL, cfg.ContractAddress, cfg.TokenAddress)
	if err != nil {
		logger.Fatalf("Failed to initialize blockchain service: %v", err)
	}

	// Set private key for transaction signing
	if err := blockchainService.SetPrivateKey(cfg.PrivateKey); err != nil {
		logger.Fatalf("Failed to set private key: %v", err)
	}

	// Initialize services
	whitelistService := services.NewWhitelistService(db, redisClient, blockchainService, logger)
	authService := services.NewAuthService(cfg.JWTSecret, logger)
	analyticsService := services.NewAnalyticsService(db, redisClient, logger)

	// Initialize handlers
	handlers := handlers.NewHandlers(
		whitelistService,
		authService,
		analyticsService,
		blockchainService,
		logger,
	)

	// Setup router
	router := setupRouter(cfg, handlers, logger)

	// Setup server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Infof("Server starting on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited")
}

func setupRouter(cfg *config.Config, h *handlers.Handlers, logger *logrus.Logger) *gin.Engine {
	router := gin.New()

	// Middleware
	router.Use(middleware.Logger(logger))
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.CORS())
	router.Use(middleware.RateLimit())

	// Health check
	router.GET("/health", h.HealthCheck)
	router.GET("/metrics", h.Metrics)

	// API routes
	v1 := router.Group("/v1")
	{
		// Public routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", h.Login)
			auth.POST("/verify", h.VerifySignature)
		}

		// Whitelist routes
		whitelist := v1.Group("/whitelist")
		{
			whitelist.GET("/status/:address", h.GetWhitelistStatus)
			whitelist.GET("/verify/:address", h.VerifyWhitelist)
		}

		// Sale routes
		sale := v1.Group("/sale")
		{
			sale.GET("/info", h.GetSaleInfo)
			sale.GET("/purchases/:address", h.GetUserPurchases)
			sale.GET("/stats", h.GetSaleStats)
		}

		// Analytics routes
		analytics := v1.Group("/analytics")
		{
			analytics.GET("/overview", h.GetAnalyticsOverview)
			analytics.GET("/sales", h.GetSalesAnalytics)
			analytics.GET("/users", h.GetUserAnalytics)
		}

		// Protected admin routes
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthRequired(cfg.JWTSecret))
		admin.Use(middleware.AdminRequired())
		{
			admin.POST("/whitelist", h.AddToWhitelist)
			admin.DELETE("/whitelist", h.RemoveFromWhitelist)
			admin.POST("/whitelist/batch", h.BatchUpdateWhitelist)
			admin.GET("/users", h.GetAllUsers)
			admin.PUT("/sale/config", h.UpdateSaleConfig)
			admin.POST("/sale/pause", h.PauseSale)
			admin.POST("/sale/unpause", h.UnpauseSale)
		}
	}

	return router
}