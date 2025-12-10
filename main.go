package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"licenseplate-plugin/internal/broker"
	"licenseplate-plugin/internal/database"
	"licenseplate-plugin/internal/handlers"
	evt "licenseplate-plugin/internal/events"
	"licenseplate-plugin/internal/services"
	"licenseplate-plugin/internal/eventbus"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	port := getEnv("PORT", "9002")
	host := getEnv("HOST", "localhost")
	baseAPIRoute := getEnv("BASE_API_ROUTE", "/api/licenseplate")
	databaseURL := getEnv("DATABASE_URL", "")

	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Initialize database
	db := database.NewDatabase(databaseURL)

	// Test database connection
	if conn, err := db.GetConnection(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	} else {
		conn.Close()
		log.Println("Successfully connected to database")
	}

	// Initialize services
	licensePlateService := services.NewLicensePlateService(db)

	// Register with broker
	go broker.RegisterWithBroker()

	// Start event listener (subscribes to 'events' channel)
	go func() {
		ctx := context.Background()
		_ = eventbus.Listen(ctx, "events", func(channel, message string) {
			// Dispatch incoming events to typed handlers
			// use the licensePlateService to let handlers call service-layer logic
			// non-blocking: dispatcher will start handlers asynchronously
			imported := licensePlateService
			// call into the events dispatcher
			// note: package imported below to avoid unused import until file edited
			eventDispatchWrapper(ctx, imported, message)
		})
	}()

	// Setup Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "licenseplate-plugin"})
	})

	// Initialize handlers
	handler := handlers.NewLicensePlateHandler(licensePlateService)
	webhookHandler := handlers.NewWebhookHandler(licensePlateService)

	// Register routes
	api := router.Group(baseAPIRoute)
	{
		api.POST("/scan", handler.ScanLicensePlate)
		api.GET("/records", handler.GetAllRecords)
		api.GET("/records/:plate", handler.GetRecord)
		api.GET("/records/:plate/events", handler.GetParkingEvents)
		api.DELETE("/records/:plate", handler.DeleteRecord)
		
		// Webhook endpoints
		api.POST("/webhook/xpots", webhookHandler.HandleXPOTSWebhook)
		api.GET("/webhook/info", webhookHandler.GetWebhookInfo)
	}

	// Start server
	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("License Plate Recognition Plugin running on http://%s%s", addr, baseAPIRoute)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// eventDispatchWrapper is a tiny indirection so we can call the dispatcher
// without adding complex logic inline in the listener callback.
func eventDispatchWrapper(ctx context.Context, svc *services.LicensePlateService, message string) {
	evt.Dispatch(ctx, svc, message)
}
