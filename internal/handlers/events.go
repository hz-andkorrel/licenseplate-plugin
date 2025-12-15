package handlers

import (
    "context"
    "encoding/json"
    "fmt"
    "log"

    "licenseplate-plugin/internal/models"
    "licenseplate-plugin/internal/services"
)

// HandleLicenseplateScanned is a typed handler for licenseplate.scanned events.
// It unmarshals the payload and delegates to the service layer.
func HandleLicenseplateScanned(service *services.LicensePlateService, ctx context.Context, raw json.RawMessage) error {
    var payload models.XPOTSWebhookPayload
    log.Printf("[handlers] HandleLicenseplateScanned: received raw payload: %s", string(raw))
    if err := json.Unmarshal(raw, &payload); err != nil {
        return fmt.Errorf("invalid licenseplate.scanned payload: %w", err)
    }

    // Minimal validation
    if payload.PlateNumber == "" {
        return fmt.Errorf("missing plate number in payload")
    }

    log.Printf("[handlers] HandleLicenseplateScanned: processing plate=%s", payload.PlateNumber)

    // Delegate to existing service logic that already handles XPOTS payloads
    if err := service.ProcessXPOTSWebhook(&payload); err != nil {
        return fmt.Errorf("service.ProcessXPOTSWebhook failed: %w", err)
    }

    log.Printf("[handlers] HandleLicenseplateScanned: processed plate=%s successfully", payload.PlateNumber)
    return nil
}
