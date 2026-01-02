package analysis

import (
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"

	"eeg-analyzer/models"
)

// ParseCSVData extracts time and amplitude data from base64 encoded CSV
func ParseCSVData(config models.EEGFileConfig) ([]float64, []float64, error) {
	if config.RawFile == nil {
		return nil, nil, fmt.Errorf("rawFile is nil")
	}

	// Decode base64
	decoded, err := base64.StdEncoding.DecodeString(*config.RawFile)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	// Parse CSV
	reader := csv.NewReader(strings.NewReader(string(decoded)))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, models.ErrInvalidCSV
	}

	if len(records) < 2 {
		return nil, nil, models.ErrInsufficientData
	}

	// Find column indices
	header := records[0]
	timeIdx := -1
	ampIdx := -1

	for i, col := range header {
		colTrimmed := strings.TrimSpace(col)
		if colTrimmed == config.TimeColumn {
			timeIdx = i
		}
		if colTrimmed == config.AmplitudeColumn {
			ampIdx = i
		}
	}

	if timeIdx == -1 || ampIdx == -1 {
		return nil, nil, models.ErrColumnNotFound
	}

	// Extract data
	timeData := make([]float64, 0, len(records)-1)
	ampData := make([]float64, 0, len(records)-1)

	for i := 1; i < len(records); i++ {
		row := records[i]
		if len(row) <= timeIdx || len(row) <= ampIdx {
			continue
		}

		time, err1 := strconv.ParseFloat(strings.TrimSpace(row[timeIdx]), 64)
		amp, err2 := strconv.ParseFloat(strings.TrimSpace(row[ampIdx]), 64)

		if err1 == nil && err2 == nil {
			timeData = append(timeData, time)
			ampData = append(ampData, amp)
		}
	}

	if len(timeData) < 10 {
		return nil, nil, models.ErrInsufficientData
	}

	return timeData, ampData, nil
}

// CalculateSamplingRate estimates sampling rate from time data
func CalculateSamplingRate(timeData []float64) float64 {
	if len(timeData) < 2 {
		return 250.0 // Default fallback
	}

	// Calculate average time difference
	totalDiff := 0.0
	validDiffs := 0

	for i := 1; i < len(timeData) && i < 100; i++ {
		diff := timeData[i] - timeData[i-1]
		if diff > 0 {
			totalDiff += diff
			validDiffs++
		}
	}

	if validDiffs == 0 {
		return 250.0
	}

	avgDiff := totalDiff / float64(validDiffs)
	return 1.0 / avgDiff // Hz
}
