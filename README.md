# Gonely - Social Feed Platform (Golang Recreation)

A complete recreation of the Bunely project using Golang, following Domain-Driven Design (DDD) principles.

## 🏗️ Architecture

This project follows Clean Architecture and DDD principles with clear separation of concerns:

- **Domain Layer**: Core business logic and entities
- **Application Layer**: Use cases and business rules
- **Infrastructure Layer**: External concerns (database, logging, config)
- **Interfaces Layer**: HTTP handlers and DTOs

## 🛠️ Technology Stack

- **Language**: Go 1.22+
- **Web Framework**: Gin
- **ORM**: GORM with PostgreSQL driver
- **Configuration**: Viper
- **Logging**: Zap
- **Authentication**: JWT with bcrypt
- **Validation**: go-playground/validator
- **Testing**: Testify

## 📋 Prerequisites

- Go 1.22 or higher
- PostgreSQL 14 or higher
- Make (optional, for convenience commands)

## 🚀 Quick Start

### 1. Clone and Navigate

```bash
cd backend-golang
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Setup Database

Create PostgreSQL database:

```bash
createdb bunely
# or using psql
psql -U salju -c "CREATE DATABASE bunely;"
```

### 4. Configure Environment

Copy the example environment file and update with your settings:

```bash
cp .env.example .env
```

Edit `.env` and set your configuration:

```env
# Application
APP_NAME=Gonely
APP_ENV=development
APP_URL=http://localhost:8080
APP_PORT=8080
BASE_PATH=/api

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=salju
DB_PASSWORD=password
DB_NAME=bunely
DB_SSL_MODE=disable
DB_TIMEZONE=Asia/Jakarta

# Authentication (IMPORTANT: Change this in production!)
JWT_SECRET=your-very-secure-secret-key-minimum-32-characters-long
JWT_EXPIRY=168h
SESSION_EXPIRY=7d

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Origin,Content-Type,Accept,Authorization
```

### 5. Run the Application

Using Make:
```bash
make run
```

Or directly with Go:
```bash
go run cmd/api/main.go
```

The server will start at `http://localhost:8080`

## 📡 API Endpoints

### Health Check
- `GET /health` - Health status

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login user
- `POST /api/auth/anonymous` - Create anonymous user
- `POST /api/auth/logout` - Logout (requires auth)
- `GET /api/auth/me` - Get current user (requires auth)

### Posts (V1 API)
- `GET /api/v1/posts` - List all posts (public)
- `GET /api/v1/posts/:id` - Get post by ID (public)
- `POST /api/v1/posts` - Create post (requires auth)
- `GET /api/v1/posts/:id/comments` - Get post comments (public)
- `POST /api/v1/posts/:id/comments` - Add comment (requires auth)

## 📝 Example API Calls

### Register User
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "securePassword123"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securePassword123"
  }'
```

### Create Post (with authentication)
```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "content": "My first post!"
  }'
```

### Get All Posts
```bash
curl http://localhost:8080/api/v1/posts?page=0&limit=20
```

## 🧪 Testing

Run all tests:
```bash
make test
```

Run tests with coverage:
```bash
make test-coverage
```

## 🔧 Development Commands

All available Make commands:

```bash
make help          # Display all available commands
make build         # Build the application
make run           # Run the application
make dev           # Run with hot reload (requires air)
make test          # Run tests
make test-coverage # Run tests with coverage report
make lint          # Run linter
make fmt           # Format code
make vet           # Run go vet
make tidy          # Tidy dependencies
make clean         # Clean build artifacts
make docker-up     # Start docker containers
make docker-down   # Stop docker containers
make db-reset      # Reset database
```

## 📁 Project Structure

```
backend-golang/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── domain/                  # Domain layer (business entities)
│   │   ├── auth/                # Authentication domain
│   │   │   ├── user.go
│   │   │   ├── account.go
│   │   │   ├── session.go
│   │   │   ├── verification.go
│   │   │   └── repository.go    # Repository interfaces
│   │   └── feeds/               # Feeds domain
│   │       ├── post.go
│   │       ├── comment.go
│   │       └── repository.go
│   ├── application/             # Application layer (use cases)
│   │   ├── auth/
│   │   │   └── auth_service.go
│   │   └── feeds/
│   │       └── feed_service.go
│   ├── infrastructure/          # Infrastructure layer
│   │   ├── config/
│   │   │   └── config.go        # Configuration management
│   │   ├── database/
│   │   │   └── database.go      # Database connection & migrations
│   │   ├── logger/
│   │   │   └── logger.go        # Logging setup
│   │   └── repository/          # Repository implementations
│   │       ├── auth/
│   │       │   ├── user_repository.go
│   │       │   ├── account_repository.go
│   │       │   └── session_repository.go
│   │       └── feeds/
│   │           ├── post_repository.go
│   │           └── comment_repository.go
│   └── interfaces/              # Interfaces layer (HTTP)
│       └── http/
│           ├── handler/         # HTTP handlers
│           │   ├── auth_handler.go
│           │   └── feed_handler.go
│           ├── middleware/      # HTTP middleware
│           │   ├── auth.go
│           │   ├── error.go
│           │   └── logger.go
│           ├── dto/             # Data Transfer Objects
│           │   ├── auth.go
│           │   ├── post.go
│           │   └── response.go
│           └── router.go        # Route definitions
├── pkg/                         # Shared packages
│   └── errors/
│       └── errors.go            # Custom error types
├── .env.example                 # Example environment variables
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
├── PROJECT_ANALYSIS.md          # Detailed project analysis
└── README.md
```

## 🎯 Key Features

- ✅ **Clean Architecture**: Clear separation of concerns with DDD principles
- ✅ **JWT Authentication**: Secure token-based authentication
- ✅ **Anonymous Users**: Support for anonymous user sessions
- ✅ **GORM Integration**: Type-safe database operations
- ✅ **Structured Logging**: Zap logger with contextual logging
- ✅ **Error Handling**: Centralized error handling middleware
- ✅ **CORS Support**: Configurable CORS for frontend integration
- ✅ **Graceful Shutdown**: Proper cleanup on application termination
- ✅ **Database Migrations**: Automatic schema migrations
- ✅ **Validation**: Request validation with struct tags
- ✅ **Configuration**: Environment-based configuration with Viper

## 🔐 Security Considerations

1. **JWT Secret**: Always use a strong, random secret (minimum 32 characters) in production
2. **Password Hashing**: Uses bcrypt with default cost factor
3. **SQL Injection**: Protected by GORM's parameterized queries
4. **CORS**: Configure allowed origins appropriately for production
5. **Rate Limiting**: Consider adding rate limiting middleware for production

## 🐳 Docker Support (Optional)

Create a `Dockerfile`:

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /gonely cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /gonely .
COPY .env.example .env
EXPOSE 8080
CMD ["./gonely"]
```

## 📊 Database Schema

### Users Table
- id (varchar, PK)
- name (varchar)
- handle (varchar, indexed)
- email (varchar, unique, indexed)
- email_verified (boolean)
- is_anonymous (boolean)
- image (varchar)
- created_at, updated_at (timestamp)

### Accounts Table
- id (varchar, PK)
- account_id (varchar)
- provider_id (varchar)
- user_id (varchar, FK)
- access_token (varchar, indexed)
- password (varchar, hashed)
- created_at, updated_at (timestamp)

### Sessions Table
- id (varchar, PK)
- user_id (varchar, FK)
- token (varchar, unique, indexed)
- user_agent, ip_address (varchar)
- revoked (boolean)
- expires_at (timestamp)
- created_at, updated_at (timestamp)

### Posts Table
- id (varchar, PK)
- content (text)
- created_at (timestamp, indexed)
- updated_at (timestamp)

### Comments Table
- id (varchar, PK)
- post_id (varchar, FK, cascade delete)
- content (text)
- created_at (timestamp, indexed)
- updated_at (timestamp)

## 🤝 Contributing

1. Follow Go best practices and conventions
2. Maintain test coverage above 80%
3. Use meaningful commit messages
4. Update documentation for new features
5. Run `make fmt` and `make lint` before committing

## 📄 License

This project is a recreation for educational purposes.
