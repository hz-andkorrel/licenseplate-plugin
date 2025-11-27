package handlers

import (
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

	c.JSON(http.StatusCreated, gin.H{
		"message": "License plate scanned successfully",
		"record":  record,
	})
}

func (h *LicensePlateHandler) GetAllRecords(c *gin.Context) {
	guestName := c.Query("guest_name")
	
	if guestName != "" {
		records := h.service.SearchByGuestName(guestName)
		c.JSON(http.StatusOK, gin.H{
			"records": records,
			"count":   len(records),
		})
		return
	}

	records := h.service.GetAllRecords()
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
