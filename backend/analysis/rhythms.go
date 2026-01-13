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
// Matches Python implementation from lab notebook:
// 1. Remove DC offset
// 2. Apply FFT pre-filter (0.5-40 Hz) to remove noise
// 3. Compute PSD using Welch's method on the pre-filtered signal
// 4. Extract power in the rhythm band
// 5. Apply bandpass filter for visualization
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

	// Step 1: Remove DC offset
	signal := RemoveDCOffset(ampData)

	// Step 2: Pre-filter with FFT to remove noise (0.5-40 Hz)
	// This matches Python: df_cleaned = clean_eeg_data(df, lower_freq=0.5, upper_freq=40, ...)
	signalPreFiltered := PreFilterSignal(signal, samplingRate, 0.5, 40)

	// Step 3: Compute PSD using Welch's method on the PRE-FILTERED signal
	// This matches Python: calculate_lambd_power(df_cleaned['A0 с FFT (В)'].values, sampling_rate)
	fftResult := ComputeWelchPSD(signalPreFiltered, samplingRate, 1024)

	// Step 4: Calculate band power by integrating PSD in the rhythm band
	absolutePower := ExtractBandPower(fftResult, band.Low, band.High)
	totalPower := CalculateTotalPower(fftResult)
	relativePower := CalculateRelativePower(absolutePower, totalPower)

	// Step 5: Apply bandpass filter for visualization (on pre-filtered signal)
	filter := NewButterworthFilter(band.Low, band.High, samplingRate)
	filtered := filter.Apply(signalPreFiltered)

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
