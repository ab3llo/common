package utils

import (
	"testing"
	"time"
)

func TestFormatTime(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "zero time",
			input:    time.Time{},
			expected: "",
		},
		{
			name:     "valid time",
			input:    time.Date(2023, 12, 15, 14, 30, 45, 0, time.UTC),
			expected: "2023-12-15T14:30:45Z",
		},
		{
			name:     "time with local conversion",
			input:    time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			expected: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC).Local().Format(time.RFC3339),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatTime(tt.input)
			if result != tt.expected {
				t.Errorf("FormatTime() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func BenchmarkFormatTime(b *testing.B) {
	t := time.Date(2023, 12, 15, 14, 30, 45, 0, time.UTC)
	for i := 0; i < b.N; i++ {
		FormatTime(t)
	}
}

func TestFormatTime_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "distant past",
			input:    time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339),
		},
		{
			name:     "distant future",
			input:    time.Date(2030, 12, 31, 23, 59, 59, 0, time.UTC),
			expected: time.Date(2030, 12, 31, 23, 59, 59, 0, time.UTC).Local().Format(time.RFC3339),
		},
		{
			name:     "leap year",
			input:    time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC),
			expected: time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC).Local().Format(time.RFC3339),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatTime(tt.input)
			if result != tt.expected {
				t.Errorf("FormatTime() = %v, want %v", result, tt.expected)
			}
		})
	}
}