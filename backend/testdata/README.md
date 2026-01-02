# EEG Test Data

This directory contains sample EEG data files for testing the analyzer.

## Files

- **sample_eeg_alpha.csv** - Simulated alpha wave activity (8-13 Hz)
  - Columns: Time, Amplitude, Channel
  - Electrode: Fp1 (frontal)
  - Sampling rate: ~250 Hz

- **sample_eeg_beta.csv** - Simulated beta wave activity (13-30 Hz)
  - Columns: Time, Amplitude, Channel
  - Electrode: Oz (occipital)
  - Sampling rate: ~250 Hz

- **sample_eeg_theta.csv** - Simulated theta wave activity (4-8 Hz)
  - Columns: Time, Signal, Electrode
  - Electrode: C3 (central)
  - Sampling rate: ~250 Hz

## Usage

When testing the API, encode these CSV files as base64 and include in the `rawFile` field of the request.

Example in Python:
```python
import base64

with open('sample_eeg_alpha.csv', 'rb') as f:
    encoded = base64.b64encode(f.read()).decode('utf-8')
```

Note: Column names vary to test the dynamic column detection feature.
