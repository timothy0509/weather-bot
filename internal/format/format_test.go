package format

import (
	"testing"
)

func TestFormatTime(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"2026-07-30T11:25:00+08:00", "11:25"},
		{"2026-01-01T00:00:00+08:00", "00:00"},
		{"2026-12-31T23:59:00+08:00", "23:59"},
		{"invalid", "invalid"},
		{"", ""},
	}
	for _, tt := range tests {
		result := FormatTime(tt.input)
		if result != tt.expected {
			t.Errorf("FormatTime(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestFormatDateTime(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"2026-07-30T11:25:00+08:00", "Jul 30, 11:25"},
		{"2026-01-01T00:00:00+08:00", "Jan 1, 00:00"},
		{"invalid", "invalid"},
	}
	for _, tt := range tests {
		result := FormatDateTime(tt.input)
		if result != tt.expected {
			t.Errorf("FormatDateTime(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestFormatWeekdayDate(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"20260728", "Tue 28 Jul"},
		{"20260101", "Thu 1 Jan"},
		{"invalid", "invalid"},
	}
	for _, tt := range tests {
		result := FormatWeekdayDate(tt.input)
		if result != tt.expected {
			t.Errorf("FormatWeekdayDate(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestFormatCoordinates(t *testing.T) {
	tests := []struct {
		lat, lon float64
		expected string
	}{
		{22.3, 114.1, "22.3°N, 114.1°E"},
		{-33.9, 151.2, "33.9°S, 151.2°E"},
		{40.7, -74.0, "40.7°N, 74.0°W"},
		{0.0, 0.0, "0.0°N, 0.0°E"},
	}
	for _, tt := range tests {
		result := FormatCoordinates(tt.lat, tt.lon)
		if result != tt.expected {
			t.Errorf("FormatCoordinates(%v, %v) = %q, want %q", tt.lat, tt.lon, result, tt.expected)
		}
	}
}

func TestFormatDate(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"20260730", "Jul 30"},
		{"20260101", "Jan 1"},
		{"invalid", "invalid"},
	}
	for _, tt := range tests {
		result := FormatDate(tt.input)
		if result != tt.expected {
			t.Errorf("FormatDate(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
