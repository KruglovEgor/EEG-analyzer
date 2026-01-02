package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Version string `json:"version"`
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Returns the health status of the API
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:  "ok",
		Message: "EEG Analyzer API is running",
		Version: "1.0.0",
	})
}
