# Go Backend Architecture Layers

This document describes the typical layers in a Go backend application and their responsibilities.

## Layer Overview

```
┌─────────────────────────────────────┐
│           Routes/Handlers           │  ← HTTP layer
├─────────────────────────────────────┤
│             Services                │  ← Business logic
├─────────────────────────────────────┤
│           Repositories              │  ← Data access
├─────────────────────────────────────┤
│             Models                  │  ← Domain entities
└─────────────────────────────────────┘
```

---

## 1. Routes/Handlers Layer

**Location:** `internal/routes/` or `internal/handlers/`

**Responsibilities:**
- Define HTTP endpoints and methods (GET, POST, PUT, DELETE)
- Parse and validate request data (path params, query params, JSON body)
- Call appropriate service methods
- Format and return HTTP responses
- Handle HTTP-specific concerns (status codes, headers, cookies)

**Should NOT:**
- Contain business logic
- Directly access the database
- Know about database schemas or SQL

**Example:**
```go
func (h *UserHandler) GetUser(c *gin.Context) {
    id := c.Param("id")
    user, err := h.userService.GetByID(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }
    c.JSON(http.StatusOK, user)
}
```

---

## 2. Services Layer

**Location:** `internal/services/`

**Responsibilities:**
- Implement business logic and rules
- Orchestrate calls to one or more repositories
- Handle transactions when multiple operations must succeed together
- Validate business constraints (not just input format)
- Transform data between layers when needed

**Should NOT:**
- Know about HTTP (no request/response objects)
- Write raw SQL or database queries
- Handle HTTP status codes

**Example:**
```go
func (s *UserService) CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
    // Business rule: check if email already exists
    existing, _ := s.userRepo.GetByEmail(ctx, input.Email)
    if existing != nil {
        return nil, ErrEmailAlreadyExists
    }

    // Business logic: hash password
    hashedPassword := hashPassword(input.Password)

    return s.userRepo.Create(ctx, &User{
        Email:    input.Email,
        Password: hashedPassword,
    })
}
```

---

## 3. Repositories Layer

**Location:** `internal/repositories/` or `internal/db/`

**Responsibilities:**
- Abstract database operations (CRUD)
- Write and execute SQL queries or ORM calls
- Map database rows to domain models
- Handle database-specific errors

**Should NOT:**
- Contain business logic
- Know about HTTP
- Call other repositories (leave orchestration to services)

**Example:**
```go
func (r *UserRepository) GetByID(ctx context.Context, id string) (*User, error) {
    var user User
    err := r.db.QueryRowContext(ctx,
        "SELECT id, email, created_at FROM users WHERE id = $1", id,
    ).Scan(&user.ID, &user.Email, &user.CreatedAt)
    if err == sql.ErrNoRows {
        return nil, ErrNotFound
    }
    return &user, err
}
```

---

## 4. Models/Domain Layer

**Location:** `internal/models/` or `internal/domain/`

**Responsibilities:**
- Define domain entities and their fields
- Define DTOs (Data Transfer Objects) for API requests/responses
- May contain simple validation or domain methods

**Example:**
```go
// Domain model
type User struct {
    ID        string
    Email     string
    Password  string
    CreatedAt time.Time
}

// Request DTO
type CreateUserRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

// Response DTO
type UserResponse struct {
    ID    string `json:"id"`
    Email string `json:"email"`
}
```

---

## Supporting Layers

### Middleware

**Location:** `internal/middleware/`

**Handles:**
- Authentication/Authorization
- Logging
- CORS
- Rate limiting
- Request ID injection

### Config

**Location:** `internal/config/`

**Handles:**
- Loading environment variables
- Application configuration
- Feature flags

### Utils/Helpers

**Location:** `internal/utils/` or `pkg/`

**Handles:**
- Shared utility functions
- Common helpers (string manipulation, date formatting)
- Custom error types

---

## Data Flow Example

```
HTTP Request
    │
    ▼
┌─────────────────┐
│  Middleware     │  Auth, logging, etc.
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Handler        │  Parse request, validate input format
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Service        │  Business logic, orchestration
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Repository     │  Database operations
└────────┬────────┘
         │
         ▼
    Database
```

---

## Key Principles

1. **Dependency Direction:** Outer layers depend on inner layers, never the reverse
2. **Single Responsibility:** Each layer has one clear purpose
3. **Testability:** Layers can be tested independently using mocks/interfaces
4. **Separation of Concerns:** HTTP logic stays in handlers, business logic in services, data access in repositories
