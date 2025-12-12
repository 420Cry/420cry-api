# Backend Refactoring: Microservices Architecture with Dependency Injection

## Overview

This document describes the refactoring of the backend codebase to implement a proper microservices architecture with improved Dependency Injection (DI).

## Goals

1. **Type-Safe DI Container**: Replace string-based service registration with type-safe interfaces
2. **Service Modularity**: Organize services into distinct microservice modules
3. **Interface-Based DI**: Ensure all dependencies are injected via interfaces
4. **Better Separation of Concerns**: Clear boundaries between services
5. **Maintainability**: Easier to test, extend, and maintain

## Architecture Changes

### Before (Old Container)

```go
// String-based registration (not type-safe)
container.Register("userService", userService)
service := container.Get("userService").(UserServiceInterface) // Type assertion needed
```

**Issues:**
- No compile-time type safety
- Services created directly in controllers (2FA, CoinMarketCap, WalletExplorer)
- Inconsistent dependency resolution
- Hard to test and mock

### After (New Container)

```go
// Type-safe getters
service := container.GetUserService() // Returns UserServiceInterface directly
```

**Benefits:**
- Compile-time type safety
- All services registered in container
- Consistent dependency resolution
- Easy to test and mock

## Service Modules

The application is organized into the following microservice modules:

### 1. **User Service Module**
- **Repository**: `UserRepository`, `UserTokenRepository`
- **Services**: `UserService`, `UserTokenService`
- **Responsibilities**: User CRUD operations, token management

### 2. **Auth Service Module**
- **Services**: `AuthService`, `PasswordService`
- **Responsibilities**: Authentication, password hashing/verification, OTP verification

### 3. **Email Service Module**
- **Services**: `EmailService`
- **Components**: `EmailSender`, `EmailCreator`
- **Responsibilities**: Email template creation and sending

### 4. **2FA Service Module**
- **Services**: `TwoFactorService`
- **Responsibilities**: TOTP generation, QR code creation, OTP verification

### 5. **CoinMarketCap Service Module**
- **Services**: `CoinMarketCapService`
- **Responsibilities**: External API integration for market data

### 6. **Wallet Explorer Service Module**
- **Services**: `TransactionService`
- **Responsibilities**: Blockchain transaction data retrieval

## Implementation Details

### New Container Structure

#### `ServiceContainer` (Type-Safe)
```go
type ServiceContainer struct {
    // Core dependencies
    db     *gorm.DB
    config *EnvTypes.EnvConfig
    
    // Repositories
    userRepo      UserRepository.UserRepository
    userTokenRepo UserRepository.UserTokenRepository
    
    // Services (all interfaces)
    passwordService      PasswordService.PasswordServiceInterface
    emailService        EmailService.EmailServiceInterface
    authService         AuthService.AuthServiceInterface
    userTokenService    UserService.UserTokenServiceInterface
    userService         UserService.UserServiceInterface
    twoFactorService    TwoFactorService.TwoFactorServiceInterface
    coinMarketCapService CoinMarketCapService.CoinMarketCapServiceInterface
    transactionService  WalletExplorerService.TransactionServiceInterface
}
```

#### Type-Safe Getters
```go
func (c *ServiceContainer) GetUserService() UserService.UserServiceInterface
func (c *ServiceContainer) GetAuthService() AuthService.AuthServiceInterface
func (c *ServiceContainer) GetEmailService() EmailService.EmailServiceInterface
// ... etc
```

### Service Provider Pattern

For better organization, a service provider pattern is available:

```go
type ServiceProvider interface {
    Register(container *ServiceContainer)
}

// Example providers:
- UserServiceProvider
- AuthServiceProvider
- EmailServiceProvider
- UserBusinessServiceProvider
- TwoFactorServiceProvider
- ExternalAPIServiceProvider
```

### Backward Compatibility

The old `Container` type is maintained for backward compatibility:

```go
type Container struct {
    *ServiceContainer
}
```

This allows existing code to continue working while new code can use the improved `ServiceContainer` directly.

## Migration Guide

### For Controllers

**Before:**
```go
func NewTwoFactorController(container *container.Container) *TwoFactorController {
    userService := container.GetUserService()
    twoFactorService := TwoFactorService.NewTwoFactorService() // Created directly!
    // ...
}
```

**After:**
```go
func NewTwoFactorController(container *container.Container) *TwoFactorController {
    return &TwoFactorController{
        UserService:      container.GetUserService(),
        TwoFactorService: container.GetTwoFactorService(), // From container
        // ...
    }
}
```

### For Services

All services are now registered in the container during initialization:

```go
container := container.InitializeContainer(cfg, db)
// All services are now available via type-safe getters
```

## Benefits

1. **Type Safety**: Compile-time checking prevents runtime errors
2. **Testability**: Easy to mock services for unit testing
3. **Maintainability**: Clear dependency graph and service boundaries
4. **Extensibility**: Easy to add new services or swap implementations
5. **Consistency**: All services follow the same registration pattern

## Testing

With the new structure, testing is easier:

```go
// Create a test container with mocked services
testContainer := &ServiceContainer{
    userService: mockUserService,
    authService: mockAuthService,
    // ...
}
```

## Future Improvements

1. **Wire Integration**: Consider using Google's Wire for code generation
2. **Service Discovery**: For true microservices, add service discovery
3. **Event Bus**: Add event-driven communication between services
4. **API Gateway**: Separate API gateway for routing
5. **Service Mesh**: For inter-service communication and observability

## File Structure

```
app/container/
├── container.go    # Legacy container (backward compatibility)
├── di.go          # New ServiceContainer (type-safe)
└── providers.go   # Service provider pattern
```

## Usage Example

```go
// Initialize container
container := container.InitializeContainer(cfg, db)

// Use in controllers
userController := controllers.NewUserController(container)
twoFactorController := controllers.NewTwoFactorController(container)

// All services are properly injected via interfaces
```

## Notes

- All services implement interfaces for better testability
- Dependencies are resolved in the correct order
- No service is created outside the container
- Backward compatibility is maintained for existing code

