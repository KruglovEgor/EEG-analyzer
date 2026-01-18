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
// 3. Extract power in the rhythm band from PSD (uses filterParams if provided)
// 4. Apply Butterworth bandpass filter on original signal for visualization
func AnalyzeRhythm(
	timeData []float64,
	ampData []float64,
	rhythm models.RhythmType,
	samplingRate float64,
	filterParams *models.EEGFilterParams,
) (*RhythmAnalysisResult, error) {
	// Get default parameters if not provided
	if filterParams == nil {
		defaults := models.GetDefaultFilterParams(rhythm)
		filterParams = &defaults
	} else {
		// Validate and apply defaults for invalid values
		if err := filterParams.Validate(rhythm); err != nil {
			return nil, err
		}
	}

	// Step 1: Remove DC offset from original signal
	signal := RemoveDCOffset(ampData)

	// Step 2: For PSD computation, apply FFT pre-filter (0.5-40 Hz)
	// This matches Python implementation: clean_eeg_data() + welch()
	signalForPSD := ApplyFFTBandpass(signal, samplingRate, 0.5, 40.0)
	fftResult := ComputeWelchPSD(signalForPSD, samplingRate, filterParams.NPerSeg, filterParams.NOverlap)

	// Step 3: Calculate band power using filterParams boundaries
	absolutePower := ExtractBandPower(fftResult, filterParams.FilterMin, filterParams.FilterMax)
	totalPower := CalculateTotalPower(fftResult)
	relativePower := CalculateRelativePower(absolutePower, totalPower)

	// Step 4: For visualization, apply Butterworth filter on the ORIGINAL signal
	// Uses filterParams boundaries and order
	filter := NewButterworthFilter(filterParams.FilterMin, filterParams.FilterMax, samplingRate, filterParams.FilterOrder)
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
// For SINGLE mode with multiple rhythms, uses default parameters for each rhythm
func AnalyzeMultipleRhythms(
	timeData []float64,
	ampData []float64,
	rhythms []models.RhythmType,
	samplingRate float64,
) (map[models.RhythmType]*RhythmAnalysisResult, error) {
	results := make(map[models.RhythmType]*RhythmAnalysisResult)

	for _, rhythm := range rhythms {
		// Use default parameters for each rhythm in multi-rhythm analysis
		result, err := AnalyzeRhythm(timeData, ampData, rhythm, samplingRate, nil)
		if err != nil {
			return nil, err
		}
		results[rhythm] = result
	}

	return results, nil
}

// AnalyzeMultipleRhythmsWithParams analyzes multiple rhythm bands with custom PSD params
// Uses default frequency boundaries for each rhythm, but applies custom nPerSeg/nOverlap
// filterMin, filterMax, and filterOrder from filterParams are IGNORED (each rhythm uses its own defaults)
func AnalyzeMultipleRhythmsWithParams(
	timeData []float64,
	ampData []float64,
	rhythms []models.RhythmType,
	samplingRate float64,
	filterParams *models.EEGFilterParams,
) (map[models.RhythmType]*RhythmAnalysisResult, error) {
	results := make(map[models.RhythmType]*RhythmAnalysisResult)

	for _, rhythm := range rhythms {
		// Get default params for this rhythm
		params := models.GetDefaultFilterParams(rhythm)

		// Override nPerSeg, nOverlap, and filterOrder from filterParams
		// filterMin and filterMax are IGNORED (use rhythm defaults for frequency boundaries)
		if filterParams != nil {
			if filterParams.NPerSeg > 0 {
				params.NPerSeg = filterParams.NPerSeg
			}
			if filterParams.NOverlap > 0 && filterParams.NOverlap < params.NPerSeg {
				params.NOverlap = filterParams.NOverlap
			}
			if filterParams.FilterOrder > 0 {
				params.FilterOrder = filterParams.FilterOrder
			}
		}

		result, err := AnalyzeRhythm(timeData, ampData, rhythm, samplingRate, &params)
		if err != nil {
			return nil, err
		}
		results[rhythm] = result
	}

	return results, nil
}
