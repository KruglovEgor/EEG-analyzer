package analysis

import (
	"github.com/mjibson/go-dsp/fft"
)

// PreFilterSignal applies FFT-based bandpass filtering to remove noise outside 0.5-40 Hz
// This matches the Python implementation: clean_eeg_data(df, lower_freq=0.5, upper_freq=40, ...)
// Used before PSD calculation to focus on EEG-relevant frequencies
func PreFilterSignal(signal []float64, samplingRate, lowFreq, highFreq float64) []float64 {
	n := len(signal)
	if n < 2 {
		return signal
	}

	// Perform FFT
	fftOutput := fft.FFTReal(signal)

	// Calculate frequency step
	freqStep := samplingRate / float64(n)

	// Zero out frequencies outside the desired band
	for i := 0; i < len(fftOutput); i++ {
		freq := float64(i) * freqStep
		if freq < lowFreq || freq > highFreq {
			fftOutput[i] = complex(0, 0)
		}
	}

	// Perform inverse FFT
	filtered := fft.IFFT(fftOutput)

	// Take only the real part and trim to original length
	result := make([]float64, n)
	for i := 0; i < n && i < len(filtered); i++ {
		result[i] = real(filtered[i])
	}

	return result
}
