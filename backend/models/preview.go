package models

// EEGPreviewRequest is the request for preview endpoint
type EEGPreviewRequest struct {
	PreviewID      string           `json:"previewId" binding:"required"`
	File           EEGFileConfig    `json:"file" binding:"required"`
	ExperimentName string           `json:"experimentName" binding:"required"`
	Rhythm         RhythmType       `json:"rhythm" binding:"required"`
	FilterParams   *EEGFilterParams `json:"filterParams"`
}

// EEGPreviewResponse is the response for preview endpoint
type EEGPreviewResponse struct {
	PreviewID      string      `json:"previewId" binding:"required"`
	ExperimentName string      `json:"experimentName" binding:"required"`
	Rhythm         RhythmType  `json:"rhythm" binding:"required"`
	Plot           EEGPlotPair `json:"plot" binding:"required"`
}
