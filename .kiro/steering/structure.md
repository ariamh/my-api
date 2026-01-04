# Project Structure

```
├── cmd/api/main.go          # Application entry point
├── internal/                 # Private application code
│   ├── config/              # Configuration loading, database setup, migrations
│   ├── handler/             # HTTP handlers (controllers)
│   ├── middleware/          # Fiber middleware (auth, logging, security)
│   ├── model/               # GORM models with Base embedding
│   ├── repository/          # Data access layer with generic BaseRepository
│   ├── router/              # Route definitions
│   └── service/             # Business logic layer
├── pkg/                     # Reusable packages
│   ├── jwt/                 # JWT token management
│   ├── logger/              # Zap logger wrapper
│   ├── response/            # Standardized API responses
│   └── validator/           # Input validation wrapper
├── docs/                    # Generated Swagger documentation
└── migrations/              # Database migrations
```

## Architecture Pattern

Layered architecture with dependency injection:

1. **Handler** → Receives HTTP requests, validates input, calls service
2. **Service** → Business logic, calls repository, returns DTOs
3. **Repository** → Data access via GORM, uses generic BaseRepository[T]

## Key Conventions

- Models embed `model.Base` for ID (UUID), timestamps, and soft delete
- Services define interfaces and domain errors (e.g., `ErrUserNotFound`)
- Handlers use `pkg/response` for consistent JSON responses
- Input/output DTOs defined in service layer with validation tags
- Swagger annotations on handler methods for API documentation
- Constructor pattern: `NewXxxHandler()`, `NewXxxService()`, `NewXxxRepository()`
