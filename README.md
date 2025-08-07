# Go Quiz API with PostgreSQL

A RESTful quiz API built with Go and PostgreSQL, featuring proper database integration with migrations, repository pattern, and environment configuration.

## Features

✅ **PostgreSQL Integration**: Full database connectivity using `lib/pq` driver  
✅ **Database Models**: Proper Go structs with JSON and database tags  
✅ **Migration System**: Automated database schema migrations  
✅ **Repository Pattern**: Clean data access layer with proper separation  
✅ **Environment Configuration**: Database connection via environment variables  
✅ **Data Seeding**: Automatic population of initial quiz questions  
✅ **CORS Support**: Cross-origin request handling  
✅ **Health Check**: Service health monitoring endpoint  

## Architecture

```
├── config/          # Environment configuration
├── database/        # Database connection and seeding
├── migrations/      # Database schema migrations
├── models/          # Data models and structs
├── repository/      # Data access layer
└── main.go         # HTTP server and API routes
```

## Database Schema

- **questions**: Quiz questions with timestamps
- **options**: Question options with correct answer flags
- **quiz_sessions**: User quiz sessions (future use)
- **quiz_results**: User answer tracking (future use)
- **schema_migrations**: Migration version tracking

## Environment Variables

```bash
# Database Configuration
DB_HOST=localhost          # Database host
DB_PORT=5432              # Database port
DB_USER=postgres          # Database user
DB_PASSWORD=postgres      # Database password
DB_NAME=postgres          # Database name
DB_SSL_MODE=disable       # SSL mode (disable/require/verify-full)

# Server Configuration
PORT=8080                 # Server port (optional, defaults to 8080)
```

## Quick Start

1. **Setup Environment** (optional):
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

2. **Start PostgreSQL** (if not running):
   ```bash
   # Using Docker
   docker run --name postgres-quiz \
     -e POSTGRES_PASSWORD=postgres \
     -p 5432:5432 \
     -d postgres:15
   ```

3. **Run the Application**:
   ```bash
   go mod tidy
   go run .
   ```

4. **Verify Setup**:
   ```bash
   # Check health
   curl http://localhost:8080/health
   
   # Get all questions
   curl http://localhost:8080/api/questions
   
   # Get specific question
   curl http://localhost:8080/api/questions/1
   ```

## API Endpoints

### Questions
- `GET /api/questions` - Get all quiz questions
- `GET /api/questions/{id}` - Get specific question by ID

### System
- `GET /health` - Health check endpoint
- `GET /` - Welcome message

## Database Operations

The application automatically:

1. **Connects** to PostgreSQL using environment configuration
2. **Runs migrations** to create required tables
3. **Seeds data** with initial quiz questions (if empty)
4. **Serves data** from PostgreSQL instead of in-memory storage

## Development

### Adding New Questions

Questions are automatically seeded on first run. To add more questions, modify the `SeedData` function in `database/seed.go`.

### Creating Migrations

Add new migrations to the `GetMigrations()` function in `migrations/migrations.go`:

```go
{
    Version: 6,
    Name:    "add_new_table",
    Up:      `CREATE TABLE new_table (...);`,
    Down:    `DROP TABLE new_table;`,
}
```

### Repository Pattern

The repository pattern separates database operations from business logic:

```go
// Get question repository
questionRepo := repository.NewQuestionRepository(db)

// Use repository methods
questions, err := questionRepo.GetAll()
question, err := questionRepo.GetByID(1)
```

## Production Deployment

1. Set proper environment variables
2. Use connection pooling for high-traffic scenarios
3. Enable SSL mode for secure connections
4. Consider using database migrations in CI/CD pipeline

## Dependencies

- **github.com/lib/pq**: PostgreSQL driver for Go
- **Standard library**: HTTP server, JSON encoding, environment variables

---

🐹 **Go Quiz API** - Built with Go 1.21+ and PostgreSQL
