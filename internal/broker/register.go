package broker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type PluginRegistration struct {
	Slug          string   `json:"slug"`
	Name          string   `json:"name"`
	Version       string   `json:"version"`
	Description   string   `json:"description"`
	Host          string   `json:"host"`
	BaseAPIRoute  string   `json:"base-api-route"`
	SettingsRoute string   `json:"settings-route,omitempty"`
	APIRoutes     []string `json:"api-routes,omitempty"`
	Enabled       bool     `json:"enabled"`
}

func RegisterWithBroker() {
	brokerURL := os.Getenv("BROKER_URL")
	if brokerURL == "" {
		brokerURL = "http://localhost:8081"
	}

	authToken := os.Getenv("BROKER_AUTH_TOKEN")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	baseAPIRoute := os.Getenv("BASE_API_ROUTE")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "8082"
	}
	if baseAPIRoute == "" {
		baseAPIRoute = "/api/licenseplate"
	}

	pluginHost := fmt.Sprintf("http://%s:%s", host, port)

	registration := PluginRegistration{
		Slug:         "licenseplate-recognition",
		Name:         "License Plate Recognition",
		Version:      "1.0.0",
		Description:  "License plate recognition and management for hotel guests",
		Host:         pluginHost,
		BaseAPIRoute: baseAPIRoute,
		APIRoutes: []string{
			"/scan",
			"/records",
			"/records/:plate",
		},
		Enabled: true,
	}

	// Retry logic
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			time.Sleep(time.Duration(i*2) * time.Second)
			log.Printf("Retrying broker registration (attempt %d/%d)...", i+1, maxRetries)
		}

		if err := registerPlugin(brokerURL, authToken, registration); err != nil {
			log.Printf("Failed to register with broker: %v", err)
			continue
		}

		log.Println("Successfully registered with broker")
		return
	}

	log.Println("Failed to register with broker after all retries. Plugin will run without broker integration.")
}

func registerPlugin(brokerURL, authToken string, registration PluginRegistration) error {
	data, err := json.Marshal(registration)
	if err != nil {
		return fmt.Errorf("failed to marshal registration: %w", err)
	}

	req, err := http.NewRequest("POST", brokerURL+"/api/v1/route", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("broker returned status %d", resp.StatusCode)
	}

	return nil
}
