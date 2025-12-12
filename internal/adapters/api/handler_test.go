package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/princem/peripage-printer/internal/core/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/print", handler.Print)
	router.GET("/health", handler.HealthCheck)
	return router
}

func TestHandler_Print_ValidTextRequest(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		mockSetup      func(*mocks.MockPrinter)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful print with text only",
			requestBody: map[string]interface{}{
				"text": "Hello, World!",
			},
			mockSetup: func(m *mocks.MockPrinter) {
				m.On("PrintText", "Hello, World!").Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"success": true,
				"message": "Print job completed successfully",
			},
		},
		{
			name: "successful print with multiline text",
			requestBody: map[string]interface{}{
				"text": "Line 1\nLine 2\nLine 3",
			},
			mockSetup: func(m *mocks.MockPrinter) {
				m.On("PrintText", "Line 1\nLine 2\nLine 3").Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"success": true,
				"message": "Print job completed successfully",
			},
		},
		{
			name: "empty text returns error",
			requestBody: map[string]interface{}{
				"text": "",
			},
			mockSetup: func(m *mocks.MockPrinter) {
				// No mock call expected
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Either 'text' or 'data' must be provided",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockPrinter := new(mocks.MockPrinter)
			tt.mockSetup(mockPrinter)
			service := &mockPrintService{printer: mockPrinter}
			handler := &Handler{service: service}
			router := setupTestRouter(handler)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/print", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Act
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			for key, expectedValue := range tt.expectedBody {
				assert.Contains(t, response, key)
				if str, ok := expectedValue.(string); ok {
					assert.Contains(t, response[key], str)
				} else {
					assert.Equal(t, expectedValue, response[key])
				}
			}

			mockPrinter.AssertExpectations(t)
		})
	}
}

func TestHandler_Print_ValidJSONDataRequest(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		mockSetup      func(*mocks.MockPrinter)
		expectedStatus int
	}{
		{
			name: "successful print with JSON data",
			requestBody: map[string]interface{}{
				"data": map[string]interface{}{
					"name": "John Doe",
					"age":  30,
				},
			},
			mockSetup: func(m *mocks.MockPrinter) {
				// We expect pretty-printed JSON
				m.On("PrintText", "{\n  \"age\": 30,\n  \"name\": \"John Doe\"\n}").Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "successful print with nested JSON",
			requestBody: map[string]interface{}{
				"data": map[string]interface{}{
					"user": map[string]interface{}{
						"name":  "Jane",
						"email": "jane@example.com",
					},
					"timestamp": "2025-12-12",
				},
			},
			mockSetup: func(m *mocks.MockPrinter) {
				// Accept any string for nested JSON (order may vary)
				m.On("PrintText", `{
  "timestamp": "2025-12-12",
  "user": {
    "email": "jane@example.com",
    "name": "Jane"
  }
}`).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "data takes precedence over text",
			requestBody: map[string]interface{}{
				"text": "This will be ignored",
				"data": map[string]interface{}{
					"message": "Only this will be printed",
				},
			},
			mockSetup: func(m *mocks.MockPrinter) {
				m.On("PrintText", "{\n  \"message\": \"Only this will be printed\"\n}").Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockPrinter := new(mocks.MockPrinter)
			tt.mockSetup(mockPrinter)
			service := &mockPrintService{printer: mockPrinter}
			handler := &Handler{service: service}
			router := setupTestRouter(handler)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/print", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Act
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)
			mockPrinter.AssertExpectations(t)
		})
	}
}

func TestHandler_Print_ErrorCases(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		contentType    string
		mockSetup      func(*mocks.MockPrinter)
		expectedStatus int
		errorContains  string
	}{
		{
			name:           "invalid JSON body",
			requestBody:    `{"text": invalid json}`,
			contentType:    "application/json",
			mockSetup:      func(m *mocks.MockPrinter) {},
			expectedStatus: http.StatusBadRequest,
			errorContains:  "Invalid request body",
		},
		{
			name:           "empty request body",
			requestBody:    `{}`,
			contentType:    "application/json",
			mockSetup:      func(m *mocks.MockPrinter) {},
			expectedStatus: http.StatusBadRequest,
			errorContains:  "Either 'text' or 'data' must be provided",
		},
		{
			name: "printer error is propagated",
			requestBody: `{"text": "Test"}`,
			contentType: "application/json",
			mockSetup: func(m *mocks.MockPrinter) {
				m.On("PrintText", "Test").Return(errors.New("printer offline")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			errorContains:  "Print failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockPrinter := new(mocks.MockPrinter)
			tt.mockSetup(mockPrinter)
			service := &mockPrintService{printer: mockPrinter}
			handler := &Handler{service: service}
			router := setupTestRouter(handler)

			req := httptest.NewRequest(http.MethodPost, "/print", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()

			// Act
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Contains(t, response, "error")
			assert.Contains(t, response["error"], tt.errorContains)

			mockPrinter.AssertExpectations(t)
		})
	}
}

func TestHandler_HealthCheck(t *testing.T) {
	// Arrange
	mockPrinter := new(mocks.MockPrinter)
	service := &mockPrintService{printer: mockPrinter}
	handler := &Handler{service: service}
	router := setupTestRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "healthy", response["status"])
}

// mockPrintService is a test helper that mimics PrintService behavior
type mockPrintService struct {
	printer *mocks.MockPrinter
}

func (m *mockPrintService) PrintText(text string) error {
	if text == "" {
		return errors.New("text cannot be empty")
	}
	return m.printer.PrintText(text)
}

func (m *mockPrintService) PrintJSON(data interface{}) error {
	if data == nil {
		return errors.New("data cannot be nil")
	}

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return m.printer.PrintText(string(jsonBytes))
}
