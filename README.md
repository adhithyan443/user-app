# User Management Application

A full-stack user management web application built with Go, Gin, GORM, PostgreSQL, JWT Authentication, and Session Management.

This project demonstrates:

* User Authentication
* JWT-based Authorization
* Admin/User Role Management
* GORM ORM Integration
* PostgreSQL Database Integration
* Logging using slog
* Middleware Handling
* HTML Template Rendering
* Session Management
* Form Validation

---

# Tech Stack

## Backend

* Go (Golang)
* Gin Web Framework
* GORM
* PostgreSQL

## Authentication & Security

* JWT (JSON Web Token)
* Sessions
* Password Hashing using bcrypt

## Frontend

* HTML Templates
* CSS

## Logging

* slog (Structured Logging)

---

# Project Structure

```bash
user-app/
│
├── config/
│   └── database.go
│
├── handlers/
│   ├── admin_handler.go
│   ├── auth_handler.go
│   └── user_handler.go
│
├── middleware/
│   ├── admin_middleware.go
│   ├── auth_middleware.go
│   ├── cache_middleware.go
│   └── jwt.go
│
├── models/
│   └── user.go
│
├── routes/
│   └── routes.go
│
├── templates/
│   ├── pages/
│   └── static/
│
├── utils/
│   ├── jwt.go
│   └── validator.go
│
├── .env
├── go.mod
├── go.sum
└── main.go
```

---

# Features

## Authentication

* User Login
* User Registration
* JWT Token Generation
* JWT Verification
* Secure Password Hashing
* Logout Functionality

## Authorization

* Admin Middleware
* Protected Routes
* Role-based Access Control

## User Management

* Add Users
* Edit Users
* Delete Users
* Change Password
* User Listing

## Database

* PostgreSQL Integration
* Auto Migration using GORM
* Connection Pooling
* Structured Queries

## Logging & Monitoring

* Structured JSON Logging using slog
* Request Logging
* Error Logging

---

# JWT Authentication Flow

1. User logs in with email and password.
2. Password is verified using bcrypt.
3. JWT token is generated.
4. Token is stored and sent with requests.
5. Middleware validates the token.
6. Protected routes are accessed only if the token is valid.

---

# Environment Variables

Create a `.env` file in the root directory.

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=userdb
DB_SSLMODE=disable

JWT_SECRET=your-secret-key
```

---

# Installation & Setup

## 1. Clone the Repository

```bash
git clone https://github.com/your-username/user-app.git
cd user-app
```

## 2. Install Dependencies

```bash
go mod tidy
```

## 3. Setup PostgreSQL Database

Create a PostgreSQL database.

Example:

```sql
CREATE DATABASE userdb;
```

---

## 4. Configure Environment Variables

Create a `.env` file and add your database credentials.

---

## 5. Run the Application

```bash
go run main.go
```

Server will start at:

```bash
http://localhost:8080
```

---

# Main Dependencies

| Package                 | Purpose                     |
| ----------------------- | --------------------------- |
| gin-gonic/gin           | Web Framework               |
| gorm.io/gorm            | ORM                         |
| gorm.io/driver/postgres | PostgreSQL Driver           |
| golang-jwt/jwt/v5       | JWT Authentication          |
| gin-contrib/sessions    | Session Management          |
| golang.org/x/crypto     | Password Hashing            |
| joho/godotenv           | Environment Variable Loader |

---

# Middleware Used

## Auth Middleware

Protects authenticated routes.

## Admin Middleware

Allows access only for admin users.

## Cache Middleware

Prevents browser caching issues.

## JWT Middleware

Validates JWT tokens.

---

# Database Migration

Auto migration is handled using:

```go
config.AutoMigrate()
```

This automatically creates and updates database tables.

---

# Logging

The project uses structured logging with `slog`.

Example:

```go
slog.Info("Starting User Application...")
```

Benefits:

* Better debugging
* Structured logs
* Production-friendly logging
* Easier monitoring

---

# Security Features

* Password hashing using bcrypt
* JWT token validation
* HTTP-only sessions
* Route protection using middleware
* Role-based authorization

---

# Future Improvements

* Refresh Tokens
* Email Verification
* Password Reset
* Docker Support
* Unit Testing
* API Documentation using Swagger
* Rate Limiting
* CSRF Protection
* Redis Session Store
* Pagination & Search

---

# Learning Concepts Covered

This project is useful for learning:

* Golang Web Development
* REST API Design
* Gin Framework
* GORM ORM
* Authentication & Authorization
* JWT
* Middleware
* PostgreSQL
* Logging
* MVC Pattern
* Session Management

---

