# FibNotes API

A production-ready REST API for notes management built with Go, Fiber web framework, and MySQL. The application provides secure user authentication, CRUD operations for notes, and comprehensive API documentation.

> [!IMPORTANT] 
>
> This Project also has postman collections as well as swagger docs powered by scalar
>


## Architecture Overview

This application follows a clean architecture pattern with clear separation of concerns:

- **Handlers**: HTTP request/response handling and business logic
- **Models**: Data structures and database models using GORM
- **Middleware**: Cross-cutting concerns like JWT authentication
- **Routes**: API endpoint definitions and routing
- **Utils**: Utility functions for JWT, hashing, and common operations
- **Database**: Connection management and configuration

## Technology Stack

- **Language**: Go 1.23+
- **Web Framework**: Fiber v2
- **Database**: MySQL 8.0
- **ORM**: GORM
- **Authentication**: JWT (golang-jwt/jwt)
- **Documentation**: Swagger/OpenAPI with Scalar
- **Containerization**: Docker & Docker Compose
- **Environment**: dotenv for configuration

## Quick Start

### Prerequisites
- Go 1.23 or higher
- Docker and Docker Compose
- MySQL 8.0 (if running locally)

### Using Docker (Recommended)

1. Clone the repository:
```bash
git clone https://github.com/kunalsinghdadhwal/fib_notes
cd fib_notes
```

2. Create environment configuration:
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. Start the application:
```bash
docker-compose up --build
```

The API will be available at `http://localhost:3000`

### Local Development

1. Install dependencies:
```bash
go mod download
```

2. Set up environment variables:
```bash
cp .env.example .env
```

3. Run the application:
```bash
go run main.go
```

## Database Seeding

The application includes a comprehensive database seeding tool for development and testing:

### Build the seeder:
```bash
go build -o bin/seed ./cmd/seeder
```

### Seed database with sample data:
```bash
./bin/seed                   # Creates 10 users with 5-10 notes each (default)
./bin/seed -n 25             # Creates 25 users
```

### Clear all data:
```bash
./bin/seed clear
```

### Using Docker:
```bash
docker exec -it notes_app ./seed help
docker exec -it notes_app ./seed -n 20
```



## API Endpoints

### Authentication
```
POST   /auth/register        - User registration
POST   /auth/login           - User login
POST   /auth/refresh         - Refresh access token
POST   /auth/logout          - User logout (protected)
GET    /auth/me              - Get current user info (protected)
PUT    /auth/change-password - Change user password (protected)
```

### Notes
```
POST   /notes                - Create a new note (protected)
GET    /notes                - List user notes with pagination (protected)
GET    /notes/:id            - Get specific note (protected)
PUT    /notes/:id            - Update note (protected)
DELETE /notes/:id            - Delete note (protected)
```

### System
```
GET    /                     - Redirect to API documentation
GET    /api/reference        - Interactive API documentation
GET    /health               - Health check endpoint
```

## Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | Server port | 3000 | No |
| `JWT_SECRET` | JWT signing secret | - | Yes |
| `DB_HOST` | Database host | localhost | Yes |
| `DB_PORT` | Database port | 3306 | Yes |
| `DB_USER` | Database username | - | Yes |
| `DB_PASS` | Database password | - | Yes |
| `DB_NAME` | Database name | - | Yes |


## API Documentation

Interactive API documentation is available at `/api/reference` when the server is running. The documentation includes:

- Complete endpoint specifications
- Request/response schemas
- Authentication requirements
- Example requests and responses
- Error code documentation

## Database Schema

### Users Table
- `id` (UUID, Primary Key)
- `name` (VARCHAR, NOT NULL)
- `email` (VARCHAR, UNIQUE, NOT NULL)
- `password` (VARCHAR, NOT NULL, hashed)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

### Notes Table
- `id` (INT, Primary Key, Auto Increment)
- `user_id` (UUID, Foreign Key, NOT NULL)
- `title` (VARCHAR, NOT NULL)
- `content` (TEXT, NOT NULL)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

## Testing

### Manual Testing
Use the included Postman collection:
```
Fib Notes.postman_collection.json
```

### Health Check
```bash
curl http://localhost:3000/health
```

Expected response:
```json
{
  "status": "ok",
  "message": "Server is running"
}
```
