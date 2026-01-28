# PPPKJson - Hexagonal Architecture Guide

## Overview
This project implements **Hexagonal Architecture** (Ports and Adapters pattern) for building a PPPK exam question management system.

## Architecture Structure

```
PPPKJson/
├── cmd/                     # Application Entry Points
│   ├── server/             # Main application server
│   │   ├── main.go         # Server entry point  
│   │   └── config.go       # Backward compatibility config
│   ├── seeder/             # Database seeding utility
│   │   └── main.go         # Seeder entry point
│   └── config/             # 🔧 SHARED CONFIG
│       └── config.go       # Re-exports for compatibility
├── internal/               # Private application code
│   ├── adapters/          # 🔌 ADAPTERS (Infrastructure)
│   │   ├── db_adapter/    # Database connection adapter
│   │   └── logger/        # Logging adapter  
│   ├── ports/             # 🚪 PORTS (Interfaces)
│   │   └── adapter_port.go # Core adapter interface
│   ├── repositories/      # 📊 REPOSITORIES (Data Layer)
│   │   └── models/        # Domain models (question.go)
│   └── utils/             # 🛠️ UTILITIES
├── migrations/            # 📋 DATABASE MIGRATIONS
│   ├── *.up.sql          # Migration up files
│   └── *.down.sql        # Migration down files
├── data/                  # 📄 DATA FILES
│   ├── manajerial/        # Management questions JSON
│   ├── sosial_kultural/   # Social cultural questions JSON
│   ├── teknis/           # Technical questions JSON  
│   └── wawancara/        # Interview questions JSON
├── compose.dev.yaml      # Docker compose for development
├── Makefile             # Build and migration commands
└── go.mod               # Go module definition
```

## Hexagonal Architecture Components

### 1. 🎯 **Core Domain** (Center of Hexagon)
- **Location**: `internal/repositories/models/`
- **Purpose**: Business entities and domain logic
- **Files**: 
  - `question.go` - Question and QuestionOption models with GORM tags
- **Features**: 
  - Soft delete support with `gorm.DeletedAt`
  - Foreign key relationships with CASCADE delete
  - JSON serialization tags

### 2. 🚪 **Ports** (Interfaces)
- **Location**: `internal/ports/`
- **Purpose**: Define contracts between core and adapters
- **Interface**:
```go
type AdapterPort interface {
    Connect(ctx context.Context) error
    Disconnect(ctx context.Context) error  
    IsReady() bool
    Value() any
}
```

### 3. 🔌 **Adapters** (Infrastructure)
- **Location**: `internal/adapters/`
- **Purpose**: Implement ports for external systems

#### Primary Adapters (Left Side - Driving)
- **Server**: `cmd/server/main.go` - Main application server (currently basic setup)
- **Seeder**: `cmd/seeder/main.go` - Data seeding interface (fully functional)
- **Config**: `cmd/config/config.go` - Shared configuration management

#### Secondary Adapters (Right Side - Driven)  
- **Database**: `internal/adapters/db_adapter/` - PostgreSQL adapter
- **Logger**: `internal/adapters/logger/` - Logging adapter

## Implementation Guide

### 🔄 **Adding New Adapters**

1. **Create Adapter Interface** (if needed):
```go
// internal/ports/new_service_port.go
type NewServicePort interface {
    DoSomething(ctx context.Context) error
}
```

2. **Implement Adapter**:
```go
// internal/adapters/new_service_adapter/adapter.go
type newServiceAdapter struct {
    // configuration fields
}

func (n *newServiceAdapter) Connect(ctx context.Context) error {
    // implementation
    return nil
}

func (n *newServiceAdapter) Disconnect(ctx context.Context) error {
    // cleanup implementation
    return nil
}

func (n *newServiceAdapter) IsReady() bool {
    // readiness check
    return true
}

func (n *newServiceAdapter) Value() any {
    // return underlying connection/value
    return n
}
```

3. **Register in Connection Manager**:
```go
// In cmd/server/main.go or cmd/seeder/main.go
adapters := []config.ConnectManager{
    {Name: "Database", Adapter: dbAdapter},
    {Name: "New Service", Adapter: newServiceAdapter},
}
```

### 📊 **Adding New Domain Models**

1. **Create Model**:
```go
// internal/repositories/models/new_model.go
type NewModel struct {
    ID        string    `gorm:"primaryKey" json:"id"`
    Name      string    `gorm:"not null" json:"name"`
    CreatedAt time.Time `json:"created_at"`
}
```

2. **Add to Migration**:
```sql
-- migrations/create_new_table.up.sql
CREATE TABLE new_models (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 🚀 **Running the Application**

#### Database Setup:
```bash
# Start containers
docker compose -f compose.dev.yaml up -d

# Run migrations  
make migrate-up

# Seed data
go run cmd/seeder/main.go
```

#### Development:
```bash
# Run server
go run cmd/server/main.go

# Or use Makefile
make dev-go-server
```

## Benefits of This Architecture

### 🎯 **Separation of Concerns**
- Core business logic independent of infrastructure
- Easy to test core logic without external dependencies
- Clear boundaries between layers

### 🔄 **Flexibility** 
- Switch database from PostgreSQL to MongoDB without changing core
- Replace logger implementation without affecting business logic
- Add new interfaces (REST, GraphQL, gRPC) easily

### 🧪 **Testability**
- Mock adapters for unit testing
- Test core logic in isolation
- Integration tests with real adapters

### 📈 **Scalability**
- Add new adapters without modifying existing code
- Horizontal scaling through adapter implementations
- Independent deployment of components

## Configuration Management

### Connection Management:
```go
// cmd/config/config.go (shared configuration)
type ConnectManager struct {
    Name    string
    Adapter ports.AdapterPort
}

// Functions:
// - ConnectAdapters(ctx context.Context, adapters ...ConnectManager) error
// - DisconnectAdapters(adapters ...ConnectManager)
```

### Backward Compatibility:
```go
// cmd/server/config.go
// Re-exports for backward compatibility:
type connectManager = config.ConnectManager
var connectAdapters = config.ConnectAdapters
var disconnectAdapters = config.DisconnectAdapters
```

### Environment Variables:
- `DB_HOST` - Database host (default: localhost)
- `DB_USER` - Database user (default: encang_cutbray)  
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name (default: togotestgo)
- `DB_PORT` - Database port (default: 5432)

## Data Flow

```
CLI/HTTP Request → Primary Adapter → Core Domain → Secondary Adapter → External System
     ↑                ↑                 ↑              ↑                    ↑
  cmd/server/      cmd/config/      repositories/  internal/adapters/   PostgreSQL
  cmd/seeder/                        models/        db_adapter/
```

### Current Implementation:
- **Seeder Flow**: `cmd/seeder/main.go` → `config.ConnectManager` → `models.Question` → `db_adapter` → PostgreSQL
- **Server Flow**: `cmd/server/main.go` → `config.ConnectManager` → (future: business logic) → `db_adapter` → PostgreSQL

## Best Practices

### ✅ **DO**
- Keep core domain free of infrastructure dependencies
- Use dependency injection through ports
- Handle errors at adapter boundaries
- Use context for cancellation and timeouts
- Keep adapters focused on single responsibility

### ❌ **DON'T**
- Import infrastructure packages in core domain
- Put business logic in adapters
- Couple adapters to specific implementations
- Ignore error handling
- Mix concerns across layers

## Next Steps for Development

1. **Add Service Layer** - Business logic orchestration
   - Create `internal/services/` package
   - Implement `QuestionService` with business rules
   - Add validation and business logic

2. **Add REST API** - HTTP adapter for external access  
   - Implement HTTP handlers in `internal/adapters/http_adapter/`
   - Add routing and middleware
   - Create API documentation

3. **Add Repository Layer** - Data access abstraction
   - Create `internal/repositories/question_repository.go`
   - Abstract database operations from models
   - Implement repository interfaces

4. **Add Validation** - Input validation adapter
   - Create validation middleware
   - Add struct validation tags
   - Implement custom business rules

5. **Add Testing** - Comprehensive test coverage
   - Unit tests for core domain
   - Integration tests with mock adapters
   - End-to-end tests

6. **Add Monitoring** - Metrics and health check adapters
   - Health check endpoints
   - Prometheus metrics
   - Logging improvements

## Current Project Status

### ✅ **Completed Components:**
- ✅ **Database Schema**: PostgreSQL with migrations
- ✅ **Models**: Question and QuestionOption with GORM
- ✅ **Seeder**: Functional data import from JSON files  
- ✅ **Configuration**: Shared config management
- ✅ **Database Adapter**: PostgreSQL connection with GORM
- ✅ **Logger Adapter**: Basic logging functionality

### 🚧 **In Progress:**
- 🚧 **Server Implementation**: Basic structure exists, needs business logic
- 🚧 **Repository Pattern**: Models exist but no repository interface yet

### 📋 **To Do:**
- 📋 **HTTP API**: REST endpoints for question management
- 📋 **Business Logic**: Service layer with validation
- 📋 **Testing**: Unit and integration tests
- 📋 **Documentation**: API documentation and examples

This hexagonal architecture provides a solid foundation for building maintainable, testable, and scalable applications.