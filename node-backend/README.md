# Node.js Backend API with Comprehensive Testing

This is a Node.js Express API server that provides endpoints for users and tasks management with comprehensive logging, caching capabilities, and complete test coverage.

## 🎯 **What's Been Done**

### ✅ **Complete Test Coverage**
- **Unit Tests**: Comprehensive testing of all components
- **Integration Tests**: API endpoint testing with real HTTP requests
- **Controller Tests**: Route handler testing with mock services
- **Middleware Tests**: Logging and error handling middleware
- **Utility Tests**: Helper function testing
- **Error Scenarios**: Comprehensive error path testing

### ✅ **API Gateway Architecture**
- **Request Proxying**: Seamless integration with Go backend
- **Structured Logging**: Consistent log format across services
- **Error Handling**: Robust error handling and propagation
- **Performance Monitoring**: Request timing and backend performance
- **Cache Integration**: Access to Go backend cache statistics

### ✅ **Logging System**
- **Structured Format**: Consistent logging with timestamps
- **Request Tracking**: Full request lifecycle logging
- **Backend Communication**: Log all backend requests and responses
- **Error Logging**: Detailed error context and stack traces
- **Performance Metrics**: Response time tracking

### ✅ **Middleware Stack**
- **Request Logging**: Automatic request logging with timing
- **Error Handling**: Centralized error handling middleware
- **CORS Support**: Cross-origin resource sharing
- **Request Validation**: Input validation and sanitization

### ✅ **Development Features**
- **Hot Reload**: Auto-reload in development mode
- **Environment Configuration**: Flexible configuration management
- **Health Monitoring**: Service health checks
- **Statistics API**: Performance and usage statistics

## 📁 **Project Structure**

```
node-backend/
├── package.json                 # Node.js dependencies and scripts
├── package-lock.json            # Dependency lock file
├── app.js                       # Main application entry point
├── server.js                    # Server startup and configuration
├── README.md                    # This file
├── controllers/                 # Route handlers
│   ├── healthController.js      # Health check endpoints
│   ├── healthController.test.js # ✅ Health controller tests
│   ├── statsController.js       # Statistics endpoints
│   ├── statsController.test.js  # ✅ Stats controller tests
│   ├── usersController.js       # User management endpoints
│   ├── usersController.test.js  # ✅ User controller tests
│   ├── tasksController.js       # Task management endpoints
│   └── tasksController.test.js  # ✅ Task controller tests
├── middleware/                  # Express middleware
│   ├── logger.js                # Request logging middleware
│   ├── logger.test.js           # ✅ Logger middleware tests
│   ├── errorHandler.js         # Error handling middleware
│   └── errorHandler.test.js    # ✅ Error handler tests
├── routes/                      # API route definitions
│   ├── health.js                # Health check routes
│   ├── stats.js                 # Statistics routes
│   ├── users.js                 # User management routes
│   └── tasks.js                 # Task management routes
├── services/                    # Business logic services
│   ├── goBackendService.js      # Go backend integration
│   ├── goBackendService.test.js # ✅ Backend service tests
│   ├── cacheService.js          # Cache statistics service
│   └── cacheService.test.js     # ✅ Cache service tests
├── utils/                       # Utility functions
│   ├── logger.js                # Logging utilities
│   ├── logger.test.js           # ✅ Logger utility tests
│   ├── validator.js             # Input validation
│   └── validator.test.js        # ✅ Validation tests
├── config/                      # Configuration management
│   ├── index.js                 # Configuration loader
│   └── index.test.js            # ✅ Configuration tests
└── tests/                       # Integration tests
    ├── integration.test.js       # ✅ API integration tests
    └── setup.js                 # Test setup utilities
```

## 🧪 **Testing Summary**

### **Test Coverage by Component**

#### **Controllers (4 test files)**
- ✅ **Health Controller**: Health check endpoints
- ✅ **Stats Controller**: Statistics and cache endpoints
- ✅ **Users Controller**: User CRUD operations
- ✅ **Tasks Controller**: Task CRUD operations

#### **Middleware (2 test files)**
- ✅ **Logger Middleware**: Request logging and timing
- ✅ **Error Handler**: Error handling and response formatting

#### **Services (2 test files)**
- ✅ **Go Backend Service**: Backend integration and proxying
- ✅ **Cache Service**: Cache statistics retrieval

#### **Utilities (2 test files)**
- ✅ **Logger Utilities**: Logging format and utilities
- ✅ **Validator**: Input validation functions

#### **Configuration (1 test file)**
- ✅ **Config Management**: Environment variable handling

#### **Integration (1 test file)**
- ✅ **API Integration**: End-to-end API testing

### **Test Statistics**
- **Total Test Files**: 12
- **Test Functions**: 50+
- **Coverage Areas**:
  - ✅ HTTP Endpoints
  - ✅ Error Handling
  - ✅ Request Validation
  - ✅ Backend Integration
  - ✅ Logging Functionality
  - ✅ Configuration Management

## 🚀 **Getting Started**

### **Prerequisites**
- Node.js 16 or higher
- npm (comes with Node.js)
- Git

### **Installation**
```bash
# Clone the repository
git clone https://github.com/sachinsharma3191/go-test.git
cd go-test/node-backend

# Install dependencies
npm install

# Start the server
npm start
```

### **Development Mode**
```bash
# Start with auto-reload
npm run dev
```

### **Running Tests**
```bash
# Run all tests
npm test

# Run tests with coverage
npm run test:coverage

# Run tests in watch mode
npm run test:watch

# Run specific test file
npm test controllers/usersController.test.js
```

## 🌐 **API Endpoints**

### **Health Check**
- `GET /health` - Check if the server is running

### **Users**
- `GET /api/users` - Get all users (proxied to Go backend)
- `GET /api/users/:id` - Get user by ID (proxied to Go backend)
- `POST /api/users` - Create a new user (proxied to Go backend)
  - Body: `{ "name": "string", "email": "string", "role": "string" }`

### **Tasks**
- `GET /api/tasks` - Get all tasks (proxied to Go backend)
- `GET /api/tasks/:id` - Get task by ID (proxied to Go backend)
- `POST /api/tasks` - Create a new task (proxied to Go backend)
  - Body: `{ "title": "string", "status": "string", "userId": number }`

### **Statistics**
- `GET /api/stats` - Get statistics about users and tasks
- `GET /api/stats/cache` - Get cache statistics from Go backend

## 📊 **Logging System**

### **Log Format**
```
YYYY-MM-DD HH:MM:SS METHOD PATH STATUS_CODE DURATION
```

### **Log Examples**
```
2026-03-03 12:30:45 GET /health 200 15ms
2026-03-03 12:30:46 POST /api/users 201 250ms
2026-03-03 12:30:46 BACKEND GET /api/users 200 180ms
2026-03-03 12:30:47 ERROR POST /api/users 400 0ms - Invalid email format
2026-03-03 12:30:48 CACHE_STATS hits=15 misses=3 evictions=0 entries=5
```

### **Log Types**
- **Request Logs**: All incoming HTTP requests
- **Backend Logs**: Requests to Go backend service
- **Error Logs**: Errors with full context and stack traces
- **Warning Logs**: 404 errors and other warnings
- **Cache Stats**: Cache performance metrics

## 🔧 **Configuration**

### **Environment Variables**
```bash
# Server configuration
PORT=3000
GO_BACKEND_URL=http://localhost:8080

# Logging configuration
LOG_LEVEL=info

# Node environment
NODE_ENV=development
```

### **Default Configuration**
- **Server**: Port 3000
- **Go Backend**: http://localhost:8080
- **Log Level**: info
- **Environment**: development

## 🏗️ **Architecture**

### **Request Flow**
```
Client → Node.js Backend → Go Backend → Node.js Backend → Client
```

### **Component Architecture**
```
┌─────────────┐    ┌──────────────┐    ┌─────────────┐
│   Client    │───▶│ Node.js API  │───▶│ Go Backend  │
│             │    │   Gateway    │    │  Service    │
│             │◀───│              │◀───│             │
└─────────────┘    └──────────────┘    └─────────────┘
       │                   │                   │
       ▼                   ▼                   ▼
   HTTP Requests    Express Routes    HTTP Requests
                      Middleware       JSON Responses
                      Controllers      File Storage
                      Services         In-Memory Cache
                      Logging          Health Checks
```

## 🧪 **Testing Strategy**

### **Unit Testing**
- **Controllers**: Test route handlers with mock services
- **Services**: Test business logic in isolation
- **Middleware**: Test request/response processing
- **Utilities**: Test helper functions

### **Integration Testing**
- **API Endpoints**: Test full HTTP request/response cycle
- **Backend Integration**: Test Go backend communication
- **Error Scenarios**: Test error handling paths
- **Performance**: Test response times and logging

### **Test Features**
- **Mock Services**: Isolated unit testing
- **HTTP Assertions**: Request/response validation
- **Error Simulation**: Test error conditions
- **Logging Verification**: Ensure proper logging

## 📈 **Performance Features**

### **Request Monitoring**
- **Response Times**: Track API response times
- **Backend Performance**: Monitor Go backend response times
- **Error Rates**: Track error frequencies
- **Request Volume**: Monitor request patterns

### **Cache Integration**
- **Cache Statistics**: Monitor Go backend cache
- **Hit/Miss Ratios**: Track cache performance
- **Cache Health**: Monitor cache status

### **Logging Performance**
- **Structured Logging**: Efficient log formatting
- **Async Logging**: Non-blocking log operations
- **Log Levels**: Configurable verbosity

## 🛡️ **Security Features**

### **Input Validation**
- **Request Validation**: Validate request bodies and parameters
- **Type Checking**: Ensure correct data types
- **Sanitization**: Clean input data

### **Error Handling**
- **Safe Errors**: Don't expose internal details
- **Error Logging**: Detailed error context
- **Graceful Failures**: Handle errors gracefully

### **CORS Support**
- **Cross-Origin**: Enable cross-origin requests
- **Configurable**: Flexible CORS configuration
- **Security**: Proper CORS headers

## 🔗 **Integration with Go Backend**

### **Service Communication**
- **HTTP Client**: Axios for backend communication
- **Error Handling**: Proper error propagation
- **Timeouts**: Request timeout handling
- **Retries**: Automatic retry logic

### **Data Consistency**
- **Request Forwarding**: Pass requests to Go backend
- **Response Proxying**: Return backend responses
- **Status Mapping**: Proper HTTP status codes
- **Error Mapping**: Consistent error format

### **Monitoring Integration**
- **Unified Logging**: Consistent log format
- **Performance Tracking**: Track end-to-end performance
- **Health Monitoring**: Monitor both services
- **Cache Statistics**: Access Go backend cache data

## 🎯 **Development Workflow**

### **Local Development**
```bash
# Terminal 1: Start Go backend
cd ../go-backend
go run server.go

# Terminal 2: Start Node.js backend
cd node-backend
npm run dev

# Terminal 3: Run tests
npm test
```

### **Testing Workflow**
```bash
# Run all tests
npm test

# Run with coverage
npm run test:coverage

# Watch mode for development
npm run test:watch

# Integration tests
npm run test:integration
```

### **Development Features**
- **Hot Reload**: Automatic server restart on changes
- **Test Watch**: Auto-run tests on file changes
- **Coverage Reports**: Detailed test coverage
- **Debug Logging**: Enhanced logging in development

## 🚀 **Next Steps**

This Node.js backend provides:
- ✅ **Complete Test Coverage**: All components thoroughly tested
- ✅ **API Gateway**: Seamless backend integration
- ✅ **Production Ready**: Error handling, logging, monitoring
- ✅ **Well Documented**: Comprehensive documentation
- ✅ **Extensible**: Clean architecture for easy expansion

The Node.js backend serves as a robust API gateway that demonstrates best practices in Node.js development with proper testing, logging, and integration patterns.
