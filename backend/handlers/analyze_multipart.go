package handlers

import (
	"eeg-analyzer/analysis"
	"eeg-analyzer/models"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// AnalyzeMultipart handles multipart/form-data requests for EEG analysis
// This endpoint accepts direct file uploads without base64 encoding
func AnalyzeMultipart(c *gin.Context) {
	// Parse multipart form (32 MB max)
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to parse multipart form: %v", err)})
		return
	}

	analysisMode := c.PostForm("analysisMode")
	if analysisMode != "SINGLE" && analysisMode != "GROUP" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "analysisMode must be 'SINGLE' or 'GROUP'"})
		return
	}

	if analysisMode == "SINGLE" {
		handleSingleMultipart(c)
	} else {
		handleGroupMultipart(c)
	}
}

func handleSingleMultipart(c *gin.Context) {
	analysisID := c.PostForm("analysisId")
	experimentName := c.PostForm("experimentName")
	timeColumn := c.PostForm("timeColumn")
	amplitudeColumn := c.PostForm("amplitudeColumn")
	rhythmsStr := c.PostForm("rhythms")

	if analysisID == "" || experimentName == "" || timeColumn == "" || amplitudeColumn == "" || rhythmsStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required fields: analysisId, experimentName, timeColumn, amplitudeColumn, rhythms"})
		return
	}

	// Parse rhythms
	rhythmStrs := strings.Split(rhythmsStr, ",")
	var rhythms []models.RhythmType
	for _, r := range rhythmStrs {
		r = strings.TrimSpace(r)
		if r != "" {
			rhythms = append(rhythms, models.RhythmType(r))
		}
	}

	if len(rhythms) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one rhythm must be specified"})
		return
	}

	// Get uploaded file
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to read file: %v", err)})
		return
	}
	defer file.Close()

	// Parse CSV from file
	timeData, ampData, err := parseCSVFromReader(file, timeColumn, amplitudeColumn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Calculate sampling rate
	samplingRate := analysis.CalculateSamplingRate(timeData)

	// Analyze multiple rhythms
	results, err := analysis.AnalyzeMultipleRhythms(timeData, ampData, rhythms, samplingRate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate total power once for relative power calculation
	var totalPower float64
	if len(results) > 0 {
		for _, result := range results {
			totalPower = analysis.CalculateTotalPowerFromArrays(result.Frequencies, result.PSD)
			break
		}
	}

	// Build response
	response := &models.EEGAnalysisResponse{
		AnalysisID:     analysisID,
		AnalysisMode:   models.ModeSingle,
		ExperimentName: &experimentName,
		Rhythms:        rhythms,
		DataByRhythm:   make(map[models.RhythmType]models.EEGPlotPair),
	}

	absolutePowers := make([][2]interface{}, 0, len(rhythms))
	relativePowers := make([][2]interface{}, 0, len(rhythms))

	for rhythm, result := range results {
		plotPair := createPlotPair(timeData, ampData, result, samplingRate)
		response.DataByRhythm[rhythm] = plotPair

		absolutePowers = append(absolutePowers, [2]interface{}{rhythm, result.AbsolutePower})
		// Calculate relative power from common total power
		relativePower := analysis.CalculateRelativePower(result.AbsolutePower, totalPower)
		relativePowers = append(relativePowers, [2]interface{}{rhythm, relativePower})
	}

	response.AbsolutePowers = absolutePowers
	response.RelativePowers = relativePowers

	c.JSON(http.StatusOK, response)
}

func handleGroupMultipart(c *gin.Context) {
	analysisID := c.PostForm("analysisId")
	rhythm := c.PostForm("rhythm")
	timeColumn := c.PostForm("timeColumn")
	amplitudeColumn := c.PostForm("amplitudeColumn")
	experimentNamesStr := c.PostForm("experimentNames")

	if analysisID == "" || rhythm == "" || timeColumn == "" || amplitudeColumn == "" || experimentNamesStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required fields: analysisId, rhythm, timeColumn, amplitudeColumn, experimentNames"})
		return
	}

	experimentNames := strings.Split(experimentNamesStr, ",")
	for i := range experimentNames {
		experimentNames[i] = strings.TrimSpace(experimentNames[i])
	}

	// Get uploaded files
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to parse form: %v", err)})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one file must be uploaded"})
		return
	}

	if len(files) != len(experimentNames) {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Number of files (%d) must match number of experiment names (%d)", len(files), len(experimentNames))})
		return
	}

	// Build response
	response := &models.EEGAnalysisResponse{
		AnalysisID:       analysisID,
		AnalysisMode:     models.ModeGroup,
		ExperimentNames:  experimentNames,
		DataByExperiment: make(map[string]models.EEGPlotPair),
	}

	absolutePowers := make([][2]interface{}, 0, len(files))
	relativePowers := make([][2]interface{}, 0, len(files))

	rhythmType := models.RhythmType(rhythm)
	response.Rhythm = &rhythmType

	// Analyze each file
	for i, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to open file %s: %v", fileHeader.Filename, err)})
			return
		}
		defer file.Close()

		timeData, ampData, err := parseCSVFromReader(file, timeColumn, amplitudeColumn)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error parsing %s: %v", fileHeader.Filename, err)})
			return
		}

		samplingRate := analysis.CalculateSamplingRate(timeData)

		// Analyze the specified rhythm with filter params
		result, err := analysis.AnalyzeRhythm(timeData, ampData, rhythmType, samplingRate, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error analyzing %s: %v", fileHeader.Filename, err)})
			return
		}

		plotPair := createPlotPair(timeData, ampData, result, samplingRate)
		expName := experimentNames[i]
		response.DataByExperiment[expName] = plotPair

		absolutePowers = append(absolutePowers, [2]interface{}{expName, result.AbsolutePower})
		relativePowers = append(relativePowers, [2]interface{}{expName, result.RelativePower})
	}

	response.AbsolutePowers = absolutePowers
	response.RelativePowers = relativePowers

	c.JSON(http.StatusOK, response)
}

// parseCSVFromReader reads CSV data from an io.Reader
func parseCSVFromReader(reader io.Reader, timeColumn, amplitudeColumn string) ([]float64, []float64, error) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read CSV: %v", err)
	}

	if len(records) < 2 {
		return nil, nil, fmt.Errorf("CSV file must have at least a header and one data row")
	}

	header := records[0]
	timeIdx := -1
	ampIdx := -1

	for i, col := range header {
		if col == timeColumn {
			timeIdx = i
		}
		if col == amplitudeColumn {
			ampIdx = i
		}
	}

	if timeIdx == -1 {
		return nil, nil, fmt.Errorf("time column '%s' not found in CSV", timeColumn)
	}
	if ampIdx == -1 {
		return nil, nil, fmt.Errorf("amplitude column '%s' not found in CSV", amplitudeColumn)
	}

	var times []float64
	var amplitudes []float64

	for i, record := range records[1:] {
		if len(record) <= timeIdx || len(record) <= ampIdx {
			return nil, nil, fmt.Errorf("row %d has insufficient columns", i+2)
		}

		t, err := strconv.ParseFloat(record[timeIdx], 64)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid time value at row %d: %v", i+2, err)
		}

		amp, err := strconv.ParseFloat(record[ampIdx], 64)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid amplitude value at row %d: %v", i+2, err)
		}

		times = append(times, t)
		amplitudes = append(amplitudes, amp)
	}

	if len(times) == 0 {
		return nil, nil, fmt.Errorf("no data rows found in CSV")
	}

	return times, amplitudes, nil
}
