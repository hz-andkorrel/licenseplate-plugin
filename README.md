# License Plate Plugin

Hotel parking management with XPOTS camera integration and visitor tracking.

## Quick Start

```bash
# 1. Start database
docker-compose up -d

# 2. Run migrations
Get-Content migrations/001_create_license_plates_table.sql | docker exec -i hotelhub-postgres psql -U broker -d broker_db
Get-Content migrations/002_add_visitor_support.sql | docker exec -i hotelhub-postgres psql -U broker -d broker_db
Get-Content migrations/003_create_parking_events.sql | docker exec -i hotelhub-postgres psql -U broker -d broker_db

# 3. Start plugin
go run main.go

# 4. Open UI
http://localhost:8082/render
```

## Features

- 🚗 Manual & automatic license plate registration
- 📹 XPOTS camera integration (webhooks)
- 📊 Entry/exit event history
- 👥 Visitor types (guest, visitor, staff, delivery, contractor, VIP)
- ⏱️ Temporary access with expiration
- 🔗 Mews PMS integration ready

## API Endpoints

```
POST   /api/licenseplate/scan              - Register plate
GET    /api/licenseplate/records           - List all
GET    /api/licenseplate/records/:plate    - Get specific
DELETE /api/licenseplate/records/:plate    - Remove
GET    /api/licenseplate/records/:plate/events - Event history
POST   /api/licenseplate/webhook/xpots     - XPOTS webhook
```

## XPOTS Setup

Configure camera system:
- URL: `http://your-server:8082/api/licenseplate/webhook/xpots`
- Header: `Authorization: Bearer your-webhook-key`

## License

MIT
