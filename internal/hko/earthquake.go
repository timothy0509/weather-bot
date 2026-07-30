package hko

import (
	"fmt"
	"net/url"
	"time"
)

const (
	ttlEarthquake = 60 * time.Second
	ttlTides      = 1 * time.Hour
	ttlLunar      = 24 * time.Hour
)

// EarthquakeInfo is a single earthquake report.
type EarthquakeInfo struct {
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"`
	Mag        float64 `json:"mag"`
	Region     string  `json:"region"`
	PTime      string  `json:"ptime"`
	UpdateTime string  `json:"updateTime"`
}

// EarthquakeResponse is the response from dataType=qem.
// The API returns either a single object or a single-element array.
type EarthquakeResponse struct {
	EarthquakeInfo
	Data []EarthquakeInfo `json:"data"`
}

// GetEarthquakeInfo fetches quick earthquake messages.
func (c *Client) GetEarthquakeInfo(lang string) (*EarthquakeResponse, error) {
	u, _ := url.Parse(baseEarthquakeURL)
	q := u.Query()
	q.Set("dataType", "qem")
	q.Set("lang", languageCode(lang))
	u.RawQuery = q.Encode()

	var res EarthquakeResponse
	if err := c.GetWithTTL(u.String(), &res, ttlEarthquake); err != nil {
		return nil, err
	}
	return &res, nil
}

// Results returns the earthquake records, handling both single-object and array forms.
func (e *EarthquakeResponse) Results() []EarthquakeInfo {
	if len(e.Data) > 0 {
		return e.Data
	}
	if e.Region != "" || e.Mag != 0 {
		return []EarthquakeInfo{e.EarthquakeInfo}
	}
	return nil
}

// Time parses the earthquake time.
func (e *EarthquakeInfo) Time() (time.Time, error) {
	return time.Parse(time.RFC3339, e.PTime)
}

// FormatMag formats the magnitude.
func (e *EarthquakeInfo) FormatMag() string {
	return fmt.Sprintf("%.1f", e.Mag)
}
