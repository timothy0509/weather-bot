package format

import (
	"fmt"
	"time"
)

var hkt = time.FixedZone("HKT", 8*60*60)

func FormatTime(iso string) string {
	t, err := time.Parse(time.RFC3339, iso)
	if err != nil {
		return iso
	}
	return t.In(hkt).Format("15:04")
}

func FormatDateTime(iso string) string {
	t, err := time.Parse(time.RFC3339, iso)
	if err != nil {
		return iso
	}
	return t.In(hkt).Format("Jan 2, 15:04")
}

func FormatWeekdayDate(yyyymmdd string) string {
	t, err := time.Parse("20060102", yyyymmdd)
	if err != nil {
		return yyyymmdd
	}
	return t.In(hkt).Format("Mon 2 Jan")
}

func FormatCoordinates(lat, lon float64) string {
	ns := "N"
	if lat < 0 {
		ns = "S"
		lat = -lat
	}
	ew := "E"
	if lon < 0 {
		ew = "W"
		lon = -lon
	}
	return fmt.Sprintf("%.1f°%s, %.1f°%s", lat, ns, lon, ew)
}

func TruncateRunes(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max-3]) + "..."
}
