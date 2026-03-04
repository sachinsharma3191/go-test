# Node.js Backend Logging

This document describes the comprehensive logging system implemented in the Node.js backend that matches the Go backend's structured logging format.

## Features

### Structured Request Logging
- **Format**: `YYYY-MM-DD HH:MM:SS METHOD PATH STATUS_CODE DURATION`
- **Examples**:
  ```
  2026-03-03 12:30:45 GET /health 200 15ms
  2026-03-03 12:30:46 POST /api/users 201 250ms
  2026-03-03 12:30:47 GET /api/tasks 200 89ms
  ```

### Backend Request Logging
- Logs all requests to the Go backend with `BACKEND` prefix
- **Format**: `YYYY-MM-DD HH:MM:SS BACKEND METHOD PATH STATUS_CODE DURATION`
- **Examples**:
  ```
  2026-03-03 12:30:45 BACKEND GET /health 200 12ms
  2026-03-03 12:30:46 BACKEND POST /api/users 201 180ms
  ```

### Error Logging
- Structured error logging with context
- **Format**: `YYYY-MM-DD HH:MM:SS ERROR METHOD PATH STATUS_CODE DURATION - MESSAGE`
- **Examples**:
  ```
  2026-03-03 12:30:45 ERROR POST /api/users 400 0ms - Invalid email format
  2026-03-03 12:30:46 ERROR BACKEND POST /api/tasks 500 0ms - Connection refused
  ```

### Warning Logging
- Logs 404 errors and other warnings
- **Format**: `YYYY-MM-DD HH:MM:SS WARNING METHOD PATH STATUS_CODE DURATION - MESSAGE`
- **Examples**:
  ```
  2026-03-03 12:30:45 WARNING GET /nonexistent 404 0ms - Route not found
  ```

### Cache Statistics Logging
- Logs cache statistics when requested
- **Format**: `YYYY-MM-DD HH:MM:SS CACHE_STATS hits=X misses=Y evictions=Z entries=N`
- **Examples**:
  ```
  2026-03-03 12:30:45 CACHE_STATS hits=15 misses=3 evictions=0 entries=5
  ```

### Server Startup Logging
- Structured startup messages
- **Examples**:
  ```
  2026-03-03 12:30:45 Node.js backend server starting on http://localhost:3000
  2026-03-03 12:30:45 Connecting to Go backend at http://localhost:8080
  2026-03-03 12:30:45 Request logging enabled
  2026-03-03 12:30:45 Health check: http://localhost:3000/health
  ```

## Implementation

### Middleware Components

#### Request Logger (`middleware/logger.js`)
- Captures incoming HTTP requests
- Measures response duration
- Logs in structured format

#### Error Handler (`middleware/errorHandler.js`)
- Logs errors with full context
- Includes stack traces for debugging
- Handles 404 warnings

#### HTTP Client Logger (`utils/httpClient.js`)
- Logs requests to Go backend
- Tracks backend performance
- Logs backend errors separately

### Logger Utility (`utils/logger.js`)
Centralized logging utility with methods:
- `logRequest(method, path, statusCode, duration)`
- `logBackend(method, path, statusCode, duration)`
- `logError(method, path, statusCode, duration, error)`
- `logWarning(method, path, statusCode, duration, message)`
- `logCacheStats(stats)`
- `logStartup(port, goBackendUrl)`

## Log Levels

The logging system supports different log levels via the `LOG_LEVEL` environment variable:
- `error` - Only error messages
- `warn` - Warnings and errors
- `info` - All messages (default)

## Environment Variables

```bash
# Set log level
LOG_LEVEL=info

# Server configuration
PORT=3000
GO_BACKEND_URL=http://localhost:8080
```

## Usage Examples

### Starting the Server
```bash
# Start with default logging
npm start

# Start with specific log level
LOG_LEVEL=warn npm start
```

### Monitoring Logs

The logs provide comprehensive monitoring capabilities:

#### Performance Monitoring
```bash
# Monitor response times
grep "GET /api/users" logs/app.log | tail -10

# Monitor backend performance
grep "BACKEND" logs/app.log | tail -10
```

#### Error Monitoring
```bash
# Monitor errors
grep "ERROR" logs/app.log | tail -10

# Monitor 404s
grep "WARNING" logs/app.log | tail -10
```

#### Cache Performance
```bash
# Monitor cache statistics
grep "CACHE_STATS" logs/app.log | tail -10
```

## Integration with Go Backend

The Node.js logging system is designed to complement the Go backend logging:

1. **Request Flow**: Client → Node.js → Go
2. **Logging Sequence**:
   - Node.js logs incoming request
   - Node.js logs backend request
   - Go logs backend processing
   - Node.js logs response completion

### Example Request Flow
```
# Node.js receives request
2026-03-03 12:30:45 GET /api/users 200 95ms

# Node.js calls Go backend
2026-03-03 12:30:45 BACKEND GET /api/users 200 85ms

# Go processes request
2026/03/03 12:30:45 GET /api/users 200 82.45µs

# Node.js completes response
2026-03-03 12:30:45 GET /api/users 200 95ms
```

## Troubleshooting

### Common Issues

#### Missing Logs
- Check if logging middleware is properly configured
- Verify log level settings
- Ensure proper order of middleware

#### Performance Impact
- Logging is designed to be minimal impact
- Use appropriate log levels in production
- Consider log rotation for high-traffic applications

#### Backend Connection Issues
- Check Go backend status
- Verify `GO_BACKEND_URL` configuration
- Monitor `BACKEND` log entries for connection errors

### Debug Mode

For enhanced debugging, set:
```bash
LOG_LEVEL=debug npm start
```

This will include additional debugging information in the logs.

## Log Analysis Tools

The structured log format enables easy analysis with standard tools:

### Using `grep`
```bash
# Find all slow requests (>100ms)
grep -E "[0-9]{3}ms$" logs/app.log | grep -E "[1-9][0-9]{2}ms$"

# Find error patterns
grep "ERROR" logs/app.log | cut -d' ' -f1-8

# Monitor cache hit ratios
grep "CACHE_STATS" logs/app.log | awk '{print $4, $5}'
```

### Using `awk`
```bash
# Calculate average response time
grep "GET /api/users" logs/app.log | awk '{sum+=$6; count++} END {print "Average:", sum/count, "ms"}'

# Extract error rates
grep -c "ERROR" logs/app.log
```

This comprehensive logging system provides full visibility into the Node.js backend operations and complements the Go backend's logging capabilities.
