# EEG Analyzer

This project consists of 2 parts:
- backend (mine)
- [frontend](https://github.com/Vad1mChK/CourseworkEEG) (Vadim)

---
# EEG Analyzer Backend

REST API for real-time EEG (electroencephalography) signal analysis with support for multiple rhythm bands and analysis modes.

## Features

- **Dual Analysis Modes**
  - SINGLE: Analyze one file across multiple rhythm bands
  - GROUP: Compare multiple files for a single rhythm band
  
- **8 Rhythm Bands**
  - DELTA (0.5-4 Hz) - Deep sleep
  - THETA (4-8 Hz) - Drowsiness, meditation
  - ALPHA (8-13 Hz) - Relaxed wakefulness
  - BETA (13-30 Hz) - Active thinking, concentration
  - GAMMA (30-100 Hz) - Higher cognitive functions
  - MU (8-13 Hz) - Motor inhibition
  - LAMBDA (4-8 Hz) - Visual scanning
  - KAPPA (8-13 Hz) - Alpha variant

- **Signal Processing**
  - DC offset removal
  - FFT pre-filter (0.5-40 Hz) for PSD computation
  - Welch's method for Power Spectral Density
  - Butterworth bandpass filtering (1st-4th order) for visualization
  - Configurable filter parameters (frequency, order, PSD segments)
  - LTTB downsampling (1000-2000 points)

- **Performance**
  - Gzip compression (70-90% size reduction)
  - Multipart file upload (no base64 overhead)
  - Efficient memory handling for large files

## Quick Start

### Run Locally

```bash
go mod download
go run main.go
```

Server starts on `http://localhost:3000`

### Run with Docker

```bash
docker-compose up --build
```

### Build Binary

```bash
go build -o eeg-analyzer
./eeg-analyzer
```

## API Endpoints

### Health Check
```
GET /health
```

### Analyze EEG Data
```
POST /analyze
Content-Type: multipart/form-data
```

**SINGLE Mode** (one file, multiple rhythms):
```bash
curl -X POST http://localhost:3000/analyze \
  -F "analysisId=uuid-here" \
  -F "analysisMode=SINGLE" \
  -F "file=@data.csv" \
  -F "experimentName=Test 1" \
  -F "timeColumn=Time" \
  -F "amplitudeColumn=Amplitude" \
  -F "rhythms=ALPHA,BETA,THETA"
```

**GROUP Mode** (multiple files, one rhythm):
```bash
curl -X POST http://localhost:3000/analyze \
  -F "analysisId=uuid-here" \
  -F "analysisMode=GROUP" \
  -F "files=@subject1.csv" \
  -F "files=@subject2.csv" \
  -F "files=@subject3.csv" \
  -F "experimentNames=Subject 1,Subject 2,Subject 3" \
  -F "timeColumn=Time" \
  -F "amplitudeColumn=Amplitude" \
  -F "rhythm=ALPHA"
```

### Swagger Documentation
```
GET /swagger/index.html
```

## CSV Format

Required columns (names configurable):
- Time column: timestamps in seconds
- Amplitude column: signal amplitude in µV

Example:
```csv
Time,Amplitude
0.000,12.45
0.004,15.23
0.008,18.67
...
```

## Testing

Test files are provided in `testdata/`:

**Small files** (~76 samples, 0.3s):
- `sample_eeg_alpha.csv`
- `sample_eeg_beta.csv`
- `sample_eeg_theta.csv`

**Large files** (15,000 samples, 60s):
- `large_mixed_rhythms.csv` - For SINGLE mode with multiple rhythms
- `large_subject1_alpha.csv` - For GROUP mode comparison
- `large_subject2_alpha.csv`
- `large_subject3_alpha.csv`

Generate new test files:
```bash
cd testdata
python generate_large_test_files.py
```

## Response Format

```json
{
  "analysisId": "uuid",
  "analysisMode": "SINGLE",
  "experimentName": "Test 1",
  "rhythms": ["ALPHA", "BETA"],
  "absolutePowers": [
    ["ALPHA", 12.45],
    ["BETA", 8.32]
  ],
  "relativePowers": [
    ["ALPHA", 60.0],
    ["BETA", 40.0]
  ],
  "dataByRhythm": {
    "ALPHA": {
      "psdPlot": {
        "data": [{"x": 0, "psd": 1.5}, ...],
        "yLogarithmic": true
      },
      "signalPlot": {
        "data": [{"x": 0, "raw": 12.5, "filtered": 11.2}, ...]
      }
    }
  }
}
```

## Environment Variables

- `PORT` - Server port (default: 3000)
- `GIN_MODE` - Gin mode: release/debug (default: release)

## Architecture

```
backend/
├── main.go              # Entry point, routing
├── handlers/            # HTTP handlers
│   ├── health.go
│   ├── analyze.go       # JSON endpoint (Swagger)
│   └── analyze_multipart.go  # Production endpoint
├── models/              # Data structures
│   ├── request.go
│   ├── response.go
│   └── types.go
├── analysis/            # Signal processing
│   ├── csv_parser.go
│   ├── filter.go
│   ├── fft.go
│   ├── rhythms.go
│   └── downsampler.go
├── docs/                # Swagger generated docs
└── testdata/            # Test CSV files
```

## Signal Processing

The backend implements a dual-path processing pipeline:

1. **Analysis Path**: FFT pre-filter (0.5-40 Hz) → Welch PSD → Power extraction
2. **Visualization Path**: Butterworth filter → Downsampling → Display

See [detailed signal processing documentation](../temp/SIGNAL_PROCESSING.md) for mathematical explanations.

**Key Features:**
- Configurable Butterworth filter order (1-4)
- Adjustable Welch PSD parameters (nperseg, noverlap)
- Automatic handling of short signals
- Single-point PSD approximation for low-resolution data

## Dependencies

- **gin-gonic/gin** - Web framework
- **gin-contrib/cors** - CORS middleware
- **gin-contrib/gzip** - Gzip compression
- **mjibson/go-dsp** - FFT implementation
- **swaggo/gin-swagger** - API documentation

## Frontend

Frontend repository: [EEG-analyzer/eeg_frontend](https://github.com/vad1mchk/CourseworkEEG)

Live demo: https://vad1mchk.github.io/CourseworkEEG/
