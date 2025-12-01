package models

import "time"

type LicensePlateRecord struct {
	PlateNumber      string    `json:"plate_number"`
	GuestName        string    `json:"guest_name"`
	RoomNumber       string    `json:"room_number,omitempty"`
	CheckIn          time.Time `json:"check_in"`
	CheckOut         time.Time `json:"check_out,omitempty"`
	VehicleMake      string    `json:"vehicle_make,omitempty"`
	VehicleModel     string    `json:"vehicle_model,omitempty"`
	Notes            string    `json:"notes,omitempty"`
	VisitorType      string    `json:"visitor_type"`           // guest, visitor, staff, delivery, contractor, vip
	AccessExpiresAt  time.Time `json:"access_expires_at,omitempty"` // When temporary access expires
	Purpose          string    `json:"purpose,omitempty"`      // Purpose of visit for non-guests
	CreatedAt        time.Time `json:"created_at"`
}

type ScanRequest struct {
	PlateNumber     string `json:"plate_number" binding:"required"`
	GuestName       string `json:"guest_name" binding:"required"`
	RoomNumber      string `json:"room_number"`
	VehicleMake     string `json:"vehicle_make"`
	VehicleModel    string `json:"vehicle_model"`
	Notes           string `json:"notes"`
	VisitorType     string `json:"visitor_type"`      // guest, visitor, staff, delivery, contractor, vip
	AccessExpiresAt string `json:"access_expires_at"` // ISO 8601 timestamp for expiration
	Purpose         string `json:"purpose"`           // Purpose of visit
}
