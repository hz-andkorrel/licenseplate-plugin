-- Migration: Add visitor support and access control
-- This adds visitor types, access expiration, and purpose tracking

ALTER TABLE license_plates 
ADD COLUMN visitor_type VARCHAR(50) DEFAULT 'guest' CHECK (visitor_type IN ('guest', 'visitor', 'staff', 'delivery', 'contractor', 'vip')),
ADD COLUMN access_expires_at TIMESTAMP,
ADD COLUMN purpose TEXT;

-- Create index for quick lookup of expired access
CREATE INDEX idx_access_expires_at ON license_plates(access_expires_at) WHERE access_expires_at IS NOT NULL;

-- Create index for visitor type filtering
CREATE INDEX idx_visitor_type ON license_plates(visitor_type);

COMMENT ON COLUMN license_plates.visitor_type IS 'Type of visitor: guest (hotel guest), visitor (temporary), staff, delivery, contractor, vip';
COMMENT ON COLUMN license_plates.access_expires_at IS 'When temporary access expires (NULL = no expiration)';
COMMENT ON COLUMN license_plates.purpose IS 'Purpose of visit for non-guests';
