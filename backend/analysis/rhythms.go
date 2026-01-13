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
// Matches Python implementation:
// 1. Clean signal (remove DC offset)
// 2. Compute PSD using Welch's method on the full signal
// 3. Extract power in the rhythm band
// 4. Filter signal for visualization
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

	// Preprocess signal (remove DC offset)
	signal := RemoveDCOffset(ampData)

	// Compute PSD using Welch's method on the FULL signal (not filtered)
	// This matches Python: calculate_lambd_power(df_cleaned['A0 с FFT (В)'].values, sampling_rate)
	fftResult := ComputeWelchPSD(signal, samplingRate, 1024)

	// Calculate band power by integrating PSD in the rhythm band
	absolutePower := ExtractBandPower(fftResult, band.Low, band.High)
	totalPower := CalculateTotalPower(fftResult)
	relativePower := CalculateRelativePower(absolutePower, totalPower)

	// Filter signal for visualization (bandpass filter)
	filter := NewButterworthFilter(band.Low, band.High, samplingRate)
	filtered := filter.Apply(signal)

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
