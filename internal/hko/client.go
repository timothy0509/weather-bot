// Package hko provides a client for the Hong Kong Observatory Open Data API.
package hko

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

const (
	baseWeatherURL    = "https://data.weather.gov.hk/weatherAPI/opendata/weather.php"
	baseEarthquakeURL = "https://data.weather.gov.hk/weatherAPI/opendata/earthquake.php"
	baseOpenDataURL   = "https://data.weather.gov.hk/weatherAPI/opendata/opendata.php"
	baseLunarURL      = "https://data.weather.gov.hk/weatherAPI/opendata/lunardate.php"

	maxBodySize  = 5 * 1024 * 1024
	maxRetries   = 3
	retryBaseMs  = 500
)

type cacheEntry struct {
	data      []byte
	expiresAt time.Time
}

type responseCache struct {
	mu      sync.RWMutex
	entries map[string]cacheEntry
}

func (c *responseCache) get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[key]
	if !ok || time.Now().After(e.expiresAt) {
		return nil, false
	}
	cp := make([]byte, len(e.data))
	copy(cp, e.data)
	return cp, true
}

func (c *responseCache) set(key string, data []byte, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.entries == nil {
		c.entries = make(map[string]cacheEntry)
	}
	cp := make([]byte, len(data))
	copy(cp, data)
	c.entries[key] = cacheEntry{data: cp, expiresAt: time.Now().Add(ttl)}
}

// Client performs HKO API requests.
type Client struct {
	httpClient *http.Client
	cache      *responseCache
	flight     singleflight.Group
	logger     *slog.Logger
}

// NewClient creates a new HKO API client.
func NewClient(timeout time.Duration) *Client {
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &Client{
		httpClient: &http.Client{Timeout: timeout},
		cache:      &responseCache{},
	}
}

// SetLogger sets the logger for the client.
func (c *Client) SetLogger(l *slog.Logger) {
	c.logger = l
}

// Get requests a JSON endpoint and decodes it into v.
func (c *Client) Get(targetURL string, v interface{}) error {
	return c.GetWithTTL(targetURL, v, 0)
}

// GetWithTTL requests a JSON endpoint with response caching.
func (c *Client) GetWithTTL(targetURL string, v interface{}, ttl time.Duration) error {
	if ttl > 0 {
		if cached, ok := c.cache.get(targetURL); ok {
			return json.Unmarshal(cached, v)
		}
	}

	result, err, _ := c.flight.Do(targetURL, func() (interface{}, error) {
		if ttl > 0 {
			if cached, ok := c.cache.get(targetURL); ok {
				return cached, nil
			}
		}

		body, err := c.doGet(targetURL)
		if err != nil {
			return nil, err
		}

		if ttl > 0 {
			c.cache.set(targetURL, body, ttl)
		}
		return body, nil
	})
	if err != nil {
		return err
	}

	raw := result.([]byte)
	return json.Unmarshal(raw, v)
}

func (c *Client) doGet(targetURL string) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(retryBaseMs*(1<<(attempt-1))) * time.Millisecond
			time.Sleep(backoff)
		}

		body, err := c.doRequest(targetURL)
		if err == nil {
			return body, nil
		}
		lastErr = err
		if !isRetryable(err) {
			return nil, err
		}
		if c.logger != nil {
			c.logger.Warn("HKO request failed, retrying",
				slog.String("url", targetURL),
				slog.Int("attempt", attempt+1),
				slog.Any("err", err))
		}
	}
	return nil, fmt.Errorf("after %d retries: %w", maxRetries, lastErr)
}

func (c *Client) doRequest(targetURL string) ([]byte, error) {
	resp, err := c.httpClient.Get(targetURL)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodySize))
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &httpError{statusCode: resp.StatusCode, body: string(body)}
	}

	return body, nil
}

type httpError struct {
	statusCode int
	body       string
}

func (e *httpError) Error() string {
	return fmt.Sprintf("unexpected status %d: %s", e.statusCode, e.body)
}

func isRetryable(err error) bool {
	var he *httpError
	if errors.As(err, &he) {
		return he.statusCode >= 500
	}
	return true
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
	Unit string            `json:"unit"`
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
