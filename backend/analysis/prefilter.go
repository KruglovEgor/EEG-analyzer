package analysis

import (
	"github.com/mjibson/go-dsp/fft"
)

// PreFilterSignal applies FFT-based band-pass filtering to remove noise outside desired frequency range
// This matches Python: remove_noise_fft(signal, sampling_rate, lower_freq, upper_freq)
func PreFilterSignal(signal []float64, samplingRate, lowFreq, highFreq float64) []float64 {
	n := len(signal)
	if n < 2 {
		return signal
	}

	// Perform real FFT (matching Python's np.fft.rfft)
	fftOutput := fft.FFTReal(signal)

	// Calculate frequency for each bin
	freqStep := samplingRate / float64(n)

	// Zero out frequencies outside the desired range [lowFreq, highFreq]
	halfN := len(fftOutput) // rfft returns n/2 + 1 complex values
	for i := 0; i < halfN; i++ {
		freq := float64(i) * freqStep
		if freq < lowFreq || freq > highFreq {
			fftOutput[i] = 0
		}
	}

	// Inverse FFT to get cleaned signal (matching Python's np.fft.irfft)
	// We need to convert back from rfft format to full FFT format for IFFT
	// Create full FFT array: [0, 1, 2, ..., N/2, -N/2+1, ..., -2, -1]
	fullFFT := make([]complex128, n)

	// Copy positive frequencies
	for i := 0; i < halfN && i < n; i++ {
		fullFFT[i] = fftOutput[i]
	}

	// Mirror for negative frequencies (conjugate symmetry for real signals)
	for i := 1; i < halfN-1 && i < n/2; i++ {
		fullFFT[n-i] = complex(real(fftOutput[i]), -imag(fftOutput[i]))
	}

	// Inverse FFT
	cleanedComplex := fft.IFFT(fullFFT)

	// Extract real part
	cleaned := make([]float64, n)
	for i := 0; i < n; i++ {
		cleaned[i] = real(cleanedComplex[i])
	}

	return cleaned
}
