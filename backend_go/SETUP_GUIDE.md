# Go Backend Setup Guide

## What's Been Completed ✅

- ✅ All 9 repository implementations with full PostgreSQL database access
- ✅ All service implementations
- ✅ All handler implementations
- ✅ Route definitions

## What's Needed to Run

### 1. Database Setup (CRITICAL)

You have two options:

#### Option A: Share Database with Kotlin Backend (RECOMMENDED)

Since the Kotlin backend already has Flyway migrations that create all tables:

1. Use the same database for both backends
2. Let the Kotlin backend create the schema (it already has)
3. Point the Go backend to the same database

**Update your `.env` file:**
```bash
cp .env.example .env
# Edit .env and set DB_NAME to match your Kotlin backend database
# Usually: menunderfire
```

**In `cmd/api/main.go`:**
```go
// Comment out or remove the InitSchema call:
// if err := database.InitSchema(db); err != nil {
//     log.Fatalf("Failed to initialize database schema: %v", err)
// }
```

#### Option B: Add Migration Tool to Go Backend

```bash
cd backend_go
go get -u github.com/golang-migrate/migrate/v4
go get -u github.com/golang-migrate/migrate/v4/database/postgres
go get -u github.com/golang-migrate/migrate/v4/source/file
```

Then copy migrations from `backend/src/main/resources/db/migration/` to `backend_go/migrations/`

### 2. Complete main.go Wiring

Replace your current `cmd/api/main.go` with the example in `cmd/api/main_complete.go.example`:

```bash
cd backend_go/cmd/api
# Backup current main.go
mv main.go main.go.backup
# Use the complete version
mv main_complete.go.example main.go
```

Or manually update main.go to initialize all repositories, services, and handlers.

### 3. Install Dependencies

```bash
cd backend_go
go mod download
go mod tidy
```

### 4. Set Up Environment Variables

```bash
cd backend_go
cp .env.example .env
# Edit .env with your actual database credentials
```

Make sure these match your Kotlin backend database:
- `DB_NAME`: Should be the same database name (e.g., `menunderfire`)
- `DB_HOST`: Usually `localhost` if running locally
- `DB_PORT`: Usually `5432` for PostgreSQL
- `DB_USER`: Your PostgreSQL username
- `DB_PASSWORD`: Your PostgreSQL password

### 5. Build and Run

```bash
cd backend_go

# Build
go build -o bin/server cmd/api/main.go

# Or run directly
go run cmd/api/main.go
```

The server will start on `http://localhost:8080` (or whatever port you configured).

## Testing the Setup

### Health Check
```bash
curl http://localhost:8080/health
# Should return: {"status":"ok"}
```

### API Endpoints

The following routes are available (see `internal/routes/routes.go` for complete list):

- `GET /api/users/me` - Get current user profile
- `GET /api/users/me/permissions` - Get user permissions
- `GET /api/organizations` - List organizations
- `GET /api/groups` - List groups
- And many more...

## Common Issues

### Issue: "Failed to connect to database"
**Solution:** Check that PostgreSQL is running and your `.env` credentials are correct.

### Issue: "Table doesn't exist"
**Solution:** Make sure the Kotlin backend has run at least once to create the schema via Flyway migrations.

### Issue: "routes.Setup expects 10 arguments"
**Solution:** You need to initialize all 10 handlers. Use the example in `main_complete.go.example`.

### Issue: Port already in use
**Solution:** Either stop the Kotlin backend or change `SERVER_PORT` in `.env` to a different port (e.g., 8081).

## Architecture

```
cmd/api/main.go
    ↓
internal/
├── config/          - Environment configuration
├── database/        - Database connection
├── repositories/    - Data access layer (✅ ALL IMPLEMENTED)
│   └── postgres/    - PostgreSQL implementations
├── services/        - Business logic layer
├── handlers/        - HTTP handlers
└── routes/          - Route definitions
```

## Next Steps

1. Set up authentication middleware (currently missing)
2. Add proper error handling and logging
3. Set up CORS for frontend integration
4. Add unit tests for repositories
5. Add integration tests
