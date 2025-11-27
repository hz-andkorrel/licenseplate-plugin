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

2. Configure environment variables in `.env`:
```
PORT=8082
HOST=localhost
BROKER_URL=http://localhost:8081
BROKER_AUTH_TOKEN=your_jwt_token_here
BASE_API_ROUTE=/api/licenseplate
```

3. Run the plugin:
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
