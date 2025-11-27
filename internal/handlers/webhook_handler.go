package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"licenseplate-plugin/internal/models"
	"licenseplate-plugin/internal/services"

	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	service *services.LicensePlateService
	apiKey  string
}

func NewWebhookHandler(service *services.LicensePlateService) *WebhookHandler {
	apiKey := os.Getenv("WEBHOOK_API_KEY")
	if apiKey == "" {
		log.Println("WARNING: WEBHOOK_API_KEY not set - webhook endpoint will be insecure!")
		apiKey = "default-insecure-key"
	}
	
	return &WebhookHandler{
		service: service,
		apiKey:  apiKey,
	}
}

// HandleXPOTSWebhook receives license plate data from XPOTS system
func (h *WebhookHandler) HandleXPOTSWebhook(c *gin.Context) {
	// Validate API key from header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, models.WebhookResponse{
			Success: false,
			Message: "Missing Authorization header",
		})
		return
	}

	// Support both "Bearer TOKEN" and "TOKEN" formats
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token != h.apiKey {
		c.JSON(http.StatusUnauthorized, models.WebhookResponse{
			Success: false,
			Message: "Invalid API key",
		})
		return
	}

	// Parse XPOTS payload
	var payload models.XPOTSWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Printf("Failed to parse XPOTS webhook payload: %v", err)
		c.JSON(http.StatusBadRequest, models.WebhookResponse{
			Success: false,
			Message: fmt.Sprintf("Invalid payload format: %v", err),
		})
		return
	}

	// Log the incoming webhook
	log.Printf("Received XPOTS webhook - Event: %s, Plate: %s, Time: %s, Location: %s",
		payload.EventType, payload.PlateNumber, payload.Timestamp, payload.Location)

	// Process the webhook through service layer
	err := h.service.ProcessXPOTSWebhook(&payload)
	if err != nil {
		log.Printf("Error processing XPOTS webhook: %v", err)
		c.JSON(http.StatusInternalServerError, models.WebhookResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to process webhook: %v", err),
			Plate:   payload.PlateNumber,
		})
		return
	}

	// Send success response back to XPOTS
	c.JSON(http.StatusOK, models.WebhookResponse{
		Success: true,
		Message: fmt.Sprintf("Successfully processed %s event for plate %s", payload.EventType, payload.PlateNumber),
		Plate:   payload.PlateNumber,
	})
}

// GetWebhookInfo provides information about the webhook endpoint
func (h *WebhookHandler) GetWebhookInfo(c *gin.Context) {
	info := gin.H{
		"endpoint": "/api/licenseplate/webhook/xpots",
		"method":   "POST",
		"authentication": gin.H{
			"type":   "API Key",
			"header": "Authorization: Bearer YOUR_API_KEY",
		},
		"payload_example": models.XPOTSWebhookPayload{
			EventType:   "entry",
			PlateNumber: "ABC-123",
			Location:    "Main Gate",
			Confidence:  0.98,
			CameraID:    "CAM-001",
			Direction:   "in",
		},
		"response_example": models.WebhookResponse{
			Success: true,
			Message: "Successfully processed entry event for plate ABC-123",
			Plate:   "ABC-123",
		},
	}
	
	c.JSON(http.StatusOK, info)
}
