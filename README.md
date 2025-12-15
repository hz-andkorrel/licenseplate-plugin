# License Plate Plugin — Short Summary

What it does
- Receives and stores license plate scans and parking events (manual API or camera webhooks).
- Records events in the database and keeps an outbox for reliable publishing.

Communication
- Publishes JSON events to the Redis event bus channel `events` (example: `{"type":"licenseplate.scanned","record":{...}}`).
- Subscribes to the same `events` channel to receive messages from other services.

HTTP endpoints (important)
- `POST /api/licenseplate/scan`  — register a scanned plate
- `POST /api/licenseplate/webhook/xpots` — XPOTS camera webhook

Operational notes
- Requires a Postgres DB (migrations must create `outbox_events` table) and Redis reachable via `HUB_BUS_ADDR`.
- Env vars: `DATABASE_URL`, `HUB_BUS_ADDR` (default `hub_bus:6379`), `PORT`.
- The outbox publisher background task delivers DB-backed events to Redis for reliable delivery.

Quick run (development)
```powershell
cd licenseplate-plugin
$env:DATABASE_URL='postgres://user:pass@localhost:5432/dbname'
go run main.go
```

Run with Docker Compose (recommended for local testing)

We included a `docker-compose.yml` and `Dockerfile` to run the plugin together with Postgres and Redis for development.

```bash
cd licenseplate-plugin
docker-compose up --build
```

The service is exposed on port `9002` and reads configuration from `.env`. The compose file provides sensible defaults for `DATABASE_URL` and `HUB_BUS_ADDR` — update `.env` if you change credentials.
