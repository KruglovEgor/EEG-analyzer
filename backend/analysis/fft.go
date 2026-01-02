package analysis

import (
	"math"
	"math/cmplx"

	"github.com/mjibson/go-dsp/fft"
)

// FFTResult contains the results of FFT analysis
type FFTResult struct {
	Frequencies []float64 // Frequency bins (Hz)
	PSD         []float64 // Power Spectral Density
	Magnitude   []float64 // Magnitude spectrum
}

// ComputeFFT performs FFT on the signal and returns frequency domain data
func ComputeFFT(signal []float64, samplingRate float64) *FFTResult {
	n := len(signal)
	if n < 2 {
		return &FFTResult{
			Frequencies: []float64{},
			PSD:         []float64{},
			Magnitude:   []float64{},
		}
	}

	// Apply Hamming window to reduce spectral leakage
	windowed := applyHammingWindow(signal)

	// Perform FFT
	fftOutput := fft.FFTReal(windowed)

	// Only take first half (positive frequencies)
	halfN := n / 2
	frequencies := make([]float64, halfN)
	psd := make([]float64, halfN)
	magnitude := make([]float64, halfN)

	// Calculate frequency bins
	freqStep := samplingRate / float64(n)
	for i := 0; i < halfN; i++ {
		frequencies[i] = float64(i) * freqStep
	}

	// Calculate PSD and magnitude
	for i := 0; i < halfN; i++ {
		mag := cmplx.Abs(fftOutput[i])
		magnitude[i] = mag

		// PSD = (magnitude^2) / N
		// Ensure PSD > 0 for logarithmic scale
		psdValue := (mag * mag) / float64(n)
		if psdValue < 1e-10 {
			psdValue = 1e-10 // Minimum value for log scale
		}
		psd[i] = psdValue
	}

	return &FFTResult{
		Frequencies: frequencies,
		PSD:         psd,
		Magnitude:   magnitude,
	}
}

// applyHammingWindow applies a Hamming window to reduce spectral leakage
func applyHammingWindow(signal []float64) []float64 {
	n := len(signal)
	windowed := make([]float64, n)

	for i := 0; i < n; i++ {
		window := 0.54 - 0.46*math.Cos(2*math.Pi*float64(i)/float64(n-1))
		windowed[i] = signal[i] * window
	}

	return windowed
}

// ExtractBandPower calculates the power in a specific frequency band
func ExtractBandPower(fftResult *FFTResult, lowFreq, highFreq float64) float64 {
	power := 0.0
	count := 0

	for i, freq := range fftResult.Frequencies {
		if freq >= lowFreq && freq <= highFreq {
			power += fftResult.PSD[i]
			count++
		}
	}

	if count == 0 {
		return 0
	}

	// Return average power in the band
	return power / float64(count)
}

// CalculateRelativePower calculates relative power as percentage of total
func CalculateRelativePower(bandPower, totalPower float64) float64 {
	if totalPower == 0 {
		return 0
	}
	return (bandPower / totalPower) * 100.0
}

// CalculateTotalPower calculates total power across all frequencies
func CalculateTotalPower(fftResult *FFTResult) float64 {
	total := 0.0
	for _, psd := range fftResult.PSD {
		total += psd
	}
	return total
}
