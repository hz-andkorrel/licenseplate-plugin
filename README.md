# License Plate Recognition Plugin

Hotel parking management plugin with automatic XPOTS integration and visitor tracking.

## Quick Start

1. **Start PostgreSQL:**
```bash
docker-compose up -d
```

2. **Run migrations:**
```bash
Get-Content migrations/001_create_license_plates_table.sql | docker exec -i hotelhub-postgres psql -U broker -d broker_db
Get-Content migrations/002_add_visitor_support.sql | docker exec -i hotelhub-postgres psql -U broker -d broker_db
```

3. **Configure `.env`:**
```env
PORT=8082
DATABASE_URL=postgres://broker:broker123@localhost:5432/broker_db?sslmode=disable
WEBHOOK_API_KEY=your-secure-key
```

4. **Start the plugin:**
```bash
go run main.go
```

5. **Open UI:** http://localhost:8082/render

## Features

- 🚗 Manual license plate registration
- 📹 XPOTS webhook integration (automatic camera detection)
- 👥 Visitor types: Guest, Visitor, Staff, Delivery, Contractor, VIP
- ⏱️ Temporary access passes with expiration
- 🎨 Color-coded badges and expired access alerts

## API Endpoints

### Core Operations
- `POST /api/licenseplate/scan` - Register a license plate
- `GET /api/licenseplate/records` - Get all records
- `GET /api/licenseplate/records/:plate` - Get specific plate
- `DELETE /api/licenseplate/records/:plate` - Delete record

### XPOTS Integration
- `POST /api/licenseplate/webhook/xpots` - Receive XPOTS events
- `GET /api/licenseplate/webhook/info` - Webhook configuration

## XPOTS Setup

Configure XPOTS to send webhooks to:
- **URL:** `http://your-server:8082/api/licenseplate/webhook/xpots`
- **Header:** `Authorization: Bearer your-webhook-key`

Events supported: `entry`, `exit`, `scan`

## Architecture

```
├── main.go                 # Server setup & routes
├── public/                 # Frontend UI (HTML/CSS/JS)
├── internal/
│   ├── handlers/          # HTTP request handlers
│   ├── services/          # Business logic
│   ├── models/            # Data structures
│   ├── database/          # PostgreSQL connection
│   └── broker/            # Plugin registration
└── migrations/            # Database schema
```

## Documentation

- [XPOTS Integration Guide](./XPOTS_INTEGRATION.md) - Detailed webhook setup
- [Database Schema](./DATABASE.md) - Table structure and indexes

## License

MIT
