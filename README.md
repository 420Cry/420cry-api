# 420cry-api

This is a Go-based API server for the 420cry application, built with modern Go best practices and enterprise-level patterns following Clean Architecture principles.

## 🏗️ Architecture Overview

The application follows a **Clean Architecture** pattern with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                    HTTP Layer (Gin)                        │
├─────────────────────────────────────────────────────────────┤
│                Middleware & Validation                      │
├─────────────────────────────────────────────────────────────┤
│                   Controllers                              │
├─────────────────────────────────────────────────────────────┤
│                    Services                                │
├─────────────────────────────────────────────────────────────┤
│                  Repositories                              │
├─────────────────────────────────────────────────────────────┤
│                   Database (GORM)                          │
└─────────────────────────────────────────────────────────────┘
```

## 📁 Project Structure

```
app/
├── cmd/                    # Application entry point
├── config/                 # Configuration management
├── container/              # Dependency injection container
├── controllers/            # HTTP request handlers
├── database/               # Database connection and utilities
├── email/                  # Email templates and sending
├── factories/              # Object creation factories
├── logger/                 # Structured logging
├── middleware/             # HTTP middleware
├── models/                 # Data models
├── repositories/           # Data access layer
├── routes/                 # Route registration
├── services/               # Business logic layer
├── types/                  # Type definitions and DTOs
├── utils/                  # Utility functions
└── validators/             # Request validation
```

## 🔧 Key Components

### 1. Dependency Injection Container

The `container` package provides a centralized way to manage dependencies:

```go
// Initialize all services and repositories
container := container.InitializeContainer(cfg, db)

// Get services from container
userService := container.GetUserService()
emailService := container.GetEmailService()
```

### 2. Centralized Error Handling

The `middleware/error_handler.go` provides consistent error responses:

```go
// Custom error types
type ValidationError struct {
    *AppError
    Field string `json:"field,omitempty"`
}

// Automatic error handling
router.Use(middleware.ErrorHandler())
```

### 3. Request Validation

The `validators` package provides structured request validation:

```go
// Validate user signup request
input, err := validators.ValidateUserSignup(c)
if err != nil {
    middleware.AbortWithError(c, err)
    return
}
```

### 4. Structured Logging

The `logger` package provides structured logging with context:

```go
logger := logger.GetLogger()
logger.WithField("user_id", userID).Info("User created successfully")
logger.WithError(err).Error("Database operation failed")
```

### 5. Security Middleware

Multiple security layers protect the application:

```go
router.Use(middleware.SecurityMiddleware())        // Security headers
router.Use(middleware.RateLimitMiddleware())       // Rate limiting
router.Use(middleware.RequestSizeMiddleware(10MB)) // Request size limits
```

## 🚀 Features

### API Versioning
All routes are versioned under `/api/v1/`:

```
POST /api/v1/users/signup
POST /api/v1/users/signin
GET  /api/v1/coin-market-cap/fear-and-greed-latest
```

### Health Check Endpoint
```
GET /health
```

### Enhanced Configuration
Configuration is validated at startup:

```go
// Configuration validation
if err := config.Validate(); err != nil {
    log.Fatalf("Configuration validation failed: %v", err)
}
```

### Database Connection Pooling
Optimized database connections with configurable pooling:

```go
// Configure connection pool
dbConn.ConfigureConnectionPool(25, 5, time.Hour)
```

### Transaction Support
Database operations can use transactions:

```go
err := db.Transaction(func(tx *gorm.DB) error {
    // Multiple operations in transaction
    return nil
})
```

## Prerequisites

- Go 1.23.4

## Preparation

1. **Add Development Hosts to the `/etc/hosts` File**:
    * On **Linux/macOS**, add the following lines to your `/etc/hosts` file.
    * On **Windows**, add them to the `C:\Windows\System32\drivers\etc\hosts` file.

    Add the following lines to the file:
    ```bash
    127.0.0.1 api.420.crypto.test
    127.0.0.1 db.420.crypto.test
    ```
2. **Copy .env.example to .env**:
    ```bash
    cp .env.example .env
    ```

## Installation

1. Clone the repository
2. Install Go dependencies:
    ```bash
    make install
    ```
3. Build the Go application and create a binary:
    ```bash
    make build
    ```
4. Migration:
    ```bash
    make migrate
    ```
5. Run the server:
    ```bash
    make dev
    ```

### Lint:

This project uses `golangci-lint`. You can install it using the following commands based on your OS:

#### macOS:
You can install it with `brew`:
```bash
brew install golangci-lint
```

#### Windows:
You can install go and golangci-lint with choco:
```bash
choco install go golangci-lint
```

## ⚠️ Installing `make` on Windows (Required for Makefile scripts)
```bash
choco install make
```

#### Run Lint:
You can run the linter with:
```bash
make lint
```
This command applies gofumpt and goimports to fix formatting and organize imports.
You need to have both gofumpt and goimports installed and available in your system's PATH.
Install them with:

```bash
go install mvdan.cc/gofumpt@latest
go install golang.org/x/tools/cmd/goimports@latest
```

After that you can run
```bash
make lint-fix
```

### Test:
```bash
make test
```

## ⚠️ Ensure Go Tools Are in Your PATH (MAC OS)

If you encounter a "command not found" error when running `make lint-fix`, it's likely because your Go-installed binaries are not in your system `PATH`.

Add this to your shell profile (`~/.zshrc`, `~/.bashrc`, or `~/.bash_profile`):

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```
Then reload your shell:
```bash
source ~/.zshrc  # or source ~/.bashrc
```

## 🧪 Testing

### Test Suites
The testing architecture has been enhanced with test suites:

```go
// Run user controller tests
func TestUserController(t *testing.T) {
    suite.Run(t, new(UserControllerTestSuite))
}
```

### Test Database
Tests use in-memory SQLite for fast, isolated testing:

```go
// In-memory database for tests
db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
```

## 🔒 Security Enhancements

### Security Headers
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Strict-Transport-Security: max-age=31536000`

### Rate Limiting
- 100 requests per minute per IP
- Configurable limits per endpoint

### Request Validation
- Content-Type validation
- Request size limits
- Input sanitization

## 📊 Monitoring & Observability

### Structured Logging
- JSON-formatted logs
- Contextual information
- Log levels (DEBUG, INFO, WARN, ERROR)

### Request Logging
- HTTP request/response logging
- Performance metrics
- Error tracking

## 🚀 Deployment

### Environment Variables
Required environment variables:

```bash
# Database
DB_HOST=localhost
DB_PORT=3306
DB_DATABASE=420cry-db
DB_USERNAME=420cry-user
DB_PASSWORD=your_password

# Application
API_PORT=8080
APP_ENV=production
LOG_LEVEL=info

# URLs
CRY_APP_URL=https://app.420cry.com
CRY_API_URL=https://api.420cry.com

# Email
NO_REPLY_EMAIL=noreply@420cry.com
SMTP_HOST=smtp.example.com
SMTP_PORT=587

# External APIs
COIN_MARKET_CAP_API=https://pro-api.coinmarketcap.com
COIN_MARKET_CAP_API_KEY=your_api_key
```

### Docker Deployment
The application is ready for Docker deployment with the existing `docker-compose.yaml`.

### With Docker
1. Shutdown the dev server docker compose for this project.
    ```bash
    docker compose down
    ```

2. Build and start application in production mode.
    ```bash
    docker compose build
    ```

3. Start the application in DEV mode.
    ```bash
    docker compose up -d
   ```

### Without Docker
1. Log into MySQL.
    ```bash
    mysql -u root
    ```

2. Create a new user (we use 420cry-user for this project): In the MySQL shell, run the following SQL command to create the new user with a password:
    ```bash
    CREATE USER '420cry-user'@'localhost' IDENTIFIED BY 'Password';
    ```

3. Grant privileges to the new user: Now, grant the necessary privileges to the new user for the 420cry-db database:
    ```bash
    GRANT ALL PRIVILEGES ON `420cry-db`.* TO '420cry-user'@'localhost';
   ```

4. Flush privileges: Apply the changes to the user privileges:
    ```bash
    FLUSH PRIVILEGES;
   ```

5. Exit MySQL: Exit the MySQL shell:
    ```bash
    EXIT;
   ```

6. Verify the new user:
    ```bash
    mysql -u 420cry-user -p
   ```

7. Create the database: Once you're logged in to the MySQL shell, run the following SQL command to create the database:
    ```bash
   CREATE DATABASE `420cry-db`;
   ```

## MailHog
You can access MailHog at 
```bash
    http://localhost:8025/#
```

## 📈 Performance Improvements

1. **Connection Pooling**: Optimized database connections
2. **Middleware Stack**: Efficient request processing
3. **Structured Logging**: Better performance monitoring
4. **Error Handling**: Reduced overhead in error scenarios
5. **Validation**: Early request validation reduces processing

## 🔄 Migration Guide

### From Old Architecture

1. **Update imports**: Use the new container instead of direct database access
2. **Error handling**: Use the new error types and middleware
3. **Validation**: Use the new validators instead of manual validation
4. **Logging**: Replace `log.Printf` with structured logging
5. **Testing**: Update tests to use the new test suites

### Breaking Changes

1. **API Routes**: All routes now require `/api/v1/` prefix
2. **Controller Constructor**: Controllers now use the container
3. **Error Responses**: Error format has changed to be more structured
4. **Configuration**: Some configuration validation is now required

## 🔮 Future Enhancements

1. **Metrics Collection**: Prometheus/Grafana integration
2. **Distributed Tracing**: OpenTelemetry support
3. **Caching Layer**: Redis integration
4. **Message Queue**: Asynchronous processing
5. **API Documentation**: OpenAPI/Swagger generation

## Docs
- [DB Structure Documentation](./docs/db_structure/DB_STRUCTURE.md)
- [API Documentation](./docs/routes/API_DOCS.md)

## Frequently asked questions
### How can I see which application uses a port?
You can easily check this with the command below.
```shell
sudo netstat -tulpn | grep -E "(80|443|3306)"
```

This is very useful if you get an error like
```
ERROR: for dev-server_mysql_1  Cannot start service mysql: Ports are not available: listen tcp 0.0.0.0:3306: bind: address already in use
```
or
```
WARNING: Host is already in use by another container
ERROR: for dev-server_proxy_1  Cannot start service proxy: driver failed program
```

### What should I do if I encounter a port issue?
If you encounter a port issue, you have two options:

1. Stop MySQL locally: If MySQL is running locally on your machine and using the port, you can stop it to free up the port.

2. Update Docker port: You can modify your Docker configuration to use a different port for MySQL.

## 📚 Additional Resources

- [Go Best Practices](https://golang.org/doc/effective_go.html)
- [Gin Framework Documentation](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [Testify Suite Documentation](https://github.com/stretchr/testify/tree/master/suite)