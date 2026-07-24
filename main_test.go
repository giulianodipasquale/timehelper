package main

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestParseUnix(t *testing.T) {
	tests := []struct {
		name     string
		value    int64
		raw      string
		unit     string
		expected time.Time
		wantErr  bool
	}{
		{"seconds", 1721845200, "1721845200", "s", time.Unix(1721845200, 0), false},
		{"milliseconds", 1721845200000, "1721845200000", "ms", time.UnixMilli(1721845200000), false},
		{"microseconds", 1721845200000000, "1721845200000000", "us", time.UnixMicro(1721845200000000), false},
		{"nanoseconds", 1721845200000000000, "1721845200000000000", "ns", time.Unix(0, 1721845200000000000), false},
		{"auto seconds", 1721845200, "1721845200", "auto", time.Unix(1721845200, 0), false},
		{"auto milliseconds", 1721845200000, "1721845200000", "auto", time.UnixMilli(1721845200000), false},
		{"auto microseconds", 1721845200000000, "1721845200000000", "auto", time.UnixMicro(1721845200000000), false},
		{"auto nanoseconds", 1721845200000000000, "1721845200000000000", "auto", time.Unix(0, 1721845200000000000), false},
		{"invalid unit", 1721845200, "1721845200", "seconds", time.Time{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseUnix(tt.value, tt.raw, tt.unit)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseUnix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.expected) {
				t.Errorf("parseUnix() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseDateTime(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{"RFC3339", "2026-07-23T18:00:00Z", time.Date(2026, 7, 23, 18, 0, 0, 0, time.UTC)},
		{"RFC3339Nano", "2026-07-23T18:00:00.123456789Z", time.Date(2026, 7, 23, 18, 0, 0, 123456789, time.UTC)},
		{"HPFDateTime", "20260723180000.123456", time.Date(2026, 7, 23, 18, 0, 0, 123456000, time.UTC)}, // Note: hpfDateTimeUTC is "20060102150405.000000"
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDateTime(tt.input)
			if err != nil {
				t.Errorf("parseDateTime() error = %v", err)
			}
			if !got.Equal(tt.expected) {
				t.Errorf("parseDateTime() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
		wantErr  bool
	}{
		{"now", "now", time.Now(), false}, // Not exactly equal but should not error
		{"unix-s", "1721845200", time.Unix(1721845200, 0), false},
		{"rfc3339", "2026-07-23T18:00:00Z", time.Date(2026, 7, 23, 18, 0, 0, 0, time.UTC), false},
		{"invalid", "not-a-date", time.Time{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseInput(tt.input, "auto")
			if (err != nil) != tt.wantErr {
				t.Errorf("parseInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.expected) && tt.name != "now" {
				t.Errorf("parseInput() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPrintFormats(t *testing.T) {
	timeVal := time.Date(2026, 7, 23, 18, 0, 0, 0, time.UTC)
	var buf bytes.Buffer
	printFormats(&buf, timeVal)

	output := buf.String()
	// Check for several key lines to ensure they are present
	expectedLines := []string{
		"Unix:",
		"Unix Milli:",
		"Unix Micro:",
		"Unix Nano:",
		"RFC3339Nano:",
		"RFC3339Nano UTC:",
		"HPFDateTime UTC:",
	}

	for _, line := range expectedLines {
		if !strings.Contains(output, line) {
			t.Errorf("Expected output to contain %q, but it didn't", line)
		}
	}
}
