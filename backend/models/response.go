package models

// EEGSeriesMetadata describes a single curve/series in a plot
type EEGSeriesMetadata struct {
	DataKey        string  `json:"dataKey" binding:"required"`
	Legend         *string `json:"legend,omitempty"`
	PreferredColor *string `json:"preferredColor,omitempty"`
}

// EEGCombinedDataPoint represents a unified data point for Recharts
type EEGCombinedDataPoint map[string]float64

// EEGLinePlot contains plot data and metadata
type EEGLinePlot struct {
	XMin            *float64                `json:"xMin,omitempty"`
	XMax            *float64                `json:"xMax,omitempty"`
	YMin            *float64                `json:"yMin,omitempty"`
	YMax            *float64                `json:"yMax,omitempty"`
	XAxisName       *string                 `json:"xAxisName,omitempty"`
	YAxisName       *string                 `json:"yAxisName,omitempty"`
	XHighlightRange *[2]float64             `json:"xHighlightRange,omitempty"`
	YLogarithmic    *bool                   `json:"yLogarithmic,omitempty"`
	ShowLegend      *bool                   `json:"showLegend,omitempty"`
	Area            *bool                   `json:"area,omitempty"`
	SeriesMetadata  []EEGSeriesMetadata     `json:"seriesMetadata,omitempty"`
	Data            []EEGCombinedDataPoint  `json:"data" binding:"required"`
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
	ExperimentName  *string                      `json:"experimentName,omitempty"`
	Rhythms         []RhythmType                 `json:"rhythms,omitempty"`
	AbsolutePowers  [][2]interface{}             `json:"absolutePowers,omitempty"` // [[RhythmType, number]]
	RelativePowers  [][2]interface{}             `json:"relativePowers,omitempty"` // [[RhythmType, number]]
	DataByRhythm    map[RhythmType]EEGPlotPair   `json:"dataByRhythm,omitempty"`
	
	// For GROUP mode
	ExperimentNames []string                     `json:"experimentNames,omitempty"`
	Rhythm          *RhythmType                  `json:"rhythm,omitempty"`
	AbsolutePowersGroup [][2]interface{}         `json:"absolutePowers,omitempty"` // [[string, number]]
	RelativePowersGroup [][2]interface{}         `json:"relativePowers,omitempty"` // [[string, number]]
	DataByExperiment    map[string]EEGPlotPair   `json:"dataByExperiment,omitempty"`
}
