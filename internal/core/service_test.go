package core

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/princem/peripage-printer/internal/core/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewPrintService(t *testing.T) {
	// Arrange
	mockPrinter := new(mocks.MockPrinter)

	// Act
	service := NewPrintService(mockPrinter)

	// Assert
	require.NotNil(t, service, "Service should not be nil")
	assert.Equal(t, mockPrinter, service.printer, "Service should use provided printer")
}

func TestPrintService_PrintText(t *testing.T) {
	tests := []struct {
		name          string
		text          string
		mockSetup     func(*mocks.MockPrinter)
		expectedError string
	}{
		{
			name: "successful print with valid text",
			text: "Hello, World!",
			mockSetup: func(m *mocks.MockPrinter) {
				m.On("PrintText", "Hello, World!").Return(nil).Once()
			},
			expectedError: "",
		},
		{
			name: "successful print with multiline text",
			text: "Line 1\nLine 2\nLine 3",
			mockSetup: func(m *mocks.MockPrinter) {
				m.On("PrintText", "Line 1\nLine 2\nLine 3").Return(nil).Once()
			},
			expectedError: "",
		},
		{
			name: "empty text returns error",
			text: "",
			mockSetup: func(m *mocks.MockPrinter) {
				// No call to PrintText expected
			},
			expectedError: "text cannot be empty",
		},
		{
			name: "printer error is propagated",
			text: "Test text",
			mockSetup: func(m *mocks.MockPrinter) {
				m.On("PrintText", "Test text").Return(errors.New("printer offline")).Once()
			},
			expectedError: "printer offline",
		},
		{
			name: "long text is printed successfully",
			text: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. " +
				"Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
			mockSetup: func(m *mocks.MockPrinter) {
				m.On("PrintText", "Lorem ipsum dolor sit amet, consectetur adipiscing elit. "+
					"Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.").Return(nil).Once()
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockPrinter := new(mocks.MockPrinter)
			tt.mockSetup(mockPrinter)
			service := NewPrintService(mockPrinter)

			// Act
			err := service.PrintText(tt.text)

			// Assert
			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			mockPrinter.AssertExpectations(t)
		})
	}
}

func TestPrintService_PrintJSON(t *testing.T) {
	tests := []struct {
		name          string
		data          interface{}
		mockSetup     func(*mocks.MockPrinter, string)
		expectedError string
	}{
		{
			name: "simple JSON object",
			data: map[string]interface{}{
				"name": "John",
				"age":  30,
			},
			mockSetup: func(m *mocks.MockPrinter, expectedJSON string) {
				m.On("PrintText", expectedJSON).Return(nil).Once()
			},
			expectedError: "",
		},
		{
			name: "nested JSON object",
			data: map[string]interface{}{
				"user": map[string]interface{}{
					"name":  "Jane",
					"email": "jane@example.com",
					"preferences": map[string]interface{}{
						"theme": "dark",
						"lang":  "en",
					},
				},
				"timestamp": "2025-12-12T10:00:00Z",
			},
			mockSetup: func(m *mocks.MockPrinter, expectedJSON string) {
				m.On("PrintText", expectedJSON).Return(nil).Once()
			},
			expectedError: "",
		},
		{
			name: "JSON array",
			data: []interface{}{
				map[string]interface{}{"id": 1, "name": "Item 1"},
				map[string]interface{}{"id": 2, "name": "Item 2"},
			},
			mockSetup: func(m *mocks.MockPrinter, expectedJSON string) {
				m.On("PrintText", expectedJSON).Return(nil).Once()
			},
			expectedError: "",
		},
		{
			name: "nil data returns error",
			data: nil,
			mockSetup: func(m *mocks.MockPrinter, expectedJSON string) {
				// No call to PrintText expected
			},
			expectedError: "data cannot be nil",
		},
		{
			name: "printer error is propagated",
			data: map[string]interface{}{
				"test": "data",
			},
			mockSetup: func(m *mocks.MockPrinter, expectedJSON string) {
				m.On("PrintText", expectedJSON).Return(errors.New("connection lost")).Once()
			},
			expectedError: "connection lost",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockPrinter := new(mocks.MockPrinter)

			// Generate expected JSON if data is not nil
			var expectedJSON string
			if tt.data != nil {
				jsonBytes, err := json.MarshalIndent(tt.data, "", "  ")
				require.NoError(t, err, "Test setup failed: could not marshal test data")
				expectedJSON = string(jsonBytes)
			}

			tt.mockSetup(mockPrinter, expectedJSON)
			service := NewPrintService(mockPrinter)

			// Act
			err := service.PrintJSON(tt.data)

			// Assert
			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			mockPrinter.AssertExpectations(t)
		})
	}
}

func TestPrintService_PrintJSON_VerifyFormatting(t *testing.T) {
	// This test verifies that JSON is pretty-printed with proper indentation
	mockPrinter := new(mocks.MockPrinter)
	service := NewPrintService(mockPrinter)

	data := map[string]interface{}{
		"name": "Test User",
		"age":  25,
	}

	// Capture the actual text sent to printer
	var capturedText string
	mockPrinter.On("PrintText", mock.AnythingOfType("string")).
		Run(func(args mock.Arguments) {
			capturedText = args.Get(0).(string)
		}).
		Return(nil).
		Once()

	// Act
	err := service.PrintJSON(data)

	// Assert
	require.NoError(t, err)
	assert.Contains(t, capturedText, "  ", "JSON should be indented")
	assert.Contains(t, capturedText, "\"name\": \"Test User\"")
	assert.Contains(t, capturedText, "\"age\": 25")

	mockPrinter.AssertExpectations(t)
}
