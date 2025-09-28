# WhitelistToken Backend

A Go-based REST API server for the WhitelistToken DApp. Provides blockchain interaction, whitelist management, analytics, and user authentication services.

## 🚀 Features

- **REST API** - Clean HTTP endpoints for frontend integration
- **Blockchain Integration** - Direct interaction with smart contracts
- **Database Management** - PostgreSQL for persistent data storage
- **Authentication** - JWT-based user authentication
- **Whitelist Management** - Admin controls for whitelist operations
- **Analytics** - Token sale and usage metrics
- **Health Monitoring** - Service health checks and monitoring

## 🛠️ Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin (HTTP web framework)
- **Database**: PostgreSQL with GORM ORM
- **Cache**: Redis for session and data caching
- **Blockchain**: go-ethereum for Web3 interaction
- **Authentication**: JWT tokens
- **Logging**: Logrus for structured logging

## 📋 Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- Running Hardhat node (port 8545)
- Deployed smart contracts

## 🚀 Quick Start

### 1. Start Infrastructure Services
```bash
# Start PostgreSQL, Redis, and Hardhat in Docker
docker-compose -f docker-compose-simple.yml up -d

# Check services are healthy
docker-compose -f docker-compose-simple.yml ps
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Environment Configuration
Create `.env` file in the backend directory:
```env
# Server Configuration
API_PORT=8080
ENVIRONMENT=development

# Database Configuration
DATABASE_URL=postgresql://whitelist_user:secure_password@localhost:5432/whitelist_token_db?sslmode=disable

# Redis Configuration
REDIS_URL=redis://localhost:6379

# Blockchain Configuration
BLOCKCHAIN_RPC_URL=http://localhost:8545
CONTRACT_ADDRESS=0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512
TOKEN_ADDRESS=0x5FbDB2315678afecb367f032d93F642f64180aa3
PRIVATE_KEY=ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRE_HOURS=24

# CORS Configuration
CORS_ORIGINS=http://localhost:3000,http://localhost:3001
```

### 3. Database Setup
```bash
# The database will be automatically initialized by docker-compose
# Or manually run the init script:
psql -U whitelist_user -d whitelist_token_db -f ../database/init.sql
```

### 4. Run the Server
```bash
# Development mode
go run cmd/server/main.go

# Build and run
go build -o bin/server cmd/server/main.go
./bin/server
```

The API will be available at `http://localhost:8080`

## 🏗️ Infrastructure Services

This repository includes infrastructure setup for the entire WhitelistToken DApp:

### Services Included
- **PostgreSQL** (port 5432) - Database
- **Redis** (port 6379) - Cache
- **Hardhat Node** (port 8545) - Local blockchain

### Service Management
```bash
# Start all services
docker-compose -f docker-compose-simple.yml up -d

# Stop all services
docker-compose -f docker-compose-simple.yml down

# View logs
docker-compose -f docker-compose-simple.yml logs

# Check service status
docker-compose -f docker-compose-simple.yml ps
```

## 📁 Project Structure

```
backend/
├── cmd/
│   └── server/          # Application entry point
│       └── main.go
├── internal/            # Internal packages
│   ├── config/         # Configuration management
│   ├── database/       # Database connection and models
│   ├── handlers/       # HTTP request handlers
│   ├── middleware/     # HTTP middleware
│   ├── models/         # Data models
│   └── services/       # Business logic services
│       ├── analytics.go
│       ├── auth.go
│       ├── blockchain.go
│       └── whitelist.go
├── migrations/         # Database migrations
├── go.mod             # Go modules file
├── go.sum             # Go modules checksum
└── .env               # Environment variables
```

## 🔌 API Endpoints

### Health Check
```
GET /health - Service health status
```

### Authentication
```
POST /api/auth/login     - User login
POST /api/auth/register  - User registration
POST /api/auth/refresh   - Refresh JWT token
```

### Whitelist Management
```
GET    /api/v1/whitelist/status/:address  - Check whitelist status
POST   /api/v1/whitelist/add             - Add address to whitelist (admin)
POST   /api/v1/whitelist/remove          - Remove address from whitelist (admin)
GET    /api/v1/whitelist/list            - Get all whitelisted addresses (admin)
```

### Token Information
```
GET /api/v1/token/info            - Get token contract information
GET /api/v1/token/balance/:address - Get token balance for address
GET /api/v1/token/supply          - Get token supply information
```

### Sale Management
```
GET /api/v1/sale/info         - Get sale contract information
GET /api/v1/sale/stats        - Get sale statistics
POST /api/v1/sale/purchase    - Purchase tokens (authenticated)
```

### Analytics
```
GET /api/v1/analytics/overview    - General analytics overview
GET /api/v1/analytics/transactions - Transaction history
GET /api/v1/analytics/users       - User statistics (admin)
```

## 📊 Example API Usage

### Check Health
```bash
curl http://localhost:8080/health
```

### Check Whitelist Status
```bash
curl http://localhost:8080/api/v1/whitelist/status/0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
```

### Get Token Information
```bash
curl http://localhost:8080/api/v1/token/info
```

### Add to Whitelist (Admin)
```bash
curl -X POST http://localhost:8080/api/v1/whitelist/add \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"address": "0x1234567890123456789012345678901234567890"}'
```

## 🗃️ Database Models

### User Model
```go
type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Address   string    `json:"address" gorm:"uniqueIndex"`
    Email     string    `json:"email" gorm:"uniqueIndex"`
    Role      string    `json:"role" gorm:"default:user"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### Whitelist Model
```go
type WhitelistEntry struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Address   string    `json:"address" gorm:"uniqueIndex"`
    IsActive  bool      `json:"is_active" gorm:"default:true"`
    AddedBy   string    `json:"added_by"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### Transaction Model
```go
type Transaction struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    TxHash      string    `json:"tx_hash" gorm:"uniqueIndex"`
    FromAddress string    `json:"from_address"`
    ToAddress   string    `json:"to_address"`
    Amount      string    `json:"amount"`
    Type        string    `json:"type"`
    Status      string    `json:"status"`
    BlockNumber uint64    `json:"block_number"`
    CreatedAt   time.Time `json:"created_at"`
}
```

## 🔐 Authentication & Authorization

### JWT Authentication
- Login returns JWT access token
- Token expires in 24 hours (configurable)
- Protected endpoints require `Authorization: Bearer <token>` header

### Role-Based Access Control
- **User**: Basic access to public endpoints
- **Admin**: Full access including whitelist management
- **Super Admin**: All permissions including user management

### Example Protected Endpoint
```go
func (h *Handler) AddToWhitelist(c *gin.Context) {
    // JWT middleware validates token and sets user context
    user := c.MustGet("user").(*models.User)
    
    if user.Role != "admin" && user.Role != "super_admin" {
        c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
        return
    }
    
    // Handle whitelist addition logic
}
```

## 🔗 Blockchain Integration

### Smart Contract Interaction
```go
// Connect to blockchain
client, err := ethers.Dial(rpcURL)
if err != nil {
    log.Fatal("Failed to connect to blockchain:", err)
}

// Load contract
tokenContract, err := contracts.NewWhitelistToken(tokenAddress, client)
if err != nil {
    log.Fatal("Failed to load token contract:", err)
}

// Call contract method
isWhitelisted, err := tokenContract.IsWhitelisted(address)
```

### Event Monitoring
The backend monitors blockchain events for:
- Token transfers
- Whitelist updates
- Sale purchases
- Mint events

## 📊 Services Architecture

### Analytics Service
- Track token transactions
- Monitor sale progress
- Generate usage reports
- Calculate ROI metrics

### Blockchain Service
- Smart contract interaction
- Event monitoring
- Transaction broadcasting
- Balance checking

### Whitelist Service
- Address validation
- Batch operations
- Admin controls
- Status tracking

### Auth Service
- User registration/login
- JWT token management
- Role-based permissions
- Session handling

## 🔧 Configuration Management

### Environment Variables
All configuration is handled through environment variables with sensible defaults:

```go
type Config struct {
    Server struct {
        Port string `default:"8080"`
        Environment string `default:"development"`
    }
    Database struct {
        URL string `required:"true"`
    }
    Blockchain struct {
        RpcURL string `required:"true"`
        TokenAddress string `required:"true"`
        SaleAddress string `required:"true"`
        PrivateKey string `required:"true"`
    }
    JWT struct {
        Secret string `required:"true"`
        ExpireHours int `default:"24"`
    }
}
```

## 🚨 Error Handling

### HTTP Error Responses
```go
type ErrorResponse struct {
    Error   string `json:"error"`
    Code    int    `json:"code"`
    Details string `json:"details,omitempty"`
}
```

### Common Error Patterns
- **400 Bad Request** - Invalid input data
- **401 Unauthorized** - Missing or invalid JWT token
- **403 Forbidden** - Insufficient permissions
- **404 Not Found** - Resource not found
- **500 Internal Server Error** - Server-side errors

## 📝 Logging

### Structured Logging with Logrus
```go
log.WithFields(log.Fields{
    "user_id": userID,
    "address": address,
    "action": "whitelist_add",
}).Info("Address added to whitelist")
```

### Log Levels
- **Debug**: Development debugging
- **Info**: General information
- **Warn**: Warning conditions
- **Error**: Error conditions
- **Fatal**: Critical errors

## 🧪 Testing

### Run Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/services/...
```

### Test Structure
```
backend/
├── internal/
│   ├── handlers/
│   │   └── handlers_test.go
│   └── services/
│       ├── auth_test.go
│       ├── blockchain_test.go
│       └── whitelist_test.go
```

## 🚀 Deployment

### Build for Production
```bash
# Build binary
go build -o bin/server cmd/server/main.go

# Build with optimizations
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/server cmd/server/main.go
```

### Docker Deployment
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
CMD ["./server"]
```

## 🤝 Contributing

### Development Guidelines
1. Follow Go conventions and best practices
2. Add tests for new functionality
3. Use structured logging with appropriate levels
4. Handle errors gracefully with proper HTTP status codes
5. Document API endpoints with proper examples
6. Validate all input data
7. Use meaningful commit messages

### Code Review Checklist
- [ ] Tests added/updated
- [ ] Error handling implemented
- [ ] Logging added for important operations
- [ ] Input validation included
- [ ] Documentation updated
- [ ] No sensitive data in logs

## 🔗 Related Documentation

- [Gin Web Framework](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [go-ethereum Documentation](https://geth.ethereum.org/docs/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)