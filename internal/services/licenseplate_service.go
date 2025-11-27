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

	checkIn := time.Now()
	query := `
		INSERT INTO license_plates (plate_number, guest_name, room_number, check_in, vehicle_make, vehicle_model, notes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (plate_number) 
		DO UPDATE SET guest_name = $2, room_number = $3, check_in = $4, vehicle_make = $5, vehicle_model = $6, notes = $7, updated_at = NOW()
		RETURNING id, created_at
	`

	var id int
	var createdAt time.Time
	row := s.db.QueryRow(query, plateNumber, req.GuestName, req.RoomNumber, checkIn, req.VehicleMake, req.VehicleModel, req.Notes, checkIn)
	
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
		CreatedAt:    createdAt,
	}

	return record, nil
}

func (s *LicensePlateService) GetAllRecords() []*models.LicensePlateRecord {
	query := `
		SELECT plate_number, guest_name, room_number, check_in, check_out, vehicle_make, vehicle_model, notes, created_at
		FROM license_plates
		ORDER BY created_at DESC
	`

	conn, err := s.db.GetConnection()
	if err != nil {
		log.Println("[LicensePlateService] Error connecting to database:", err)
		return []*models.LicensePlateRecord{}
	}
	defer conn.Close()

	rows, err := conn.Query(query)
	if err != nil {
		log.Println("[LicensePlateService] Error querying records:", err)
		return []*models.LicensePlateRecord{}
	}
	defer rows.Close()

	records := make([]*models.LicensePlateRecord, 0)
	for rows.Next() {
		var record models.LicensePlateRecord
		var checkOut sql.NullTime
		var roomNumber, vehicleMake, vehicleModel, notes sql.NullString

		err := rows.Scan(
			&record.PlateNumber,
			&record.GuestName,
			&roomNumber,
			&record.CheckIn,
			&checkOut,
			&vehicleMake,
			&vehicleModel,
			&notes,
			&record.CreatedAt,
		)
		if err != nil {
			log.Println("[LicensePlateService] Error scanning row:", err)
			continue
		}

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

		records = append(records, &record)
	}

	return records
}

func (s *LicensePlateService) GetRecord(plateNumber string) (*models.LicensePlateRecord, error) {
	plateNumber = strings.ToUpper(strings.ReplaceAll(plateNumber, " ", ""))

	query := `
		SELECT plate_number, guest_name, room_number, check_in, check_out, vehicle_make, vehicle_model, notes, created_at
		FROM license_plates
		WHERE plate_number = $1
	`

	var record models.LicensePlateRecord
	var checkOut sql.NullTime
	var roomNumber, vehicleMake, vehicleModel, notes sql.NullString

	row := s.db.QueryRow(query, plateNumber)
	err := row.Scan(
		&record.PlateNumber,
		&record.GuestName,
		&roomNumber,
		&record.CheckIn,
		&checkOut,
		&vehicleMake,
		&vehicleModel,
		&notes,
		&record.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("record not found")
	}
	if err != nil {
		log.Println("[LicensePlateService] Error querying record:", err)
		return nil, errors.New("failed to retrieve record")
	}

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

	return &record, nil
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

func (s *LicensePlateService) SearchByGuestName(guestName string) []*models.LicensePlateRecord {
	query := `
		SELECT plate_number, guest_name, room_number, check_in, check_out, vehicle_make, vehicle_model, notes, created_at
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
		var record models.LicensePlateRecord
		var checkOut sql.NullTime
		var roomNumber, vehicleMake, vehicleModel, notes sql.NullString

		err := rows.Scan(
			&record.PlateNumber,
			&record.GuestName,
			&roomNumber,
			&record.CheckIn,
			&checkOut,
			&vehicleMake,
			&vehicleModel,
			&notes,
			&record.CreatedAt,
		)
		if err != nil {
			log.Println("[LicensePlateService] Error scanning row:", err)
			continue
		}

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

		records = append(records, &record)
	}

	return records
}

// ProcessXPOTSWebhook handles incoming webhook data from XPOTS system
func (s *LicensePlateService) ProcessXPOTSWebhook(payload *models.XPOTSWebhookPayload) error {
	// Normalize plate number
	plateNumber := strings.ToUpper(strings.ReplaceAll(payload.PlateNumber, " ", ""))
	
	if plateNumber == "" {
		return errors.New("plate number is required")
	}

	// Check if record exists
	existingRecord, _ := s.GetRecord(plateNumber)

	// Handle different event types
	switch payload.EventType {
	case "entry", "scan", "in":
		// Vehicle entering - create or update check-in
		if existingRecord != nil {
			// Update existing record with new check-in time
			query := `
				UPDATE license_plates 
				SET check_in = $1, updated_at = NOW(), notes = $2
				WHERE plate_number = $3
			`
			notes := fmt.Sprintf("Auto-detected by XPOTS at %s (confidence: %.2f%%)", 
				payload.Location, payload.Confidence*100)
			
			_, err := s.db.Execute(query, payload.Timestamp, notes, plateNumber)
			if err != nil {
				log.Printf("[LicensePlateService] Error updating check-in for %s: %v", plateNumber, err)
				return err
			}
			log.Printf("Updated check-in time for plate %s", plateNumber)
		} else {
			// Create new record for unknown vehicle
			query := `
				INSERT INTO license_plates (plate_number, guest_name, check_in, notes, created_at)
				VALUES ($1, $2, $3, $4, $5)
			`
			guestName := "Unknown Guest (Auto-detected)"
			notes := fmt.Sprintf("Auto-detected by XPOTS at %s (confidence: %.2f%%, camera: %s)", 
				payload.Location, payload.Confidence*100, payload.CameraID)
			
			_, err := s.db.Execute(query, plateNumber, guestName, payload.Timestamp, notes, time.Now())
			if err != nil {
				log.Printf("[LicensePlateService] Error creating record for %s: %v", plateNumber, err)
				return err
			}
			log.Printf("Created new record for unknown plate %s", plateNumber)
		}

	case "exit", "out":
		// Vehicle exiting - update check-out time
		if existingRecord != nil {
			query := `
				UPDATE license_plates 
				SET check_out = $1, updated_at = NOW()
				WHERE plate_number = $2
			`
			_, err := s.db.Execute(query, payload.Timestamp, plateNumber)
			if err != nil {
				log.Printf("[LicensePlateService] Error updating check-out for %s: %v", plateNumber, err)
				return err
			}
			log.Printf("Updated check-out time for plate %s", plateNumber)
		} else {
			// Unknown vehicle leaving - still log it
			query := `
				INSERT INTO license_plates (plate_number, guest_name, check_out, notes, created_at)
				VALUES ($1, $2, $3, $4, $5)
			`
			guestName := "Unknown Guest (Exit only)"
			notes := fmt.Sprintf("Exit detected by XPOTS at %s (no entry record found)", payload.Location)
			
			_, err := s.db.Execute(query, plateNumber, guestName, payload.Timestamp, notes, time.Now())
			if err != nil {
				log.Printf("[LicensePlateService] Error creating exit record for %s: %v", plateNumber, err)
				return err
			}
			log.Printf("Created exit record for unknown plate %s", plateNumber)
		}

	default:
		log.Printf("Unknown event type '%s' for plate %s - treating as entry", payload.EventType, plateNumber)
		// Treat unknown events as entry by default
		if existingRecord == nil {
			query := `
				INSERT INTO license_plates (plate_number, guest_name, check_in, notes, created_at)
				VALUES ($1, $2, $3, $4, $5)
			`
			guestName := "Unknown Guest (Auto-detected)"
			notes := fmt.Sprintf("Event: %s, Location: %s, Camera: %s", 
				payload.EventType, payload.Location, payload.CameraID)
			
			_, err := s.db.Execute(query, plateNumber, guestName, payload.Timestamp, notes, time.Now())
			return err
		}
	}

	return nil
}
