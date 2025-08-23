# API Documentation

## Overview
This document provides detailed information about all API endpoints available in the practice-go application. The API follows REST conventions and uses JSON for data exchange.

## Base URL
```
http://localhost:3000/api
```

## Authentication
Most endpoints require JWT authentication. Include the token in the Authorization header:
```
Authorization: Bearer <JWT_TOKEN>
```

## Standard Response Format
All responses follow this structure:
```json
{
  "status": "success|error",
  "message": "Description message",
  "data": "Response data (varies by endpoint)"
}
```

## Error Codes
- `200` - Success
- `201` - Created successfully  
- `400` - Bad Request (validation errors)
- `401` - Unauthorized (missing or invalid token)
- `404` - Not Found
- `500` - Internal Server Error

---

## Authentication Endpoints

### 1. User Registration
**POST** `/api/auth/register`

Register a new user account.

**Request Body:**
```json
{
  "username": "string (required, 3-50 characters)",
  "email": "string (required, valid email format)",
  "password": "string (required, min 6 characters)",
  "names": "string (optional, display name)"
}
```

**Response (201 - Success):**
```json
{
  "status": "success",
  "message": "User created successfully",
  "data": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "names": "John Doe"
  }
}
```

**Response (400 - Validation Error):**
```json
{
  "status": "error",
  "message": "Username already exists"
}
```

### 2. User Login
**POST** `/api/auth/login`

Authenticate user and receive JWT token.

**Request Body:**
```json
{
  "username": "string (required)",
  "password": "string (required)"
}
```

**Response (200 - Success):**
```json
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "user": {
      "id": 1,
      "username": "john_doe",
      "email": "john@example.com",
      "names": "John Doe"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Response (401 - Invalid Credentials):**
```json
{
  "status": "error",
  "message": "Invalid username or password"
}
```

---

## User Management Endpoints

### 3. Get All Users
**GET** `/api/user`

Retrieve a list of all users.

**Authentication:** Not required

**Response (200 - Success):**
```json
{
  "status": "success",
  "message": "All users",
  "data": [
    {
      "id": 1,
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z",
      "deleted_at": null,
      "username": "john_doe",
      "email": "john@example.com",
      "names": "John Doe"
    }
  ]
}
```

### 4. Get User by ID
**GET** `/api/user/:id`

Retrieve a specific user by their ID.

**Authentication:** Not required

**Path Parameters:**
- `id` (integer, required) - User ID

**Response (200 - Success):**
```json
{
  "status": "success",
  "message": "User found",
  "data": {
    "id": 1,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z",
    "deleted_at": null,
    "username": "john_doe",
    "email": "john@example.com",
    "names": "John Doe"
  }
}
```

**Response (404 - Not Found):**
```json
{
  "status": "error",
  "message": "User not found"
}
```

### 5. Create User
**POST** `/api/user`

Create a new user (alternative to registration).

**Authentication:** Not required

**Request Body:**
```json
{
  "username": "string (required, 3-50 characters)",
  "email": "string (required, valid email format)",
  "password": "string (required, min 6 characters)",
  "names": "string (optional, display name)"
}
```

**Response (201 - Success):**
```json
{
  "status": "success",
  "message": "User created successfully",
  "data": {
    "id": 2,
    "username": "jane_doe",
    "email": "jane@example.com",
    "names": "Jane Doe"
  }
}
```

### 6. Update User
**PATCH** `/api/user/:id`

Update an existing user.

**Authentication:** Required (JWT Token)

**Path Parameters:**
- `id` (integer, required) - User ID

**Request Body:**
```json
{
  "username": "string (optional, 3-50 characters)",
  "email": "string (optional, valid email format)",
  "password": "string (optional, min 6 characters)",
  "names": "string (optional, display name)"
}
```

**Response (200 - Success):**
```json
{
  "status": "success",
  "message": "User updated successfully",
  "data": {
    "id": 1,
    "username": "john_updated",
    "email": "john.updated@example.com",
    "names": "John Updated"
  }
}
```

### 7. Delete User
**DELETE** `/api/user/:id`

Delete a user (soft delete).

**Authentication:** Required (JWT Token)

**Path Parameters:**
- `id` (integer, required) - User ID

**Response (200 - Success):**
```json
{
  "status": "success",
  "message": "User deleted successfully"
}
```

---

## Resource Management Endpoints

### 8. Get All Resources
**GET** `/api/resource`

Retrieve a list of all resources.

**Authentication:** Not required

**Response (200 - Success):**
```json
{
  "status": "success",
  "message": "All resources",
  "data": [
    {
      "id": 1,
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z",
      "deleted_at": null,
      "name": "Сталь",
      "description": "Конструкционная сталь высокого качества",
      "unit": "кг",
      "quantity": 1000
    }
  ]
}
```

### 9. Get Resource by ID
**GET** `/api/resource/:id`

Retrieve a specific resource by its ID.

**Authentication:** Not required

**Path Parameters:**
- `id` (integer, required) - Resource ID

**Response (200 - Success):**
```json
{
  "status": "success",
  "message": "Resource found",
  "data": {
    "id": 1,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z",
    "deleted_at": null,
    "name": "Сталь",
    "description": "Конструкционная сталь высокого качества",
    "unit": "кг",
    "quantity": 1000
  }
}
```

**Response (404 - Not Found):**
```json
{
  "status": "error",
  "message": "Resource not found"
}
```

### 10. Create Resource
**POST** `/api/resource`

Create a new resource.

**Authentication:** Required (JWT Token)

**Request Body:**
```json
{
  "name": "string (required, 2-100 characters, must be unique)",
  "description": "string (optional)",
  "unit": "string (required, 1-20 characters, e.g., кг, л, шт)",
  "quantity": "integer (required, >= 0)"
}
```

**Response (201 - Success):**
```json
{
  "status": "success",
  "message": "Resource created successfully",
  "data": {
    "id": 2,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z",
    "deleted_at": null,
    "name": "Алюминий",
    "description": "Алюминиевые листы для производства",
    "unit": "кг",
    "quantity": 500
  }
}
```

**Response (400 - Validation Error):**
```json
{
  "status": "error",
  "message": "Resource name already exists"
}
```

### 11. Update Resource
**PUT** `/api/resource/:id`

Update an existing resource.

**Authentication:** Required (JWT Token)

**Path Parameters:**
- `id` (integer, required) - Resource ID

**Request Body:**
```json
{
  "name": "string (optional, 2-100 characters)",
  "description": "string (optional)",
  "unit": "string (optional, 1-20 characters)",
  "quantity": "integer (optional, >= 0)"
}
```

**Response (200 - Success):**
```json
{
  "status": "success",
  "message": "Resource updated successfully",
  "data": {
    "id": 1,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T10:30:00Z",
    "deleted_at": null,
    "name": "Сталь обновленная",
    "description": "Обновленное описание стали",
    "unit": "кг",
    "quantity": 1200
  }
}
```

### 12. Delete Resource
**DELETE** `/api/resource/:id`

Delete a resource (soft delete).

**Authentication:** Required (JWT Token)

**Path Parameters:**
- `id` (integer, required) - Resource ID

**Response (200 - Success):**
```json
{
  "status": "success",
  "message": "Resource deleted successfully"
}
```

### 13. Get Resource History
**GET** `/api/resource/:id/history`

Retrieve the change history for a specific resource.

**Authentication:** Not required

**Path Parameters:**
- `id` (integer, required) - Resource ID

**Response (200 - Success):**
```json
{
  "status": "success",
  "message": "Resource history",
  "data": [
    {
      "id": 1,
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z",
      "deleted_at": null,
      "resource_id": 1,
      "action": "CREATE",
      "user_id": 1,
      "old_data": null,
      "new_data": "{\"id\":1,\"name\":\"Сталь\",\"description\":\"Конструкционная сталь\",\"unit\":\"кг\",\"quantity\":1000}",
      "timestamp": "2023-01-01T00:00:00Z",
      "description": "Ресурс 'Сталь' создан",
      "resource": {
        "id": 1,
        "name": "Сталь",
        "description": "Конструкционная сталь высокого качества",
        "unit": "кг",
        "quantity": 1000
      },
      "user": {
        "id": 1,
        "username": "admin",
        "email": "admin@ipr-op.org",
        "names": "Администратор"
      }
    }
  ]
}
```

---

## General Information Endpoint

### 14. API Health Check
**GET** `/api/`

Basic health check endpoint.

**Authentication:** Not required

**Response (200 - Success):**
```json
{
  "message": "Hello, World!"
}
```

---

## HTTP Status Codes Details

### Success Codes
- **200 OK** - Request successful (GET, PUT, DELETE operations)
- **201 Created** - Resource created successfully (POST operations)

### Client Error Codes
- **400 Bad Request** - Invalid request body, validation errors, or malformed data
- **401 Unauthorized** - Missing, invalid, or expired JWT token
- **404 Not Found** - Requested resource does not exist
- **409 Conflict** - Resource already exists (e.g., duplicate username/email)

### Server Error Codes
- **500 Internal Server Error** - Database errors, server configuration issues

---

## Data Validation Rules

### User Fields
- **username**: 3-50 characters, alphanumeric and underscores only, must be unique
- **email**: Valid email format, must be unique
- **password**: Minimum 6 characters, will be hashed before storage
- **names**: Optional display name, up to 255 characters

### Resource Fields
- **name**: 2-100 characters, must be unique across all resources
- **description**: Optional text description, up to 500 characters
- **unit**: 1-20 characters (examples: кг, л, шт, м², м³, т)
- **quantity**: Non-negative integer

### History Fields
- **action**: Automatically set to CREATE, UPDATE, or DELETE
- **timestamp**: Automatically set to current time
- **user_id**: Extracted from JWT token
- **old_data/new_data**: JSON representation of resource before/after change

---

## Rate Limiting
Currently no rate limiting is implemented. In production, consider implementing rate limiting to prevent abuse.

## CORS
CORS is not explicitly configured. Ensure proper CORS settings for frontend integration.

## Security Notes
1. JWT tokens expire after 72 hours (3 days)
2. Passwords are hashed using bcrypt
3. All protected routes require valid JWT authentication
4. Soft deletion is used for users and resources (records are marked as deleted but not physically removed)
5. Resource history provides complete audit trail of all changes

---

## Example Usage

### Complete User Registration and Resource Management Flow

1. **Register a new user:**
```bash
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "email": "test@example.com", "password": "password123", "names": "Test User"}'
```

2. **Login to get JWT token:**
```bash
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "password123"}'
```

3. **Create a resource (using token from login):**
```bash
curl -X POST http://localhost:3000/api/resource \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"name": "Новый ресурс", "description": "Описание ресурса", "unit": "шт", "quantity": 100}'
```

4. **Get resource history:**
```bash
curl -X GET http://localhost:3000/api/resource/1/history
```

This API provides comprehensive resource management with full audit capabilities, suitable for inventory tracking, asset management, and similar applications requiring detailed change history.