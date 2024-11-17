# Tender Management System API Documentation

This document provides detailed information about the RESTful API for a tender management system. The API utilizes PostgreSQL for data persistence, Redis for caching, Casbin for authorization, and WebSockets for real-time notifications.

## Table of Contents
- [Introduction](#introduction)
- [Run Project](#Run)
- [Authentication](#authentication)
- [Tenders](#tenders)
    - [Creating a Tender](#creating-a-tender)
    - [Listing Tenders](#listing-tenders)
    - [Getting a Tender](#getting-a-tender)
    - [Updating a Tender](#updating-a-tender)
    - [Deleting a Tender](#deleting-a-tender)
    - [Filtering Tenders](#filtering-tenders)
    - [Getting Bids for a Tender (Client)](#getting-bids-for-a-tender-client)
    - [Awarding a Bid](#awarding-a-bid)
- [Bids](#bids)
    - [Creating a Bid](#creating-a-bid)
    - [Listing Bids (Contractor)](#listing-bids-contractor)
    - [Deleting a Bid (Contractor)](#deleting-a-bid-contractor)
- [History](#history)
    - [Tender History](#tender-history)
    - [Bid History](#bid-history)
- [Real-time Notifications (WebSockets)](#real-time-notifications-websockets)
- [Data Model](#data-model)
- [Error Handling](#error-handling)
- [Rate Limiting](#rate-limiting)
- [Authorization](#authorization)
- [Deployment (Docker)](#deployment-docker)

## Introduction
The API provides endpoints for managing tenders and bids within a tendering process. It supports client and contractor roles with distinct permissions. Real-time updates are facilitated through WebSockets.

All API endpoints are prefixed with `/api` and require JSON formatted requests and responses, except for the WebSocket endpoint. Authentication is required for all endpoints except `/register` and `/login`, using JSON Web Tokens (JWTs).

## Run

1. Run db
```
    make run_db
```
2. Migrate db (required migrate cli)
```
    make migrate-up
```

3. Run application
```
go run cmd/server/main.go
```


## Authentication
### `/register` (POST)
Registers a new user.

Request Body:
```json
{
    "username": "string",
    "email": "string",
    "password": "string",
    "role": "string"
}
```

Response (201 Created):
A successful response returns an `AuthResponse` object with a JWT and user details.

Error Responses:
- 400 Bad Request: Invalid input data.
- 409 Conflict: Username or email already exists.
- 500 Internal Server Error: Server error during registration.

### `/login` (POST)
Logs in an existing user.

Request Body:
```json
{
    "username": "string",
    "password": "string"
}
```

Response (200 OK):
A successful response returns an `AuthResponse` object with a JWT and user details.

Error Responses:
- 400 Bad Request: Invalid input data.
- 401 Unauthorized: Invalid credentials.
- 404 Not Found: User not found.
- 500 Internal Server Error: Server error during login.

## Tenders
### Creating a Tender
`/api/client/tenders` (POST): Creates a new tender. Requires client authentication.

Request Body:
```json
{
    "title": "string",
    "description": "string",
    "deadline": "string",
    "budget": "number",
    "attachment": "string"
}
```

Response (201 Created):
The ID and title of the newly created tender.

Error Responses:
- 400 Bad Request: Invalid input data. Deadline in the past is not allowed.
- 500 Internal Server Error: Server error during tender creation.

### Listing Tenders
`/api/client/tenders` (GET): Lists all tenders for the authenticated client. Requires client authentication.

Response (200 OK):
An array of Tender objects.

Error Responses:
- 500 Internal Server Error: Server error during tender retrieval.

### Getting a Tender
`/api/client/tenders/:id` (GET): Gets a tender by ID. Requires client authentication.

Path Parameters:
- id (string): The ID of the tender.

Response (200 OK):
A Tender object.

Error Responses:
- 400 Bad Request: Invalid tender ID.
- 404 Not Found: Tender not found.
- 500 Internal Server Error: Server error during tender retrieval.

### Updating a Tender
`/api/client/tenders/:id` (PUT): Updates the status of a tender. Requires client authentication.

Path Parameters:
- id (string): The ID of the tender.

Request Body:
```json
{
    "status": "string"
}
```

Response (200 OK):
A success message.

Error Responses:
- 400 Bad Request: Invalid input data.
- 404 Not Found: Tender not found.
- 500 Internal Server Error: Server error during tender update.

### Deleting a Tender
`/api/client/tenders/:id` (DELETE): Deletes a tender. Requires client authentication.

Path Parameters:
- id (string): The ID of the tender.

Response (200 OK):
A success message.

Error Responses:
- 400 Bad Request: Invalid tender ID.
- 404 Not Found: Tender not found.
- 403 Forbidden: Unauthorized access.
- 500 Internal Server Error: Server error during tender deletion.

### Filtering Tenders
`/api/client/tenders/filter` (GET): Lists tenders with filtering options. Requires client authentication.

Query Parameters:
- status (string, optional): Filter by tender status ("open", "closed", "awarded").
- search (string, optional): Search tenders by keyword (case-insensitive).

Response (200 OK):
An array of Tender objects.

Error Responses:
- 500 Internal Server Error: Server error during tender retrieval.

### Getting Bids for a Tender (Client)
`/api/client/tenders/:tender_id/bids` (GET): Gets all bids for a specific tender. Requires client authentication.

Path Parameters:
- tender_id (string): The ID of the tender.

Query Parameters (Optional filtering and sorting):
- status (string): Filter by bid status.
- search (string): Search bids by comments.
- min_price (number): Minimum bid price.
- max_price (number): Maximum bid price.
- min_delivery_time (integer): Minimum delivery time.
- max_delivery_time (integer): Maximum delivery time.
- sort_by (string): Sort by field (e.g., price, delivery_time).
- sort_order (string): Sort order (asc or desc).

Response (200 OK):
An array of Bid objects.

Error Responses:
- 400 Bad Request: Invalid tender ID.
- 500 Internal Server Error: Server error.

### Awarding a Bid
`/api/client/tenders/:tender_id/award/:bid_id` (POST): Awards a bid to a contractor. Requires client authentication.

Path Parameters:
- tender_id (string): The ID of the tender.
- bid_id (string): The ID of the bid to award.

Response (200 OK):
A success message.

Error Responses:
- 400 Bad Request: Invalid tender ID or bid ID.
- 403 Forbidden: Unauthorized access or bid does not belong to the tender.
- 404 Not Found: Tender or bid not found.
- 500 Internal Server Error: Server error.

## Bids
### Creating a Bid
`/api/contractor/tenders/:tender_id/bid` (POST): Creates a new bid for a tender. Requires contractor authentication. Subject to rate limiting.

Path Parameters:
- tender_id (string): The ID of the tender.

Request Body:
```json
{
    "price": "number",
    "delivery_time": "integer",
    "comments": "string"
}
```

Response (201 Created):
The created Bid object.

Error Responses:
- 400 Bad Request: Invalid input data, or tender not open for bids.
- 429 Too Many Requests: Rate limit exceeded.
- 500 Internal Server Error: Server error.

### Listing Bids (Contractor)
`/api/contractor/bids` (GET): Lists all bids submitted by the authenticated contractor. Requires contractor authentication.

Response (200 OK):
An array of Bid objects.

Error Responses:
- 500 Internal Server Error: Server error.

### Deleting a Bid (Contractor)
`/api/contractor/bids/:bid_id` (DELETE): Deletes a bid. Requires contractor authentication.

Path Parameters:
- bid_id (string): The ID of the bid to delete.

Response (200 OK):
A success message.

Error Responses:
- 400 Bad Request: Invalid bid ID.
- 403 Forbidden: Unauthorized access (contractor does not own the bid).
- 404 Not Found: Bid not found.
- 500 Internal Server Error: Server error.

## History
### Tender History
`/api/users/:id/tenders` (GET): Retrieves tender history for a given user ID. Requires authentication.

Path Parameters:
- id (string): The ID of the user.

Response (200 OK):
An array of Tender objects.

Error Responses:
- 400 Bad Request: Invalid user ID.
- 500 Internal Server Error: Server error.

### Bid History
`/api/users/:id/bids` (GET): Retrieves bid history for a given user ID. Requires authentication.

Path Parameters:
- id (string): The ID of the user.

Response (200 OK):
An array of Bid objects.

Error Responses:
- 400 Bad Request: Invalid user ID.
- 500 Internal Server Error: Server error.

## Real-time Notifications (WebSockets)
`/api/ws` (GET): Establishes a WebSocket connection for real-time notifications. Requires authentication. The client should send a keep-alive message periodically to keep the connection open.

The WebSocket connection will send JSON notifications:
- `new_bid`: Notifies clients of new bids on their tenders.
- `bid_awarded`: Notifies contractors that their bid has been awarded.

Error Responses:
- 400 Bad Request: Invalid input data.
- 401 Unauthorized: Authentication required.

## Data Model
- User: id (UUID), username (string), email (string), password_hash (string), role (string, "client" or "contractor"), created_at (timestamp), updated_at (timestamp)
- Tender: id (UUID), client_id (UUID), title (string), description (string), deadline (timestamp), budget (decimal), status (string, "open", "closed", "awarded"), attachment (string, optional), created_at (timestamp), updated_at (timestamp)
- Bid: id (UUID), tender_id (UUID), contractor_id (UUID), price (decimal), delivery_time (integer), comments (string), status (string, "open", "awarded"), created_at (timestamp), updated_at (timestamp)
- Notification: id (UUID), user_id (UUID), message (string), relation_id (UUID, optional), type (string), read (boolean), created_at (timestamp)

## Error Handling
Error responses will generally be in JSON format with a "message" field indicating the error. HTTP status codes will indicate the type of error (e.g., 400 Bad Request, 404 Not Found, 500 Internal Server Error).

## Rate Limiting
Rate limiting is implemented to prevent abuse of the bid creation endpoint. Contractor users are limited to a specific number of bid submissions within a time window. The rate limit is configurable. If exceeded, a 429 Too Many Requests response is returned.

## Authorization
The API uses Casbin for role-based access control (RBAC). Client users can only create and manage their tenders. Contractor users can only submit bids. Specific permissions are defined in a Casbin policy file.

## Deployment (Docker)
The application can be deployed using Docker Compose. The `docker-compose.yml` file defines the services for the application, PostgreSQL database, and Redis cache. Remember to replace "your-secret-key" in `docker-compose.yml` with your actual JWT secret. To run:

```bash
docker-compose up -d --build
```

This will build the application and start all the services in detached mode. You can then access the API through `localhost:8080`.

This documentation provides a comprehensive overview of the API. For detailed information on individual endpoints, refer to the Swagger documentation, accessible via `/swagger` after the application starts.

