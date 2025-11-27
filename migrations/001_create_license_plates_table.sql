-- Create license_plates table for storing vehicle records
CREATE TABLE IF NOT EXISTS license_plates (
    id SERIAL PRIMARY KEY,
    plate_number VARCHAR(20) NOT NULL UNIQUE,
    guest_name VARCHAR(255) NOT NULL,
    room_number VARCHAR(10),
    check_in TIMESTAMP NOT NULL,
    check_out TIMESTAMP,
    vehicle_make VARCHAR(100),
    vehicle_model VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_license_plates_plate_number ON license_plates(plate_number);
CREATE INDEX IF NOT EXISTS idx_license_plates_guest_name ON license_plates(guest_name);
CREATE INDEX IF NOT EXISTS idx_license_plates_room_number ON license_plates(room_number);
CREATE INDEX IF NOT EXISTS idx_license_plates_check_in ON license_plates(check_in);

-- Insert sample data (optional, for testing)
-- INSERT INTO license_plates (plate_number, guest_name, room_number, check_in, vehicle_make, vehicle_model, notes)
-- VALUES ('ABC123', 'John Doe', '101', NOW(), 'Toyota', 'Camry', 'Regular guest');
