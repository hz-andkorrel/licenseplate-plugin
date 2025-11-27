# License Plate Recognition Plugin

A simple license plate recognition plugin for the Hotel Hub modular platform.

## Features

- Scan and recognize license plates
- Store license plate records with guest information
- Search and retrieve license plate data
- Integration with Hotel Hub broker
- **XPOTS webhook integration** for automatic license plate detection

## API Endpoints

### Manual Entry
- `POST /api/licenseplate/scan` - Scan a license plate (accepts image or text)
- `GET /api/licenseplate/records` - Get all license plate records
- `GET /api/licenseplate/records/:plate` - Get specific plate record
- `DELETE /api/licenseplate/records/:plate` - Delete a plate record

### XPOTS Webhook Integration
- `POST /api/licenseplate/webhook/xpots` - Receive automatic plate scans from XPOTS
- `GET /api/licenseplate/webhook/info` - Get webhook configuration information

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
WEBHOOK_API_KEY=your-secure-webhook-key-here
```

5. Run the plugin:
```bash
go run main.go
```

## XPOTS Integration

This plugin supports automatic license plate detection through XPOTS webhooks.

### Configuring XPOTS

1. **Generate a secure API key**:
   ```bash
   # On Linux/Mac
   openssl rand -base64 32
   
   # On Windows PowerShell
   -join ((48..57) + (65..90) + (97..122) | Get-Random -Count 32 | ForEach-Object {[char]$_})
   ```

2. **Set the API key** in your `.env` file:
   ```
   WEBHOOK_API_KEY=your-generated-secure-key
   ```

3. **Configure XPOTS** to send webhooks to:
   - **URL**: `http://your-server:8082/api/licenseplate/webhook/xpots`
   - **Method**: `POST`
   - **Authentication**: Add header `Authorization: Bearer your-generated-secure-key`

### Webhook Payload

XPOTS should send JSON data in this format:
```json
{
  "event_type": "entry",
  "plate_number": "ABC-123",
  "timestamp": "2025-11-27T10:30:00Z",
  "location": "Main Gate",
  "confidence": 0.98,
  "camera_id": "CAM-001",
  "direction": "in"
}
```

### Event Types

- **entry/scan/in**: Vehicle entering - creates new record or updates check-in time
- **exit/out**: Vehicle exiting - updates check-out time
- Unknown event types are treated as entry events

### Automatic Behavior

When XPOTS detects a license plate:
- If the plate exists in the system → updates check-in/check-out time
- If the plate is unknown → creates a new record with "Unknown Guest (Auto-detected)"
- Guest info can be manually updated later through the API or hotel management system
- All auto-detected records include XPOTS metadata (location, confidence, camera ID) in notes

### Testing Webhooks

Send a test webhook using curl:
```bash
curl -X POST http://localhost:8082/api/licenseplate/webhook/xpots \
  -H "Authorization: Bearer your-webhook-key" \
  -H "Content-Type: application/json" \
  -d '{
    "event_type": "entry",
    "plate_number": "TEST-123",
    "timestamp": "2025-11-27T10:30:00Z",
    "location": "Main Gate",
    "confidence": 0.98,
    "camera_id": "CAM-001",
    "direction": "in"
  }'
```

Get webhook configuration info:
```bash
curl http://localhost:8082/api/licenseplate/webhook/info
```

## Integration

The plugin automatically registers with the broker on startup. Ensure the broker is running before starting this plugin.

## Development

This plugin integrates with XPOTS for automatic license plate recognition. XPOTS is a professional parking management system with ANPR (Automatic Number Plate Recognition) technology.

For other implementations, you could integrate with services like:
- OpenALPR
- Plate Recognizer
- Google Cloud Vision API
- AWS Rekognition
