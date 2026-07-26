package hko

import (
	"testing"
	"time"
)

func TestClientRealWeather(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live API test")
	}
	c := NewClient(10 * time.Second)

	w, err := c.GetCurrentWeather("en")
	if err != nil {
		t.Fatalf("get current weather: %v", err)
	}
	if w.UpdateTime == "" {
		t.Error("expected update time")
	}
	if len(w.Temperature.Data) == 0 {
		t.Error("expected temperature data")
	}
}

func TestClientRealForecast(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live API test")
	}
	c := NewClient(10 * time.Second)

	f, err := c.GetForecast("en")
	if err != nil {
		t.Fatalf("get forecast: %v", err)
	}
	if len(f.WeatherForecast) == 0 {
		t.Error("expected forecast data")
	}
}

func TestClientRealTides(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live API test")
	}
	c := NewClient(10 * time.Second)

	tides, err := c.GetTides("QUB", "en")
	if err != nil {
		t.Fatalf("get tides: %v", err)
	}
	if len(tides.Fields) == 0 {
		t.Error("expected fields")
	}
	if len(tides.Records()) == 0 {
		t.Error("expected tide records")
	}
}

func TestClientRealLunar(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live API test")
	}
	c := NewClient(10 * time.Second)

	lunar, err := c.GetTodayLunarCalendar()
	if err != nil {
		t.Fatalf("get lunar: %v", err)
	}
	if lunar.LunarDate == "" {
		t.Error("expected lunar date")
	}
}
