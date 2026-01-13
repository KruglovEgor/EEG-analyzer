package analysis

import (
	"math"
	"math/cmplx"

	"github.com/mjibson/go-dsp/fft"
)

// FFTResult contains the results of FFT analysis
type FFTResult struct {
	Frequencies []float64 // Frequency bins (Hz)
	PSD         []float64 // Power Spectral Density (µV²/Hz)
	Magnitude   []float64 // Magnitude spectrum
}

// ComputeWelchPSD computes Power Spectral Density using Welch's method
// This matches the Python implementation: welch(data, fs=sampling_rate, nperseg=1024)
func ComputeWelchPSD(signal []float64, samplingRate float64, nperseg int) *FFTResult {
	n := len(signal)
	if n < nperseg {
		// Fall back to simple FFT if signal is too short
		return ComputeFFT(signal, samplingRate)
	}

	// Welch parameters (matching scipy defaults)
	if nperseg <= 0 {
		nperseg = 1024
	}
	noverlap := nperseg / 2 // 50% overlap (scipy default)

	// Calculate number of segments
	step := nperseg - noverlap
	numSegments := (n - noverlap) / step

	if numSegments < 1 {
		return ComputeFFT(signal, samplingRate)
	}

	// Frequency bins
	halfN := nperseg / 2
	frequencies := make([]float64, halfN)
	freqStep := samplingRate / float64(nperseg)
	for i := 0; i < halfN; i++ {
		frequencies[i] = float64(i) * freqStep
	}

	// Calculate window scaling factor (sum of squared window values)
	// For Hamming window, this compensates for the energy loss due to windowing
	windowScale := 0.0
	for i := 0; i < nperseg; i++ {
		w := 0.54 - 0.46*math.Cos(2*math.Pi*float64(i)/float64(nperseg-1))
		windowScale += w * w
	}

	// Accumulate PSD from all segments
	psdAccum := make([]float64, halfN)

	for seg := 0; seg < numSegments; seg++ {
		start := seg * step
		end := start + nperseg
		if end > n {
			break
		}

		// Extract segment
		segment := signal[start:end]

		// Apply Hamming window
		windowed := applyHammingWindow(segment)

		// Perform FFT on segment
		fftOutput := fft.FFTReal(windowed)

		// Accumulate power
		// scipy.signal.welch uses: (|FFT|² * 2) / (fs * sum(window²))
		// Factor of 2 accounts for negative frequencies (one-sided PSD)
		for i := 0; i < halfN; i++ {
			mag := cmplx.Abs(fftOutput[i])
			psdAccum[i] += (mag * mag)
		}
	}

	// Average over all segments and normalize
	psd := make([]float64, halfN)
	scale := samplingRate * windowScale
	for i := 0; i < halfN; i++ {
		psd[i] = (psdAccum[i] / float64(numSegments)) * 2.0 / scale
		// DC and Nyquist components should not be doubled
		if i == 0 || (nperseg%2 == 0 && i == halfN-1) {
			psd[i] = psd[i] / 2.0
		}
		// Ensure minimum value for log scale
		if psd[i] < 1e-10 {
			psd[i] = 1e-10
		}
	}

	return &FFTResult{
		Frequencies: frequencies,
		PSD:         psd,
		Magnitude:   nil, // Not computed in Welch method
	}
}

// ComputeFFT performs FFT on the signal and returns frequency domain data
// Kept for compatibility, but ComputeWelchPSD should be preferred for PSD analysis
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

// ExtractBandPower calculates the power in a specific frequency band using trapezoidal integration
func ExtractBandPower(fftResult *FFTResult, lowFreq, highFreq float64) float64 {
	// Find indices in the frequency band
	var freqsInBand []float64
	var psdInBand []float64

	for i, freq := range fftResult.Frequencies {
		if freq >= lowFreq && freq <= highFreq {
			freqsInBand = append(freqsInBand, freq)
			psdInBand = append(psdInBand, fftResult.PSD[i])
		}
	}

	if len(freqsInBand) < 2 {
		return 0
	}

	// Integrate using trapezoidal rule
	return trapezoidalIntegration(freqsInBand, psdInBand)
}

// CalculateRelativePower calculates relative power as percentage of total
func CalculateRelativePower(bandPower, totalPower float64) float64 {
	if totalPower == 0 {
		return 0
	}
	return (bandPower / totalPower) * 100.0
}

// CalculateTotalPower calculates total power across all frequencies using trapezoidal integration
func CalculateTotalPower(fftResult *FFTResult) float64 {
	if len(fftResult.Frequencies) < 2 {
		return 0
	}
	return trapezoidalIntegration(fftResult.Frequencies, fftResult.PSD)
}

// trapezoidalIntegration performs numerical integration using the trapezoidal rule
// Equivalent to NumPy's trapezoid function: np.trapezoid(y, x)
func trapezoidalIntegration(x, y []float64) float64 {
	if len(x) != len(y) || len(x) < 2 {
		return 0
	}

	sum := 0.0
	for i := 0; i < len(x)-1; i++ {
		dx := x[i+1] - x[i]
		avgY := (y[i] + y[i+1]) / 2.0
		sum += dx * avgY
	}

	return sum
}
