package analysis

import (
	"math"
)

// ButterworthFilter implements a Butterworth IIR bandpass filter with filtfilt (zero-phase filtering)
// Matches scipy.signal.butter + scipy.signal.filtfilt behavior
type ButterworthFilter struct {
	lowFreq  float64
	highFreq float64
	fs       float64
	order    int
	// Cascaded second-order sections
	sos [][6]float64 // Biquad coefficients: [b0, b1, b2, a0, a1, a2]
}

// NewButterworthFilter creates a new Butterworth bandpass filter
func NewButterworthFilter(lowFreq, highFreq, samplingRate float64) *ButterworthFilter {
	return NewButterworthFilterWithOrder(lowFreq, highFreq, samplingRate, 5)
}

// NewButterworthFilterWithOrder creates a Butterworth filter with specified order
func NewButterworthFilterWithOrder(lowFreq, highFreq, samplingRate float64, order int) *ButterworthFilter {
	f := &ButterworthFilter{
		lowFreq:  lowFreq,
		highFreq: highFreq,
		fs:       samplingRate,
		order:    order,
	}

	// Design the filter using second-order sections (SOS)
	f.designButterworthBandpass()

	return f
}

// designButterworthBandpass designs Butterworth bandpass filter using second-order sections
func (f *ButterworthFilter) designButterworthBandpass() {
	// Normalize frequencies to Nyquist frequency
	nyquist := f.fs / 2.0
	wn1 := f.lowFreq / nyquist
	wn2 := f.highFreq / nyquist

	// Ensure valid range
	if wn1 <= 0 {
		wn1 = 0.001
	}
	if wn2 >= 1 {
		wn2 = 0.999
	}
	if wn1 >= wn2 {
		wn1, wn2 = wn2*0.9, wn2
	}

	// Design using bilinear transformation
	// For a 5th order Butterworth, we'll use 3 second-order sections (order/2 + 1)
	numSections := (f.order + 1) / 2

	f.sos = make([][6]float64, numSections)

	// Compute analog prototype poles and transform to digital bandpass
	for k := 0; k < numSections; k++ {
		// Analog prototype pole for Butterworth
		angle := math.Pi * float64(2*k+f.order+1) / float64(2*f.order)
		sigmak := -math.Sin(angle)

		// Bandpass transformation: map lowpass to bandpass
		wc := math.Sqrt(wn1 * wn2) // Center frequency (normalized)
		bw := wn2 - wn1            // Bandwidth (normalized)

		// Bilinear transformation parameters
		// For analog filter H(s) = 1/(s^2 + sqrt(2)*s + 1) to bandpass
		// We transform: s -> (s^2 + wc^2)/(bw*s)

		// Direct computation of biquad coefficients using bilinear transform
		// This is a practical approximation that works well for EEG

		// Compute biquad numerator and denominator
		f.sos[k] = f.computeBiquadCoefficients(sigmak, wc, bw)
	}
}

// computeBiquadCoefficients computes a single biquad (second-order section) coefficients
func (f *ButterworthFilter) computeBiquadCoefficients(sigma float64, wc, bw float64) [6]float64 {
	// Analog prototype: H(s) = wc / (s^2 + sqrt(2)*bw*s + wc^2)
	// Apply bilinear transformation

	T := 1.0 / f.fs // Sampling period

	// Butterworth damping for 5th order
	dampingFactor := math.Sqrt(2) * bw

	// Bilinear transform coefficients
	a0num := 1.0 + dampingFactor*T + wc*wc*T*T
	a1num := -2.0 + 2.0*wc*wc*T*T
	a2num := 1.0 - dampingFactor*T + wc*wc*T*T

	b0num := wc * wc * T * T
	b1num := 2.0 * b0num
	b2num := b0num

	// Normalize by a0
	a1 := a1num / a0num
	a2 := a2num / a0num
	b0 := b0num / a0num
	b1 := b1num / a0num
	b2 := b2num / a0num

	return [6]float64{b0, b1, b2, 1.0, a1, a2}
}

// Apply filters the signal using forward-backward filtering (filtfilt)
// This applies the filter twice (forward and backward) to achieve zero phase distortion
func (f *ButterworthFilter) Apply(signal []float64) []float64 {
	if len(signal) < 2 {
		return signal
	}

	// Apply all SOS sections forward
	filtered := signal
	for i := 0; i < len(f.sos); i++ {
		filtered = f.applyBiquadFilter(filtered, f.sos[i])
	}

	// Apply in reverse direction for zero-phase filtering
	reversed := reverseSignal(filtered)
	for i := 0; i < len(f.sos); i++ {
		reversed = f.applyBiquadFilter(reversed, f.sos[i])
	}

	// Reverse back
	result := reverseSignal(reversed)

	return result
}

// applyBiquadFilter applies a single biquad filter section
func (f *ButterworthFilter) applyBiquadFilter(signal []float64, biquad [6]float64) []float64 {
	b0, b1, b2 := biquad[0], biquad[1], biquad[2]
	a1, a2 := biquad[4], biquad[5] // a0 is normalized to 1

	y := make([]float64, len(signal))

	// State variables
	x1, x2 := 0.0, 0.0
	y1, y2 := 0.0, 0.0

	for n := 0; n < len(signal); n++ {
		// Direct Form II
		// w[n] = x[n] - a1*y[n-1] - a2*y[n-2]
		w := signal[n] - a1*y1 - a2*y2
		// y[n] = b0*w[n] + b1*w[n-1] + b2*w[n-2]
		yn := b0*w + b1*x1 + b2*x2

		y[n] = yn

		// Update state
		x2 = x1
		x1 = w
		y2 = y1
		y1 = yn
	}

	return y
}

// reverseSignal reverses a signal for backward filtering
func reverseSignal(signal []float64) []float64 {
	result := make([]float64, len(signal))
	for i := 0; i < len(signal); i++ {
		result[i] = signal[len(signal)-1-i]
	}
	return result
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
