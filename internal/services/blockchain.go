package services

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

// BlockchainService handles blockchain interactions
type BlockchainService struct {
	client          *ethclient.Client
	contractAddress common.Address
	tokenAddress    common.Address
	privateKey      *ecdsa.PrivateKey
	saleABI         abi.ABI
	tokenABI        abi.ABI
	logger          *logrus.Logger
}

// NewBlockchainService creates a new blockchain service
func NewBlockchainService(rpcURL, contractAddr, tokenAddr string) (*BlockchainService, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to blockchain: %w", err)
	}

	// Parse contract addresses
	var contractAddress, tokenAddress common.Address
	if contractAddr != "" {
		contractAddress = common.HexToAddress(contractAddr)
	}
	if tokenAddr != "" {
		tokenAddress = common.HexToAddress(tokenAddr)
	}

	// Validate network connectivity
	ctx := context.Background()
	if _, err := client.NetworkID(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to blockchain network: %w", err)
	}
	
	// Note: Contract validation is handled at runtime during actual operations

	// Parse ABIs
	saleABI, err := abi.JSON(strings.NewReader(WhitelistSaleABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse sale ABI: %w", err)
	}

	tokenABI, err := abi.JSON(strings.NewReader(WhitelistTokenABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse token ABI: %w", err)
	}

	return &BlockchainService{
		client:          client,
		contractAddress: contractAddress,
		tokenAddress:    tokenAddress,
		saleABI:         saleABI,
		tokenABI:        tokenABI,
		logger:          logrus.New(),
	}, nil
}

// SetPrivateKey sets the private key for signing transactions
func (bs *BlockchainService) SetPrivateKey(privateKeyHex string) error {
	if privateKeyHex == "" {
		return nil // No private key provided, read-only mode
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	bs.privateKey = privateKey
	return nil
}

// GetSaleInfo retrieves current sale information from the smart contract
func (bs *BlockchainService) GetSaleInfo(ctx context.Context) (*SaleInfo, error) {
	if bs.contractAddress == (common.Address{}) {
		return nil, fmt.Errorf("contract address not set")
	}

	// Create a call context
	callOpts := &bind.CallOpts{Context: ctx}

	// Get sale configuration
	saleConfig, err := bs.callContract(callOpts, "saleConfig")
	if err != nil {
		return nil, fmt.Errorf("failed to get sale config: %w", err)
	}

	// Get additional sale data
	totalSold, err := bs.callContract(callOpts, "totalSold")
	if err != nil {
		return nil, fmt.Errorf("failed to get total sold: %w", err)
	}

	totalEthRaised, err := bs.callContract(callOpts, "totalEthRaised")
	if err != nil {
		return nil, fmt.Errorf("failed to get total ETH raised: %w", err)
	}

	isSaleActive, err := bs.callContract(callOpts, "isSaleActive")
	if err != nil {
		return nil, fmt.Errorf("failed to get sale active status: %w", err)
	}

	return &SaleInfo{
		TokenPrice:       saleConfig[0].(*big.Int),
		MinPurchase:      saleConfig[1].(*big.Int),
		MaxPurchase:      saleConfig[2].(*big.Int),
		MaxSupply:        saleConfig[3].(*big.Int),
		StartTime:        time.Unix(saleConfig[4].(*big.Int).Int64(), 0),
		EndTime:          time.Unix(saleConfig[5].(*big.Int).Int64(), 0),
		WhitelistRequired: saleConfig[6].(bool),
		TotalSold:        totalSold[0].(*big.Int),
		TotalEthRaised:   totalEthRaised[0].(*big.Int),
		IsActive:         isSaleActive[0].(bool),
	}, nil
}

// GetUserPurchases retrieves purchase information for a user
func (bs *BlockchainService) GetUserPurchases(ctx context.Context, userAddress string) (*UserPurchaseInfo, error) {
	if bs.contractAddress == (common.Address{}) {
		return nil, fmt.Errorf("contract address not set")
	}

	address := common.HexToAddress(userAddress)
	callOpts := &bind.CallOpts{Context: ctx}

	// Get purchase info
	purchaseInfo, err := bs.callContractWithParams(callOpts, "getPurchaseInfo", address)
	if err != nil {
		return nil, fmt.Errorf("failed to get purchase info: %w", err)
	}

	// Get total purchased
	totalPurchased, err := bs.callContractWithParams(callOpts, "totalPurchased", address)
	if err != nil {
		return nil, fmt.Errorf("failed to get total purchased: %w", err)
	}

	return &UserPurchaseInfo{
		Address:        userAddress,
		Amount:         purchaseInfo[0].(*big.Int),
		EthSpent:       purchaseInfo[1].(*big.Int),
		Timestamp:      time.Unix(purchaseInfo[2].(*big.Int).Int64(), 0),
		Claimed:        purchaseInfo[3].(bool),
		TotalPurchased: totalPurchased[0].(*big.Int),
	}, nil
}

// IsWhitelisted checks if an address is whitelisted
func (bs *BlockchainService) IsWhitelisted(ctx context.Context, userAddress string, merkleProof []string) (bool, error) {
	if bs.tokenAddress == (common.Address{}) {
		return false, fmt.Errorf("token address not set")
	}

	address := common.HexToAddress(userAddress)
	callOpts := &bind.CallOpts{Context: ctx}

	// Call the whitelist mapping on the token contract
	contract := bind.NewBoundContract(bs.tokenAddress, bs.tokenABI, bs.client, bs.client, bs.client)
	var result []interface{}
	err := contract.Call(callOpts, &result, "whitelist", address)
	if err != nil {
		return false, fmt.Errorf("failed to check whitelist status: %w", err)
	}

	return result[0].(bool), nil
}

// AddToWhitelist adds addresses to the whitelist (requires admin privileges)
func (bs *BlockchainService) AddToWhitelist(ctx context.Context, addresses []string) (*types.Transaction, error) {
	if bs.privateKey == nil {
		return nil, fmt.Errorf("private key not set")
	}

	// For single address, call updateWhitelist on token contract
	if len(addresses) == 1 {
		address := common.HexToAddress(addresses[0])
		return bs.executeTokenTransaction(ctx, "updateWhitelist", address, true)
	}

	// For batch operations, call updateWhitelistBatch on token contract
	addrs := make([]common.Address, len(addresses))
	for i, addr := range addresses {
		addrs[i] = common.HexToAddress(addr)
	}
	
	return bs.executeTokenTransaction(ctx, "updateWhitelistBatch", addrs, true)
}

// RemoveFromWhitelist removes addresses from the whitelist
func (bs *BlockchainService) RemoveFromWhitelist(ctx context.Context, addresses []string) (*types.Transaction, error) {
	if bs.privateKey == nil {
		return nil, fmt.Errorf("private key not set")
	}

	// For single address, call updateWhitelist on token contract
	if len(addresses) == 1 {
		address := common.HexToAddress(addresses[0])
		return bs.executeTokenTransaction(ctx, "updateWhitelist", address, false)
	}

	// For batch operations, call updateWhitelistBatch on token contract
	addrs := make([]common.Address, len(addresses))
	for i, addr := range addresses {
		addrs[i] = common.HexToAddress(addr)
	}
	
	return bs.executeTokenTransaction(ctx, "updateWhitelistBatch", addrs, false)
}

// PauseSale pauses the token sale
func (bs *BlockchainService) PauseSale(ctx context.Context) (*types.Transaction, error) {
	if bs.privateKey == nil {
		return nil, fmt.Errorf("private key not set")
	}

	return bs.executeTransaction(ctx, "pause")
}

// UnpauseSale unpauses the token sale
func (bs *BlockchainService) UnpauseSale(ctx context.Context) (*types.Transaction, error) {
	if bs.privateKey == nil {
		return nil, fmt.Errorf("private key not set")
	}

	return bs.executeTransaction(ctx, "unpause")
}

// GetTokenBalance gets token balance for an address
func (bs *BlockchainService) GetTokenBalance(ctx context.Context, address string) (*big.Int, error) {
	if bs.tokenAddress == (common.Address{}) {
		return nil, fmt.Errorf("token address not set")
	}

	addr := common.HexToAddress(address)
	callOpts := &bind.CallOpts{Context: ctx}

	// Call balanceOf
	contract := bind.NewBoundContract(bs.tokenAddress, bs.tokenABI, bs.client, bs.client, bs.client)
	var result []interface{}
	err := contract.Call(callOpts, &result, "balanceOf", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to get token balance: %w", err)
	}

	return result[0].(*big.Int), nil
}

// WatchPurchaseEvents watches for token purchase events
func (bs *BlockchainService) WatchPurchaseEvents(ctx context.Context, eventChan chan<- PurchaseEvent) error {
	if bs.contractAddress == (common.Address{}) {
		return fmt.Errorf("contract address not set")
	}

	// Create filter query
	query := ethereum.FilterQuery{
		Addresses: []common.Address{bs.contractAddress},
		Topics: [][]common.Hash{
			{crypto.Keccak256Hash([]byte("TokenPurchase(address,uint256,uint256,uint256)"))},
		},
	}

	// Subscribe to logs
	logs := make(chan types.Log)
	sub, err := bs.client.SubscribeFilterLogs(ctx, query, logs)
	if err != nil {
		return fmt.Errorf("failed to subscribe to logs: %w", err)
	}

	go func() {
		defer sub.Unsubscribe()
		for {
			select {
			case err := <-sub.Err():
				bs.logger.Errorf("Subscription error: %v", err)
				return
			case vLog := <-logs:
				// Parse the log
				event, err := bs.parsePurchaseEvent(vLog)
				if err != nil {
					bs.logger.Errorf("Failed to parse purchase event: %v", err)
					continue
				}
				
				select {
				case eventChan <- *event:
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

// Helper methods

func (bs *BlockchainService) callContract(opts *bind.CallOpts, method string) ([]interface{}, error) {
	contract := bind.NewBoundContract(bs.contractAddress, bs.saleABI, bs.client, bs.client, bs.client)
	var result []interface{}
	err := contract.Call(opts, &result, method)
	return result, err
}

func (bs *BlockchainService) callContractWithParams(opts *bind.CallOpts, method string, params ...interface{}) ([]interface{}, error) {
	contract := bind.NewBoundContract(bs.contractAddress, bs.saleABI, bs.client, bs.client, bs.client)
	var result []interface{}
	err := contract.Call(opts, &result, method, params...)
	return result, err
}

func (bs *BlockchainService) executeTransaction(ctx context.Context, method string, params ...interface{}) (*types.Transaction, error) {
	// Get chain ID
	chainID, err := bs.client.NetworkID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Create transactor
	auth, err := bind.NewKeyedTransactorWithChainID(bs.privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	// Set gas parameters
	auth.GasLimit = uint64(300000) // Adjust based on method
	gasPrice, err := bs.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}
	auth.GasPrice = gasPrice

	// Execute transaction on sale contract
	contract := bind.NewBoundContract(bs.contractAddress, bs.saleABI, bs.client, bs.client, bs.client)
	tx, err := contract.Transact(auth, method, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute transaction: %w", err)
	}

	return tx, nil
}

func (bs *BlockchainService) executeTokenTransaction(ctx context.Context, method string, params ...interface{}) (*types.Transaction, error) {
	// Get chain ID
	chainID, err := bs.client.NetworkID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Create transactor
	auth, err := bind.NewKeyedTransactorWithChainID(bs.privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	// Set gas parameters
	auth.GasLimit = uint64(300000)
	gasPrice, err := bs.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}
	auth.GasPrice = gasPrice

	// Validate method exists in ABI
	if _, exists := bs.tokenABI.Methods[method]; !exists {
		return nil, fmt.Errorf("method %s not found in ABI", method)
	}

	// Execute transaction on token contract
	contract := bind.NewBoundContract(bs.tokenAddress, bs.tokenABI, bs.client, bs.client, bs.client)
	tx, err := contract.Transact(auth, method, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute transaction: %w", err)
	}

	// Wait for the transaction to be mined
	receipt, err := bind.WaitMined(ctx, bs.client, tx)
	if err != nil {
		return tx, err // Return transaction even if waiting fails
	}

	// Check transaction status
	if receipt.Status == 0 {
		return tx, fmt.Errorf("transaction failed")
	}

	bs.logger.Infof("Transaction %s completed successfully (block %d)", tx.Hash().Hex(), receipt.BlockNumber.Uint64())
	return tx, nil
}

func (bs *BlockchainService) parsePurchaseEvent(vLog types.Log) (*PurchaseEvent, error) {
	// This is a simplified version - in practice, you'd use the ABI to properly decode
	return &PurchaseEvent{
		Buyer:       common.BytesToAddress(vLog.Topics[1].Bytes()),
		TokenAmount: new(big.Int).SetBytes(vLog.Data[0:32]),
		EthAmount:   new(big.Int).SetBytes(vLog.Data[32:64]),
		Timestamp:   new(big.Int).SetBytes(vLog.Data[64:96]),
		TxHash:      vLog.TxHash,
		BlockNumber: vLog.BlockNumber,
	}, nil
}

// Data structures

type SaleInfo struct {
	TokenPrice        *big.Int  `json:"token_price"`
	MinPurchase       *big.Int  `json:"min_purchase"`
	MaxPurchase       *big.Int  `json:"max_purchase"`
	MaxSupply         *big.Int  `json:"max_supply"`
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
	WhitelistRequired bool      `json:"whitelist_required"`
	TotalSold         *big.Int  `json:"total_sold"`
	TotalEthRaised    *big.Int  `json:"total_eth_raised"`
	IsActive          bool      `json:"is_active"`
}

type UserPurchaseInfo struct {
	Address        string    `json:"address"`
	Amount         *big.Int  `json:"amount"`
	EthSpent       *big.Int  `json:"eth_spent"`
	Timestamp      time.Time `json:"timestamp"`
	Claimed        bool      `json:"claimed"`
	TotalPurchased *big.Int  `json:"total_purchased"`
}

type PurchaseEvent struct {
	Buyer       common.Address `json:"buyer"`
	TokenAmount *big.Int       `json:"token_amount"`
	EthAmount   *big.Int       `json:"eth_amount"`
	Timestamp   *big.Int       `json:"timestamp"`
	TxHash      common.Hash    `json:"tx_hash"`
	BlockNumber uint64         `json:"block_number"`
}

// ABI definitions (simplified - in practice, load from files or generate with abigen)
const WhitelistSaleABI = `[
	{
		"inputs": [],
		"name": "saleConfig",
		"outputs": [
			{"internalType": "uint256", "name": "tokenPrice", "type": "uint256"},
			{"internalType": "uint256", "name": "minPurchase", "type": "uint256"},
			{"internalType": "uint256", "name": "maxPurchase", "type": "uint256"},
			{"internalType": "uint256", "name": "maxSupply", "type": "uint256"},
			{"internalType": "uint256", "name": "startTime", "type": "uint256"},
			{"internalType": "uint256", "name": "endTime", "type": "uint256"},
			{"internalType": "bool", "name": "whitelistRequired", "type": "bool"}
		],
		"stateMutability": "view",
		"type": "function"
	}
]`

const WhitelistTokenABI = `[
	{
		"inputs": [{"internalType": "address", "name": "account", "type": "address"}],
		"name": "balanceOf",
		"outputs": [{"internalType": "uint256", "name": "", "type": "uint256"}],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [{"internalType": "address", "name": "", "type": "address"}],
		"name": "whitelist",
		"outputs": [{"internalType": "bool", "name": "", "type": "bool"}],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{"internalType": "address", "name": "user", "type": "address"},
			{"internalType": "bool", "name": "status", "type": "bool"}
		],
		"name": "updateWhitelist",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{"internalType": "address[]", "name": "users", "type": "address[]"},
			{"internalType": "bool", "name": "status", "type": "bool"}
		],
		"name": "updateWhitelistBatch",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	}
]`