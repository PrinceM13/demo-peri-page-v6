package printer

import (
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMockPrinter(t *testing.T) {
	tests := []struct {
		name   string
		logger *log.Logger
	}{
		{
			name:   "with nil logger uses default",
			logger: nil,
		},
		{
			name:   "with provided logger",
			logger: log.New(&bytes.Buffer{}, "[test] ", log.LstdFlags),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			printer := NewMockPrinter(tt.logger)

			// Assert
			require.NotNil(t, printer)
			assert.NotNil(t, printer.logger)
		})
	}
}

func TestMockPrinter_PrintText(t *testing.T) {
	tests := []struct {
		name         string
		text         string
		expectedText string
	}{
		{
			name:         "prints simple text",
			text:         "Hello, World!",
			expectedText: "Hello, World!",
		},
		{
			name:         "prints multiline text",
			text:         "Line 1\nLine 2\nLine 3",
			expectedText: "Line 1\nLine 2\nLine 3",
		},
		{
			name:         "prints empty string",
			text:         "",
			expectedText: "",
		},
		{
			name:         "prints JSON",
			text:         `{"name":"John","age":30}`,
			expectedText: `{"name":"John","age":30}`,
		},
		{
			name: "prints long text",
			text: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. " +
				"Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
			expectedText: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. " +
				"Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			var buf bytes.Buffer
			logger := log.New(&buf, "", 0)
			printer := NewMockPrinter(logger)

			// Act
			err := printer.PrintText(tt.text)

			// Assert
			require.NoError(t, err)

			// Verify log output contains the markers
			logOutput := buf.String()
			assert.Contains(t, logOutput, "=== MOCK PRINTER OUTPUT ===")
			assert.Contains(t, logOutput, "=== END MOCK PRINTER OUTPUT ===")
		})
	}
}

func TestMockPrinter_PrintText_AlwaysSucceeds(t *testing.T) {
	// The mock printer should never return an error
	printer := NewMockPrinter(nil)

	testCases := []string{
		"",
		"short",
		"very long text " + string(make([]byte, 10000)),
		"\n\n\n",
		"special chars: !@#$%^&*()",
	}

	for _, text := range testCases {
		err := printer.PrintText(text)
		assert.NoError(t, err, "MockPrinter should never fail")
	}
}
