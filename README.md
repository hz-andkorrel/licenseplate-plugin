# License Plate Recognition Plugin

A simple license plate recognition plugin for the Hotel Hub modular platform.

## Features

- Scan and recognize license plates
- Store license plate records with guest information
- Search and retrieve license plate data
- Integration with Hotel Hub broker

## API Endpoints

- `POST /api/licenseplate/scan` - Scan a license plate (accepts image or text)
- `GET /api/licenseplate/records` - Get all license plate records
- `GET /api/licenseplate/records/:plate` - Get specific plate record
- `DELETE /api/licenseplate/records/:plate` - Delete a plate record

## Setup

1. Install dependencies:
```bash
go mod download
```

2. Set up PostgreSQL database (if not already set up):
```bash
# Using the broker's database
# Make sure PostgreSQL is running on localhost:5432
# Default credentials: broker/broker123
# Default database: broker_db
```

3. Run database migrations:
```bash
# Connect to PostgreSQL and run the migration
psql -U broker -d broker_db -f migrations/001_create_license_plates_table.sql

# Or on Windows with PowerShell:
.\setup-database.ps1
```

4. Configure environment variables in `.env`:
```
PORT=8082
HOST=localhost
BROKER_URL=http://localhost:8081
BROKER_AUTH_TOKEN=your_jwt_token_here
BASE_API_ROUTE=/api/licenseplate
DATABASE_URL=postgres://broker:broker123@localhost:5432/broker_db?sslmode=disable
```

5. Run the plugin:
```bash
go run main.go
```

## Integration

The plugin automatically registers with the broker on startup. Ensure the broker is running before starting this plugin.

## Development

This is a simple implementation for demonstration purposes. In production, you would integrate with actual license plate recognition services like:
- OpenALPR
- Plate Recognizer
- Google Cloud Vision API
- AWS Rekognition
