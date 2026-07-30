package hko

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// TideResponse is the response from dataType=HHOT.
// HKO returns a tabular format: fields=["MM","DD","01","02",...] and data=[rows].
type TideResponse struct {
	Fields []string   `json:"fields"`
	Data   [][]string `json:"data"`
}

// TideRecord is a parsed hourly tide record.
type TideRecord struct {
	Time   string
	Hour   int
	Height float64
}

// Records parses the tabular data into a slice of records for today.
func (t *TideResponse) Records() []TideRecord {
	if len(t.Fields) < 3 || len(t.Data) == 0 {
		return nil
	}

	var records []TideRecord
	row := t.Data[0]
	month := row[0]
	day := row[1]
	for i, field := range t.Fields[2:] {
		if i+2 >= len(row) {
			break
		}
		value := row[i+2]
		if value == "" || value == "-" || value == "N/A" {
			continue
		}
		h, err := strconv.ParseFloat(value, 64)
		if err != nil {
			continue
		}
		hour, _ := strconv.Atoi(field)
		records = append(records, TideRecord{
			Time:   fmt.Sprintf("%s-%s %02d:00", month, day, hour),
			Hour:   hour,
			Height: h,
		})
	}
	return records
}

// GetTides fetches hourly astronomical tide heights for a station.
func (c *Client) GetTides(station, lang string) (*TideResponse, error) {
	if station == "" {
		station = "QUB"
	}
	u, _ := url.Parse(baseOpenDataURL)
	q := u.Query()
	q.Set("dataType", "HHOT")
	q.Set("station", station)
	q.Set("lang", languageCode(lang))
	q.Set("rformat", "json")
	now := time.Now()
	q.Set("year", fmt.Sprintf("%d", now.Year()))
	q.Set("month", fmt.Sprintf("%d", int(now.Month())))
	q.Set("day", fmt.Sprintf("%d", now.Day()))
	u.RawQuery = q.Encode()

	var res TideResponse
	if err := c.GetWithTTL(u.String(), &res, ttlTides); err != nil {
		return nil, err
	}
	return &res, nil
}

// LunarResponse is the response from the lunar calendar API.
type LunarResponse struct {
	LunarYear string `json:"LunarYear"`
	LunarDate string `json:"LunarDate"`
}

// GetLunarCalendar fetches the lunar calendar entry for a date.
func (c *Client) GetLunarCalendar(date string) (*LunarResponse, error) {
	u, _ := url.Parse(baseLunarURL)
	q := u.Query()
	q.Set("date", date)
	u.RawQuery = q.Encode()

	var res LunarResponse
	if err := c.GetWithTTL(u.String(), &res, ttlLunar); err != nil {
		return nil, err
	}
	return &res, nil
}

// GetTodayLunarCalendar fetches today's lunar calendar entry.
func (c *Client) GetTodayLunarCalendar() (*LunarResponse, error) {
	return c.GetLunarCalendar(time.Now().Format("2006-01-02"))
}

// ValidTideStations returns the list of supported tide station codes.
func ValidTideStations() []string {
	return []string{
		"CCH", "CLK", "CMW", "KCT", "KLW", "LOP", "MWC", "QUB",
		"SPW", "TAO", "TBT", "TMW", "TPK", "WAG",
	}
}

// IsValidTideStation checks if a station code is valid.
func IsValidTideStation(station string) bool {
	for _, s := range ValidTideStations() {
		if s == station {
			return true
		}
	}
	return false
}
