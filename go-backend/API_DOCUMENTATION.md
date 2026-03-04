# API Documentation

This document provides comprehensive documentation for all REST API endpoints in the Go backend service.

## Base URL

```
http://localhost:8080
```

## Overview

The API provides RESTful endpoints for managing users and tasks, along with health monitoring and statistics. All endpoints return JSON responses and follow HTTP status code conventions.

## Authentication

Currently, the API does not require authentication. All endpoints are publicly accessible.

## Content Type

All requests should use `Content-Type: application/json` for request bodies.

## Response Format

### Success Response
```json
{
  "success": true,
  "data": { ... },
  "message": "Operation completed successfully"
}
```

### Error Response
```json
{
  "success": false,
  "error": "Error message",
  "code": "ERROR_CODE",
  "details": "Additional error details"
}
```

## Endpoints

### Health Check Endpoints

#### GET /health
Returns a comprehensive health report including all system components.

**Response:**
```json
{
  "status": "healthy",
  "message": "All systems operational",
  "version": "1.0.0",
  "timestamp": "2026-03-04T12:00:00Z",
  "uptime": "2h30m15s",
  "checks": {
    "memory": {
      "name": "memory",
      "status": "healthy",
      "message": "Memory usage within limits",
      "latency": "1ms"
    },
    "datastore": {
      "name": "datastore",
      "status": "healthy",
      "message": "Data store accessible",
      "latency": "5ms"
    }
  },
  "readiness": {
    "ready": true,
    "reason": "All systems ready"
  },
  "liveness": {
    "alive": true,
    "reason": "Service is alive"
  }
}
```

#### GET /health/ready
Kubernetes readiness probe. Returns whether the service is ready to accept traffic.

**Response:**
```json
{
  "ready": true,
  "reason": "All systems ready"
}
```

#### GET /health/live
Kubernetes liveness probe. Returns whether the service is alive.

**Response:**
```json
{
  "alive": true,
  "reason": "Service is alive"
}
```

---

### User Management Endpoints

#### GET /api/users
Retrieves a list of all users.

**Response:**
```json
{
  "success": true,
  "data": {
    "users": [
      {
        "id": 1,
        "name": "John Doe",
        "email": "john@example.com",
        "role": "developer"
      },
      {
        "id": 2,
        "name": "Jane Smith",
        "email": "jane@example.com",
        "role": "designer"
      }
    ],
    "count": 2
  }
}
```

#### GET /api/users/{id}
Retrieves a specific user by ID.

**Parameters:**
- `id` (path): User ID (integer)

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "role": "developer"
  }
}
```

**Error Responses:**
- `404 Not Found`: User with specified ID does not exist
- `400 Bad Request`: Invalid user ID format

#### POST /api/users
Creates a new user.

**Request Body:**
```json
{
  "name": "Alice Johnson",
  "email": "alice@example.com",
  "role": "manager"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 3,
    "name": "Alice Johnson",
    "email": "alice@example.com",
    "role": "manager"
  },
  "message": "User created successfully"
}
```

**Validation Requirements:**
- `name`: Required, max 100 characters
- `email`: Required, valid email format, max 255 characters
- `role`: Required, max 50 characters

**Error Responses:**
- `400 Bad Request`: Invalid request body or validation errors
- `409 Conflict`: User with email already exists

#### PUT /api/users/{id}
Updates an existing user completely.

**Parameters:**
- `id` (path): User ID (integer)

**Request Body:**
```json
{
  "name": "Alice Johnson Updated",
  "email": "alice.updated@example.com",
  "role": "senior_manager"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 3,
    "name": "Alice Johnson Updated",
    "email": "alice.updated@example.com",
    "role": "senior_manager"
  },
  "message": "User updated successfully"
}
```

**Error Responses:**
- `404 Not Found`: User with specified ID does not exist
- `400 Bad Request`: Invalid request body or validation errors

#### DELETE /api/users/{id}
Deletes a user by ID.

**Parameters:**
- `id` (path): User ID (integer)

**Response:**
```json
{
  "success": true,
  "message": "User deleted successfully"
}
```

**Error Responses:**
- `404 Not Found`: User with specified ID does not exist
- `400 Bad Request`: Invalid user ID format

---

### Task Management Endpoints

#### GET /api/tasks
Retrieves a list of all tasks.

**Response:**
```json
{
  "success": true,
  "data": {
    "tasks": [
      {
        "id": 1,
        "title": "Implement authentication",
        "status": "pending",
        "userId": 1
      },
      {
        "id": 2,
        "title": "Design user interface",
        "status": "in-progress",
        "userId": 2
      }
    ],
    "count": 2
  }
}
```

#### GET /api/tasks/{id}
Retrieves a specific task by ID.

**Parameters:**
- `id` (path): Task ID (integer)

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "Implement authentication",
    "status": "pending",
    "userId": 1
  }
}
```

**Error Responses:**
- `404 Not Found`: Task with specified ID does not exist
- `400 Bad Request`: Invalid task ID format

#### POST /api/tasks
Creates a new task.

**Request Body:**
```json
{
  "title": "Write API documentation",
  "status": "pending",
  "userId": 1
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 3,
    "title": "Write API documentation",
    "status": "pending",
    "userId": 1
  },
  "message": "Task created successfully"
}
```

**Validation Requirements:**
- `title`: Required, max 200 characters
- `status`: Required, must be one of: "pending", "in-progress", "completed"
- `userId`: Required, must reference an existing user

**Error Responses:**
- `400 Bad Request`: Invalid request body or validation errors
- `404 Not Found`: Referenced user does not exist

#### PUT /api/tasks/{id}
Updates an existing task completely.

**Parameters:**
- `id` (path): Task ID (integer)

**Request Body:**
```json
{
  "title": "Write comprehensive API documentation",
  "status": "in-progress",
  "userId": 2
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 3,
    "title": "Write comprehensive API documentation",
    "status": "in-progress",
    "userId": 2
  },
  "message": "Task updated successfully"
}
```

**Error Responses:**
- `404 Not Found`: Task with specified ID does not exist
- `400 Bad Request`: Invalid request body or validation errors
- `404 Not Found`: Referenced user does not exist

#### DELETE /api/tasks/{id}
Deletes a task by ID.

**Parameters:**
- `id` (path): Task ID (integer)

**Response:**
```json
{
  "success": true,
  "message": "Task deleted successfully"
}
```

**Error Responses:**
- `404 Not Found`: Task with specified ID does not exist
- `400 Bad Request`: Invalid task ID format

---

### Statistics Endpoints

#### GET /api/stats
Retrieves system statistics including user and task counts.

**Response:**
```json
{
  "success": true,
  "data": {
    "users": {
      "total": 3
    },
    "tasks": {
      "total": 5,
      "pending": 2,
      "inProgress": 2,
      "completed": 1
    }
  }
}
```

#### GET /api/stats/cache
Retrieves cache performance statistics.

**Response:**
```json
{
  "success": true,
  "data": {
    "hits": 150,
    "misses": 25,
    "evictions": 0,
    "totalEntries": 10
  }
}
```

---

## Error Codes

| Error Code | Description | HTTP Status |
|------------|-------------|-------------|
| `VALIDATION_ERROR` | Request validation failed | 400 |
| `INVALID_JSON` | Invalid JSON format | 400 |
| `FIELD_TOO_LONG` | Field exceeds maximum length | 400 |
| `NOT_FOUND` | Resource not found | 404 |
| `DUPLICATE` | Resource already exists | 409 |
| `DATA_STORE_ERROR` | Data store operation failed | 500 |
| `INTERNAL_ERROR` | Internal server error | 500 |
| `INVALID_METHOD` | HTTP method not allowed | 405 |
| `REQUEST_TOO_LARGE` | Request body too large | 413 |

## Rate Limiting

Currently, no rate limiting is implemented. All endpoints have unlimited access.

## Pagination

Currently, pagination is not implemented. List endpoints return all available data.

## CORS

The API supports Cross-Origin Resource Sharing (CORS) with the following headers:

```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization
```

## Examples

### Create a User and Task

```bash
# Create a user
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bob Wilson",
    "email": "bob@example.com",
    "role": "developer"
  }'

# Create a task for the user
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Setup development environment",
    "status": "pending",
    "userId": 1
  }'
```

### Get All Users and Their Tasks

```bash
# Get all users
curl http://localhost:8080/api/users

# Get all tasks
curl http://localhost:8080/api/tasks

# Get statistics
curl http://localhost:8080/api/stats
```

### Update and Delete

```bash
# Update a user
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bob Wilson Updated",
    "email": "bob.updated@example.com",
    "role": "senior_developer"
  }'

# Delete a task
curl -X DELETE http://localhost:8080/api/tasks/1
```

## SDK and Client Libraries

Currently, no official SDKs are provided. The API can be accessed using any HTTP client library.

## Versioning

The current API version is 1.0.0. Version information is available in the health endpoint response.

## Support

For API support and questions, please refer to the project documentation or create an issue in the repository.
