package services

import (
	"errors"
	"licenseplate-plugin/internal/models"
	"strings"
	"sync"
	"time"
)

type LicensePlateService struct {
	records map[string]*models.LicensePlateRecord
	mu      sync.RWMutex
}

func NewLicensePlateService() *LicensePlateService {
	return &LicensePlateService{
		records: make(map[string]*models.LicensePlateRecord),
	}
}

func (s *LicensePlateService) ScanAndStore(req models.ScanRequest) (*models.LicensePlateRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Normalize plate number (uppercase, remove spaces)
	plateNumber := strings.ToUpper(strings.ReplaceAll(req.PlateNumber, " ", ""))

	if plateNumber == "" {
		return nil, errors.New("plate number is required")
	}

	record := &models.LicensePlateRecord{
		PlateNumber:  plateNumber,
		GuestName:    req.GuestName,
		RoomNumber:   req.RoomNumber,
		CheckIn:      time.Now(),
		VehicleMake:  req.VehicleMake,
		VehicleModel: req.VehicleModel,
		Notes:        req.Notes,
		CreatedAt:    time.Now(),
	}

	s.records[plateNumber] = record
	return record, nil
}

func (s *LicensePlateService) GetAllRecords() []*models.LicensePlateRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()

	records := make([]*models.LicensePlateRecord, 0, len(s.records))
	for _, record := range s.records {
		records = append(records, record)
	}
	return records
}

func (s *LicensePlateService) GetRecord(plateNumber string) (*models.LicensePlateRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	plateNumber = strings.ToUpper(strings.ReplaceAll(plateNumber, " ", ""))
	record, exists := s.records[plateNumber]
	if !exists {
		return nil, errors.New("record not found")
	}
	return record, nil
}

func (s *LicensePlateService) DeleteRecord(plateNumber string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	plateNumber = strings.ToUpper(strings.ReplaceAll(plateNumber, " ", ""))
	if _, exists := s.records[plateNumber]; !exists {
		return errors.New("record not found")
	}

	delete(s.records, plateNumber)
	return nil
}

func (s *LicensePlateService) SearchByGuestName(guestName string) []*models.LicensePlateRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []*models.LicensePlateRecord
	searchTerm := strings.ToLower(guestName)

	for _, record := range s.records {
		if strings.Contains(strings.ToLower(record.GuestName), searchTerm) {
			results = append(results, record)
		}
	}
	return results
}
