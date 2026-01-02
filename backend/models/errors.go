package models

import "errors"

var (
	ErrMissingFile       = errors.New("file is required for SINGLE mode")
	ErrMissingFiles      = errors.New("files are required for GROUP mode")
	ErrMissingRhythms    = errors.New("rhythms are required for SINGLE mode")
	ErrMissingRhythm     = errors.New("rhythm is required for GROUP mode")
	ErrInvalidCSV        = errors.New("invalid CSV format")
	ErrColumnNotFound    = errors.New("specified column not found in CSV")
	ErrInsufficientData  = errors.New("insufficient data points for analysis")
	ErrInvalidRhythmBand = errors.New("invalid rhythm band frequency range")
)
