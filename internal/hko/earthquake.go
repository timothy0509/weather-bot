package hko

import (
	"fmt"
	"net/url"
	"time"
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
	if err := c.Get(u.String(), &res); err != nil {
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

// FeltEarthquakeInfo is a felt earthquake report.
type FeltEarthquakeInfo struct {
	Mag        float64 `json:"mag"`
	Region     string  `json:"region"`
	Intensity  string  `json:"intensity"`
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"`
	Details    string  `json:"details"`
	PTime      string  `json:"ptime"`
	UpdateTime string  `json:"updateTime"`
}

// FeltEarthquakeResponse is the response from dataType=feltearthquake.
type FeltEarthquakeResponse struct {
	Data []FeltEarthquakeInfo `json:"data"`
}

// GetFeltEarthquakeInfo fetches locally felt earthquakes.
func (c *Client) GetFeltEarthquakeInfo(lang string) (*FeltEarthquakeResponse, error) {
	u, _ := url.Parse(baseEarthquakeURL)
	q := u.Query()
	q.Set("dataType", "feltearthquake")
	q.Set("lang", languageCode(lang))
	u.RawQuery = q.Encode()

	var res FeltEarthquakeResponse
	if err := c.Get(u.String(), &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// ValidateTime validates a time for earthquake API.
func (e *EarthquakeInfo) Time() time.Time {
	t, _ := time.Parse(time.RFC3339, e.PTime)
	return t
}

// FormatMag formats the magnitude.
func (e *EarthquakeInfo) FormatMag() string {
	return fmt.Sprintf("%.1f", e.Mag)
}
