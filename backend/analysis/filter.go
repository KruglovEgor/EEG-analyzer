package analysis

import (
	"math"
)

// ButterworthFilter applies a simple 2nd order Butterworth bandpass filter
type ButterworthFilter struct {
	lowFreq  float64
	highFreq float64
	fs       float64
}

// NewButterworthFilter creates a new bandpass filter
func NewButterworthFilter(lowFreq, highFreq, samplingRate float64) *ButterworthFilter {
	return &ButterworthFilter{
		lowFreq:  lowFreq,
		highFreq: highFreq,
		fs:       samplingRate,
	}
}

// Apply filters the signal using a simple moving average bandpass approximation
// For production, consider using a proper IIR/FIR filter implementation
func (f *ButterworthFilter) Apply(signal []float64) []float64 {
	filtered := make([]float64, len(signal))
	copy(filtered, signal)

	// High-pass filter (removes DC and very low frequencies)
	if f.lowFreq > 0 {
		filtered = f.highPass(filtered)
	}

	// Low-pass filter (removes high frequencies)
	if f.highFreq < f.fs/2 {
		filtered = f.lowPass(filtered)
	}

	return filtered
}

func (f *ButterworthFilter) highPass(signal []float64) []float64 {
	// Simple first-order high-pass filter
	alpha := 1.0 / (1.0 + f.fs/(2*math.Pi*f.lowFreq))
	filtered := make([]float64, len(signal))
	filtered[0] = signal[0]

	for i := 1; i < len(signal); i++ {
		filtered[i] = alpha*(filtered[i-1]+signal[i]-signal[i-1])
	}

	return filtered
}

func (f *ButterworthFilter) lowPass(signal []float64) []float64 {
	// Simple first-order low-pass filter
	alpha := (2 * math.Pi * f.highFreq) / (f.fs + 2*math.Pi*f.highFreq)
	filtered := make([]float64, len(signal))
	filtered[0] = signal[0]

	for i := 1; i < len(signal); i++ {
		filtered[i] = alpha*signal[i] + (1-alpha)*filtered[i-1]
	}

	return filtered
}

// RemoveDCOffset removes the mean from the signal
func RemoveDCOffset(signal []float64) []float64 {
	if len(signal) == 0 {
		return signal
	}

	// Calculate mean
	sum := 0.0
	for _, v := range signal {
		sum += v
	}
	mean := sum / float64(len(signal))

	// Remove mean
	result := make([]float64, len(signal))
	for i, v := range signal {
		result[i] = v - mean
	}

	return result
}

// Normalize normalizes the signal to [-1, 1] range
func Normalize(signal []float64) []float64 {
	if len(signal) == 0 {
		return signal
	}

	// Find max absolute value
	maxAbs := 0.0
	for _, v := range signal {
		abs := math.Abs(v)
		if abs > maxAbs {
			maxAbs = abs
		}
	}

	if maxAbs == 0 {
		return signal
	}

	// Normalize
	result := make([]float64, len(signal))
	for i, v := range signal {
		result[i] = v / maxAbs
	}

	return result
}
