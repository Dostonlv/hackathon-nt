# Tender Management System API Documentation

## Overview
This document outlines the RESTful API endpoints for a tender management system. The API supports tender creation, bidding, and management with distinct roles for clients and contractors.

## Base URL
All API endpoints are prefixed with `/api` unless specified otherwise.

## Authentication
Authentication is required for all endpoints except registration and login. The API uses JSON Web Tokens (JWT) for authentication.

### Public Endpoints

#### Register User
```
POST /register
```

Creates a new user account.

**Request Body:**
```json
{
    "username": "string",
    "email": "string",
    "password": "string",
    "role": "string"  // "client" or "contractor"
}
```

**Responses:**
- `201 Created`: Registration successful
- `400 Bad Request`: Invalid input data
- `409 Conflict`: Username/email already exists
- `500 Internal Server Error`: Server error

#### Login
```
POST /login
```

Authenticates a user and returns a JWT token.

**Request Body:**
```json
{
    "username": "string",
    "password": "string"
}
```

**Responses:**
- `200 OK`: Login successful, returns JWT token
- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Invalid credentials
- `404 Not Found`: User not found
- `500 Internal Server Error`: Server error

### Protected Endpoints

## Client Endpoints

### Tender Management

#### Create Tender
```
POST /api/client/tenders
```

Creates a new tender. Requires client role.

**Request Body:**
```json
{
    "title": "string",
    "description": "string",
    "deadline": "string",
    "budget": "number",
    "attachment": "string"
}
```

**Responses:**
- `201 Created`: Tender created successfully
- `400 Bad Request`: Invalid input data
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized as client
- `500 Internal Server Error`: Server error

#### List Tenders
```
GET /api/client/tenders
```

Returns all tenders for the authenticated client.

**Responses:**
- `200 OK`: List of tenders
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized as client
- `500 Internal Server Error`: Server error

#### Update Tender Status
```
PUT /api/client/tenders/:id
```

Updates the status of a specific tender.

**Path Parameters:**
- `id`: Tender ID

**Request Body:**
```json
{
    "status": "string"  // "open", "closed", "awarded"
}
```

**Responses:**
- `200 OK`: Status updated successfully
- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized to modify tender
- `404 Not Found`: Tender not found
- `500 Internal Server Error`: Server error

#### Delete Tender
```
DELETE /api/client/tenders/:id
```

Deletes a specific tender.

**Path Parameters:**
- `id`: Tender ID

**Responses:**
- `200 OK`: Tender deleted successfully
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized to delete tender
- `404 Not Found`: Tender not found
- `500 Internal Server Error`: Server error

#### Filter Tenders
```
GET /api/client/tenders/filter
```

Returns filtered list of tenders.

**Query Parameters:**
- `status`: Filter by tender status
- `search`: Search keyword in tender details

**Responses:**
- `200 OK`: Filtered list of tenders
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized as client
- `500 Internal Server Error`: Server error

#### Get Bids for Tender
```
GET /api/client/tenders/:tender_id/bids
```

Returns all bids for a specific tender.

**Path Parameters:**
- `tender_id`: Tender ID

**Query Parameters:**
- `status`: Filter by bid status
- `search`: Search in bid comments
- `min_price`: Minimum bid price
- `max_price`: Maximum bid price
- `min_delivery_time`: Minimum delivery time
- `max_delivery_time`: Maximum delivery time
- `sort_by`: Field to sort by
- `sort_order`: "asc" or "desc"

**Responses:**
- `200 OK`: List of bids
- `400 Bad Request`: Invalid tender ID
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized to view bids
- `500 Internal Server Error`: Server error

#### Award Bid
```
POST /api/client/tenders/:tender_id/award/:bid_id
```

Awards a bid to a contractor.

**Path Parameters:**
- `tender_id`: Tender ID
- `bid_id`: Bid ID

**Responses:**
- `200 OK`: Bid awarded successfully
- `400 Bad Request`: Invalid IDs
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized to award bid
- `404 Not Found`: Tender or bid not found
- `500 Internal Server Error`: Server error

## Contractor Endpoints

### Bid Management

#### Create Bid
```
POST /api/contractor/tenders/:tender_id/bid
```

Creates a new bid for a tender. Subject to rate limiting.

**Path Parameters:**
- `tender_id`: Tender ID

**Request Body:**
```json
{
    "price": "number",
    "delivery_time": "integer",
    "comments": "string"
}
```

**Responses:**
- `201 Created`: Bid created successfully
- `400 Bad Request`: Invalid input or tender not open
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized as contractor
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Server error

#### List Contractor's Bids
```
GET /api/contractor/bids
```

Returns all bids submitted by the authenticated contractor.

**Responses:**
- `200 OK`: List of bids
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized as contractor
- `500 Internal Server Error`: Server error

#### Delete Bid
```
DELETE /api/contractor/bids/:bid_id
```

Deletes a specific bid.

**Path Parameters:**
- `bid_id`: Bid ID

**Responses:**
- `200 OK`: Bid deleted successfully
- `400 Bad Request`: Invalid bid ID
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized to delete bid
- `404 Not Found`: Bid not found
- `500 Internal Server Error`: Server error

## History Endpoints

#### Get Tender History
```
GET /api/users/:id/tenders
```

Returns tender history for a specific user.

**Path Parameters:**
- `id`: User ID

**Responses:**
- `200 OK`: List of tenders
- `400 Bad Request`: Invalid user ID
- `401 Unauthorized`: Not authenticated
- `500 Internal Server Error`: Server error

#### Get Bid History
```
GET /api/users/:id/bids
```

Returns bid history for a specific user.

**Path Parameters:**
- `id`: User ID

**Responses:**
- `200 OK`: List of bids
- `400 Bad Request`: Invalid user ID
- `401 Unauthorized`: Not authenticated
- `500 Internal Server Error`: Server error

## WebSocket Endpoint

#### Real-time Notifications
```
GET /api/ws
```

Establishes WebSocket connection for real-time notifications.

**Events:**
- `new_bid`: Notification when new bid is placed
- `bid_awarded`: Notification when bid is awarded

**Responses:**
- `101 Switching Protocols`: Connection established
- `401 Unauthorized`: Not authenticated

## Error Handling
All error responses follow this format:
```json
{
    "message": "Error description"
}
```

## Rate Limiting
Rate limiting is applied to the bid creation endpoint to prevent abuse. Exceeding the rate limit results in a 429 response.

## Authorization
The API implements role-based access control (RBAC) using Casbin. Each endpoint requires specific roles and permissions as detailed above.
