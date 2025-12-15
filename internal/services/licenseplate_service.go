package services

import (
	"database/sql"
	"errors"
	"fmt"
	"licenseplate-plugin/internal/database"
	"licenseplate-plugin/internal/models"
	"log"
	"strings"
	"time"
)

type LicensePlateService struct {
	db *database.Database
}

func NewLicensePlateService(db *database.Database) *LicensePlateService {
	return &LicensePlateService{
		db: db,
	}
}

func (s *LicensePlateService) ScanAndStore(req models.ScanRequest) (*models.LicensePlateRecord, error) {
	// Normalize plate number (uppercase, remove spaces)
	plateNumber := strings.ToUpper(strings.ReplaceAll(req.PlateNumber, " ", ""))

	if plateNumber == "" {
		return nil, errors.New("plate number is required")
	}

	// Default visitor type to "guest" if not specified
	visitorType := req.VisitorType
	if visitorType == "" {
		visitorType = "guest"
	}

	// Validate visitor type
	validTypes := map[string]bool{"guest": true, "visitor": true, "staff": true, "delivery": true, "contractor": true, "vip": true}
	if !validTypes[visitorType] {
		return nil, errors.New("invalid visitor type")
	}

	// Parse expiration time if provided
	var expiresAt sql.NullTime
	if req.AccessExpiresAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, req.AccessExpiresAt)
		if err != nil {
			return nil, errors.New("invalid access_expires_at format, use ISO 8601")
		}
		expiresAt = sql.NullTime{Time: parsedTime, Valid: true}
	}

	checkIn := time.Now()
	query := `
		INSERT INTO license_plates (plate_number, guest_name, room_number, check_in, vehicle_make, vehicle_model, notes, visitor_type, access_expires_at, purpose, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (plate_number) 
		DO UPDATE SET guest_name = $2, room_number = $3, check_in = $4, vehicle_make = $5, vehicle_model = $6, notes = $7, visitor_type = $8, access_expires_at = $9, purpose = $10, updated_at = NOW()
		RETURNING id, created_at
	`

	var id int
	var createdAt time.Time
	row := s.db.QueryRow(query, plateNumber, req.GuestName, req.RoomNumber, checkIn, req.VehicleMake, req.VehicleModel, req.Notes, visitorType, expiresAt, req.Purpose, checkIn)
	
	if err := row.Scan(&id, &createdAt); err != nil {
		log.Println("[LicensePlateService] Error inserting/updating record:", err)
		return nil, errors.New("failed to store license plate record")
	}

	record := &models.LicensePlateRecord{
		PlateNumber:  plateNumber,
		GuestName:    req.GuestName,
		RoomNumber:   req.RoomNumber,
		CheckIn:      checkIn,
		VehicleMake:  req.VehicleMake,
		VehicleModel: req.VehicleModel,
		Notes:        req.Notes,
		VisitorType:  visitorType,
		Purpose:      req.Purpose,
		CreatedAt:    createdAt,
	}
	
	if expiresAt.Valid {
		record.AccessExpiresAt = expiresAt.Time
	}

	return record, nil
}

// scanLicensePlateRecord is a helper function to reduce duplicate code
func scanLicensePlateRecord(scanner interface {
	Scan(dest ...interface{}) error
}) (*models.LicensePlateRecord, error) {
	record := &models.LicensePlateRecord{}
	var checkOut, expiresAt sql.NullTime
	var roomNumber, vehicleMake, vehicleModel, notes, purpose sql.NullString

	err := scanner.Scan(
		&record.PlateNumber,
		&record.GuestName,
		&roomNumber,
		&record.CheckIn,
		&checkOut,
		&vehicleMake,
		&vehicleModel,
		&notes,
		&record.VisitorType,
		&expiresAt,
		&purpose,
		&record.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Handle nullable fields
	if checkOut.Valid {
		record.CheckOut = checkOut.Time
	}
	if roomNumber.Valid {
		record.RoomNumber = roomNumber.String
	}
	if vehicleMake.Valid {
		record.VehicleMake = vehicleMake.String
	}
	if vehicleModel.Valid {
		record.VehicleModel = vehicleModel.String
	}
	if notes.Valid {
		record.Notes = notes.String
	}
	if expiresAt.Valid {
		record.AccessExpiresAt = expiresAt.Time
	}
	if purpose.Valid {
		record.Purpose = purpose.String
	}

	return record, nil
}

// LogParkingEvent creates an entry/exit event record
func (s *LicensePlateService) LogParkingEvent(plateNumber, eventType string, location, cameraID string, confidence float64, notes string) error {
	query := `
		INSERT INTO parking_events (plate_number, event_type, event_time, location, camera_id, confidence, notes)
		VALUES ($1, $2, NOW(), $3, $4, $5, $6)
	`
	
	_, err := s.db.Execute(query, plateNumber, eventType, location, cameraID, confidence, notes)
	if err != nil {
		log.Printf("[LicensePlateService] Error logging parking event: %v", err)
		return err
	}
	
	log.Printf("Logged %s event for plate %s", eventType, plateNumber)
	return nil
}

// GetParkingEvents retrieves all events for a specific license plate
func (s *LicensePlateService) GetParkingEvents(plateNumber string) ([]models.ParkingEvent, error) {
	plateNumber = strings.ToUpper(strings.ReplaceAll(plateNumber, " ", ""))
	
	query := `
		SELECT id, plate_number, event_type, event_time, location, camera_id, confidence, notes, created_at
		FROM parking_events
		WHERE plate_number = $1
		ORDER BY event_time DESC
	`
	
	conn, err := s.db.GetConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	
	rows, err := conn.Query(query, plateNumber)
	if err != nil {
		log.Printf("[LicensePlateService] Error querying parking events: %v", err)
		return nil, err
	}
	defer rows.Close()
	
	events := make([]models.ParkingEvent, 0)
	for rows.Next() {
		var event models.ParkingEvent
		var location, cameraID, notes sql.NullString
		var confidence sql.NullFloat64
		
		err := rows.Scan(
			&event.ID,
			&event.PlateNumber,
			&event.EventType,
			&event.EventTime,
			&location,
			&cameraID,
			&confidence,
			&notes,
			&event.CreatedAt,
		)
		if err != nil {
			log.Printf("[LicensePlateService] Error scanning event row: %v", err)
			continue
		}
		
		if location.Valid {
			event.Location = location.String
		}
		if cameraID.Valid {
			event.CameraID = cameraID.String
		}
		if confidence.Valid {
			event.Confidence = confidence.Float64
		}
		if notes.Valid {
			event.Notes = notes.String
		}
		
		events = append(events, event)
	}
	
	return events, nil
}

// SearchFilters contains all search and filter parameters
type SearchFilters struct {
	Search       string // Search in plate_number or guest_name
	VisitorType  string // Filter by visitor type
	DateFrom     string // Filter by check_in >= date
	DateTo       string // Filter by check_in <= date
}

func (s *LicensePlateService) GetAllRecords(filters SearchFilters) []*models.LicensePlateRecord {
	// Build dynamic query based on filters
	query := `
		SELECT plate_number, guest_name, room_number, check_in, check_out, vehicle_make, vehicle_model, notes, visitor_type, access_expires_at, purpose, created_at
		FROM license_plates
		WHERE 1=1
	`
	
	args := make([]interface{}, 0)
	argIndex := 1
	
	// Add search filter (plate number or guest name)
	if filters.Search != "" {
		query += fmt.Sprintf(" AND (UPPER(plate_number) LIKE $%d OR UPPER(guest_name) LIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+strings.ToUpper(filters.Search)+"%")
		argIndex++
	}
	
	// Add visitor type filter
	if filters.VisitorType != "" {
		query += fmt.Sprintf(" AND visitor_type = $%d", argIndex)
		args = append(args, filters.VisitorType)
		argIndex++
	}
	
	// Add date from filter
	if filters.DateFrom != "" {
		query += fmt.Sprintf(" AND check_in >= $%d", argIndex)
		args = append(args, filters.DateFrom)
		argIndex++
	}
	
	// Add date to filter
	if filters.DateTo != "" {
		query += fmt.Sprintf(" AND check_in <= $%d", argIndex)
		args = append(args, filters.DateTo)
		argIndex++
	}
	
	query += " ORDER BY created_at DESC"

	conn, err := s.db.GetConnection()
	if err != nil {
		log.Println("[LicensePlateService] Error connecting to database:", err)
		return []*models.LicensePlateRecord{}
	}
	defer conn.Close()

	rows, err := conn.Query(query, args...)
	if err != nil {
		log.Println("[LicensePlateService] Error querying records:", err)
		return []*models.LicensePlateRecord{}
	}
	defer rows.Close()

	records := make([]*models.LicensePlateRecord, 0)
	for rows.Next() {
		record, err := scanLicensePlateRecord(rows)
		if err != nil {
			log.Println("[LicensePlateService] Error scanning row:", err)
			continue
		}
		records = append(records, record)
	}

	return records
}

func (s *LicensePlateService) GetRecord(plateNumber string) (*models.LicensePlateRecord, error) {
	plateNumber = strings.ToUpper(strings.ReplaceAll(plateNumber, " ", ""))

	query := `
		SELECT plate_number, guest_name, room_number, check_in, check_out, vehicle_make, vehicle_model, notes, visitor_type, access_expires_at, purpose, created_at
		FROM license_plates
		WHERE plate_number = $1
	`

	row := s.db.QueryRow(query, plateNumber)
	record, err := scanLicensePlateRecord(row)

	if err == sql.ErrNoRows {
		return nil, errors.New("record not found")
	}
	if err != nil {
		log.Println("[LicensePlateService] Error querying record:", err)
		return nil, errors.New("failed to retrieve record")
	}

	return record, nil
}

func (s *LicensePlateService) DeleteRecord(plateNumber string) error {
	plateNumber = strings.ToUpper(strings.ReplaceAll(plateNumber, " ", ""))

	query := `DELETE FROM license_plates WHERE plate_number = $1`

	rowsAffected, err := s.db.Execute(query, plateNumber)
	if err != nil {
		log.Println("[LicensePlateService] Error deleting record:", err)
		return errors.New("failed to delete record")
	}

	if rowsAffected == 0 {
		return errors.New("record not found")
	}

	return nil
}

// Outbox: insert an event to be published reliably
func (s *LicensePlateService) InsertOutboxEvent(channel, payload string) (int64, error) {
	query := `
		INSERT INTO outbox_events (channel, payload, attempts, created_at)
		VALUES ($1, $2, 0, NOW())
		RETURNING id
	`
	conn, err := s.db.GetConnection()
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	var id int64
	err = conn.QueryRow(query, channel, payload).Scan(&id)
	if err != nil {
		log.Printf("[LicensePlateService] InsertOutboxEvent error: %v", err)
		return 0, err
	}
	return id, nil
}

// FetchPendingOutboxEvents fetches unsent outbox events up to limit
func (s *LicensePlateService) FetchPendingOutboxEvents(limit int) ([]models.OutboxEvent, error) {
	query := `
		SELECT id, channel, payload, attempts, last_error, created_at, sent_at
		FROM outbox_events
		WHERE sent_at IS NULL
		ORDER BY created_at ASC
		LIMIT $1
	`
	conn, err := s.db.GetConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []models.OutboxEvent{}
	for rows.Next() {
		var e models.OutboxEvent
		var lastError sql.NullString
		var sentAt sql.NullTime
		if err := rows.Scan(&e.ID, &e.Channel, &e.Payload, &e.Attempts, &lastError, &e.CreatedAt, &sentAt); err != nil {
			log.Printf("[LicensePlateService] FetchPendingOutboxEvents scan error: %v", err)
			continue
		}
		if lastError.Valid {
			e.LastError = lastError.String
		}
		if sentAt.Valid {
			t := sentAt.Time
			e.SentAt = &t
		}
		events = append(events, e)
	}
	return events, nil
}

// MarkOutboxSent marks the given outbox event as sent (sets sent_at)
func (s *LicensePlateService) MarkOutboxSent(id int64) error {
	query := `UPDATE outbox_events SET sent_at = NOW() WHERE id = $1`
	_, err := s.db.Execute(query, id)
	return err
}

// IncrementOutboxAttempts increments attempts and sets last_error
func (s *LicensePlateService) IncrementOutboxAttempts(id int64, errMsg string) error {
	query := `UPDATE outbox_events SET attempts = attempts + 1, last_error = $2 WHERE id = $1`
	_, err := s.db.Execute(query, id, errMsg)
	return err
}


func (s *LicensePlateService) SearchByGuestName(guestName string) []*models.LicensePlateRecord {
	query := `
		SELECT plate_number, guest_name, room_number, check_in, check_out, vehicle_make, vehicle_model, notes, visitor_type, access_expires_at, purpose, created_at
		FROM license_plates
		WHERE LOWER(guest_name) LIKE LOWER($1)
		ORDER BY created_at DESC
	`

	conn, err := s.db.GetConnection()
	if err != nil {
		log.Println("[LicensePlateService] Error connecting to database:", err)
		return []*models.LicensePlateRecord{}
	}
	defer conn.Close()

	searchTerm := "%" + guestName + "%"
	rows, err := conn.Query(query, searchTerm)
	if err != nil {
		log.Println("[LicensePlateService] Error searching records:", err)
		return []*models.LicensePlateRecord{}
	}
	defer rows.Close()

	records := make([]*models.LicensePlateRecord, 0)
	for rows.Next() {
		record, err := scanLicensePlateRecord(rows)
		if err != nil {
			log.Println("[LicensePlateService] Error scanning row:", err)
			continue
		}
		records = append(records, record)
	}

	return records
}

// ProcessXPOTSWebhook handles incoming webhook data from XPOTS system
// Now logs events in parking_events table instead of overwriting check_in/check_out
func (s *LicensePlateService) ProcessXPOTSWebhook(payload *models.XPOTSWebhookPayload) error {
	// Normalize plate number
	plateNumber := strings.ToUpper(strings.ReplaceAll(payload.PlateNumber, " ", ""))
	
	if plateNumber == "" {
		return errors.New("plate number is required")
	}

	// Determine event type
	eventType := "entry"
	switch payload.EventType {
	case "exit", "out":
		eventType = "exit"
	case "entry", "scan", "in":
		eventType = "entry"
	default:
		log.Printf("Unknown event type '%s' for plate %s - treating as entry", payload.EventType, plateNumber)
	}

	// Log the parking event
	notes := fmt.Sprintf("Auto-detected by XPOTS (confidence: %.2f%%)", payload.Confidence*100)
	err := s.LogParkingEvent(plateNumber, eventType, payload.Location, payload.CameraID, payload.Confidence, notes)
	if err != nil {
		return err
	}

	// Check if vehicle is registered in license_plates table
	existingRecord, _ := s.GetRecord(plateNumber)
	
	if existingRecord == nil {
		// Unknown vehicle - create a record for tracking
		query := `
			INSERT INTO license_plates (plate_number, guest_name, check_in, notes, visitor_type, created_at)
			VALUES ($1, $2, $3, $4, $5, NOW())
		`
		guestName := "Unknown Guest (Auto-detected)"
		notes := fmt.Sprintf("First detected at %s by camera %s", payload.Location, payload.CameraID)
		
		_, err := s.db.Execute(query, plateNumber, guestName, payload.Timestamp, notes, "visitor")
		if err != nil {
			log.Printf("[LicensePlateService] Error creating record for unknown vehicle %s: %v", plateNumber, err)
		}
	}

	return nil
}
