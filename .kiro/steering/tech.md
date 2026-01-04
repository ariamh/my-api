# Tech Stack

## Language & Framework

- Go 1.24
- Fiber v2 (web framework)
- GORM (ORM with PostgreSQL driver)

## Key Dependencies

- `github.com/gofiber/fiber/v2` - HTTP framework
- `gorm.io/gorm` + `gorm.io/driver/postgres` - Database ORM
- `github.com/golang-jwt/jwt/v5` - JWT authentication
- `github.com/go-playground/validator/v10` - Input validation
- `github.com/swaggo/swag` + `github.com/gofiber/swagger` - API documentation
- `go.uber.org/zap` - Structured logging
- `github.com/google/uuid` - UUID generation
- `github.com/joho/godotenv` - Environment configuration
- `golang.org/x/crypto/bcrypt` - Password hashing
- `github.com/stretchr/testify` - Testing assertions/mocks

## Common Commands

```bash
# Run the API
make run

# Run tests
make test

# Run tests with coverage
make test-cover

# Build binary
make build

# Generate Swagger docs
make swagger

# Start dev database (PostgreSQL via Docker)
make dev-db

# Stop dev database
make dev-db-down

# Docker operations
make docker-build    # Build container
make docker-up       # Start services
make docker-down     # Stop services
make docker-logs     # View logs

# Lint (requires golangci-lint)
make lint
```

## Configuration

Environment variables loaded from `.env` file:
- `APP_ENV` - Environment (development/production)
- `APP_PORT` - Server port (default: 3000)
- `APP_NAME` - Application name
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` - PostgreSQL config
- `JWT_SECRET` - JWT signing secret
- `JWT_EXPIRE_HOURS` - Token expiration (default: 24)
