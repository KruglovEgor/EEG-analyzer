package models

// EEGPreviewRequest is the request for preview endpoint (multipart/form-data)
// Note: File is uploaded as multipart file, not in this struct
type EEGPreviewRequest struct {
	PreviewID       string     `json:"previewId" binding:"required"`
	ExperimentName  string     `json:"experimentName" binding:"required"`
	Rhythm          RhythmType `json:"rhythm" binding:"required"`
	TimeColumn      string     `json:"timeColumn"`      // default: "time"
	AmplitudeColumn string     `json:"amplitudeColumn"` // default: "amplitude"
	// Filter parameters (all optional)
	FilterMin   *float64 `json:"filterMin"`   // bandpass min freq (Hz)
	FilterMax   *float64 `json:"filterMax"`   // bandpass max freq (Hz)
	FilterOrder *int     `json:"filterOrder"` // filter order (1-4)
	NPerSeg     *int     `json:"nPerSeg"`     // Welch window size
	NOverlap    *int     `json:"nOverlap"`    // Welch overlap
}

// EEGPreviewResponse is the response for preview endpoint
type EEGPreviewResponse struct {
	PreviewID      string      `json:"previewId" binding:"required"`
	ExperimentName string      `json:"experimentName" binding:"required"`
	Rhythm         RhythmType  `json:"rhythm" binding:"required"`
	Plot           EEGPlotPair `json:"plot" binding:"required"`
}
