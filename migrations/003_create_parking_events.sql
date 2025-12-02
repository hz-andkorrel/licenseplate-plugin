-- Migration 003: Create parking events table for entry/exit history
-- This allows tracking multiple entries/exits instead of overwriting check_in/check_out

CREATE TABLE IF NOT EXISTS parking_events (
    id SERIAL PRIMARY KEY,
    plate_number VARCHAR(20) NOT NULL,
    event_type VARCHAR(10) NOT NULL CHECK (event_type IN ('entry', 'exit')),
    event_time TIMESTAMP NOT NULL DEFAULT NOW(),
    location VARCHAR(100),
    camera_id VARCHAR(50),
    confidence DECIMAL(3,2),
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for fast queries
CREATE INDEX idx_parking_events_plate ON parking_events(plate_number);
CREATE INDEX idx_parking_events_time ON parking_events(event_time DESC);
CREATE INDEX idx_parking_events_type ON parking_events(event_type);
CREATE INDEX idx_parking_events_plate_time ON parking_events(plate_number, event_time DESC);

-- Add comments for documentation
COMMENT ON TABLE parking_events IS 'Audit trail of all vehicle entry/exit events';
COMMENT ON COLUMN parking_events.event_type IS 'Type of event: entry or exit';
COMMENT ON COLUMN parking_events.location IS 'Gate/entrance where event occurred (from XPOTS)';
COMMENT ON COLUMN parking_events.camera_id IS 'Camera that detected the vehicle (from XPOTS)';
COMMENT ON COLUMN parking_events.confidence IS 'XPOTS confidence score (0.00-1.00)';

-- Add guest_id column to license_plates for future reservation integration
ALTER TABLE license_plates 
ADD COLUMN IF NOT EXISTS guest_id VARCHAR(50),
ADD COLUMN IF NOT EXISTS reservation_id VARCHAR(50);

-- Create indexes for guest linking
CREATE INDEX IF NOT EXISTS idx_license_plates_guest ON license_plates(guest_id);
CREATE INDEX IF NOT EXISTS idx_license_plates_reservation ON license_plates(reservation_id);

-- Add comments
COMMENT ON COLUMN license_plates.guest_id IS 'Reference to guest in booking system (e.g., Mews guest ID)';
COMMENT ON COLUMN license_plates.reservation_id IS 'Reference to reservation in booking system';
