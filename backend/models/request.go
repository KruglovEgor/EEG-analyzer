package models

// EEGFileConfig contains metadata about a single EEG file
type EEGFileConfig struct {
	ID              string  `json:"id" binding:"required"`
	Filename        string  `json:"filename" binding:"required"`
	ExperimentName  string  `json:"experimentName" binding:"required"`
	TimeColumn      string  `json:"timeColumn" binding:"required"`
	AmplitudeColumn string  `json:"amplitudeColumn" binding:"required"`
	RawFile         *string `json:"rawFile"`  // Base64 encoded CSV content
	ServerID        *string `json:"serverId"` // Can be ignored
}

// EEGAnalysisRequest is the base request structure
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
