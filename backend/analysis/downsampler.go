package analysis

import (
	"math"
)

// DownsampleStrategy defines the downsampling method
type DownsampleStrategy int

const (
	// StrategyLTTB uses Largest-Triangle-Three-Buckets algorithm (best for visualization)
	StrategyLTTB DownsampleStrategy = iota
	// StrategyDecimate uses simple decimation (faster but less accurate)
	StrategyDecimate
	// StrategyAverage uses averaging within buckets
	StrategyAverage
)

// DownsampleData reduces the number of data points while preserving shape
func DownsampleData(x, y []float64, targetPoints int, strategy DownsampleStrategy) ([]float64, []float64) {
	if len(x) <= targetPoints || targetPoints < 3 {
		return x, y
	}

	switch strategy {
	case StrategyLTTB:
		return downsampleLTTB(x, y, targetPoints)
	case StrategyAverage:
		return downsampleAverage(x, y, targetPoints)
	default:
		return downsampleDecimate(x, y, targetPoints)
	}
}

// downsampleDecimate performs simple decimation (every Nth point)
func downsampleDecimate(x, y []float64, targetPoints int) ([]float64, []float64) {
	n := len(x)
	step := float64(n) / float64(targetPoints)

	xOut := make([]float64, 0, targetPoints)
	yOut := make([]float64, 0, targetPoints)

	for i := 0; i < targetPoints; i++ {
		idx := int(float64(i) * step)
		if idx >= n {
			idx = n - 1
		}
		xOut = append(xOut, x[idx])
		yOut = append(yOut, y[idx])
	}

	return xOut, yOut
}

// downsampleAverage averages points within each bucket
func downsampleAverage(x, y []float64, targetPoints int) ([]float64, []float64) {
	n := len(x)
	bucketSize := n / targetPoints

	xOut := make([]float64, 0, targetPoints)
	yOut := make([]float64, 0, targetPoints)

	for i := 0; i < targetPoints; i++ {
		start := i * bucketSize
		end := start + bucketSize
		if i == targetPoints-1 {
			end = n
		}

		if end > n {
			end = n
		}

		// Calculate averages
		xSum := 0.0
		ySum := 0.0
		count := 0

		for j := start; j < end; j++ {
			xSum += x[j]
			ySum += y[j]
			count++
		}

		if count > 0 {
			xOut = append(xOut, xSum/float64(count))
			yOut = append(yOut, ySum/float64(count))
		}
	}

	return xOut, yOut
}

// downsampleLTTB implements Largest-Triangle-Three-Buckets algorithm
// This preserves visual characteristics better than simple decimation
func downsampleLTTB(x, y []float64, targetPoints int) ([]float64, []float64) {
	n := len(x)
	if n <= targetPoints {
		return x, y
	}

	xOut := make([]float64, targetPoints)
	yOut := make([]float64, targetPoints)

	// Always include first point
	xOut[0] = x[0]
	yOut[0] = y[0]

	// Always include last point
	xOut[targetPoints-1] = x[n-1]
	yOut[targetPoints-1] = y[n-1]

	// Bucket size
	bucketSize := float64(n-2) / float64(targetPoints-2)

	// Index of point in the current bucket
	a := 0

	for i := 0; i < targetPoints-2; i++ {
		// Calculate point average for next bucket (for area calculation)
		avgRangeStart := int(math.Floor(float64(i+1)*bucketSize)) + 1
		avgRangeEnd := int(math.Floor(float64(i+2)*bucketSize)) + 1

		if avgRangeEnd >= n {
			avgRangeEnd = n
		}

		avgX := 0.0
		avgY := 0.0
		avgRangeLength := avgRangeEnd - avgRangeStart

		for ; avgRangeStart < avgRangeEnd; avgRangeStart++ {
			avgX += x[avgRangeStart]
			avgY += y[avgRangeStart]
		}
		avgX /= float64(avgRangeLength)
		avgY /= float64(avgRangeLength)

		// Get the range for this bucket
		rangeStart := int(math.Floor(float64(i)*bucketSize)) + 1
		rangeEnd := int(math.Floor(float64(i+1)*bucketSize)) + 1

		// Point a
		pointAX := x[a]
		pointAY := y[a]

		maxArea := -1.0
		maxAreaPoint := rangeStart

		for ; rangeStart < rangeEnd; rangeStart++ {
			// Calculate triangle area
			area := math.Abs((pointAX-avgX)*(y[rangeStart]-pointAY) -
				(pointAX-x[rangeStart])*(avgY-pointAY)) * 0.5

			if area > maxArea {
				maxArea = area
				maxAreaPoint = rangeStart
			}
		}

		xOut[i+1] = x[maxAreaPoint]
		yOut[i+1] = y[maxAreaPoint]
		a = maxAreaPoint
	}

	return xOut, yOut
}

// SuggestTargetPoints suggests a reasonable number of points for visualization
func SuggestTargetPoints(dataLength int) int {
	// Aim for 1000-2000 points for good visualization performance
	if dataLength <= 2000 {
		return dataLength
	}
	if dataLength <= 10000 {
		return 1500
	}
	return 1000
}
