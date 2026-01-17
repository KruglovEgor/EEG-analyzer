package analysis

import (
	"math"
)

// DownsampleData reduces the number of data points while preserving visual shape
// Uses the LTTB (Largest-Triangle-Three-Buckets) algorithm for optimal visualization
func DownsampleData(x, y []float64, targetPoints int) ([]float64, []float64) {
	if len(x) <= targetPoints || targetPoints < 3 {
		return x, y
	}

	return downsampleLTTB(x, y, targetPoints)
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
			area := math.Abs((pointAX-avgX)*(y[rangeStart]-pointAY)-
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
