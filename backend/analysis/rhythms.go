package analysis

import (
	"eeg-analyzer/models"
)

// RhythmAnalysisResult contains the analysis results for a specific rhythm
type RhythmAnalysisResult struct {
	Rhythm         models.RhythmType
	AbsolutePower  float64
	RelativePower  float64
	FilteredSignal []float64
	PSD            []float64
	Frequencies    []float64
}

// AnalyzeRhythm performs complete analysis for a specific rhythm band
func AnalyzeRhythm(
	timeData []float64,
	ampData []float64,
	rhythm models.RhythmType,
	samplingRate float64,
) (*RhythmAnalysisResult, error) {
	// Get rhythm band
	band, ok := models.DefaultRhythmBands[rhythm]
	if !ok {
		return nil, models.ErrInvalidRhythmBand
	}

	// Preprocess signal
	signal := RemoveDCOffset(ampData)

	// Filter signal for this rhythm band
	filter := NewButterworthFilter(band.Low, band.High, samplingRate)
	filtered := filter.Apply(signal)

	// Compute FFT on filtered signal
	fftResult := ComputeFFT(filtered, samplingRate)

	// Calculate band power
	absolutePower := ExtractBandPower(fftResult, band.Low, band.High)
	totalPower := CalculateTotalPower(fftResult)
	relativePower := CalculateRelativePower(absolutePower, totalPower)

	return &RhythmAnalysisResult{
		Rhythm:         rhythm,
		AbsolutePower:  absolutePower,
		RelativePower:  relativePower,
		FilteredSignal: filtered,
		PSD:            fftResult.PSD,
		Frequencies:    fftResult.Frequencies,
	}, nil
}

// AnalyzeMultipleRhythms analyzes multiple rhythm bands
func AnalyzeMultipleRhythms(
	timeData []float64,
	ampData []float64,
	rhythms []models.RhythmType,
	samplingRate float64,
) (map[models.RhythmType]*RhythmAnalysisResult, error) {
	results := make(map[models.RhythmType]*RhythmAnalysisResult)

	for _, rhythm := range rhythms {
		result, err := AnalyzeRhythm(timeData, ampData, rhythm, samplingRate)
		if err != nil {
			return nil, err
		}
		results[rhythm] = result
	}

	return results, nil
}
