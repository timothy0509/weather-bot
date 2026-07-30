package format

import (
	"fmt"
	"math"
	"time"
)

func FormatTime(iso string) string {
	t, err := time.Parse(time.RFC3339, iso)
	if err != nil {
		return iso
	}
	return t.Format("15:04")
}

func FormatDateTime(iso string) string {
	t, err := time.Parse(time.RFC3339, iso)
	if err != nil {
		return iso
	}
	return t.Format("Jan 2, 15:04")
}

func FormatWeekdayDate(yyyymmdd string) string {
	t, err := time.Parse("20060102", yyyymmdd)
	if err != nil {
		return yyyymmdd
	}
	return t.Format("Mon 2 Jan")
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

func FormatRelative(iso string) string {
	t, err := time.Parse(time.RFC3339, iso)
	if err != nil {
		return iso
	}
	d := time.Since(t)
	if d < 0 {
		return FormatDateTime(iso)
	}
	if d < time.Minute {
		return "just now"
	}
	if d < time.Hour {
		mins := int(math.Round(d.Minutes()))
		return fmt.Sprintf("%d min ago", mins)
	}
	if d < 48*time.Hour {
		hours := int(math.Round(d.Hours()))
		return fmt.Sprintf("%d hours ago", hours)
	}
	return FormatDateTime(iso)
}

func FormatDate(yyyymmdd string) string {
	t, err := time.Parse("20060102", yyyymmdd)
	if err != nil {
		return yyyymmdd
	}
	return t.Format("Jan 2")
}
