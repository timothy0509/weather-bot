// Package hko provides a client for the Hong Kong Observatory Open Data API.
package hko

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	baseWeatherURL = "https://data.weather.gov.hk/weatherAPI/opendata/weather.php"
	baseEarthquakeURL = "https://data.weather.gov.hk/weatherAPI/opendata/earthquake.php"
	baseOpenDataURL = "https://data.weather.gov.hk/weatherAPI/opendata/opendata.php"
	baseLunarURL = "https://data.weather.gov.hk/weatherAPI/opendata/lunardate.php"
)

// Client performs HKO API requests.
type Client struct {
	httpClient *http.Client
}

// NewClient creates a new HKO API client.
func NewClient(timeout time.Duration) *Client {
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &Client{
		httpClient: &http.Client{Timeout: timeout},
	}
}

// Get requests a JSON endpoint and decodes it into v.
func (c *Client) Get(targetURL string, v interface{}) error {
	resp, err := c.httpClient.Get(targetURL)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("decode json: %w", err)
	}
	return nil
}

// buildWeatherURL builds a URL for the weather API endpoint.
func buildWeatherURL(dataType, lang string) string {
	u, _ := url.Parse(baseWeatherURL)
	q := u.Query()
	q.Set("dataType", dataType)
	q.Set("lang", lang)
	u.RawQuery = q.Encode()
	return u.String()
}

// languageCode normalizes language input to HKO-supported values.
func languageCode(lang string) string {
	switch lang {
	case "en", "tc", "sc":
		return lang
	case "bilingual":
		return "en"
	default:
		return "en"
	}
}

// secondaryLanguage returns the second language for bilingual mode.
func secondaryLanguage(lang string) string {
	if lang == "bilingual" {
		return "tc"
	}
	return ""
}

// StringValue represents a value that can be either a string or number.
type StringValue struct {
	Value string `json:"value"`
	Unit  string `json:"unit"`
}

// NumberValue represents a numeric value with unit.
type NumberValue struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

// Reading is a place/value/unit reading.
type Reading struct {
	Place string  `json:"place"`
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
	Max   float64 `json:"max,omitempty"`
	Min   float64 `json:"min,omitempty"`
}

// PlaceValue is a simple place/value pair.
type PlaceValue struct {
	Place string `json:"place"`
	Value string `json:"value"`
	Unit  string `json:"unit,omitempty"`
}

// ReadingGroup groups readings by record time.
type ReadingGroup struct {
	RecordTime string    `json:"recordTime"`
	Data       []Reading `json:"data"`
}

// RainfallGroup groups rainfall data.
type RainfallGroup struct {
	Unit string          `json:"unit"`
	Data []RainfallReading `json:"data"`
}

// RainfallReading is a rainfall reading.
type RainfallReading struct {
	Place string  `json:"place"`
	Max   float64 `json:"max"`
	Min   float64 `json:"min,omitempty"`
	Unit  string  `json:"unit"`
	Main  string  `json:"main,omitempty"`
}
