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
// Processing steps:
// 1. Remove DC offset (mean)
// 2. Apply FFT bandpass (0.5-40 Hz) and compute PSD (matches Python implementation)
// 3. Extract power in the rhythm band from PSD
// 4. Apply Butterworth 1st order bandpass filter on original signal for visualization
func AnalyzeRhythm(
	timeData []float64,
	ampData []float64,
	rhythm models.RhythmType,
	samplingRate float64,
) (*RhythmAnalysisResult, error) {
	// Get rhythm frequency band
	band, ok := models.DefaultRhythmBands[rhythm]
	if !ok {
		return nil, models.ErrInvalidRhythmBand
	}

	// Step 1: Remove DC offset from original signal
	signal := RemoveDCOffset(ampData)

	// Step 2: For PSD computation, apply FFT pre-filter (0.5-40 Hz)
	// This matches Python implementation: clean_eeg_data() + welch()
	signalForPSD := ApplyFFTBandpass(signal, samplingRate, 0.5, 40.0)
	fftResult := ComputeWelchPSD(signalForPSD, samplingRate, 1024)

	// Step 3: Calculate band power by integrating PSD in the rhythm band
	absolutePower := ExtractBandPower(fftResult, band.Low, band.High)
	totalPower := CalculateTotalPower(fftResult)
	relativePower := CalculateRelativePower(absolutePower, totalPower)

	// Step 4: For visualization, apply Butterworth filter on the ORIGINAL signal
	// This provides cleaner visual representation without harsh FFT artifacts
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
