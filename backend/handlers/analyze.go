package handlers

import (
	"fmt"
	"net/http"

	"eeg-analyzer/analysis"
	"eeg-analyzer/models"

	"github.com/gin-gonic/gin"
)

// AnalyzeEEG godoc
// @Summary Analyze EEG data
// @Description Performs EEG analysis in SINGLE or GROUP mode
// @Tags analysis
// @Accept json
// @Produce json
// @Param request body models.EEGAnalysisRequest true "Analysis Request"
// @Success 200 {object} models.EEGAnalysisResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /analyze [post]
func AnalyzeEEG(c *gin.Context) {
	var req models.EEGAnalysisRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Route to appropriate handler
	var response *models.EEGAnalysisResponse
	var err error

	if req.AnalysisMode == models.ModeSingle {
		response, err = handleSingleAnalysis(&req)
	} else {
		response, err = handleGroupAnalysis(&req)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// handleSingleAnalysis processes a single file with multiple rhythms
func handleSingleAnalysis(req *models.EEGAnalysisRequest) (*models.EEGAnalysisResponse, error) {
	// Parse CSV data
	timeData, ampData, err := analysis.ParseCSVData(*req.File)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV: %w", err)
	}

	// Calculate sampling rate
	samplingRate := analysis.CalculateSamplingRate(timeData)

	// Analyze all requested rhythms
	rhythmResults, err := analysis.AnalyzeMultipleRhythms(timeData, ampData, req.Rhythms, samplingRate)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze rhythms: %w", err)
	}

	// Build response
	response := &models.EEGAnalysisResponse{
		AnalysisID:     req.AnalysisID,
		AnalysisMode:   models.ModeSingle,
		ExperimentName: &req.File.ExperimentName,
		Rhythms:        req.Rhythms,
		DataByRhythm:   make(map[models.RhythmType]models.EEGPlotPair),
	}

	// Calculate absolute and relative powers
	absolutePowers := make([][2]interface{}, 0, len(req.Rhythms))
	relativePowers := make([][2]interface{}, 0, len(req.Rhythms))

	for _, rhythm := range req.Rhythms {
		result := rhythmResults[rhythm]
		absolutePowers = append(absolutePowers, [2]interface{}{rhythm, result.AbsolutePower})
		relativePowers = append(relativePowers, [2]interface{}{rhythm, result.RelativePower})

		// Create plot pair for this rhythm
		plotPair := createPlotPair(timeData, ampData, result, samplingRate)
		response.DataByRhythm[rhythm] = plotPair
	}

	response.AbsolutePowers = absolutePowers
	response.RelativePowers = relativePowers

	return response, nil
}

// handleGroupAnalysis processes multiple files with a single rhythm
func handleGroupAnalysis(req *models.EEGAnalysisRequest) (*models.EEGAnalysisResponse, error) {
	if req.Rhythm == nil {
		return nil, fmt.Errorf("rhythm is required for GROUP mode")
	}

	response := &models.EEGAnalysisResponse{
		AnalysisID:       req.AnalysisID,
		AnalysisMode:     models.ModeGroup,
		Rhythm:           req.Rhythm,
		ExperimentNames:  make([]string, 0, len(req.Files)),
		DataByExperiment: make(map[string]models.EEGPlotPair),
	}

	absolutePowers := make([][2]interface{}, 0, len(req.Files))
	relativePowers := make([][2]interface{}, 0, len(req.Files))

	// Process each file
	for _, fileConfig := range req.Files {
		// Parse CSV
		timeData, ampData, err := analysis.ParseCSVData(fileConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to parse CSV for %s: %w", fileConfig.ExperimentName, err)
		}

		// Calculate sampling rate
		samplingRate := analysis.CalculateSamplingRate(timeData)

		// Analyze the specified rhythm
		result, err := analysis.AnalyzeRhythm(timeData, ampData, *req.Rhythm, samplingRate)
		if err != nil {
			return nil, fmt.Errorf("failed to analyze rhythm for %s: %w", fileConfig.ExperimentName, err)
		}

		// Add to response
		response.ExperimentNames = append(response.ExperimentNames, fileConfig.ExperimentName)
		absolutePowers = append(absolutePowers, [2]interface{}{fileConfig.ExperimentName, result.AbsolutePower})
		relativePowers = append(relativePowers, [2]interface{}{fileConfig.ExperimentName, result.RelativePower})

		// Create plot pair
		plotPair := createPlotPair(timeData, ampData, result, samplingRate)
		response.DataByExperiment[fileConfig.ExperimentName] = plotPair
	}

	response.AbsolutePowers = absolutePowers
	response.RelativePowers = relativePowers

	return response, nil
}

// createPlotPair creates a plot pair (PSD + signal) from analysis results
func createPlotPair(
	timeData []float64,
	rawSignal []float64,
	result *analysis.RhythmAnalysisResult,
	samplingRate float64,
) models.EEGPlotPair {
	// Downsample data for efficient visualization
	targetPoints := analysis.SuggestTargetPoints(len(timeData))

	// Downsample signal plot data
	timeDownsampled, rawDownsampled := analysis.DownsampleData(
		timeData, rawSignal, targetPoints, analysis.StrategyLTTB,
	)
	_, filteredDownsampled := analysis.DownsampleData(
		timeData, result.FilteredSignal, targetPoints, analysis.StrategyLTTB,
	)

	// Downsample PSD data
	freqDownsampled, psdDownsampled := analysis.DownsampleData(
		result.Frequencies, result.PSD, targetPoints/2, analysis.StrategyLTTB,
	)

	// Create signal plot data (combined format for Recharts)
	signalData := make([]models.EEGCombinedDataPoint, len(timeDownsampled))
	for i := range timeDownsampled {
		signalData[i] = models.EEGCombinedDataPoint{
			"x":        timeDownsampled[i],
			"raw":      rawDownsampled[i],
			"filtered": filteredDownsampled[i],
		}
	}

	// Create PSD plot data
	psdData := make([]models.EEGCombinedDataPoint, len(freqDownsampled))
	for i := range freqDownsampled {
		psdData[i] = models.EEGCombinedDataPoint{
			"x":   freqDownsampled[i],
			"psd": psdDownsampled[i],
		}
	}

	// Create metadata
	// Metadata fields commented out as per frontend request - frontend will add localized labels
	// rawLegend := "Raw Signal"
	// filteredLegend := "Filtered Signal"
	// psdLegend := "PSD"
	// secondaryColor := "secondary"
	// primaryColor := "primary"

	// yLogTrue := true
	// xAxisFreq := "Frequency (Hz)"
	// yAxisPower := "Power (µV²/Hz)"
	// xAxisTime := "Time (s)"
	// yAxisAmp := "Amplitude (µV)"

	// Get rhythm band for highlight
	band := models.DefaultRhythmBands[result.Rhythm]
	xHighlight := [2]float64{band.Low, band.High}

	return models.EEGPlotPair{
		PSDPlot: models.EEGLinePlot{
			// SeriesMetadata: []models.EEGSeriesMetadata{
			//   {DataKey: "psd", Legend: &psdLegend, PreferredColor: &secondaryColor},
			// },
			Data: psdData,
			// YLogarithmic:    &yLogTrue,
			// XAxisName:       &xAxisFreq,
			// YAxisName:       &yAxisPower,
			XHighlightRange: &xHighlight,
		},
		SignalPlot: models.EEGLinePlot{
			// SeriesMetadata: []models.EEGSeriesMetadata{
			//   {DataKey: "raw", Legend: &rawLegend, PreferredColor: &primaryColor},
			//   {DataKey: "filtered", Legend: &filteredLegend, PreferredColor: &secondaryColor},
			// },
			Data: signalData,
			// XAxisName: &xAxisTime,
			// YAxisName: &yAxisAmp,
		},
	}
}
