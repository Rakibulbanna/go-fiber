# Go Fiber GORM API with Authentication

A production-ready REST API built with Go Fiber, GORM, and PostgreSQL, featuring JWT-based authentication and protected routes.

## Features

- ✅ User authentication (Sign Up, Login/Sign In)
- ✅ JWT-based authentication
- ✅ Protected routes with middleware
- ✅ Book CRUD operations
- ✅ User ownership tracking for books
- ✅ Clean architecture (controllers, services, repositories)
- ✅ Hot reloading with Air
- ✅ Environment-based configuration

## Project Structure

```
.
├── config/          # Configuration management
├── controllers/     # HTTP handlers
├── dtos/            # Data Transfer Objects
├── middleware/     # HTTP middleware (auth, etc.)
├── models/          # Database models
├── repositories/    # Data access layer
├── routes/          # Route definitions
├── services/        # Business logic
├── storage/         # Database connection
├── utils/           # Utility functions (JWT, password hashing)
├── main.go          # Application entry point
├── .air.toml        # Air configuration for hot reload
└── .env.dev.example # Environment variables template
```

## Prerequisites

- Go 1.23.5 or higher
- PostgreSQL
- Make (optional, for using Makefile commands)
  - **Windows**: Install [Make for Windows](http://gnuwin32.sourceforge.net/packages/make.htm) or use Git Bash/WSL
  - **Linux/Mac**: Usually pre-installed

## Quick Start

### Using Makefile (Recommended)

1. **Initial Setup**

   ```bash
   make setup
   ```

   This will install tools, download dependencies, and check your environment file.

2. **Configure environment variables**

   - Copy `.env.dev.example` to `.env.dev`
   - Update the database credentials and JWT secret

3. **Create the database**

   ```sql
   CREATE DATABASE fiber_demo;
   ```

4. **Run in development mode (with hot reload)**
   ```bash
   make air
   # or
   make dev
   ```

### Manual Setup

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd go-api-with-fiber-gorm
   ```

2. **Install dependencies**

   ```bash
   go mod download
   ```

3. **Configure environment variables**

   - Copy `.env.dev.example` to `.env.dev`
   - Update the database credentials and JWT secret:
     ```env
     DB_HOST=localhost
     DB_PORT=5432
     DB_USER=postgres
     DB_PASSWORD=your_password
     DB_NAME=fiber_demo
     DB_SSLMODE=disable
     JWT_SECRET=your-strong-secret-key
     PORT=8080
     ```

4. **Create the database**

   ```sql
   CREATE DATABASE fiber_demo;
   ```

5. **Run migrations**
   The application automatically runs migrations on startup.

## Running the Application

### Using Makefile

```bash
# Development mode (hot reload with Air)
make air
# or
make dev

# Production mode (build and run)
make prod

# Run with go run
make run

# Build only
make build
```

### Manual Commands

#### Development Mode (with Air hot reload)

```bash
air
```

#### Production Mode

```bash
go run main.go
```

The server will start on `http://localhost:8080` (or the port specified in your `.env.dev` file).

### Using Docker Compose

#### Production Mode

```bash
# Start all services (API + PostgreSQL)
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down
```

#### Development Mode (with hot reload)

```bash
# Start development environment
make docker-dev

# View logs
make docker-logs-dev

# Stop development services
make docker-down-dev
```

#### Manual Docker Compose Commands

```bash
# Production
docker-compose up -d
docker-compose logs -f
docker-compose down

# Development
docker-compose -f docker-compose.dev.yml up -d
docker-compose -f docker-compose.dev.yml logs -f
docker-compose -f docker-compose.dev.yml down
```

## Makefile Commands

The project includes a comprehensive Makefile with the following commands:

### Setup & Dependencies

- `make setup` - Complete initial setup (install tools, deps, check env)
- `make deps` - Download and install dependencies
- `make tidy` - Clean up dependencies
- `make install-tools` - Install required tools (air, etc.)

### Database

- `make migrate` - Run database migrations

### Build & Run

- `make build` - Build the application
- `make run` - Run with `go run`
- `make air` - Run with air (hot reload)
- `make dev` - Development mode (check env + air)
- `make prod` - Production mode (build + run)
- `make run-prod` - Run the built binary

### Cross-platform Builds

- `make build-linux` - Build for Linux
- `make build-windows` - Build for Windows
- `make build-mac` - Build for macOS

### Testing

- `make test` - Run all tests
- `make test-verbose` - Run tests with verbose output
- `make test-coverage` - Run tests with coverage report

### Code Quality

- `make fmt` - Format code
- `make vet` - Run go vet
- `make lint` - Run linter (requires golangci-lint)

### Cleanup

- `make clean` - Remove build artifacts
- `make clean-all` - Remove all generated files and logs
- `make stop` - Stop running processes

### Docker

- `make docker-build` - Build Docker image
- `make docker-run` - Run Docker container
- `make docker-stop` - Stop Docker container
- `make docker-up` - Start services with docker-compose (production)
- `make docker-dev` - Start development environment with docker-compose (hot reload)
- `make docker-down` - Stop services with docker-compose
- `make docker-down-dev` - Stop development services
- `make docker-logs` - View docker-compose logs
- `make docker-logs-dev` - View development docker-compose logs
- `make docker-restart` - Restart services
- `make docker-clean` - Remove containers, volumes, and networks
- `make docker-ps` - Show running containers
- `make docker-setup` - Setup with Docker (start services)

### Utilities

- `make check-env` - Check if .env.dev file exists
- `make help` - Show all available commands

For detailed information, run `make help`.

## API Endpoints

### Authentication

#### Sign Up

```http
POST /api/auth/signup
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe"
}
```

**Response:**

```json
{
  "message": "User created successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "name": "John Doe"
    }
  }
}
```

#### Login / Sign In

```http
POST /api/auth/login
# or
POST /api/auth/signin
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**

```json
{
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "name": "John Doe"
    }
  }
}
```

### Books

#### Get All Books (Public)

```http
GET /api/books
```

#### Get Book by ID (Public)

```http
GET /api/books/:id
```

#### Create Book (Protected)

```http
POST /api/books
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "The Go Programming Language",
  "author": "Alan Donovan",
  "publisher": "Addison-Wesley",
  "year": 2015
}
```

#### Update Book (Protected - Owner only)

```http
PUT /api/books/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Updated Title",
  "author": "Updated Author",
  "publisher": "Updated Publisher",
  "year": 2024
}
```

#### Delete Book (Protected - Owner only)

```http
DELETE /api/books/:id
Authorization: Bearer <token>
```

## Authentication

All protected routes require a JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

Tokens are valid for 24 hours from issuance.

## Security Features

- Password hashing using bcrypt
- JWT token-based authentication
- Protected routes with middleware
- User ownership verification for book operations
- Input validation

## Development

### Project Architecture

The project follows clean architecture principles:

- **Controllers**: Handle HTTP requests and responses
- **Services**: Contain business logic
- **Repositories**: Handle database operations
- **Models**: Define database schemas
- **DTOs**: Data transfer objects for API requests/responses
- **Middleware**: Authentication and other cross-cutting concerns

### Adding New Features

1. Create model in `models/`
2. Add migration in `models/migrations.go`
3. Create repository in `repositories/`
4. Create service in `services/`
5. Create controller in `controllers/`
6. Add routes in `routes/routes.go`

## Testing

Example API calls using curl:

```bash
# Sign Up
curl -X POST http://localhost:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","name":"Test User"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Create Book (replace TOKEN with actual token)
curl -X POST http://localhost:8080/api/books \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{"title":"My Book","author":"John Doe","publisher":"Publisher","year":2024}'

# Get All Books
curl http://localhost:8080/api/books
```

## Production Considerations

- Use a strong, random JWT secret in production
- Enable SSL/TLS for database connections (`DB_SSLMODE=require`)
- Use environment-specific configuration files
- Implement rate limiting
- Add logging and monitoring
- Set up proper CORS policies
- Use HTTPS in production
- Implement database connection pooling
- Add request validation middleware
- Set up health check endpoints

## License

MIT
