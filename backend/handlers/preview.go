package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"eeg-analyzer/analysis"
	"eeg-analyzer/models"

	"github.com/gin-gonic/gin"
)

// PreviewEEG godoc
// @Summary Preview EEG signal with filter parameters
// @Description Shows how filters will affect the signal before full analysis (multipart/form-data)
// @Tags analysis
// @Accept multipart/form-data
// @Produce json
// @Param previewId formData string true "Preview ID"
// @Param file formData file true "CSV file"
// @Param experimentName formData string true "Experiment name"
// @Param rhythm formData string true "Rhythm type (ALPHA, BETA, etc.)"
// @Param filterMin formData number false "Bandpass filter min frequency (Hz)"
// @Param filterMax formData number false "Bandpass filter max frequency (Hz)"
// @Param filterOrder formData integer false "Filter order (1-4)"
// @Param nPerSeg formData integer false "Welch window size"
// @Param nOverlap formData integer false "Welch overlap size"
// @Param timeColumn formData string false "Time column name (default: time)"
// @Param amplitudeColumn formData string false "Amplitude column name (default: amplitude)"
// @Success 200 {object} models.EEGPreviewResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /preview [post]
func PreviewEEG(c *gin.Context) {
	// Parse multipart form (32 MB max)
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to parse multipart form: %v", err)})
		return
	}

	// Get form fields
	previewID := c.PostForm("previewId")
	experimentName := c.PostForm("experimentName")
	rhythmStr := c.PostForm("rhythm")
	timeColumn := c.PostForm("timeColumn")
	amplitudeColumn := c.PostForm("amplitudeColumn")

	if previewID == "" || experimentName == "" || rhythmStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required fields: previewId, experimentName, rhythm"})
		return
	}

	// Default column names
	if timeColumn == "" {
		timeColumn = "time"
	}
	if amplitudeColumn == "" {
		amplitudeColumn = "amplitude"
	}

	// Parse rhythm
	rhythm := models.RhythmType(strings.ToUpper(rhythmStr))

	// Parse filter params from flat fields
	filterParams := parseFilterParams(c)

	// Get uploaded file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// Open file
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to open file: %v", err)})
		return
	}
	defer file.Close()

	// Parse CSV
	timeData, ampData, err := parseCSVFromReader(file, timeColumn, amplitudeColumn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to parse CSV: %v", err)})
		return
	}

	// Calculate sampling rate
	samplingRate := analysis.CalculateSamplingRate(timeData)

	// Analyze with filter params
	result, err := analysis.AnalyzeRhythm(timeData, ampData, rhythm, samplingRate, filterParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create plot pair
	plotPair := createPlotPair(timeData, ampData, result, samplingRate)

	// Build response
	response := &models.EEGPreviewResponse{
		PreviewID:      previewID,
		ExperimentName: experimentName,
		Rhythm:         rhythm,
		Plot:           plotPair,
	}

	c.JSON(http.StatusOK, response)
}
