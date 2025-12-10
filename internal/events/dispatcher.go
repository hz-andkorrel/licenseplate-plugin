package events

import (
    "context"
    "encoding/json"
    "log"

    "licenseplate-plugin/internal/handlers"
    "licenseplate-plugin/internal/services"
)

// Event is the minimal wrapper expected on the bus
type Event struct {
    Type   string          `json:"type"`
    Record json.RawMessage `json:"record"`
}

// Dispatch parses a raw message and routes it to the correct handler.
// It runs handler calls asynchronously so the caller (listener) is not blocked.
func Dispatch(ctx context.Context, svc *services.LicensePlateService, rawMessage string) {
    var ev Event
    if err := json.Unmarshal([]byte(rawMessage), &ev); err != nil {
        log.Printf("[events] invalid event JSON: %v", err)
        return
    }

    switch ev.Type {
    case "licenseplate.scanned":
        go func() {
            if err := handlers.HandleLicenseplateScanned(svc, ctx, ev.Record); err != nil {
                log.Printf("[events] licenseplate.scanned handler error: %v", err)
            }
        }()
    default:
        log.Printf("[events] unknown event type: %s", ev.Type)
    }
}
