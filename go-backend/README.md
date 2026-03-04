# Go Backend Architecture

This document describes the improved architecture and code organization of the Go backend server.

## Project Structure

```
go-backend/
├── server.go                    # Application entry point and HTTP server
├── server_test.go               # Server integration tests
├── data_test.go                 # Data store unit tests
├── data/
│   └── data.json                # Default JSON data file
├── pkg/
│   ├── model/                   # Shared domain models and DTOs
│   │   ├── user.go
│   │   ├── task.go
│   │   ├── request.go
│   │   ├── response.go
│   │   └── health.go
│   ├── store/                   # Data store and persistence
│   │   ├── data.go              # In-memory data + read APIs
│   │   ├── load.go              # JSON load/save implementation
│   │   ├── users.go             # User CRUD on DataStore
│   │   ├── tasks.go             # Task CRUD on DataStore
│   │   └── cache.go             # In-memory cache implementation
│   ├── controller/
│   │   ├── users/               # User HTTP handlers
│   │   └── tasks/               # Task HTTP handlers
│   ├── middleware/
│   │   ├── response.go          # HTTP logging, CORS, error & response helpers
│   │   └── validation.go        # Request validation middleware
│   ├── errors/
│   │   └── errors.go            # Custom error types and handling
│   ├── health/
│   │   └── health.go            # Health check system
│   └── validation/
│       └── validator.go         # Input validation logic
├── go.mod                       # Go module definition
└── README.md                    # Project documentation
```

## Architecture Principles

### 1. Package Organization

- **pkg/**: Reusable library code
- **Separation of Concerns**: Each package has a single responsibility
- **Dependency Injection**: Dependencies are injected rather than hardcoded

### 2. Error Handling

- **Structured Errors**: Custom error types with error codes
- **Error Wrapping**: Using `fmt.Errorf` with `%w` for error chaining
- **HTTP Status Mapping**: Appropriate HTTP status codes for different error types
- **Client-Safe Errors**: Internal errors are not exposed to clients

### 3. Validation

- **Reusable Validators**: Modular validation rules
- **Field-Level Validation**: Detailed validation error messages
- **Early Validation**: Input validation at middleware level
- **Type Safety**: Strong typing for all validation rules

### 4. Caching

- **Thread-Safe Cache**: Using sync.RWMutex for concurrent access
- **Automatic Expiration**: Time-based cache invalidation
- **Cache Statistics**: Monitoring cache performance
- **Selective Invalidation**: Cache invalidation on data mutations

### 5. Health Monitoring

- **Component Checks**: Individual health checks for each component
- **Readiness/Liveness**: Kubernetes-style probes
- **Detailed Reporting**: Comprehensive health information
- **Performance Metrics**: Latency and resource usage tracking

## Error Handling Strategy

### Error Types

```go
// Custom error types with specific error codes
type AppError struct {
    Code       ErrorCode `json:"code"`
    Message    string    `json:"message"`
    HTTPStatus int       `json:"-"`
    Internal   error     `json:"-"`
}
```

### Error Categories

- **Validation Errors**: `VALIDATION_ERROR`, `INVALID_JSON`, `FIELD_TOO_LONG`
- **Data Store Errors**: `NOT_FOUND`, `DUPLICATE`, `DATA_STORE_ERROR`
- **System Errors**: `INTERNAL_ERROR`, `REQUEST_TOO_LARGE`, `INVALID_METHOD`

### Error Handling Pattern

```go
// Create specific errors with context
return errors.NewValidationError("field is required", nil)

// Wrap errors with additional context
return errors.NewDataStoreError("Failed to save data", fmt.Errorf("write error: %w", err))

// Check error types
if errors.IsErrorCode(err, errors.ErrCodeValidation) {
    // Handle validation error
}
```

## Validation System

### Validation Rules

```go
type ValidationRule interface {
    Validate(value interface{}) error
}
```

### Built-in Rules

- **RequiredRule**: Validates non-empty values
- **MaxLengthRule**: Validates string length limits
- **EmailRule**: Validates email format
- **EnumRule**: Validates against allowed values

### Usage Example

```go
validator := validation.NewValidator()
validator.AddField("name", name,
    validation.RequiredRule{},
    validation.MaxLengthRule{MaxLength: 100},
)
```

## Middleware Chain

### Request Processing Pipeline

1. **CORS Middleware**: Adds CORS headers
2. **Error Middleware**: Handles panics and errors
3. **Logging Middleware**: Logs all requests
4. **Validation Middleware**: Validates JSON and request size
5. **Route Handler**: Processes the specific request

### Middleware Benefits

- **Composability**: Easy to add/remove middleware
- **Reusability**: Same middleware across all routes
- **Separation of Concerns**: Each middleware has one responsibility

## Caching Strategy

### Cache Implementation

```go
type CacheEntry struct {
    Data      interface{}
    ExpiresAt time.Time
}
```

### Cache Features

- **TTL Support**: Automatic expiration
- **Pattern Invalidation**: Invalidate cache entries by pattern
- **Statistics**: Track hits, misses, and evictions
- **Thread Safety**: Safe for concurrent access

### Cache Keys

- `users`: All users
- `tasks_status_pending`: Filtered tasks
- `stats`: Calculated statistics

## Health Monitoring

### Health Checks

- **Data Store**: Verifies data accessibility
- **Memory**: Monitors memory usage
- **Cache**: Verifies cache functionality

### Health Endpoints

- `/health`: Comprehensive health report
- `/health/ready`: Readiness probe (Kubernetes)
- `/health/live`: Liveness probe (Kubernetes)

### Health Report Structure

```go
type HealthReport struct {
    Status    string                   `json:"status"`
    Checks    map[string]CheckResult   `json:"checks"`
    Readiness ReadinessStatus          `json:"readiness"`
    Liveness  LivenessStatus           `json:"liveness"`
}
```

## Code Quality Practices

### Naming Conventions

- **Package Names**: Lowercase, single words
- **Exported Names**: PascalCase
- **Private Names**: camelCase
- **Interface Names**: -er suffix (e.g., Validator, Checker)

### Function Design

- **Single Purpose**: Each function does one thing
- **Error Handling**: All functions return errors
- **Documentation**: Comments for complex logic
- **Testing**: Comprehensive test coverage

### Dependencies

- **Minimal Dependencies**: Only necessary external packages
- **Interface-Based**: Depend on interfaces, not implementations
- **Dependency Injection**: Pass dependencies as parameters

## Performance Considerations

### Concurrency

- **RWMutex**: Read-write locks for data access
- **Goroutines**: Background cleanup tasks
- **Atomic Operations**: Cache statistics updates

### Memory Management

- **Copy-on-Read**: Return copies of data structures
- **Cache Limits**: Automatic cleanup of expired entries
- **Memory Monitoring**: Track memory usage in health checks

### I/O Operations

- **Atomic Writes**: Use temp files and rename
- **Buffered Operations**: Efficient file I/O
- **Error Recovery**: Handle corrupted data gracefully

## Testing Strategy

### Unit Tests

- **Package-Level Tests**: Test each package independently
- **Mock Interfaces**: Use interfaces for testability
- **Error Scenarios**: Test all error paths

### Integration Tests

- **API Tests**: Test HTTP endpoints
- **Data Persistence**: Test file operations
- **Cache Behavior**: Test caching logic

### Health Tests

- **Component Tests**: Test individual health checks
- **System Tests**: Test overall health reporting

## Deployment Considerations

### Configuration

- **Environment Variables**: Configurable settings
- **Default Values**: Sensible defaults for development
- **Validation**: Validate configuration at startup

### Graceful Shutdown

- **Signal Handling**: Handle SIGINT/SIGTERM
- **Connection Draining**: Complete in-flight requests
- **Resource Cleanup**: Close files and connections

### Monitoring

- **Structured Logging**: JSON-formatted logs
- **Health Endpoints**: Kubernetes-ready probes
- **Performance Metrics**: Request latency and throughput

## Security Considerations

### Input Validation

- **JSON Validation**: Validate all input JSON
- **Size Limits**: Prevent oversized requests
- **Type Validation**: Ensure correct data types

### Error Exposure

- **Safe Errors**: Don't expose internal details
- **Sanitized Messages**: User-friendly error messages
- **Logging Details**: Log full errors internally

### Data Protection

- **File Permissions**: Appropriate file access rights
- **Atomic Operations**: Prevent data corruption
- **Backup Strategy**: Handle data recovery

This architecture provides a solid foundation for a production-ready Go backend service with proper error handling,
validation, caching, and monitoring capabilities.
