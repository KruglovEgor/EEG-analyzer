#!/usr/bin/env python3
"""Generate large realistic EEG test files for stress testing"""
import numpy as np
import csv

def generate_eeg_signal(duration_sec, sampling_rate, dominant_freq, noise_level=0.3):
    """Generate synthetic EEG signal with dominant frequency"""
    n_samples = int(duration_sec * sampling_rate)
    time = np.linspace(0, duration_sec, n_samples)
    
    # Base signal with dominant frequency
    signal = np.sin(2 * np.pi * dominant_freq * time)
    
    # Add harmonics
    signal += 0.3 * np.sin(2 * np.pi * (dominant_freq * 2) * time)
    signal += 0.15 * np.sin(2 * np.pi * (dominant_freq * 0.5) * time)
    
    # Add noise
    noise = np.random.normal(0, noise_level, n_samples)
    signal += noise
    
    # Scale to microvolts range
    signal = signal * 15 + 12  # Range approximately 0-30 µV
    
    return time, signal

def save_csv(filename, time, amplitude):
    """Save EEG data to CSV file"""
    with open(filename, 'w', newline='') as f:
        writer = csv.writer(f)
        writer.writerow(['Time', 'Amplitude'])
        for t, amp in zip(time, amplitude):
            writer.writerow([f'{t:.6f}', f'{amp:.6f}'])
    print(f"Created {filename}: {len(time)} samples, {time[-1]:.1f} seconds")

# Configuration
sampling_rate = 250  # Hz (typical EEG sampling rate)
duration = 60  # seconds

print("Generating large test files for EEG analysis...")
print(f"Sampling rate: {sampling_rate} Hz")
print(f"Duration: {duration} seconds")
print(f"Total samples per file: {duration * sampling_rate}")
print()

# Generate files with different dominant rhythms for GROUP mode testing
# All files have same structure but different signals

# Subject 1 - Strong Alpha (10 Hz - relaxed state)
time, signal = generate_eeg_signal(duration, sampling_rate, 10.0, noise_level=0.2)
save_csv('large_subject1_alpha.csv', time, signal)

# Subject 2 - Moderate Alpha (9.5 Hz)
time, signal = generate_eeg_signal(duration, sampling_rate, 9.5, noise_level=0.3)
save_csv('large_subject2_alpha.csv', time, signal)

# Subject 3 - Weak Alpha (11 Hz)
time, signal = generate_eeg_signal(duration, sampling_rate, 11.0, noise_level=0.4)
save_csv('large_subject3_alpha.csv', time, signal)

# Subject 4 - Beta dominant (20 Hz - alert state)
time, signal = generate_eeg_signal(duration, sampling_rate, 20.0, noise_level=0.25)
save_csv('large_subject4_beta.csv', time, signal)

# Single file with mixed rhythms for SINGLE mode testing
print("\nGenerating mixed rhythm file for SINGLE mode testing...")
time = np.linspace(0, duration, duration * sampling_rate)
signal = np.zeros_like(time)

# Mix multiple rhythms
signal += 2.0 * np.sin(2 * np.pi * 10 * time)  # Alpha (10 Hz)
signal += 1.5 * np.sin(2 * np.pi * 6 * time)   # Theta (6 Hz)
signal += 1.0 * np.sin(2 * np.pi * 18 * time)  # Beta (18 Hz)
signal += 0.5 * np.sin(2 * np.pi * 2 * time)   # Delta (2 Hz)
signal += np.random.normal(0, 0.5, len(time))  # Noise

signal = signal * 3 + 15  # Scale to realistic range
save_csv('large_mixed_rhythms.csv', time, signal)

print("\n✓ All test files generated successfully!")
print("\nFiles created:")
print("  - large_subject1_alpha.csv (15,000 samples, 60s)")
print("  - large_subject2_alpha.csv (15,000 samples, 60s)")
print("  - large_subject3_alpha.csv (15,000 samples, 60s)")
print("  - large_subject4_beta.csv  (15,000 samples, 60s)")
print("  - large_mixed_rhythms.csv  (15,000 samples, 60s)")
print("\nUse these files to test:")
print("  GROUP mode: Use subject1-3 alpha files with ALPHA rhythm")
print("  SINGLE mode: Use large_mixed_rhythms.csv with multiple rhythms")
