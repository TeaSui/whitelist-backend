-- Simple Database Setup for Whitelist Token Project

-- Create the main tables we need

-- Whitelist table - stores who can buy tokens and how much
CREATE TABLE IF NOT EXISTS whitelist (
    id SERIAL PRIMARY KEY,
    address VARCHAR(42) UNIQUE NOT NULL,  -- Ethereum address
    max_allocation DECIMAL(20,0) NOT NULL,  -- Max tokens they can buy
    current_allocation DECIMAL(20,0) DEFAULT 0,  -- How many they've bought
    created_at TIMESTAMP DEFAULT NOW()
);

-- Purchases table - stores all token purchases
CREATE TABLE IF NOT EXISTS purchases (
    id SERIAL PRIMARY KEY,
    tx_hash VARCHAR(66) UNIQUE NOT NULL,  -- Transaction hash
    buyer_address VARCHAR(42) NOT NULL,   -- Who bought
    token_amount DECIMAL(20,0) NOT NULL,  -- How many tokens
    eth_amount DECIMAL(20,0) NOT NULL,    -- How much ETH paid
    created_at TIMESTAMP DEFAULT NOW()
);

-- Admin users table - for admin login
CREATE TABLE IF NOT EXISTS admin_users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_whitelist_address ON whitelist(address);
CREATE INDEX IF NOT EXISTS idx_purchases_buyer ON purchases(buyer_address);
CREATE INDEX IF NOT EXISTS idx_purchases_tx_hash ON purchases(tx_hash);

-- Insert a default admin user (password: admin123)
INSERT INTO admin_users (username, email, password_hash) 
VALUES ('admin', 'admin@whitelist.com', '$2b$10$rQxLx4xHnqfN8H1qV7Zu2.pGzVYsZfP1Q8v7xH2yfHrW9Z3qX4Y5e')
ON CONFLICT (username) DO NOTHING;

-- Insert some sample whitelist entries for testing
INSERT INTO whitelist (address, max_allocation) VALUES 
('0x742d35Cc6638Bb532c7F4316E7C56C8f69e5D0b5', 1000000000000000000000),  -- 1000 tokens
('0x8ba1f109551bD432803012645Hac136c22C57592', 500000000000000000000),   -- 500 tokens
('0x1234567890123456789012345678901234567890', 2000000000000000000000)   -- 2000 tokens
ON CONFLICT (address) DO NOTHING;