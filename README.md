# License Plate Plugin (short)

This service provides an API for registering and managing license plate scans and related parking events. The plugin is API-only (no built-in frontend).

Key responsibilities
- Receive and store license plate scans (manual or via XPOTS webhooks)
- Provide endpoints to list, query and delete records
- Publish scan events to a central event bus so other services can react
- Subscribe to the same event bus for incoming events (currently logs incoming events)

Quick start (development)
1. Start required services (database, redis) e.g. via your docker-compose for the workspace.
2. Ensure `DATABASE_URL` and other env vars are set.
3. Run the plugin:

```powershell
# from repository root
cd 'C:\Users\larss\OneDrive\Documenten\ICT - Year 4\Cadzand hotel project\licenseplate-plugin'
$env:DATABASE_URL='postgres://user:pass@localhost:5432/dbname'
go run main.go
```

API (important endpoints)
- `POST /api/licenseplate/scan` - register a scanned plate
- `GET  /api/licenseplate/records` - list records
- `GET  /api/licenseplate/records/:plate` - get details
- `DELETE /api/licenseplate/records/:plate` - remove
- `GET  /api/licenseplate/records/:plate/events` - event history
- `POST /api/licenseplate/webhook/xpots` - XPOTS camera webhook

Event bus (what changed)
- Redis Pub/Sub is used as a simple event bus. The plugin publishes JSON events to the `events` channel and also subscribes to that channel on startup.
- Published event example (JSON):
# License Plate Plugin

API-only service for license plate scans and parking events.

Quick run (PowerShell):
```powershell
cd 'C:\Users\larss\OneDrive\Documenten\ICT - Year 4\Cadzand hotel project\licenseplate-plugin'
$env:DATABASE_URL='postgres://user:pass@localhost:5432/dbname'
go run main.go
```

Important env vars:
- `DATABASE_URL` - Postgres connection
- `HUB_BUS_ADDR` - Redis address (default `hub_bus:6379`)
- `PORT`, `HOST` - HTTP server

Event bus:
- Publishes JSON events to channel `events`, e.g. `{"type":"licenseplate.scanned","record":{...}}`
- Test tooling: `cmd/publish_test` and `cmd/listen_test`

Endpoints (summary): `POST /api/licenseplate/scan`, `GET /api/licenseplate/records`, `GET /api/licenseplate/records/:plate`, `DELETE /api/licenseplate/records/:plate`, `POST /api/licenseplate/webhook/xpots`

Notes: listener currently logs events; consider adding an event dispatcher for handling specific event types.

License: MIT
