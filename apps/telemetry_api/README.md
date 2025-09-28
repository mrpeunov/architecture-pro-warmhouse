# Telemetry API

Telemetry API for Smart Home System that provides endpoints for managing devices and collecting telemetry data.

## Features

- Device management (get devices by ID, home ID, or type)
- Telemetry data collection and storage
- Integration with WarmHouse provider
- Kafka messaging for telemetry events
- Swagger documentation
- PostgreSQL database integration

## API Endpoints

### Devices
- `GET /api/v1/devices/{device_id}` - Get device by ID
- `GET /api/v1/devices/home/{home_id}` - Get devices by home ID
- `GET /api/v1/devices/type/{device_type}` - Get devices by device type

### Telemetry
- `GET /api/v1/telemetry/device/{device_id}` - Get telemetry data for device
- `GET /api/v1/telemetry/device/{device_id}/latest` - Get latest telemetry for device
- `POST /api/v1/telemetry/device/{device_id}/collect` - Collect telemetry from provider
- `POST /api/v1/telemetry` - Create new telemetry record

## Configuration

Copy `env.example` to `.env` and configure the following variables:

- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_USER` - Database user (default: postgres)
- `DB_PASSWORD` - Database password (default: postgres)
- `DB_NAME` - Database name (default: telemetry_db)
- `KAFKA_BROKER` - Kafka broker address (default: localhost:29092)
- `WARMHOUSE_API_URL` - WarmHouse provider URL (default: http://localhost:8080)
- `PORT` - Server port (default: :8083)

## Running the Application

### Using Docker

```bash
docker build -t telemetry-api .
docker run -p 8083:8083 telemetry-api
```

### Using Go

```bash
go mod download
go run main.go
```

## Swagger Documentation

Once the application is running, visit:
- http://localhost:8083/swagger/index.html

## Database Schema

The application uses PostgreSQL with the following tables:

- `devices` - Stores device information
- `telemetries` - Stores telemetry data

Run the `init.sql` script to create the required tables and indexes.

## Architecture

The application follows a layered architecture:

- **Handlers** - HTTP request handlers
- **Services** - Business logic layer
- **Repositories** - Data access layer
- **Providers** - External system integration
- **Models** - Data structures

## Kafka Integration

Telemetry data is automatically sent to Kafka topic `telemetry` when:
- New telemetry is created via API
- Telemetry is collected from external providers

The Kafka message format:
```json
{
  "device_id": "uuid",
  "value": 22.5,
  "timestamp": "2023-12-01T10:00:00Z"
}
```
