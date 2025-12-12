// Package main provides the entry point for the Peripage printer server.
//
// @title Peripage Printer API
// @version 1.0
// @description API for printing text and JSON data to a Peripage A6 thermal printer
// @termsOfService http://swagger.io/terms/
//
// @contact.name API Support
// @contact.email support@example.com
//
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
//
// @host localhost:8080
// @BasePath /
//
// @schemes http https
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/princem/peripage-printer/internal/adapters/api"
	_ "github.com/princem/peripage-printer/internal/adapters/docs"
	"github.com/princem/peripage-printer/internal/adapters/printer"
	"github.com/princem/peripage-printer/internal/config"
	"github.com/princem/peripage-printer/internal/core"
)

func main() {
	// Setup logger
	logger := log.New(os.Stdout, "[peripage] ", log.LstdFlags|log.Lshortfile)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	logger.Printf("Starting Peripage Printer Server")
	logger.Printf("Printer type: %s", cfg.Printer.Type)

	// Initialize the appropriate printer adapter
	var printerAdapter core.Printer
	var cleanup func()

	switch cfg.Printer.Type {
	case "mock":
		logger.Println("Using mock printer adapter")
		printerAdapter = printer.NewMockPrinter(logger)
		cleanup = func() {}

	case "ble":
		logger.Println("Using BLE printer adapter")
		blePrinter, err := printer.NewBLEPrinter(printer.BLEPrinterConfig{
			DeviceName:  cfg.Printer.DeviceName,
			ScanTimeout: cfg.BLE.ScanTimeout,
			Logger:      logger,
		})
		if err != nil {
			logger.Fatalf("Failed to initialize BLE printer: %v", err)
		}

		// Connect to printer
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Printer.Timeout)
		defer cancel()

		if err := blePrinter.Connect(ctx); err != nil {
			logger.Fatalf("Failed to connect to printer: %v", err)
		}

		printerAdapter = blePrinter
		cleanup = func() {
			logger.Println("Disconnecting from printer...")
			if err := blePrinter.Disconnect(); err != nil {
				logger.Printf("Error disconnecting from printer: %v", err)
			}
		}

	default:
		logger.Fatalf("Unknown printer type: %s", cfg.Printer.Type)
	}

	// Initialize core service
	printService := core.NewPrintService(printerAdapter)

	// Initialize API handler
	handler := api.NewHandler(printService)

	// Setup router
	router := api.SetupRouter(handler)

	// Start server in a goroutine
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	go func() {
		logger.Printf("Server starting on %s", serverAddr)
		logger.Printf("Swagger documentation available at http://localhost%s/swagger/index.html", serverAddr)
		if err := router.Run(serverAddr); err != nil {
			logger.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Shutting down server...")

	// Cleanup
	cleanup()

	logger.Println("Server stopped")
}
