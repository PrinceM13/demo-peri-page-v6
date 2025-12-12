package config

import (
	"fmt"
	"os"
	"time"
)

// Config holds all application configuration.
type Config struct {
	Server  ServerConfig
	Printer PrinterConfig
	BLE     BLEConfig
}

// ServerConfig holds server-specific configuration.
type ServerConfig struct {
	Port string
}

// PrinterConfig holds printer-specific configuration.
type PrinterConfig struct {
	Type       string // "mock" or "ble"
	DeviceName string
	Timeout    time.Duration
}

// BLEConfig holds Bluetooth LE configuration.
type BLEConfig struct {
	ScanTimeout time.Duration
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
		Printer: PrinterConfig{
			Type:       getEnv("PRINTER_TYPE", "mock"),
			DeviceName: getEnv("PRINTER_DEVICE_NAME", "Peripage"),
			Timeout:    parseDuration(getEnv("PRINTER_TIMEOUT", "30s")),
		},
		BLE: BLEConfig{
			ScanTimeout: parseDuration(getEnv("BLE_SCAN_TIMEOUT", "10s")),
		},
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.Printer.Type != "mock" && c.Printer.Type != "ble" {
		return fmt.Errorf("invalid printer type: %s (must be 'mock' or 'ble')", c.Printer.Type)
	}

	if c.Printer.Type == "ble" && c.Printer.DeviceName == "" {
		return fmt.Errorf("device name is required for BLE printer")
	}

	return nil
}

// getEnv gets an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseDuration parses a duration string, returning 0 on error.
func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0
	}
	return d
}
