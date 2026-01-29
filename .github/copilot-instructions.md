# PPPKJson AI Development Guide

Project inside at folder Backend


## Architecture Overview
This is a **Hexagonal Architecture** Go application for PPPK exam management with randomized questions. The system ensures each user receives different questions across 4 categories (MANAJERIAL, SOSIAL_KULTURAL, TEKNIS, WAWANCARA) with 5 questions each.

### Core Directory Structure
- `cmd/server/` - Main application entry with Swagger docs
- `cmd/seeder/` - Database seeder that loads JSON question banks
- `internal/ports/` - AdapterPort interface defining Connect/Disconnect/IsReady/Value methods
- `internal/adapters/` - Infrastructure implementations (db_adapter, gin_adapter, logger)
- `internal/repositories/models/` - GORM domain models with soft delete support
- `internal/handlers/` - Gin HTTP handlers with Swagger annotations
- `migrations/` - SQL migrations with `migrate` CLI integration

## Key Development Patterns

### Database Models
All models follow GORM conventions with:
- Soft deletes: `gorm.DeletedAt` field with `gorm:"index"`
- CASCADE foreign keys: `constraint:OnDelete:CASCADE`
- Custom table names: implement `TableName() string`
- Example: [`models/question.go`](internal/repositories/models/question.go) and [`models/exam.go`](internal/repositories/models/exam.go)

### Adapter Pattern Implementation
- All adapters implement [`ports.AdapterPort`](internal/ports/adapter_port.go)
- Database connection: [`db_adapter.go`](internal/adapters/db_adapter/db_adapter.go) with connection pooling
- HTTP server: [`gin_adapter.go`](internal/adapters/gin_adapter/gin_adapter.go)
- Initialization pattern in [`main.go`](cmd/server/main.go) uses ConnectManager slice

### API Development
- Handlers use Swagger annotations: `@Summary`, `@Description`, `@Param`, `@Success`
- Route groups: `/api/v1/exam/{userID}` pattern and `/api/v1/dashboard/*` for admin routes
- Response format: `dto.APIResponse{Success, Message, Data}` structure
- Example: [`gin_exam_handler.go`](internal/handlers/gin_exam_handler.go)

## Essential Development Commands

### Database Operations
```bash
# Start PostgreSQL container
docker compose -f compose.dev.yaml up postgres -d

# Run migrations
make migrate-up           # Apply all pending migrations
make migrate-down         # Rollback last migration
make migrate-create name=feature_name  # Create new migration files

# Seed database from JSON files
make db-seed             # Load questions from migrations/data/*/
```

### Development Workflow
```bash
# Development server with auto-reload
make dev-go-server       # Uses gow for file watching

# Generate Swagger documentation
make swag-gen           # Updates docs/ directory

# Container operations
docker compose -f compose.dev.yaml up postgres  # Just database
docker compose -f compose.dev.yaml down         # Clean shutdown
```

## Project-Specific Conventions

### Question Bank Management
- JSON files organized in `migrations/data/{category}/` directories
- Seeder automatically discovers and loads all JSON files
- Each question has `id`, `category`, `question_text`, and `options[]` with scores 1-4
- Categories must match: MANAJERIAL, SOSIAL_KULTURAL, TEKNIS, WAWANCARA
- Current configuration: 1 question per category (4 total) for testing

### Dashboard API Endpoints
- **Individual Dashboard**: `/exam/{userID}/dashboard` - Shows user's exam status, progress, and results
- **Admin Dashboard**: `/dashboard/users` - Lists all users with their exam status and results
- Dashboard responses include exam status (NO_EXAM, NOT_STARTED, IN_PROGRESS, COMPLETED, EXPIRED)
- For completed exams: shows detailed scores, percentages, grades, and pass/fail status
- For in-progress exams: shows answered count and remaining time
- Auto-updates expired sessions before returning data

### Exam Session Logic
- User ID is hardcoded from URL path (e.g., `/exam/1234`)
- Session codes: `EXAM_{userID}_{timestamp}` format
- Random question selection: 1 per category (configurable via `questionsPerCategory`), stored in `exam_questions` table
- Status flow: NOT_STARTED → IN_PROGRESS → COMPLETED/EXPIRED
- Dashboard endpoints: individual (`/exam/{userID}/dashboard`) and admin (`/dashboard/users`)

### Environment Configuration
Key environment variables with defaults:
- `PORT=8080`, `GIN_MODE=release`
- Database: `DB_HOST=localhost`, `DB_USER=encang_cutbray`, `DB_NAME=togotestgo`
- Configuration in [`utils.go`](internal/utils/utils.go) using `GetEnvOrDefault()`

### Migration Best Practices
- Use sequential naming: `000001_`, `000002_`
- Include both `.up.sql` and `.down.sql` files
- Add proper indexes for GORM soft deletes
- Use `BIGSERIAL` for primary keys to match GORM uint type

## Integration Points

### External Dependencies
- **GORM v1.31**: ORM with PostgreSQL driver
- **Gin v1.11**: HTTP framework with middleware support  
- **Swaggo**: Auto-generates OpenAPI/Swagger documentation
- **migrate/migrate:4**: Database migration tool in Docker

### Critical File Dependencies
- [`compose.dev.yaml`](compose.dev.yaml): PostgreSQL 18 with health checks
- [`Makefile`](Makefile): All development commands delegate to Docker/migrate
- [`docs/swagger.json`](docs/swagger.json): Auto-generated from handler annotations
- [`migrations/data/`](migrations/data/): Question bank JSON structure

When modifying the codebase, always consider the hexagonal architecture boundaries and ensure proper interface implementation for adapters.