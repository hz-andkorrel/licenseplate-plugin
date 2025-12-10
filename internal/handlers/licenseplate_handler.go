package handlers

import (
	"encoding/json"
	"licenseplate-plugin/internal/eventbus"
	"licenseplate-plugin/internal/models"
	"licenseplate-plugin/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LicensePlateHandler struct {
	service *services.LicensePlateService
}

func NewLicensePlateHandler(service *services.LicensePlateService) *LicensePlateHandler {
	return &LicensePlateHandler{
		service: service,
	}
}

func (h *LicensePlateHandler) ScanLicensePlate(c *gin.Context) {
	var req models.ScanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record, err := h.service.ScanAndStore(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Publish event to event bus asynchronously
	go func(rec interface{}) {
		payload := map[string]interface{}{
			"type":   "licenseplate.scanned",
			"record": rec,
		}
		if b, err := json.Marshal(payload); err == nil {
			// use request context so cancellation propagates
			_ = eventbus.Publish(c.Request.Context(), "events", string(b))
		}
	}(record)

	c.JSON(http.StatusCreated, gin.H{
		"message": "License plate scanned successfully",
		"record":  record,
	})
}

func (h *LicensePlateHandler) GetAllRecords(c *gin.Context) {
	// Parse query parameters for search and filters
	filters := services.SearchFilters{
		Search:      c.Query("search"),       // Search in plate or name
		VisitorType: c.Query("visitor_type"), // Filter by type
		DateFrom:    c.Query("date_from"),    // Filter from date
		DateTo:      c.Query("date_to"),      // Filter to date
	}

	records := h.service.GetAllRecords(filters)
	c.JSON(http.StatusOK, gin.H{
		"records": records,
		"count":   len(records),
	})
}

func (h *LicensePlateHandler) GetRecord(c *gin.Context) {
	plate := c.Param("plate")
	record, err := h.service.GetRecord(plate)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, record)
}

func (h *LicensePlateHandler) DeleteRecord(c *gin.Context) {
	plate := c.Param("plate")
	if err := h.service.DeleteRecord(plate); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
}

// GetParkingEvents retrieves the event history for a specific license plate
func (h *LicensePlateHandler) GetParkingEvents(c *gin.Context) {
	plate := c.Param("plate")
	events, err := h.service.GetParkingEvents(plate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve parking events"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"plate_number": plate,
		"events":       events,
		"count":        len(events),
	})
}
