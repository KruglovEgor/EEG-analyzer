# Test Files

## Overview

This directory contains test CSV files for EEG analysis validation.

## Small Test Files (~1.3 KB, 76 samples, 0.3s)

Basic validation files:
- `sample_eeg_alpha.csv` - Alpha rhythm dominant (10 Hz)
- `sample_eeg_beta.csv` - Beta rhythm dominant (20 Hz)
- `sample_eeg_theta.csv` - Theta rhythm dominant (6 Hz)

**Use case**: Quick API validation and Swagger testing

## Large Test Files (~309 KB, 15,000 samples, 60s)

Realistic stress test files:

### SINGLE Mode Testing
- `large_mixed_rhythms.csv` - Multiple rhythms mixed (ALPHA, BETA, THETA, DELTA)

**Example**:
```bash
curl -X POST http://localhost:3000/analyze \
  -F "analysisId=test-001" \
  -F "analysisMode=SINGLE" \
  -F "file=@large_mixed_rhythms.csv" \
  -F "experimentName=Mixed Rhythms Test" \
  -F "timeColumn=Time" \
  -F "amplitudeColumn=Amplitude" \
  -F "rhythms=ALPHA,BETA,THETA,DELTA"
```

### GROUP Mode Testing
- `large_subject1_alpha.csv` - Subject 1 (strong Alpha, 10 Hz)
- `large_subject2_alpha.csv` - Subject 2 (moderate Alpha, 9.5 Hz)
- `large_subject3_alpha.csv` - Subject 3 (weak Alpha, 11 Hz)
- `large_subject4_beta.csv` - Subject 4 (Beta dominant, 20 Hz)

**Example**:
```bash
curl -X POST http://localhost:3000/analyze \
  -F "analysisId=test-002" \
  -F "analysisMode=GROUP" \
  -F "files=@large_subject1_alpha.csv" \
  -F "files=@large_subject2_alpha.csv" \
  -F "files=@large_subject3_alpha.csv" \
  -F "experimentNames=Subject 1,Subject 2,Subject 3" \
  -F "timeColumn=Time" \
  -F "amplitudeColumn=Amplitude" \
  -F "rhythm=ALPHA"
```

## Generate New Files

```bash
python generate_large_test_files.py
```

Customize parameters in the script:
- `sampling_rate` - Hz (default: 250)
- `duration` - seconds (default: 60)
- `dominant_freq` - Hz for each rhythm
- `noise_level` - signal noise (default: 0.2-0.4)

## Expected Results

**Performance**:
- Small files: ~10ms processing time
- Large files: ~80-100ms processing time
- Response gzip compressed: ~260KB (from ~800KB uncompressed)

**Validation**:
- All PSD values > 0 (logarithmic scale compatible)
- Signal plots: 1000-2000 points (downsampled from 15,000)
- Absolute powers: positive values in µV²/Hz
- Relative powers: percentages summing to 100%
