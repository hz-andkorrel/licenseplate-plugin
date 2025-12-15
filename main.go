package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"licenseplate-plugin/internal/broker"
	"licenseplate-plugin/internal/database"
	"licenseplate-plugin/internal/handlers"
	evt "licenseplate-plugin/internal/events"
	"licenseplate-plugin/internal/services"
	"licenseplate-plugin/internal/eventbus"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

// Global Redis client for event bus
var redisClient *redis.Client

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

	// Initialize Redis event bus
	initRedis()

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
		_ = eventbus.Listen(ctx, redisClient, "events", func(channel, message string) {
			// Dispatch incoming events to typed handlers
			// use the licensePlateService to let handlers call service-layer logic
			// non-blocking: dispatcher will start handlers asynchronously
			imported := licensePlateService
			// call into the events dispatcher
			// note: package imported below to avoid unused import until file edited
			eventDispatchWrapper(ctx, imported, message)
		})
	}()

	// Start background outbox publisher to reliably deliver events from DB to Redis
	ctx := context.Background()
	// run every 10s, process up to 50 events per tick
	startOutboxPublisher(ctx, licensePlateService, redisClient, 10*time.Second, 50)

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
	handler := handlers.NewLicensePlateHandler(licensePlateService, redisClient)
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

// initRedis initializes the global Redis client for the event bus
func initRedis() {
	redisAddr := getEnv("HUB_BUS_ADDR", "localhost:6379")

	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	// Test connection
	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis at %s: %v", redisAddr, err)
		log.Println("Plugin will continue without event bus functionality")
	} else {
		log.Printf("✓ Connected to Redis event bus at %s", redisAddr)
	}
}

// eventDispatchWrapper is a tiny indirection so we can call the dispatcher
// without adding complex logic inline in the listener callback.
func eventDispatchWrapper(ctx context.Context, svc *services.LicensePlateService, message string) {
	evt.Dispatch(ctx, svc, message)
}

// startOutboxPublisher runs a background goroutine that periodically
// fetches pending outbox events from the database and publishes them
// to the Redis event bus. On success the event is marked as sent; on
// failure the attempts counter is incremented and the error recorded.
func startOutboxPublisher(ctx context.Context, svc *services.LicensePlateService, client *redis.Client, interval time.Duration, batchSize int) {
	if client == nil {
		log.Println("[OutboxPublisher] redis client is nil; outbox publisher disabled")
		return
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("[OutboxPublisher] context canceled, stopping publisher")
				return
			case <-ticker.C:
				events, err := svc.FetchPendingOutboxEvents(batchSize)
				if err != nil {
					log.Printf("[OutboxPublisher] fetch error: %v", err)
					continue
				}

				for _, e := range events {
					// attempt publish
					if err := eventbus.Publish(ctx, client, e.Channel, e.Payload); err != nil {
						log.Printf("[OutboxPublisher] publish failed id=%d: %v", e.ID, err)
						_ = svc.IncrementOutboxAttempts(e.ID, err.Error())
						continue
					}

					// mark as sent
					if err := svc.MarkOutboxSent(e.ID); err != nil {
						log.Printf("[OutboxPublisher] mark sent failed id=%d: %v", e.ID, err)
						_ = svc.IncrementOutboxAttempts(e.ID, "mark sent failed: "+err.Error())
						continue
					}

					log.Printf("[OutboxPublisher] published and marked sent id=%d channel=%s", e.ID, e.Channel)
				}
			}
		}
	}()
}
