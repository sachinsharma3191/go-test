# Node.js Backend API

This is a simple Node.js Express API server that provides endpoints for users and tasks management with comprehensive logging and caching capabilities.

## Features

- **RESTful API** for users and tasks management
- **Request Logging** with structured format matching Go backend
- **Backend Proxy** to Go microservice
- **Error Handling** with detailed logging
- **Cache Statistics** monitoring
- **Performance Tracking** with response times

## Setup

1. Install dependencies:
```bash
npm install
```

2. Start the server:
```bash
npm start
```

Or for development with auto-reload:
```bash
npm run dev
```

The server will run on `http://localhost:3000` by default.

## API Endpoints

### Health Check
- `GET /health` - Check if the server is running

### Users
- `GET /api/users` - Get all users
- `GET /api/users/:id` - Get user by ID
- `POST /api/users` - Create a new user
  - Body: `{ "name": "string", "email": "string", "role": "string" }`

### Tasks
- `GET /api/tasks` - Get all tasks (supports query params: `status`, `userId`)
- `GET /api/tasks/:id` - Get task by ID
- `POST /api/tasks` - Create a new task
  - Body: `{ "title": "string", "status": "string", "userId": number }`

### Statistics
- `GET /api/stats` - Get statistics about users and tasks
- `GET /api/stats/cache` - Get cache statistics from Go backend

## Logging

The Node.js backend includes comprehensive structured logging that matches the Go backend format:

### Log Format
```
YYYY-MM-DD HH:MM:SS METHOD PATH STATUS_CODE DURATION
```

### Examples
```
2026-03-03 12:30:45 GET /health 200 15ms
2026-03-03 12:30:46 POST /api/users 201 250ms
2026-03-03 12:30:46 BACKEND GET /api/users 200 180ms
2026-03-03 12:30:47 ERROR POST /api/users 400 0ms - Invalid email format
2026-03-03 12:30:48 CACHE_STATS hits=15 misses=3 evictions=0 entries=5
```

### Log Types
- **Request Logs**: All incoming HTTP requests
- **Backend Logs**: Requests to Go backend service
- **Error Logs**: Errors with full context and stack traces
- **Warning Logs**: 404 errors and other warnings
- **Cache Stats**: Cache performance metrics

### Environment Variables
```bash
# Log level (error, warn, info)
LOG_LEVEL=info

# Server configuration
PORT=3000
GO_BACKEND_URL=http://localhost:8080
```

For detailed logging documentation, see [LOGGING.md](./LOGGING.md).

## Example Requests

```bash
# Health check
curl http://localhost:3000/health

# Get all users
curl http://localhost:3000/api/users

# Get user by ID
curl http://localhost:3000/api/users/1

# Get tasks by status
curl http://localhost:3000/api/tasks?status=pending

# Get statistics
curl http://localhost:3000/api/stats

# Get cache statistics
curl http://localhost:3000/api/stats/cache
```

## Architecture

The Node.js backend serves as an API gateway that:

1. **Receives HTTP requests** from clients
2. **Logs all requests** with timing information
3. **Proxies requests** to Go backend microservice
4. **Logs backend communication** with performance metrics
5. **Returns responses** to clients with appropriate status codes

### Request Flow
```
Client → Node.js Backend → Go Backend → Node.js Backend → Client
```

### Logging Flow
```
Request → Node.js Logs → Backend Request → Go Logs → Backend Response → Node.js Logs → Response
```

## Error Handling

The application includes comprehensive error handling:

- **Input Validation**: Validates request bodies and parameters
- **Backend Errors**: Handles Go backend errors appropriately
- **404 Handling**: Proper 404 responses for unknown routes
- **Error Logging**: Detailed error logging with context

## Performance Monitoring

The logging system enables performance monitoring:

- **Response Times**: Track API response times
- **Backend Performance**: Monitor Go backend response times
- **Error Rates**: Track error frequencies
- **Cache Performance**: Monitor cache hit/miss ratios

## Development

### Project Structure
```
node-backend/
├── controllers/     # Route handlers
├── middleware/      # Express middleware
├── routes/          # API routes
├── utils/           # Utility functions
├── config/          # Configuration
└── app.js           # Main application file
```

### Adding New Endpoints

1. Create controller in `controllers/`
2. Add routes in `routes/`
3. Update main routes in `routes/index.js`
4. Logging is automatically handled by middleware

### Testing

```bash
# Start both backends
cd go-backend && go run . &
cd node-backend && npm start &

# Test endpoints
curl http://localhost:3000/api/users
curl http://localhost:3000/api/stats/cache
```

## Integration with Go Backend

The Node.js backend is designed to work seamlessly with the Go backend:

- **Same Log Format**: Consistent logging across services
- **Request Correlation**: Track requests across service boundaries
- **Error Propagation**: Proper error handling and logging
- **Cache Integration**: Access to Go backend cache statistics

This provides a complete view of the entire request flow from client through both services.
