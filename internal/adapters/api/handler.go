package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/princem/peripage-printer/internal/core"
)

// Handler manages HTTP requests for the printer API.
type Handler struct {
	service *core.PrintService
}

// NewHandler creates a new API handler.
func NewHandler(service *core.PrintService) *Handler {
	return &Handler{
		service: service,
	}
}

// PrintRequest represents the request body for the print endpoint.
type PrintRequest struct {
	Text string                 `json:"text" example:"Hello, World!"`
	Data map[string]interface{} `json:"data,omitempty"`
}

// PrintResponse represents the response for the print endpoint.
type PrintResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Print job completed successfully"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request"`
}

// Print handles the POST /print endpoint.
// @Summary Print text or JSON data
// @Description Prints text or formatted JSON data to the Peripage printer
// @Tags print
// @Accept json
// @Produce json
// @Param request body PrintRequest true "Print request"
// @Success 200 {object} PrintResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /print [post]
func (h *Handler) Print(c *gin.Context) {
	var req PrintRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body: " + err.Error(),
		})
		return
	}

	// If data is provided, print JSON; otherwise print text
	var err error
	if req.Data != nil && len(req.Data) > 0 {
		err = h.service.PrintJSON(req.Data)
	} else if req.Text != "" {
		err = h.service.PrintText(req.Text)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Either 'text' or 'data' must be provided",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Print failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PrintResponse{
		Success: true,
		Message: "Print job completed successfully",
	})
}

// HealthCheck handles the GET /health endpoint.
// @Summary Health check
// @Description Returns the health status of the service
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}
