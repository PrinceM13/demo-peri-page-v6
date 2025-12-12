package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockPrinter is a mock implementation of the Printer interface for testing.
type MockPrinter struct {
	mock.Mock
}

// PrintText is a mock implementation of the Printer.PrintText method.
func (m *MockPrinter) PrintText(text string) error {
	args := m.Called(text)
	return args.Error(0)
}
