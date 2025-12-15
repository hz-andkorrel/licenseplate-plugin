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