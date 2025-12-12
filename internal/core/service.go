package core

import (
	"encoding/json"
	"fmt"
)

// PrintService orchestrates printing operations.
// It follows the hexagonal architecture pattern by depending only on the Printer port.
type PrintService struct {
	printer Printer
}

// NewPrintService creates a new print service with the given printer implementation.
func NewPrintService(printer Printer) *PrintService {
	return &PrintService{
		printer: printer,
	}
}

// PrintText sends plain text to the printer.
func (s *PrintService) PrintText(text string) error {
	if text == "" {
		return fmt.Errorf("text cannot be empty")
	}
	return s.printer.PrintText(text)
}

// PrintJSON formats JSON data and sends it to the printer.
// The JSON is pretty-printed with indentation for better readability.
func (s *PrintService) PrintJSON(data interface{}) error {
	if data == nil {
		return fmt.Errorf("data cannot be nil")
	}

	// Pretty-print the JSON with indentation
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return s.printer.PrintText(string(jsonBytes))
}
