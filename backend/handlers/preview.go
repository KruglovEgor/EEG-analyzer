package handlers

import (
	"fmt"
	"net/http"

	"eeg-analyzer/analysis"
	"eeg-analyzer/models"

	"github.com/gin-gonic/gin"
)

// PreviewEEG godoc
// @Summary Preview EEG signal with filter parameters
// @Description Shows how filters will affect the signal before full analysis
// @Tags analysis
// @Accept json
// @Produce json
// @Param request body models.EEGPreviewRequest true "Preview Request"
// @Success 200 {object} models.EEGPreviewResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /preview [post]
func PreviewEEG(c *gin.Context) {
	var req models.EEGPreviewRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse CSV data
	timeData, ampData, err := analysis.ParseCSVData(req.File)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to parse CSV: %v", err)})
		return
	}

	// Calculate sampling rate
	samplingRate := analysis.CalculateSamplingRate(timeData)

	// Analyze with filter params
	result, err := analysis.AnalyzeRhythm(timeData, ampData, req.Rhythm, samplingRate, req.FilterParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create plot pair
	plotPair := createPlotPair(timeData, ampData, result, samplingRate)

	// Build response
	response := &models.EEGPreviewResponse{
		PreviewID:      req.PreviewID,
		ExperimentName: req.ExperimentName,
		Rhythm:         req.Rhythm,
		Plot:           plotPair,
	}

	c.JSON(http.StatusOK, response)
}
