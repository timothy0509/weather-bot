package hko

import "testing"

func TestValidTideStations(t *testing.T) {
	for _, s := range ValidTideStations() {
		if !IsValidTideStation(s) {
			t.Errorf("expected %s to be a valid station", s)
		}
	}
	if IsValidTideStation("XYZ") {
		t.Error("expected XYZ to be invalid")
	}
}

func TestReadingByPlace(t *testing.T) {
	rg := ReadingGroup{
		Data: []Reading{
			{Place: "HKO", Value: 28, Unit: "C"},
			{Place: "King's Park", Value: 27, Unit: "C"},
		},
	}
	r, ok := ReadingByPlace(rg, "HKO")
	if !ok || r.Value != 28 {
		t.Errorf("expected HKO reading 28, got %v, ok=%v", r, ok)
	}
	_, ok = ReadingByPlace(rg, "Missing")
	if ok {
		t.Error("expected missing place to not be found")
	}
}

func TestFormatTemperature(t *testing.T) {
	r := Reading{Value: 28.5, Unit: "C"}
	if got := r.FormatTemperature(); got != "28.5°C" {
		t.Errorf("FormatTemperature() = %q, want 28.5°C", got)
	}
}

func TestMaxRainfallReading(t *testing.T) {
	g := RainfallGroup{
		Data: []RainfallReading{
			{Place: "A", Max: 1, Unit: "mm"},
			{Place: "B", Max: 5, Unit: "mm"},
			{Place: "C", Max: 2, Unit: "mm"},
		},
	}
	max, ok := MaxRainfallReading(g)
	if !ok || max.Place != "B" {
		t.Errorf("expected B to be max, got %v, ok=%v", max, ok)
	}
}
