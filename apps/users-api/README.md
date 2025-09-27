# Users API

A FastAPI-based microservice for user management and home creation, built according to the UserAPI Class Diagram specification.

## Features

- User registration and authentication
- JWT-based authentication
- Home creation and management
- User-home relationship management
- PostgreSQL database integration with asyncpg

## Architecture

The application follows a clean architecture pattern with the following layers:

- **Controllers**: Handle HTTP requests and responses
- **Services**: Business logic layer
- **Repositories**: Data access layer
- **Models**: Data models and DTOs

## API Endpoints

### Authentication
- `POST /auth/register` - Register a new user
- `POST /auth/login` - Login and get access token

### Homes
- `POST /homes` - Create a new home (requires authentication)
- `GET /homes?user_id={id}` - Get homes for a specific user (requires authentication)

## Setup

1. Install dependencies:
```bash
pip install -r requirements.txt
```

2. Set up environment variables:
```bash
cp env.example .env
# Edit .env with your configuration
```

3. Set up PostgreSQL database and run the init.sql script

4. Run the application:
```bash
uvicorn main:app --reload
```

## Docker

Build and run with Docker:

```bash
docker build -t users-api .
docker run -p 8000:8000 users-api
```

## Environment Variables

- `DATABASE_URL`: PostgreSQL connection string
- `SECRET_KEY`: JWT secret key
- `DEBUG`: Enable debug mode
