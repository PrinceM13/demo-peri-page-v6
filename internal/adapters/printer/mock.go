package printer

import (
	"fmt"
	"log"
)

// MockPrinter is a test implementation that prints to stdout.
// Useful for development and testing without actual hardware.
type MockPrinter struct {
	logger *log.Logger
}

// NewMockPrinter creates a new mock printer instance.
func NewMockPrinter(logger *log.Logger) *MockPrinter {
	if logger == nil {
		logger = log.Default()
	}
	return &MockPrinter{
		logger: logger,
	}
}

// PrintText outputs text to stdout, simulating a real printer.
func (m *MockPrinter) PrintText(text string) error {
	m.logger.Println("=== MOCK PRINTER OUTPUT ===")
	fmt.Println(text)
	m.logger.Println("=== END MOCK PRINTER OUTPUT ===")
	return nil
}
