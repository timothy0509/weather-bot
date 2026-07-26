package hko

import (
	"encoding/json"
	"testing"
)

func TestEarthquakeResponseSingle(t *testing.T) {
	payload := `{
		"lat": -14.04,
		"lon": 166.66,
		"mag": 6,
		"region": "near Vanuatu Islands",
		"ptime": "2026-07-25T05:37:00+08:00",
		"updateTime": "2026-07-25T05:48:00+08:00"
	}`

	var res EarthquakeResponse
	if err := json.Unmarshal([]byte(payload), &res); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	results := res.Results()
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Mag != 6 {
		t.Errorf("mag = %v, want 6", results[0].Mag)
	}
}

func TestEarthquakeResponseEmpty(t *testing.T) {
	payload := `{}`
	var res EarthquakeResponse
	if err := json.Unmarshal([]byte(payload), &res); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(res.Results()) != 0 {
		t.Errorf("expected 0 results, got %d", len(res.Results()))
	}
}
