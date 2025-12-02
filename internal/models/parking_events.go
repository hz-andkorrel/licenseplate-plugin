package models

import "time"

// ParkingEvent represents a single entry or exit event for a vehicle
type ParkingEvent struct {
	ID          int       `json:"id"`
	PlateNumber string    `json:"plate_number"`
	EventType   string    `json:"event_type"` // "entry" or "exit"
	EventTime   time.Time `json:"event_time"`
	Location    string    `json:"location,omitempty"`
	CameraID    string    `json:"camera_id,omitempty"`
	Confidence  float64   `json:"confidence,omitempty"`
	Notes       string    `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// GuestReservation represents guest booking system data (placeholder for future integration)
// This will be populated when you connect to Mews or another PMS
type GuestReservation struct {
	GuestID        string    `json:"guest_id"`
	ReservationID  string    `json:"reservation_id"`
	GuestName      string    `json:"guest_name"`
	RoomNumber     string    `json:"room_number,omitempty"`
	CheckInDate    time.Time `json:"check_in_date"`
	CheckOutDate   time.Time `json:"check_out_date"`
	LicensePlates  []string  `json:"license_plates,omitempty"`
	Email          string    `json:"email,omitempty"`
	Phone          string    `json:"phone,omitempty"`
}

// EventHistoryResponse contains a license plate record with its event history
type EventHistoryResponse struct {
	LicensePlateRecord
	Events []ParkingEvent `json:"events"`
}
