package core

// Printer defines the port for interacting with printer devices.
// This is the core domain interface that all printer adapters must implement.
type Printer interface {
	// PrintText sends text to the printer for printing.
	// Returns an error if the printing operation fails.
	PrintText(text string) error
}
