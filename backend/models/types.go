package models

// RhythmType represents EEG brain wave types
type RhythmType string

const (
	RhythmAlpha  RhythmType = "ALPHA"  // 8–13 Hz
	RhythmBeta   RhythmType = "BETA"   // 13–30 Hz
	RhythmGamma  RhythmType = "GAMMA"  // 30–100 Hz
	RhythmDelta  RhythmType = "DELTA"  // 0.5–4 Hz
	RhythmTheta  RhythmType = "THETA"  // 4–8 Hz
	RhythmKappa  RhythmType = "KAPPA"  // 8–13 Hz
	RhythmLambda RhythmType = "LAMBDA" // 4–8 Hz
	RhythmMu     RhythmType = "MU"     // 8–13 Hz
)

// BrainZone represents brain regions
type BrainZone string

const (
	BrainZoneFrontal  BrainZone = "FRONTAL"
	BrainZoneParietal BrainZone = "PARIETAL"
	BrainZoneCentral  BrainZone = "CENTRAL"
	BrainZoneTemporal BrainZone = "TEMPORAL"
	BrainZoneOccipital BrainZone = "OCCIPITAL"
)

// AnalysisMode represents the type of analysis
type AnalysisMode string

const (
	ModeSingle AnalysisMode = "SINGLE"
	ModeGroup  AnalysisMode = "GROUP"
)

// RhythmBand defines frequency ranges for rhythm types
type RhythmBand struct {
	Low  float64
	High float64
}

// DefaultRhythmBands maps rhythm types to their frequency ranges
var DefaultRhythmBands = map[RhythmType]RhythmBand{
	RhythmDelta:  {Low: 0.5, High: 4},
	RhythmTheta:  {Low: 4, High: 8},
	RhythmAlpha:  {Low: 8, High: 13},
	RhythmBeta:   {Low: 13, High: 30},
	RhythmGamma:  {Low: 30, High: 100},
	RhythmMu:     {Low: 8, High: 13},
	RhythmLambda: {Low: 4, High: 8},
	RhythmKappa:  {Low: 8, High: 13},
}
