package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	// Server configuration
	Port        string
	Environment string

	// Database configuration
	DatabaseURL string
	RedisURL    string

	// Blockchain configuration
	BlockchainRPCURL string
	BlockchainWSURL  string
	ContractAddress  string
	TokenAddress     string
	PrivateKey       string

	// JWT configuration
	JWTSecret    string
	JWTExpiryHrs int

	// External services
	EtherscanAPIKey  string
	CoinGeckoAPIKey  string

	// Monitoring
	SentryDSN string
	LogLevel  string

	// Rate limiting
	RateLimitRPS   int
	RateLimitBurst int

	// CORS settings
	AllowedOrigins []string
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found, using environment variables")
	}

	config := &Config{
		// Server
		Port:        getEnv("API_PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),

		// Database
		DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost:5432/whitelist_token_db?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),

		// Blockchain
		BlockchainRPCURL: getEnv("BLOCKCHAIN_RPC_URL", "http://localhost:8545"),
		BlockchainWSURL:  getEnv("BLOCKCHAIN_WS_URL", "ws://localhost:8545"),
		ContractAddress:  getEnv("CONTRACT_ADDRESS", ""),
		TokenAddress:     getEnv("TOKEN_ADDRESS", ""),
		PrivateKey:       getEnv("PRIVATE_KEY", ""),

		// JWT
		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpiryHrs: getEnvAsInt("JWT_EXPIRY_HOURS", 24),

		// External services
		EtherscanAPIKey: getEnv("ETHERSCAN_API_KEY", ""),
		CoinGeckoAPIKey: getEnv("COINGECKO_API_KEY", ""),

		// Monitoring
		SentryDSN: getEnv("SENTRY_DSN", ""),
		LogLevel:  getEnv("LOG_LEVEL", "info"),

		// Rate limiting
		RateLimitRPS:   getEnvAsInt("RATE_LIMIT_RPS", 10),
		RateLimitBurst: getEnvAsInt("RATE_LIMIT_BURST", 20),

		// CORS
		AllowedOrigins: getEnvAsSlice("ALLOWED_ORIGINS", []string{"http://localhost:3000", "http://localhost:3001"}),
	}

	// Validate required configuration
	config.validate()

	return config
}

func (c *Config) validate() {
	required := map[string]string{
		"DATABASE_URL":       c.DatabaseURL,
		"BLOCKCHAIN_RPC_URL": c.BlockchainRPCURL,
		"JWT_SECRET":         c.JWTSecret,
	}

	for key, value := range required {
		if value == "" {
			logrus.Fatalf("Required environment variable %s is not set", key)
		}
	}

	// Warn about missing optional but recommended variables
	optional := map[string]string{
		"CONTRACT_ADDRESS": c.ContractAddress,
		"TOKEN_ADDRESS":    c.TokenAddress,
		"PRIVATE_KEY":      c.PrivateKey,
	}

	for key, value := range optional {
		if value == "" {
			logrus.Warnf("Optional environment variable %s is not set", key)
		}
	}
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	
	// Simple comma-separated parsing
	// In production, consider using a more robust parsing method
	result := []string{}
	for _, v := range []string{valueStr} {
		if v != "" {
			result = append(result, v)
		}
	}
	return result
}