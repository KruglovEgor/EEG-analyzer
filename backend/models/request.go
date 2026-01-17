package models

// EEGFileConfig contains metadata about a single EEG file (LEGACY - used only for /analyze-json endpoint)
type EEGFileConfig struct {
	ID              string  `json:"id" binding:"required"`
	Filename        string  `json:"filename" binding:"required"`
	ExperimentName  string  `json:"experimentName" binding:"required"`
	TimeColumn      string  `json:"timeColumn" binding:"required"`
	AmplitudeColumn string  `json:"amplitudeColumn" binding:"required"`
	RawFile         *string `json:"rawFile"`  // Base64 encoded CSV content
	ServerID        *string `json:"serverId"` // Can be ignored
}

// EEGAnalysisRequest is the request structure for LEGACY /analyze-json endpoint (JSON with base64 files)
type EEGAnalysisRequest struct {
	AnalysisID   string           `json:"analysisId" binding:"required"`
	AnalysisMode AnalysisMode     `json:"analysisMode" binding:"required"`
	BrainZone    *BrainZone       `json:"brainZone"`
	FilterParams *EEGFilterParams `json:"filterParams"`

	// For SINGLE mode
	File    *EEGFileConfig `json:"file"`
	Rhythms []RhythmType   `json:"rhythms"`

	// For GROUP mode
	Files  []EEGFileConfig `json:"files"`
	Rhythm *RhythmType     `json:"rhythm"`
}

// Validate checks if the request is valid based on analysis mode
func (r *EEGAnalysisRequest) Validate() error {
	if r.AnalysisMode == ModeSingle {
		if r.File == nil {
			return ErrMissingFile
		}
		if len(r.Rhythms) == 0 {
			return ErrMissingRhythms
		}
	} else if r.AnalysisMode == ModeGroup {
		if len(r.Files) == 0 {
			return ErrMissingFiles
		}
		if r.Rhythm == nil {
			return ErrMissingRhythm
		}
	}
	return nil
}

// EEGAnalysisRequestSingle for SINGLE mode (multipart/form-data)
// Note: File is uploaded as multipart file, not in this struct
type EEGAnalysisRequestSingle struct {
	AnalysisID      string       `json:"analysisId" binding:"required"`
	ExperimentName  string       `json:"experimentName" binding:"required"`
	TimeColumn      string       `json:"timeColumn" binding:"required"`
	AmplitudeColumn string       `json:"amplitudeColumn" binding:"required"`
	Rhythms         []RhythmType `json:"rhythms" binding:"required"` // comma-separated in multipart
	// Filter parameters (all optional, only nPerSeg and nOverlap used in SINGLE mode)
	// FilterMin/FilterMax/FilterOrder can be sent but will be ignored in SINGLE mode
	FilterMin   *float64 `json:"filterMin"`   // bandpass min freq (Hz) - IGNORED in SINGLE mode
	FilterMax   *float64 `json:"filterMax"`   // bandpass max freq (Hz) - IGNORED in SINGLE mode
	FilterOrder *int     `json:"filterOrder"` // filter order (1-4) - IGNORED in SINGLE mode
	NPerSeg     *int     `json:"nPerSeg"`     // Welch window size
	NOverlap    *int     `json:"nOverlap"`    // Welch overlap
}

// EEGAnalysisRequestGroup for GROUP mode (multipart/form-data)
// Note: Files are uploaded as multipart files, not in this struct
type EEGAnalysisRequestGroup struct {
	AnalysisID      string     `json:"analysisId" binding:"required"`
	Rhythm          RhythmType `json:"rhythm" binding:"required"`
	TimeColumn      string     `json:"timeColumn" binding:"required"`
	AmplitudeColumn string     `json:"amplitudeColumn" binding:"required"`
	ExperimentNames []string   `json:"experimentNames" binding:"required"` // comma-separated in multipart
	// Filter parameters (all optional, all used in GROUP mode)
	FilterMin   *float64 `json:"filterMin"`   // bandpass min freq (Hz)
	FilterMax   *float64 `json:"filterMax"`   // bandpass max freq (Hz)
	FilterOrder *int     `json:"filterOrder"` // filter order (1-4)
	NPerSeg     *int     `json:"nPerSeg"`     // Welch window size
	NOverlap    *int     `json:"nOverlap"`    // Welch overlap
}

// Note: Validation is done in handlers since files come from multipart form
