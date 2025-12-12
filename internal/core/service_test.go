package core

import (
	"testing"
)

// mockPrinter is a test implementation of the Printer interface
type mockPrinter struct {
	lastText string
	err      error
}

func (m *mockPrinter) PrintText(text string) error {
	if m.err != nil {
		return m.err
	}
	m.lastText = text
	return nil
}

func TestNewPrintService(t *testing.T) {
	printer := &mockPrinter{}
	service := NewPrintService(printer)

	if service == nil {
		t.Fatal("Expected service to be created")
	}

	if service.printer != printer {
		t.Error("Service printer does not match provided printer")
	}
}

func TestPrintText(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		wantErr bool
	}{
		{
			name:    "valid text",
			text:    "Hello, World!",
			wantErr: false,
		},
		{
			name:    "empty text",
			text:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printer := &mockPrinter{}
			service := NewPrintService(printer)

			err := service.PrintText(tt.text)

			if (err != nil) != tt.wantErr {
				t.Errorf("PrintText() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && printer.lastText != tt.text {
				t.Errorf("Expected printer to receive %q, got %q", tt.text, printer.lastText)
			}
		})
	}
}

func TestPrintJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    interface{}
		wantErr bool
	}{
		{
			name: "valid JSON object",
			data: map[string]interface{}{
				"name": "John",
				"age":  30,
			},
			wantErr: false,
		},
		{
			name:    "nil data",
			data:    nil,
			wantErr: true,
		},
		{
			name: "nested JSON",
			data: map[string]interface{}{
				"user": map[string]interface{}{
					"name":  "Jane",
					"email": "jane@example.com",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printer := &mockPrinter{}
			service := NewPrintService(printer)

			err := service.PrintJSON(tt.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("PrintJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && printer.lastText == "" {
				t.Error("Expected printer to receive text")
			}
		})
	}
}
