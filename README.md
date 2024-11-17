# Tender Management System Documentation

This project implements a RESTful API for managing tenders, bids, and user authentication. It leverages PostgreSQL for data persistence, Redis for caching, and Casbin for authorization. Real-time notifications are implemented using WebSockets.

## Table of Contents

- [Architecture](#architecture)
- [API Endpoints](#api-endpoints)
- [Authentication](#authentication)
- [Tenders](#tenders)
- [Bids](#bids)
- [Notifications (WebSockets)](#notifications-websockets)
- [Data Model](#data-model)
- [Authorization (Casbin)](#authorization-casbin)
- [Rate Limiting](#rate-limiting)
- [Caching (Redis)](#caching-redis)
- [Database Migrations](#database-migrations)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [Dependencies](#dependencies)

## Architecture

The application follows a layered architecture:

- **API (internal/api)**: Exposes HTTP endpoints using Gin. Handles request routing, validation, and responses. Swagger documentation is integrated.
- **Service (internal/service)**: Contains the business logic for managing tenders, bids, and authentication.
- **Repository (internal/repository)**: Provides an abstraction layer for data access. Implements interfaces for interacting with the database. PostgreSQL is the chosen database.
- **Models (internal/models)**: Defines the data structures for tenders, bids, users, and notifications.
- **Utils (internal/utils)**: Contains utility functions like JWT generation/validation, rate limiting, and notification services.

## API Endpoints

All API endpoints are prefixed with `/api` and require authentication using JWT (Bearer token) unless otherwise specified.

### Authentication

- `POST /register`: Registers a new user (client or contractor).
- `POST /login`: Logs in an existing user and returns a JWT.

### Tenders

- `POST /api/client/tenders`: Creates a new tender (client role only).
- `GET /api/client/tenders`: Lists all tenders created by the authenticated client.
- `GET /api/client/tenders/filter`: Lists all tenders filtered by status or with search.
- `PUT /api/client/tenders/:id`: Updates the status of a tender (client role only).
- `DELETE /api/client/tenders/:id`: Deletes a tender (client role only).
- `GET /api/client/tenders/:tender_id/bids`: Gets all bids for a specific tender created by the client.
- `POST /api/client/tenders/:tender_id/award/:bid_id`: Award a bid for a specific tender (client role only).

### Bids

- `POST /api/contractor/tenders/:tender_id/bid`: Creates a new bid for a tender (contractor role only).
- `GET /api/contractor/bids`: Lists bids submitted by the authenticated contractor.
- `DELETE /api/contractor/bids/:bid_id`: Deletes a specific bid by the authenticated contractor.

## Notifications (WebSockets)

- `GET /api/ws`: Establishes a WebSocket connection for real-time notifications (Authentication required).
    - Clients receive notifications about new bids on their tenders.
    - Contractors receive notifications about their bid awards.

## Data Model

- **Users**: `id`, `username`, `email`, `password_hash`, `role`, `created_at`, `updated_at`
- **Tenders**: `id`, `client_id`, `title`, `description`, `deadline`, `budget`, `status`, `attachment`, `created_at`, `updated_at`
- **Bids**: `id`, `tender_id`, `contractor_id`, `price`, `delivery_time`, `comments`, `status`, `created_at`, `updated_at`
- **Notifications**: `id`, `user_id`, `message`, `relation_id`, `type`, `read`, `created_at`

## Authorization (Casbin)

Casbin is used for authorization. The model and policy are defined in `config/model.conf` and `config/policy.csv`, respectively. The API middleware (AuthorizationMiddleware) enforces the policies based on the user's role and the requested resource.

## Rate Limiting

A simple rate limiter is implemented using `golang.org/x/time/rate`. It limits requests per user to prevent abuse. This is applied in the API middleware.

## Caching (Redis)

Redis is used for caching tender lists to improve performance. Cached data is stored with an expiration time.

## Database Migrations

Database migrations are managed using `golang-migrate`. The migration files are located in the `migrations` directory.

## Configuration

Database connection details and other configurations are hardcoded in the `cmd/server/main.go` file for simplicity. In a production environment, this should be externalized (e.g., environment variables, configuration files).

## Running the Application

1. Ensure you have PostgreSQL and Redis running.
2. Run the migrations: `migrate -source file://migrations -database "postgres://postgres:postgres@localhost:5433/tender_db?sslmode=disable" up`
3. Build and run the application: `go run cmd/server/main.go`

## Dependencies

- Gin
- GORM
- PostgreSQL driver
- Redis client
- Casbin
- JWT library
- Swagger
- Golang-migrate
