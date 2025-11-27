package models

import "time"

// XPOTSWebhookPayload represents the data structure sent by XPOTS
// Adjust fields based on actual XPOTS API documentation
type XPOTSWebhookPayload struct {
	EventType   string    `json:"event_type"`   // e.g., "entry", "exit", "scan"
	PlateNumber string    `json:"plate_number"` // License plate detected
	Timestamp   time.Time `json:"timestamp"`    // When the plate was scanned
	Location    string    `json:"location"`     // Camera/gate location
	Confidence  float64   `json:"confidence"`   // Recognition confidence (0-1)
	ImageURL    string    `json:"image_url"`    // URL to plate image (if available)
	CameraID    string    `json:"camera_id"`    // ID of the camera that detected the plate
	
	// Additional fields that might be provided
	VehicleType  string `json:"vehicle_type"`  // car, motorcycle, truck, etc.
	Direction    string `json:"direction"`     // in, out
	LaneNumber   int    `json:"lane_number"`   // Which lane/gate
}

// WebhookResponse is sent back to XPOTS to acknowledge receipt
type WebhookResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Plate   string `json:"plate_number,omitempty"`
}
