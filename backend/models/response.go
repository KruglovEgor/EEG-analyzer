package models

// EEGCombinedDataPoint represents a unified data point for Recharts
// Keys are: "x", "raw", "filtered", "psd" depending on plot type
type EEGCombinedDataPoint map[string]float64

// EEGLinePlot contains plot data and optional highlight range
type EEGLinePlot struct {
	Data            []EEGCombinedDataPoint `json:"data" binding:"required"`
	XHighlightRange *[2]float64            `json:"xHighlightRange,omitempty"` // For highlighting frequency bands in PSD
}

// EEGPlotPair contains both PSD and signal plots
type EEGPlotPair struct {
	PSDPlot    EEGLinePlot `json:"psdPlot" binding:"required"`
	SignalPlot EEGLinePlot `json:"signalPlot" binding:"required"`
}

// EEGAnalysisResponse is the response structure
type EEGAnalysisResponse struct {
	AnalysisID   string       `json:"analysisId" binding:"required"`
	AnalysisMode AnalysisMode `json:"analysisMode" binding:"required"`

	// For SINGLE mode
	ExperimentName *string                    `json:"experimentName,omitempty"`
	Rhythms        []RhythmType               `json:"rhythms,omitempty"`
	AbsolutePowers [][2]interface{}           `json:"absolutePowers,omitempty"` // [[RhythmType, number]]
	RelativePowers [][2]interface{}           `json:"relativePowers,omitempty"` // [[RhythmType, number]]
	DataByRhythm   map[RhythmType]EEGPlotPair `json:"dataByRhythm,omitempty"`

	// For GROUP mode
	ExperimentNames  []string               `json:"experimentNames,omitempty"`
	Rhythm           *RhythmType            `json:"rhythm,omitempty"`
	DataByExperiment map[string]EEGPlotPair `json:"dataByExperiment,omitempty"`
}
