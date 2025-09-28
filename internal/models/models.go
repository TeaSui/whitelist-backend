package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Address     string         `json:"address" gorm:"uniqueIndex;not null"`
	Nonce       string         `json:"nonce" gorm:"not null"`
	IsAdmin     bool           `json:"is_admin" gorm:"default:false"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	LastLoginAt *time.Time     `json:"last_login_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	WhitelistEntry *WhitelistEntry `json:"whitelist_entry,omitempty"`
	Purchases      []Purchase      `json:"purchases,omitempty"`
	ActivityLogs   []ActivityLog   `json:"activity_logs,omitempty"`
}

// WhitelistEntry represents a whitelist entry
type WhitelistEntry struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	UserID        uint           `json:"user_id" gorm:"not null"`
	Address       string         `json:"address" gorm:"uniqueIndex;not null"`
	IsWhitelisted bool           `json:"is_whitelisted" gorm:"default:false"`
	MaxAllocation string         `json:"max_allocation" gorm:"type:decimal(78,0)"` // Using string for big numbers
	UsedAllocation string        `json:"used_allocation" gorm:"type:decimal(78,0);default:0"`
	TxHash        string         `json:"tx_hash"`
	BlockNumber   uint64         `json:"block_number"`
	AddedBy       string         `json:"added_by"`
	AddedAt       time.Time      `json:"added_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// Purchase represents a token purchase
type Purchase struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	UserID          uint           `json:"user_id" gorm:"not null"`
	BuyerAddress    string         `json:"buyer_address" gorm:"not null;index"`
	TokenAmount     string         `json:"token_amount" gorm:"type:decimal(78,0);not null"`
	EthAmount       string         `json:"eth_amount" gorm:"type:decimal(78,0);not null"`
	TokenPrice      string         `json:"token_price" gorm:"type:decimal(78,0);not null"`
	TxHash          string         `json:"tx_hash" gorm:"uniqueIndex;not null"`
	BlockNumber     uint64         `json:"block_number" gorm:"not null"`
	BlockTimestamp  time.Time      `json:"block_timestamp" gorm:"not null"`
	Status          string         `json:"status" gorm:"default:'pending'"`          // pending, confirmed, failed
	ClaimStatus     string         `json:"claim_status" gorm:"default:'unclaimed'"` // unclaimed, claimed
	ClaimedAt       *time.Time     `json:"claimed_at"`
	ClaimTxHash     string         `json:"claim_tx_hash"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// SaleConfig represents the token sale configuration
type SaleConfig struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	TokenPrice        string    `json:"token_price" gorm:"type:decimal(78,0);not null"`
	MinPurchase       string    `json:"min_purchase" gorm:"type:decimal(78,0);not null"`
	MaxPurchase       string    `json:"max_purchase" gorm:"type:decimal(78,0);not null"`
	MaxSupply         string    `json:"max_supply" gorm:"type:decimal(78,0);not null"`
	StartTime         time.Time `json:"start_time" gorm:"not null"`
	EndTime           time.Time `json:"end_time" gorm:"not null"`
	WhitelistRequired bool      `json:"whitelist_required" gorm:"default:true"`
	IsActive          bool      `json:"is_active" gorm:"default:true"`
	IsPaused          bool      `json:"is_paused" gorm:"default:false"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// ActivityLog represents user activity logging
type ActivityLog struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id"`
	Address    string    `json:"address" gorm:"not null;index"`
	Action     string    `json:"action" gorm:"not null"`     // login, purchase, claim, etc.
	Details    string    `json:"details" gorm:"type:text"`   // JSON details
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	TxHash     string    `json:"tx_hash"`
	CreatedAt  time.Time `json:"created_at"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// SystemLog represents system-level logging
type SystemLog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Level     string    `json:"level" gorm:"not null"`        // info, warning, error
	Component string    `json:"component" gorm:"not null"`    // blockchain, api, database
	Message   string    `json:"message" gorm:"not null"`
	Details   string    `json:"details" gorm:"type:text"`     // JSON details
	TxHash    string    `json:"tx_hash"`
	CreatedAt time.Time `json:"created_at"`
}

// DailyStats represents daily statistics
type DailyStats struct {
	ID                    uint      `json:"id" gorm:"primaryKey"`
	Date                  time.Time `json:"date" gorm:"uniqueIndex;not null"`
	TotalUsers            int64     `json:"total_users"`
	NewUsers              int64     `json:"new_users"`
	ActiveUsers           int64     `json:"active_users"`
	WhitelistedUsers      int64     `json:"whitelisted_users"`
	TotalPurchases        int64     `json:"total_purchases"`
	DailyPurchases        int64     `json:"daily_purchases"`
	TotalTokensSold       string    `json:"total_tokens_sold" gorm:"type:decimal(78,0);default:0"`
	DailyTokensSold       string    `json:"daily_tokens_sold" gorm:"type:decimal(78,0);default:0"`
	TotalEthRaised        string    `json:"total_eth_raised" gorm:"type:decimal(78,0);default:0"`
	DailyEthRaised        string    `json:"daily_eth_raised" gorm:"type:decimal(78,0);default:0"`
	AverageTokenPrice     string    `json:"average_token_price" gorm:"type:decimal(78,18);default:0"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// DTO structures for API responses

// UserDTO represents user data transfer object
type UserDTO struct {
	ID          uint       `json:"id"`
	Address     string     `json:"address"`
	IsAdmin     bool       `json:"is_admin"`
	IsActive    bool       `json:"is_active"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

// WhitelistStatusDTO represents whitelist status response
type WhitelistStatusDTO struct {
	Address        string    `json:"address"`
	IsWhitelisted  bool      `json:"is_whitelisted"`
	MaxAllocation  string    `json:"max_allocation"`
	UsedAllocation string    `json:"used_allocation"`
	RemainingAllocation string `json:"remaining_allocation"`
	AddedAt        time.Time `json:"added_at,omitempty"`
}

// PurchaseDTO represents purchase data transfer object
type PurchaseDTO struct {
	ID             uint      `json:"id"`
	BuyerAddress   string    `json:"buyer_address"`
	TokenAmount    string    `json:"token_amount"`
	EthAmount      string    `json:"eth_amount"`
	TokenPrice     string    `json:"token_price"`
	TxHash         string    `json:"tx_hash"`
	BlockNumber    uint64    `json:"block_number"`
	BlockTimestamp time.Time `json:"block_timestamp"`
	Status         string    `json:"status"`
	ClaimStatus    string    `json:"claim_status"`
	ClaimedAt      *time.Time `json:"claimed_at"`
	CreatedAt      time.Time `json:"created_at"`
}

// SaleInfoDTO represents sale information response
type SaleInfoDTO struct {
	TokenPrice        string    `json:"token_price"`
	MinPurchase       string    `json:"min_purchase"`
	MaxPurchase       string    `json:"max_purchase"`
	MaxSupply         string    `json:"max_supply"`
	TotalSold         string    `json:"total_sold"`
	RemainingSupply   string    `json:"remaining_supply"`
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
	WhitelistRequired bool      `json:"whitelist_required"`
	IsActive          bool      `json:"is_active"`
	IsPaused          bool      `json:"is_paused"`
	Progress          float64   `json:"progress"` // Percentage of tokens sold
}

// AnalyticsOverviewDTO represents analytics overview response
type AnalyticsOverviewDTO struct {
	TotalUsers       int64  `json:"total_users"`
	ActiveUsers      int64  `json:"active_users"`
	WhitelistedUsers int64  `json:"whitelisted_users"`
	TotalPurchases   int64  `json:"total_purchases"`
	TotalTokensSold  string `json:"total_tokens_sold"`
	TotalEthRaised   string `json:"total_eth_raised"`
	AverageTokenPrice string `json:"average_token_price"`
	SaleProgress     float64 `json:"sale_progress"`
}

// SalesAnalyticsDTO represents sales analytics response
type SalesAnalyticsDTO struct {
	DailyStats   []DailyStats `json:"daily_stats"`
	TopBuyers    []TopBuyerDTO `json:"top_buyers"`
	HourlyStats  []HourlyStatDTO `json:"hourly_stats"`
}

// TopBuyerDTO represents top buyer information
type TopBuyerDTO struct {
	Address     string `json:"address"`
	TokenAmount string `json:"token_amount"`
	EthAmount   string `json:"eth_amount"`
	PurchaseCount int64 `json:"purchase_count"`
}

// HourlyStatDTO represents hourly statistics
type HourlyStatDTO struct {
	Hour        time.Time `json:"hour"`
	Purchases   int64     `json:"purchases"`
	TokensSold  string    `json:"tokens_sold"`
	EthRaised   string    `json:"eth_raised"`
}