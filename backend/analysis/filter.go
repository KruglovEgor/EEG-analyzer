package analysis

import (
	"math"
)

// ButterworthFilter implements a simple 1st order Butterworth bandpass filter
// This is a classic, standard implementation that provides gentle filtering
type ButterworthFilter struct {
	lowFreq  float64
	highFreq float64
	fs       float64
	order    int
}

// NewButterworthFilter creates a new Butterworth bandpass filter with specified order
func NewButterworthFilter(lowFreq, highFreq, samplingRate float64, order int) *ButterworthFilter {
	if order < 1 {
		order = 1
	}
	return &ButterworthFilter{
		lowFreq:  lowFreq,
		highFreq: highFreq,
		fs:       samplingRate,
		order:    order,
	}
}

// Apply applies the bandpass filter to the signal
func (f *ButterworthFilter) Apply(signal []float64) []float64 {
	if len(signal) < 2 {
		return signal
	}

	filtered := make([]float64, len(signal))
	copy(filtered, signal)

	// Apply high-pass and low-pass filters sequentially, order times
	for i := 0; i < f.order; i++ {
		// Apply high-pass filter (removes DC and low frequencies)
		if f.lowFreq > 0 {
			filtered = f.highPass(filtered)
		}

		// Apply low-pass filter (removes high frequencies)
		if f.highFreq < f.fs/2 {
			filtered = f.lowPass(filtered)
		}
	}

	return filtered
}

// highPass implements a 1st order high-pass Butterworth filter
func (f *ButterworthFilter) highPass(signal []float64) []float64 {
	// RC coefficient for 1st order high-pass
	RC := 1.0 / (2.0 * math.Pi * f.lowFreq)
	dt := 1.0 / f.fs
	alpha := RC / (RC + dt)

	filtered := make([]float64, len(signal))
	filtered[0] = signal[0]

	for i := 1; i < len(signal); i++ {
		filtered[i] = alpha * (filtered[i-1] + signal[i] - signal[i-1])
	}

	return filtered
}

// lowPass implements a 1st order low-pass Butterworth filter
func (f *ButterworthFilter) lowPass(signal []float64) []float64 {
	// RC coefficient for 1st order low-pass
	RC := 1.0 / (2.0 * math.Pi * f.highFreq)
	dt := 1.0 / f.fs
	alpha := dt / (RC + dt)

	filtered := make([]float64, len(signal))
	filtered[0] = signal[0]

	for i := 1; i < len(signal); i++ {
		filtered[i] = alpha*signal[i] + (1.0-alpha)*filtered[i-1]
	}

	return filtered
}

// RemoveDCOffset removes the DC component (mean) from the signal
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

	// Subtract mean from each sample
	result := make([]float64, len(signal))
	for i, v := range signal {
		result[i] = v - mean
	}

	return result
}
